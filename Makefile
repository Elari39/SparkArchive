.PHONY: help install dev-web dev-api ingest build test docker docker-up docker-down smoke clean

PORT ?= 12026
CONTENT ?= content

help: ## 显示可用命令
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}'

install: ## 安装前端依赖
	cd web && pnpm install

dev-web: ## 启动前端开发服务器（/api 代理到 12026）
	cd web && pnpm dev

ingest: ## 从 content/ 重建数据库
	@mkdir -p server/data
	cd server && go run ./cmd/archive ingest -content ../$(CONTENT) -db data/archive.db

dev-api: ## 启动后端（需先 ingest）
	cd server && go run ./cmd/archive -addr :$(PORT) -db data/archive.db -site ../$(CONTENT)/site.yaml

build: ## 构建前端并内嵌到后端，产出单二进制
	cd web && pnpm install && pnpm build
	@rm -rf server/cmd/archive/webdist/assets
	@mkdir -p server/cmd/archive/webdist
	@cp -r web/dist/* server/cmd/archive/webdist/
	cd server && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o archive ./cmd/archive

test: ## 后端 vet + 测试，前端类型检查
	cd server && go vet ./... && go test ./...
	cd web && pnpm type-check

docker: ## 构建 Docker 镜像
	docker build -t sparkarchive:latest .

docker-up: ## 一键部署（构建并启动）
	docker compose up -d --build

docker-down: ## 停止服务
	docker compose down

smoke: ## 端到端冒烟测试（需已 docker compose up）
	pwsh -File scripts/smoke.ps1

contract: ## 校验前端 zod schema 与后端响应是否匹配（需已 docker compose up）
	cd web && pnpm contract-check

clean: ## 清理构建产物
	@rm -rf web/dist server/data server/archive server/archive.exe
	@rm -rf server/cmd/archive/webdist/assets
	@rm -f server/cmd/archive/webdist/index.html
