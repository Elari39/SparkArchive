# syntax=docker/dockerfile:1.7
# ---------- 前端构建 ----------
FROM node:24-alpine AS web
WORKDIR /src/web
# npm registry 在部分网络下较慢且易 ECONNRESET，放宽超时并多重试几次
ENV npm_config_fetch_timeout=600000 \
    npm_config_fetch_retries=10 \
    npm_config_fetch_retry_mintimeout=5000 \
    npm_config_fetch_retry_maxtimeout=120000
RUN corepack enable
COPY web/package.json web/pnpm-lock.yaml ./
# 锁定 lockfile 保证可重复构建；网络抖动时最多重试 3 次
RUN --mount=type=cache,target=/root/.local/share/pnpm/store \
    for attempt in 1 2 3; do \
        pnpm install --frozen-lockfile --prefer-offline && exit 0; \
        echo "pnpm install 第 $attempt 次失败，5 秒后重试…"; \
        sleep 5; \
    done; \
    exit 1
COPY web/ ./
RUN pnpm run build

# ---------- 后端构建 ----------
FROM golang:1.27-alpine AS build
WORKDIR /src/server
# proxy.golang.org 在部分网络下不可达，改用国内镜像；go.sum 已锁定校验和
ENV GOPROXY=https://goproxy.cn,direct \
    GOSUMDB=sum.golang.google.cn \
    CGO_ENABLED=0
COPY server/go.mod server/go.sum ./
RUN --mount=type=cache,target=/root/go/pkg/mod \
    go mod download
COPY server/ ./
# 把前端产物放进 embed 目录
COPY --from=web /src/web/dist ./cmd/archive/webdist/
RUN --mount=type=cache,target=/root/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -ldflags="-s -w" -o /out/archive ./cmd/archive

# 构建阶段顺带产出数据库：内容在构建期编译，运行时只读
COPY content/ /content/
RUN /out/archive ingest -content /content -db /out/archive.db

# ---------- 运行时 ----------
FROM alpine:3.22
RUN apk add --no-cache ca-certificates tzdata wget && \
    adduser -D -u 10001 archive
WORKDIR /app
COPY --from=build /out/archive /app/archive
COPY --from=build /out/archive.db /data/archive.db
COPY content/site.yaml /app/content/site.yaml
USER archive
ENV ADDR=:12026 \
    DB_PATH=/data/archive.db \
    SITE_PATH=/app/content/site.yaml \
    LOG_LEVEL=info
EXPOSE 12026
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -qO- http://127.0.0.1:12026/api/health >/dev/null 2>&1 || exit 1
ENTRYPOINT ["/app/archive"]
