# 星火档案馆 · SPARK ARCHIVE

> 全世界无产者，联合起来！

[![CI](https://github.com/Elari39/SparkArchive/actions/workflows/ci.yml/badge.svg)](https://github.com/Elari39/SparkArchive/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-d62828.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.27-00ADD8.svg)](server/go.mod)
[![Vue](https://img.shields.io/badge/Vue-3.5-42b883.svg)](web/package.json)
[![Image](https://img.shields.io/badge/ghcr.io-sparkarchive-2496ED.svg?logo=docker&logoColor=white)](https://github.com/Elari39/SparkArchive/pkgs/container/sparkarchive)

无产阶级革命理论与实践的文献档案。收录**马克思、恩格斯、普列汉诺夫、蔡特金、片山潜、列宁、卢森堡、柯伦泰、李大钊、胡志明、葛兰西、毛泽东、卡斯特罗、切·格瓦拉、桑卡拉**
十五位革命者的生平、著作与思想，记录国际共产主义运动的理论源流与实践历程。

前端 **Vue 3 + Vite + Tailwind CSS 4**，后端 **Go + SQLite（FTS5 全文索引）**，**Docker 一键部署**，端口 **12026**。

---

## 快速开始

```bash
git clone https://github.com/Elari39/SparkArchive.git
cd SparkArchive
docker compose up -d --build
```

打开 <http://localhost:12026>。

不想本地构建的话，可直接拉取 CI 发布到 GitHub Packages 的镜像：

```bash
docker run -d --name sparkarchive -p 12026:12026 \
  --read-only --tmpfs /tmp \
  --security-opt no-new-privileges:true \
  ghcr.io/elari39/sparkarchive:latest
```

镜像已把前端产物与数据库烘焙在内，运行时只读，不挂任何卷。

> **Windows 提示**：在 Git Bash 下执行上面的 `docker run` 时，`--tmpfs /tmp` 会被 MSYS 改写成本地路径
> 并报 `invalid mount path`。加 `MSYS_NO_PATHCONV=1` 前缀即可（PowerShell / CMD 无此问题）。

> **提示 1（端口）**：本项目端口统一为 **12026**。刻意避开过常见的 1226 —— 该端口容易被本机
> 其他软件（如微信输入法 `wetype_server.exe`）占用。一旦宿主机目标端口被占，Docker 会**静默跳过
> 端口绑定**：容器照常启动且显示 `healthy`，但访问时却是「拒绝连接」。若遇到此现象，
> 先用 `netstat -ano | findstr :12026` 确认端口空闲，或在 `docker-compose.yml` 中改用其他端口。
>
> **提示 2（localhost 解析）**：若访问出现约 20 秒延迟，请改用 <http://127.0.0.1:12026>。
> Docker Desktop 的端口代理只监听 IPv4，而 `localhost` 会优先解析到 IPv6 的 `::1`；
> 浏览器有 Happy Eyeballs 机制会瞬间回退，但部分命令行工具（如 PowerShell 的 `Invoke-WebRequest`）
> 会等完整的 TCP 超时。本仓库的测试脚本因此默认使用 `127.0.0.1`。

```bash
docker compose logs -f      # 查看日志
docker compose down         # 停止
```

镜像内置健康检查，`docker compose ps` 显示 `healthy` 即表示就绪。

---

## 功能模块

| 模块 | 路径 | 说明 |
|---|---|---|
| 人物档案 | `/people/:slug` | 生平年表、简介、政治／经济／文化／思想四域贡献、争议与评价 |
| 著作原文库 | `/works` | 精选原文摘录（标注出处）+ 书目索引 + 外部全文链接 |
| 国际共运年表 | `/timeline` | 1848—2016 大事年表，可按人物与范畴筛选 |
| 术语·概念卡 | `/glossary` | 核心概念释义，含人物、关联术语、出处著作互链 |
| 人物关系图 | `/graph` | 按出生世代分栏的思想传承谱系：节点带头像、姓名与国籍生年，连线按关系类型分族编码，配图例、聚焦详情与等价列表 |
| 全文检索 | `/search` | 跨人物／著作／事件／术语检索，结果高亮 |
| 凡例 | `/about` | 内容体例、史料来源、争议标注规则、版权说明 |

---

## 视觉风格：构成主义 × 新粗野主义

以苏联**构成主义**的历史语汇（斜向构图、红黑对比、超大字号、几何拼贴、章节编号）为骨，
以**新粗野主义**的现代手法（硬边框、硬阴影、无渐变、克制的圆角）为皮。

设计令牌集中在 [`web/src/styles/theme.css`](web/src/styles/theme.css) 的 `@theme` 块（Tailwind 4 不使用 `tailwind.config.js`）：

| 令牌 | 值 | 用途 |
|---|---|---|
| `--color-ink` | `#111111` | 主文字与边框 |
| `--color-paper` | `#f2ede4` | 档案纸张底色 |
| `--color-red` / `--color-red-deep` | `#d62828` / `#a4161a` | 革命红，强调与印章 |
| `--color-brass` | `#f4b400` | 工业黄，仅用于深底或大字号 |
| `--shadow-hard` | `6px 6px 0 #111` | 无模糊硬阴影 |
| `--border-brutal` | `3px solid #111` | 统一硬边框 |

**约束**：禁止渐变、禁止柔和阴影、禁止玻璃拟态；圆角 ≤ 4px。
**可访问性**：正文对比度 ≥ 4.5:1；`focus-visible` 硬描边焦点环；全局遵守 `prefers-reduced-motion`。
**史料照片**：人物肖像为**公有领域历史照片**，以硬边框／硬阴影嵌入版面，并以 `.portrait-photo`
的 `grayscale + sepia` 滤镜统一为单色，与红／黑／米白色板协调（滤镜不属渐变）。

---

## 技术架构

```
content/            单一内容源（版本化，无需数据库后台）
  ├─ site.yaml        站点文案
  ├─ people/*.md      人物档案（YAML front-matter + Markdown 正文）
  ├─ works.yaml       著作与摘录
  ├─ events/          大事年表
  ├─ terms.yaml       术语卡
  └─ relations.yaml   关系图边
        │
        │  archive ingest（构建期执行）
        ▼
  /data/archive.db    SQLite + FTS5 全文索引
        │
        │  go:embed（前端 dist + 数据库）
        ▼
  单二进制容器  →  :12026
        │
        │  docker compose ports 映射
        ▼
  宿主机  http://localhost:12026
```

内容在**构建期**编译进镜像，运行时数据库以**只读**方式打开（`mode=ro&immutable=1`），
因此容器无需数据卷、可设 `read_only: true`，攻击面最小。

### 目录结构

```
server/
  ├─ cmd/archive/       main.go（服务）· ingest.go（构建期入库）
  ├─ internal/content/  内容解析与引用完整性校验
  ├─ internal/store/    SQLite 写入、查询与 FTS5 检索
  └─ internal/api/      HTTP 处理器与中间件
web/src/
  ├─ api/               fetch 封装 + zod 运行时校验
  ├─ components/        新粗野主义基础组件
  │   ├─ home/          首页分区子组件（门楣 / 模块导航 / 收录人物）
  │   └─ portrait-art/  构成主义矢量肖像（照片缺图／失败时的降级分支）
  ├─ config/portraits.ts 肖像资产映射表（portrait_key → 照片路径与署名）
  ├─ composables/       useAsync / useDebouncedRef / useSiteCopy
  ├─ views/             11 个路由页面
  ├─ stores/            Pinia
  └─ styles/theme.css   设计令牌
web/scripts/
  ├─ lib/cdp.mts        无头 Chrome + CDP 公共库（启动 / 建页 / 求值 / 截图）
  ├─ render-check.mts   渲染全部路由并检查挂载与控制台报错
  ├─ og-cover.mts       生成社交分享封面
  └─ portrait-sheet.mts 生成肖像对比图
web/public/
  ├─ favicon.svg        站点图标
  ├─ og-cover.png       社交分享封面（pnpm og-cover 生成）
  └─ assets/portraits/  公有领域历史照片（WebP + JPEG）与 CREDITS.md
```

> 无头 Chrome 路径优先取 `CHROME_PATH` 环境变量，未设置时按平台自动探测常见安装位置；
> 三个渲染脚本共用 `scripts/lib/cdp.mts`，不再各自复制整套 CDP 样板。

---

## 中文全文检索的实现要点

这是本项目最容易踩坑的地方，特此说明。

SQLite FTS5 **默认的 `unicode61` 分词器不切分中文**——整段中文会被当作**一个 token**，
导致中文检索完全失效。因此：

1. 建表时显式指定 **`tokenize='trigram'`**（[`schema.sql`](server/internal/store/schema.sql)），trigram 分词器支持 CJK 子串匹配；
2. FTS5 trigram 要求查询**至少 3 个字符**，故 **2 字及以下的查询（如「列宁」「游击」）自动降级为 `LIKE` 子串匹配**（[`search.go`](server/internal/store/search.go)）；
3. 用户输入以**双引号包裹为短语**并转义内部引号，杜绝 FTS5 语法错误与注入；
4. 结果按 `bm25`（`ORDER BY rank`）排序，并用 `snippet()` 生成摘要。

[`search_test.go`](server/internal/store/search_test.go) 用真实内容覆盖了长短词、跨类型、类型过滤、
分页、注入串与空查询共 9 组中文用例。

---

## 开发

```bash
make install     # 安装前端依赖
make ingest      # 从 content/ 重建数据库
make dev-api     # 启动后端 :12026
make dev-web     # 启动前端 :5173（/api 代理到 12026）

make test        # go vet + go test + vue-tsc
make build       # 前端构建 + 内嵌 + 单二进制
make docker-up   # 一键部署
make smoke       # 端到端冒烟测试（需先 up）
make contract    # 前后端契约校验（需先 up）

cd web && pnpm render-check   # 无头 Chrome 渲染 11 条路由，检查挂载与控制台报错
cd web && pnpm og-cover       # 重新生成社交分享封面 public/og-cover.png
```

### 环境变量

| 变量 | 默认 | 说明 |
|---|---|---|
| `ADDR` | `:12026` | 监听地址 |
| `DB_PATH` | `/data/archive.db` | 数据库路径 |
| `SITE_PATH` | `/app/content/site.yaml` | 站点文案（可挂载覆盖，无需重新发版） |
| `LOG_LEVEL` | `info` | `debug\|info\|warn\|error`（JSON 结构化日志） |

改端口：编辑 [`docker-compose.yml`](docker-compose.yml) 的 `ports` 与 `ADDR`。

---

## 持续集成与发布

工作流见 [`.github/workflows/ci.yml`](.github/workflows/ci.yml)，触发于 push / PR / 手动。

**为什么 CD 是发镜像而不是发 Pages**：站点是 SPA，全部数据来自 `/api/*`；
静态托管只会得到一个没有数据的空壳页面。本项目真正自包含的产物是**镜像**
——前端在构建期 `embed` 进二进制、数据库也在构建期烘焙，运行时只读。
所以发布 = 构建镜像 → 冒烟 → 推送 GHCR，且**不需要任何保密凭据**（用内置 `GITHUB_TOKEN`）。

| job | 触发条件 | 内容 |
|---|---|---|
| `go` | push / PR | `gofmt -l`（非空即失败）、`go vet ./...`、`go test ./...`（Go `1.27`） |
| `web` | push / PR | `pnpm install --frozen-lockfile`、`type-check`、`build`；Node `22` 与 `24` 两条 LTS 线矩阵 |
| `publish` | 仅 `main`，且 `needs: [go, web]` 全绿 | buildx 构建镜像 → **实跑容器冒烟** → 推送 `ghcr.io/elari39/sparkarchive`（`latest` / `main` / `sha-<短哈希>`） |

两处设计上的取舍：

- **矩阵为什么不写 `20`**：`vite@8` 与 `@vitejs/plugin-vue@6` 声明 `engines: ^20.19.0 || >=22.12.0`，
  取两条 LTS 线即可 —— `24` 与 Dockerfile 的 `node:24-alpine` 一致（实际交付版本），`22` 是受支持的下限护栏。
- **发布前必须先冒烟**：`go test` 用的是空的 `cmd/archive/webdist/`（仓库里只有一个 `.gitkeep`），
  它**证明不了**前端真的被内嵌。所以 `publish` 先以 `load: true` 构建到本地、实跑容器断言
  `/api/meta` 计数与静态资源状态码，通过后才 `push`，坏镜像不会发出去。


## 内容编写

新增人物：在 `content/people/` 下新建 `.md`，front-matter 含 `slug/name/birth/death/...`，
正文用 `## 政治贡献`、`## 经济贡献`、`## 文化贡献`、`## 思想体系`、`## 争议与评价` 二级标题
——标题即分类依据（[`parse.go`](server/internal/content/parse.go) 的 `kindFor`），
**新增人物必须补齐这五个章节**，否则 `go test` 的 `TestLoadRealContent` 会失败。

几处容易被漏掉的地方：

- **人物顺序由 `birth` 决定**。`Load` 在读完所有档案后按 `birth`（ISO 8601 字符串）升序重排，
  人物列表、事件／术语中的人物引用顺序都取自它；关系图也按 `birth` 断代分栏。因此新人的 `birth` 必须填写完整，
  否则会排到最前（关系图里则会自成一个「未系年」带）。
- **每人至少要有 1 部著作**，否则 `search_test.go` 的 `TestPersonDetail` 会失败（它校验每人都有
  四域贡献、生平年表与关联著作）。
- **肖像**：把 `<slug>.jpg` / `<slug>.webp`（480×480）放入 `web/public/assets/portraits/`，
  并在 [`portraits.ts`](web/src/config/portraits.ts) 登记来源、摄影者与**真实的许可状态**
  ——**不要一律写「公有领域」**，那是对权利状态的陈述，写错反而更麻烦。**不登记也不会报错**——
  `Portrait.vue` 会静默回退为矢量插画，所以要主动确认。若确定长期无照片，
  请在 [`PortraitArt.vue`](web/src/components/portrait-art/PortraitArt.vue) 补一个 `<g v-else-if>` 分支，
  否则该人物会套用其他样式的人物插画。
- **首页门楣照片墙**只陈列 6 位代表性人物，名单是 [`site.ts`](web/src/config/site.ts) 中的
  `HERO_WALL`（编辑取舍，不是数据切片，故需显式声明）。其余人物在「收录人物」分区完整列出。
- 涉及「收录范围」「累计人数」的文案散落在 `content/site.yaml`、`site.ts`、`ModuleGrid.vue`、
  `PeopleView.vue`、`AboutView.vue`、`TimelineView.vue`、`README.md` 与三个 `web/scripts/*.mts`
  及 `scripts/smoke.ps1` 中，新增人物时需一并核对。后端的 `parse_test.go` 与 `api_test.go`
  各有一处数量断言；`internal/store` 的测试则从内容源推导人数，无需改动。

入库时会做**引用完整性校验**：著作引用的人物、术语引用的著作与关联术语、事件引用的、
关系边引用的两端人物，只要有一个不存在，构建即失败并在错误信息中指出具体条目。

---

## 内容立场与版权

- **写法**：歌颂性基调 + **史实可核**。事实、时间、引文均尽量标注出处；对存在学术争论或历史评价
  分歧之处，一律设「争议与评价」专节如实说明（含对失误与错误的记载）。这既是史实要求，也是档案可信度的来源。
- **著作**：仅收录**精选摘录 + 书目索引 + 外部全文链接**，不收录受版权保护的全文。
  **毛泽东著作在中国境内的版权保护期至 2026 年底**，故尤其不收录全文。
- **肖像**：人物肖像取自公开的历史照片与史料图片，**逐张**标注来源、摄影者与许可状态，
  清单见 [`portraits.ts`](web/src/config/portraits.ts) 与
  [`CREDITS.md`](web/public/assets/portraits/CREDITS.md)，并同步展示于《凡例》页。照片仅用于学习与研究，
  不作商业用途。原构成主义几何矢量插画保留为**加载失败／缺图时的降级方案**
  （[Portrait.vue](web/src/components/Portrait.vue)）。
  **两组来源的许可状态不同，请勿混同**：其中 7 张取自维基共享资源，为公有领域史料
  （含切·格瓦拉所用 Korda《英勇的游击队员》，维基共享资源标注为公有领域，惟该照片历史上曾存在道德权争议，
  本站如实记录、仅作史料展示）；另 8 张（蔡特金、普列汉诺夫、片山潜、柯伦泰、李大钊、葛兰西、卡斯特罗、桑卡拉）
  取自百度百科词条摘要图，**原图未标注许可、不属于公有领域**，依站点运营方的声明使用。
  其中卡斯特罗、桑卡拉的照片仍在著作权保护期内，片山潜、葛兰西两张为现代数字修复／上色版本。
  详见 [`CREDITS.md`](web/public/assets/portraits/CREDITS.md)。
- **授权**：本站原创文字采用 CC BY-NC-SA 4.0；所引原著摘录版权归各自权利人所有。

---

## 开源许可

本项目以 **MIT License** 发布，详见 [LICENSE](LICENSE)。

> 注意：MIT 仅覆盖本仓库的**代码**。站内收录的原著摘录版权归各自权利人所有，
> 肖像照片为公有领域史料（清单与署名见 [`CREDITS.md`](web/public/assets/portraits/CREDITS.md)）；
> 原创**文字内容**沿用 CC BY-NC-SA 4.0。三者范围不同，详见「内容立场与版权」一节。

---

## 明确不做（本期）

后台管理与鉴权（内容以仓库内文件版本化）、多语言 i18n、SSR/SSG、评论与用户系统。

---

## 技术决策记录

| 决策 | 原因 |
|---|---|
| 锁定 `typescript@5.9.3` 而非最新的 7.0.2 | **实测** `vue-tsc@3.3.12` 与 TypeScript 7 不兼容：TS 7 的 `exports` 不再暴露 `./lib/tsc`，`vue-tsc` 启动即抛 `ERR_PACKAGE_PATH_NOT_EXPORTED`。锁 5.9.3 后类型检查全绿。待 `vue-tsc` 支持 TS 7 后可升级。 |
| 后端用 `encoding/json` v1 而非 v2 | Go 1.27.1 中 `encoding/json/v2` 被 `//go:build goexperiment.jsonv2` 门控，而 `GOEXPERIMENT` 官方声明「不支持用于生产」。已在代码注释中标注该偏离。 |
| Dockerfile 固定 `GOPROXY=https://goproxy.cn,direct` | 实测本机与部分网络环境**无法访问 `proxy.golang.org`**；`go.sum` 已提交以锁定校验和。 |
| 用 `modernc.org/sqlite` 而非 `mattn/go-sqlite3` | 纯 Go 实现，`CGO_ENABLED=0`，可在 alpine 上产出静态二进制，镜像更小、构建更快。 |
| 内容在构建期入库而非运行时 | 数据集小（< 2MB）且只读，省去数据卷与写权限，容器可完全只读。 |
| 前端用 `fetch` + `zod` 校验响应 | 后端返回畸形数据时立即报错，而不是把 undefined 渲染进页面。 |
| 关系图附等价列表 | SVG 图谱对键盘与读屏用户不友好，同时提供按关系族分组的纯文本关系列表；窄屏（< 640px）默认收起图形、只给列表。 |
| 人物肖像改用公有领域历史照片，矢量插画降级 | 几何插画辨识度低、缺少史料质感。改用公有领域历史照片并逐张署名，保留原矢量插画作加载失败／缺图的降级。图片是构建期静态资产，故以 `portrait_key` 在前端映射（[`portraits.ts`](web/src/config/portraits.ts)），后端与契约零改动。 |
| 人物肖像的 face 色固定为中性米白 | 初版让面部复用背景/装饰色，导致出现"黄脸""红脸"等有歧义的观感，配色角色必须分离。该规则现仅适用于矢量降级插画。 |
| 测试脚本默认 `127.0.0.1` 而非 `localhost` | Docker Desktop 端口代理只监听 IPv4，`localhost` 优先解析到 `::1` 会让 .NET HttpClient 空等约 21 秒；改用 IPv4 后冒烟测试从 13 分钟降到 2 秒。 |
| 人物按 `birth` 升序排列 | 原先取文件名序，收录到 9 人后变成 `dimitrov → … → marx`（马克思排第 8），列表与关系图都显得随意。改按出生日期排序（ISO 字符串直接比较）后顺序符合档案馆的历史直觉。`ord` 是列表、关系图与各处人物引用的统一顺序来源，故只需在 `Load` 里排一次。 |
| 首页门楣照片墙用显式名单 `HERO_WALL` | 照片墙是**策展**而非数据切片：只陈列 6 位代表性人物，其余在「收录人物」分区完整列出。任何自然排序取前 6 位都得不到这一组合（按出生日期会漏掉毛泽东与切·格瓦拉）。名单放在 `site.ts`，与既有的 `PEOPLE_FILTERS` 同类。 |
| 缺图人物的矢量插画需显式分支 | `PortraitArt.vue` 的兜底分支是最后那个 `<g v-else>`（格瓦拉）。若新人物只加进 `SPECS` 而不加判断分支，缺图时会套用其他人物的插画（例如女性人物戴上格瓦拉的贝雷帽）。因此每位无照片人物都要补一个 `<g v-else-if>`。 |
| 关系图改「出生世代分栏」，不用单环也不用力导向 | 单环布局在 15 人时相邻圆盘净空只剩 6px，且任意两点的直线弦必穿圆内——29 条边在中心叠成一团，边标签（全部落在弦中点）随之糊成色块。力导向则首帧随机、与截图回归冲突，且仓库无 d3。改为按出生年断代（相邻间隔 ≥ 15 年）分栏、带内以「邻居在相邻带中的平均序号」做 8 轮重心排序，横轴即世代、左源右流。节点直径、带内节距、带间距与画布尺寸全部由人数推导，新增人物无需手工调坐标。 |
| 连线端点按「圆盘 ∪ 姓名牌 ∪ 角色行」的轮廓缩进 | 只按圆盘半径缩进时，朝下的连线起点会落在圆盘正下方的姓名牌与角色行上，箭头压在文字上；同带（竖直）连线更是必然如此。改为按朝外方向求轮廓出口距离，两端再按弦长等比收缩；同带连线则直接从圆盘**侧面**出入、向「未来」一侧鼓出，彻底绕开文字区。 |
| 关系图的边标签默认不显示，只在聚焦态出现 | 29 条边各配一个标签必然重叠，这正是改版前的中心色块。现在点击人物时只显示与其相关的边，并在曲线上试 135 个候选位置（15 个参数点 × 9 个法向偏移），取第一个不压圆盘、不压姓名牌／角色行、不与其他标签相撞的位置。 |
| 姓名牌取「最后一个 `·` 之后」而非尾三字 | 原先 `name.slice(-3)` 在 15 人里错了 3 处：「菲德尔·卡斯特罗」→「斯特罗」、「格奥尔基·普列汉诺夫」→「汉诺夫」、「弗拉基米尔·伊里奇·列宁」→「·列宁」（连分隔符一起截进去）。改取 `name.split('·').pop()` 后 15 人全部正确，且无需改内容层。 |
| 关系图的渲染断言按「几何不变量」写 | 节点重叠、标签遮挡这类问题肉眼容易漏。`render-check` 对 `/graph` 额外断言：节点数等于接口人物数、姓名牌等于「最后一个 `·` 之后」的显示名、圆盘中心距 ≥ 直径 + 8、连线端点不落在任何姓名牌／角色行矩形内、聚焦态边标签不压节点且互不重叠。写完后用变异验证（把 `split('·').pop()` 换回 `slice(-3)`）确认断言确实会失败，而非空转。 |
| 首页照片墙加渲染断言 | 曾出现「5 张照片 + 1 个工业黄色块」的失衡版式，且色块是唯一的亮饱和色、反而成为视觉焦点。除改逻辑外，在 `render-check` 中对 `[data-hero-wall] [data-portrait]` 断言为 6，避免同类问题再次静默发生。 |

---

## 验证情况

全部在本机实测通过（Go 1.27.1 / Node 22.22.2 / Docker 29.8.2，Windows）：

| 套件 | 结果 |
|---|---|
| `go vet ./...` / `gofmt -l .` | 通过 |
| `go test ./...` | **通过** — content 解析与边界校验、store 查询、FTS5 中文检索（9 组用例）、API 契约与限流共 4 个包 |
| `pnpm type-check` | 通过（vue-tsc + TypeScript 5.9.3） |
| `pnpm build` | 通过 |
| `pnpm contract-check` | **29/29** — 前端 zod schema 与后端响应逐字段匹配 |
| `pnpm render-check` | **11/11** — 无头 Chrome 渲染全部路由，零控制台错误；首页额外断言门楣照片墙为 6 格；`/graph` 额外断言节点数、姓名牌显示名、圆盘净空、连线端点、聚焦态标签五项几何不变量 |
| `scripts/smoke.ps1` | **63 项** — API、中文检索、错误处理、SPA 回退、缓存头、肖像资产（15 位人物 × WebP+JPEG、immutable） |
| `pnpm og-cover` / `pnpm portrait-sheet` | 通过，输出 1200×630 封面与 15 张肖像对比图，均已人工核对无裁切 |
| 肖像资产 | **15 张全部配图**（7 张维基共享资源公有领域 + 8 张百度百科、许可未标注）；WebP 合计约 344 KB + JPEG 约 480 KB，经 `:12026` 以长期缓存提供 |
| 镜像 | 单二进制容器，健康检查 `healthy` |
| GitHub Actions CI | `go` + `web`（Node 22/24 矩阵）全绿；`publish` 构建镜像并实跑容器冒烟后推送 GHCR |

> **一处环境相关的注意**：在部分沙箱环境中，`smoke.ps1` 跑到后半段（肖像资产与静态资源）时
> PowerShell 的 `Invoke-WebRequest` 会抛出「解析远程名称失败」——即对 `127.0.0.1` 的名称解析间歇性失败。
> 该现象与本项目无关：同一批 URL 用 `curl` 逐个复核全部返回 `200` 且内容类型正确
> （`image/jpeg` / `image/webp` / `image/svg+xml` / `image/png`），`Cache-Control` 亦为
> `public, max-age=31536000, immutable`。两次连续运行失败项并不相同，可确认是环境抖动而非缺陷。
> 若在本机复现失败，建议先用 `curl` 复核再判断。
