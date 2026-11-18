# whoami Backend Go 迁移练习计划

> 目标：用 Go 重写 `backend/`（Python + FastAPI），放入 `backend-go/`，**不改动 Python 版**。  
> 用途：个人 Go 练习，**不投入生产**。  
> 关键约束：原腾讯云 COS / STS **全部改为本地文件存储**，不引入任何云对象存储 SDK。

---

## 1. 现状盘点（Python 版）

### 1.1 技术栈

| 能力 | Python 依赖 |
|------|-------------|
| HTTP | FastAPI + Uvicorn |
| 配置 | pydantic-settings + `.env` |
| DB | SQLAlchemy asyncio + SQLite (`aiosqlite`) |
| 鉴权 JWT | python-jose |
| OAuth | authlib + httpx |
| 限流 | slowapi |
| OAuth state 缓存 | cachetools TTLCache |
| 对象存储 | cos-python-sdk-v5 + qcloud-python-sts |

### 1.2 业务模块与路由

| 模块 | 前缀 | 鉴权 | 职责 |
|------|------|------|------|
| health | `GET /health` | 无 | 探活 |
| auth | `/auth` | 混合 | GitHub / Google / Microsoft / Apple / Twitter OAuth；`/me`；`/logout`；JWT Cookie |
| keys | `/keys` | JWT Cookie / Bearer | API Key 创建 / 列表 / 删除；setup 状态；一次性 store URL |
| storeapi | `/storeapi` | 一次性 token | 消费 token，返回明文 `WHOAMI_API_KEY` + endpoint |
| profile | `/profile` | API Key Bearer | Agent 侧读写档案、历史、回滚（原 COS） |
| web/profile | `/web/profile` | JWT | 控制台侧读写档案（原 COS） |
| cos | `/cos/sts` | API Key | 下发临时 STS（**Go 版废弃**） |

### 1.3 数据模型

- `users`：`id`, `email`, `name`, `avatar_url`, `provider`, `provider_id`, `created_at`, `is_active`；`(provider, provider_id)` 唯一
- `api_keys`：`id`, `user_id`, `key_hash`, `key_prefix`, `key_suffix`, `label`, `created_at`, `last_used`, `is_active`
- `store_api_tokens`：`id`, `user_id`, `token`, `api_key_plain`, `api_key_id`, `is_used`, `created_at`（15 分钟过期；消费后删除）

### 1.4 COS 行为（需本地等价实现）

路径约定：

```text
profiles/{user_id}/current.md
profiles/{user_id}/history/{YYYYMMDDTHHMMSSZ}.md
```

行为：

1. `get_profile`：读 `current.md`，不存在返回空
2. `save_profile`：先把现有 `current.md` 拷到 `history/`，再写新内容；历史最多保留 3 份
3. `list_history`：列出历史文件（timestamp / size / last_modified）
4. `rollback_profile`：校验历史存在 → 备份当前 → 用历史覆盖 current → 再清理到 3 份

Go 版用本地目录（如 `./data/profiles/...`）完整复刻上述语义。

---

## 2. 设计原则（练习导向）

1. **API 尽量兼容**：路径、方法、JSON 字段与 Python 版对齐，方便本地用同一前端 / skill 联调。
2. **刻意用本地替代云**：不实现 STS；不依赖 COS / S3 SDK。
3. **优先高效、可扩展的高层库**：本项目不刻意练习 `net/http` / 手写 SQL；选型贴近「能快速做业务」的常见 Go Web 组合（Gin + GORM）。
4. **架构清晰、不过度工程**：分层保留，但不做完整 DDD / 微服务。
5. **Python 版零改动**：两套后端可并存；通过不同端口 / 环境变量切换。

---

## 3. 推荐架构

采用轻量 **分层 + 接口隔离**（类似简化版 Clean Architecture），目录建议：

```text
backend-go/
├── PLAN.md                 # 本文件
├── README.md               # 本地启动说明（后续补）
├── go.mod
├── go.sum
├── .env.example
├── cmd/
│   └── server/
│       └── main.go         # 组装依赖、启动 HTTP
├── internal/
│   ├── config/             # 环境变量加载
│   ├── domain/             # 纯结构体：User / ApiKey / StoreToken / Profile
│   ├── store/              # 持久化接口 + 实现
│   │   ├── gorm/           # users / api_keys / store_tokens（GORM）
│   │   └── localfs/        # profile 本地文件
│   ├── service/            # 业务用例（auth / key / profile / oauth）
│   ├── httpapi/            # Gin handler、路由注册、DTO、中间件
│   │   ├── middleware/     # JWT、API Key、CORS、限流、日志
│   │   └── handlers/
│   └── auth/               # JWT 签发校验、OAuth provider 适配
├── data/                   # 运行时目录（gitignore）
│   ├── whoami.db
│   └── profiles/
└── scripts/                # 可选：迁移 / 冒烟脚本
```

### 3.1 分层职责

| 层 | 做什么 | 不做什么 |
|----|--------|----------|
| `httpapi` | 解析请求、写响应、鉴权中间件 | 不写 SQL / 不直接碰文件系统细节 |
| `service` | 业务规则（找或建用户、Key 哈希、历史保留策略） | 不绑具体 HTTP 框架类型 |
| `store` | SQL / 文件 IO | 不包含业务策略 |
| `domain` | 实体与错误类型 | 无 IO |

### 3.2 关键接口（建议先写接口再实现）

```go
type UserRepository interface {
    FindByProvider(ctx context.Context, provider, providerID string) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
    FindByID(ctx context.Context, id string) (*User, error)
    Create(ctx context.Context, user *User) error
    Update(ctx context.Context, user *User) error
}

type ApiKeyRepository interface { /* create / list / delete / verify by hash */ }

type StoreTokenRepository interface { /* create / consume / cleanup */ }

type ProfileStore interface {
    Get(ctx context.Context, userID string) (content string, ok bool, err error)
    Save(ctx context.Context, userID, content string) error
    ListHistory(ctx context.Context, userID string) ([]HistoryItem, error)
    Rollback(ctx context.Context, userID, timestamp string) (bool, error)
}
```

`ProfileStore` 的本地实现替换原 `cos_service.py`；**不再有 STS 接口实现**。

### 3.3 请求链路（示意）

```text
HTTP Request
  → middleware (CORS / rate limit / recover / request log)
  → auth middleware (JWT cookie|Bearer 或 API Key)
  → handler
  → service
  → sqlite repo 和/或 localfs profile store
  → JSON / redirect / plaintext response
```

---

## 4. 第三方库选型（推荐）

原则：**优先成熟、好扩展的高层库**，少写胶水与样板代码；每个库写清「为什么选它」。

### 4.0 决策更新：Gin + GORM（不再用 chi / sqlx）

原计划偏向 chi + sqlx，是为了练标准库与 SQL。你的目标改为「高效、便于扩展、不练底层」，因此定稿为：

| 原推荐 | 现推荐 | 原因 |
|--------|--------|------|
| chi | **Gin** | 开箱即用：路由组、JSON 绑定、中间件、社区插件多，业务推进更快 |
| sqlx | **GORM** | 接近 SQLAlchemy 心智；CRUD / AutoMigrate 少写 SQL |

**chi vs Gin 一句话**：chi 更「薄」、贴近 `net/http`；Gin 更「厚」、更像 FastAPI——本项目选 Gin。

**什么时候还值得用 chi？** 强依赖标准库生态、要自己拼中间件、或团队规范禁止重度框架时。  
**什么时候 Gin 更合适？** CRUD / 业务 API、希望快速出活、团队熟悉「框架式」开发（正是本项目）。

**sqlx 还需要吗？** **不需要。** 不练 SQL 就直接 GORM；复杂报表以后真需要原生 SQL 时，GORM 也可 `db.Raw(...)`，不必先上 sqlx。

### 4.1 核心推荐

| 用途 | 推荐库 | 说明 |
|------|--------|------|
| 路由 / HTTP | [`gin-gonic/gin`](https://github.com/gin-gonic/gin) | 高性能、生态大；`ShouldBindJSON`、路由组接近 FastAPI 体验 |
| CORS | [`gin-contrib/cors`](https://github.com/gin-contrib/cors) | Gin 官方贡献库 |
| 限流 | [`ulule/limiter`](https://github.com/ulule/limiter) 或 Gin 中间件封装 `golang.org/x/time/rate` | 对齐 slowapi 按 IP 限流 |
| 配置 | [`joho/godotenv`](https://github.com/joho/godotenv) + [`caarlos0/env`](https://github.com/caarlos0/env) | 轻量；不必上 Viper |
| ORM | [`gorm.io/gorm`](https://gorm.io) | 模型、关联、AutoMigrate；心智接近 SQLAlchemy |
| SQLite 驱动 | [`gorm.io/driver/sqlite`](https://gorm.io/docs/connecting_to_the_database.html)（底层可用 modernc 纯 Go） | 练习机免 CGO 更省事 |
| 迁移 | GORM `AutoMigrate`（首版） | 表少；后期可再换 goose |
| JWT | [`golang-jwt/jwt/v5`](https://github.com/golang-jwt/jwt) | 行业事实标准 |
| OAuth2 | [`golang.org/x/oauth2`](https://pkg.go.dev/golang.org/x/oauth2) | 官方扩展；GitHub/Google 有现成 endpoint |
| HTTP Client | 标准库 `net/http` | 拉 userinfo 足够 |
| 校验（基础版） | Gin 内置 binding + [`validator/v10`](https://github.com/go-playground/validator) | 对齐 label 长度、profile 5000；**首版不用 Huma** |
| 文档 / 校验（可选增强） | [`huma`](https://huma.rocks/)（见 §4.5） | OpenAPI 文档 + 结构体驱动校验；**体验基础版后再决定是否接入** |
| UUID | [`google/uuid`](https://github.com/google/uuid) | 生成用户 / key id |
| 日志 | 标准库 `log/slog` 或 Gin 的 logger 中间件 | 结构化日志 |
| 密码/随机 | 标准库 `crypto/rand` + `crypto/sha256` | API Key 与 Python 一致：SHA-256 |
| 进程内 TTL 缓存 | 自研小 map，或 [`patrickmn/go-cache`](https://github.com/patrickmn/go-cache) | 存 OAuth `state`（见下方 TTL 说明） |

### 4.2 明确不推荐（本练习）

| 库 / 方案 | 原因 |
|-----------|------|
| 腾讯云 COS / STS Go SDK、AWS S3 SDK | 违背「去云依赖」目标 |
| chi / sqlx（本项目） | 更适合练底层；与当前「高效出业务」目标不符 |
| 首版就上 Huma | 想先体验纯 Gin 基础版，再对比决定（见 §4.5） |
| ent / sqlc（首版） | 强类型很好，但引入成本高；非必须 |
| Redis（首版） | Python 也是内存 TTL；单进程练习够用 |
| passlib / bcrypt | Python API Key 实际用 SHA-256，Go 保持一致 |

### 4.3 与 Python 依赖对照

| Python | Go 替代（基础版） | 可选增强 |
|--------|-------------------|----------|
| FastAPI / Starlette | **Gin** | Huma（可挂在 Gin/chi 等适配器上） |
| Pydantic 请求校验 | Gin `binding` + validator | Huma 结构体 schema / 校验 |
| FastAPI `/docs` OpenAPI | 暂无（或后补 swag） | **Huma 自动 OpenAPI** |
| uvicorn | 进程内 `http.Server` / Gin `Run`（或 `air` 热重载） | — |
| SQLAlchemy + aiosqlite | **GORM + SQLite** | — |
| pydantic-settings | godotenv + env struct | — |
| authlib | x/oauth2 + 各 provider callback | — |
| python-jose | golang-jwt | — |
| slowapi | Gin 限流中间件 | — |
| cachetools.TTLCache | 带 TTL 的内存 map / go-cache | — |
| cos-python-sdk / sts | **`internal/store/localfs`** | — |

### 4.5 Huma 策略：先无，后评估

目标：希望最终可能用 **Huma** 做「文档 + 校验」（更接近 FastAPI + Pydantic 的体验），但**不在一开始就绑定**。

#### 两阶段约定

| 阶段 | 做法 | 目的 |
|------|------|------|
| **基础版（Phase 0–7）** | 纯 **Gin + binding/validator**；不引入 Huma 依赖 | 先跑通业务，感受手写 handler / DTO / 无自动文档的成本 |
| **评估门禁** | 主路径可用后，你亲自体验与检查 | 再决定「加 Huma / 保持 Gin / 或只补 swag 文档」 |
| **可选 Phase 8** | 若决定采用，再接入 Huma | 用自动 OpenAPI + 统一校验替换/收敛现有 binding |

#### 基础版怎么做校验与文档

- **校验**：请求体/query 用 Gin `ShouldBindJSON` / `binding` tags（`required`、`min`、`max` 等），错误响应尽量对齐 `{"detail":"..."}`。
- **文档**：基础版**不强制** OpenAPI；需要时用 curl / 前端联调即可。不要为了「先有个 docs」提前引入 Huma。

#### 评估时看什么（通过后再上 Huma）

1. 手写 DTO + binding 是否已经够用、烦不烦。
2. 有没有强烈需要浏览器里看 `/docs`、把 schema 当契约给前端。
3. 引入 Huma 后，现有 Gin 中间件（JWT、API Key、CORS、限流）迁移成本是否可接受。
4. OAuth redirect / `PlainText` storeapi 等「非 JSON CRUD」路由，Huma 是否别扭——若别扭，可 **JSON API 走 Huma，特殊路由留纯 Gin**。

#### 若决定加入 Huma（Phase 8 草案）

- [x] 引入 `huma` + Gin adapter（或官方推荐适配方式）
- [x] 选 1～2 个只读 JSON 路由试点（如 `GET /auth/me`、`GET /keys`）
- [x] 再迁写操作（`POST /keys`、`POST /profile` 等）
- [x] 特殊路由（OAuth callback redirect、`/storeapi` plaintext）评估保留 Gin 或 Huma 自定义
- [x] 打开自动 OpenAPI UI，与 Python FastAPI docs 体验对比
- [x] 确认后再删掉重复的纯 binding 样板（避免两套校验并存过久）

#### 明确不做的事

- **不要**在 Phase 0 脚手架就 `go get huma`。
- **不要**为「将来可能用 Huma」提前把目录/接口扭成 Huma 专用形状；保持普通 Gin handler + service 即可，便于以后局部替换。

### 4.4 TTL 是什么？（本项目会用到）

**TTL = Time To Live（存活时间）**：一条数据「最多活多久」，到期自动作废/删除。

本项目里最典型的用途是 **OAuth `state`**：

1. 用户点「GitHub 登录」时，服务器生成随机 `state`，存进内存，并设 **TTL ≈ 10 分钟**。
2. GitHub 回调时带上同一个 `state`；服务器校验：存在且未过期 → 合法；否则拒绝（防 CSRF / 重放）。
3. 10 分钟没人回调，这条 `state` 过期消失，避免内存里堆满垃圾。

同理还有：

- **Store token**：约 15 分钟有效（Python 里主要靠 `created_at` 判断，再加定时清理）。
- **JWT**：也有过期时间（`exp`），本质同样是「带过期时间的凭证」，只是存在客户端 Cookie，不是服务端 TTLCache。

一句话：**TTL 就是给临时数据设闹钟，到点扔掉。** 不是独立中间件产品，而是一种过期策略；可用 Redis 的 `EXPIRE`，也可用进程内 cache。

---

## 5. COS → 本地文件 替换方案

### 5.1 配置变更

删除：

- `COS_SECRET_ID` / `COS_SECRET_KEY` / `COS_REGION` / `COS_BUCKET`

新增：

```env
# 本地档案根目录（相对或绝对路径）
PROFILE_DATA_DIR=./data/profiles
```

### 5.2 实现要点

- 用 `filepath.Join` + 强制 `userID` / `timestamp` 白名单（只允许 `[A-Za-z0-9_-]`），防止路径穿越。
- 写文件用「临时文件 + `Rename`」保证原子性。
- `Save` / `Rollback` 内加**按 userID 的互斥锁**（`sync.Mutex` map），避免并发写坏历史。
- 历史清理：按文件名排序，只留最新 3 个。
- MIME / Content-Type 对本地无意义；HTTP 层继续按 Markdown / JSON 返回即可。

### 5.3 `/cos/sts` 处理

练习版三选一（推荐 **A**）：

- **A. 直接不注册该路由**（文档标明已废弃）
- B. 返回 `501 Not Implemented` + 说明改用服务端 profile API
- C. 返回伪造字段仅供前端占位（不推荐，易误导）

Agent / 控制台应继续走 `/profile` 与 `/web/profile`，与本地存储一致。

---

## 6. 鉴权与安全（对齐 Python）

1. **JWT**：HS256；`sub` = user id；过期天数可配；Cookie 名 `access_token`（HttpOnly / Secure / SameSite=Lax / Path=/）。
2. **双重取 token**：先 Cookie，再 `Authorization: Bearer`。
3. **API Key**：前缀 `wai_` + 48 hex；存 `sha256`；列表只回 prefix/suffix。
4. **Store token**：`token_urlsafe` 等价随机串；15 分钟；消费即删；后台定时清理。
5. **OAuth state**：内存 TTL（约 10 分钟）；callback 校验；redirect 白名单 `/dashboard` `/settings` `/keys`。
6. **Twitter**：PKCE S256。
7. **Apple**：ES256 client_secret JWT（可用 `jwt` + `crypto/ecdsa`）。
8. **限流**：全局约 60/min；login 约 5/min；storeapi 约 10/min。
9. **CORS**：仅允许 `FRONTEND_URL`。
10. **Docs**：仅 `ENV=dev` 时考虑后续加 OpenAPI（练习可后置）。

---

## 7. 分阶段迁移

整体按「可运行垂直切片」推进，每阶段结束都应能 `go run ./cmd/server` 冒烟。

### Phase 0 — 脚手架（0.5 天）

- [x] `go mod init`（模块名建议 `whoami/backend-go` 或 `github.com/<you>/whoami/backend-go`）
- [x] `cmd/server/main.go` + Gin 路由 + `GET /health`
- [x] `internal/config` 读 `.env`（不含 COS）
- [x] `slog` 日志到 stdout（可选每日文件轮转后置）
- [x] `.env.example`、`.gitignore`（`data/`、`.env`）
- [x] `README` 最小启动说明

**验收**：`curl localhost:8080/api/health` → `{"status":"ok"}`

---

### Phase 1 — 本地 ProfileStore（1 天）

- [x] 定义 `ProfileStore` 接口
- [x] 实现 `localfs`：Get / Save / ListHistory / Rollback + 保留 3 份历史
- [x] 单元测试：覆盖新建、覆盖备份、超 3 清理、回滚、非法路径拒绝

**验收**：纯 Go test 通过，无需 HTTP / DB。

> 先做存储是因为这是相对 Python 最大的行为分叉点，且不依赖 OAuth。

---

### Phase 2 — SQLite 与 GORM Repository（1–1.5 天）

- [x] 定义 GORM model：`users` / `api_keys` / `store_api_tokens`（字段对齐 Python）
- [x] Repository 用 GORM 做基础 CRUD（不必手写 SQL）
- [x] `FindOrCreateUser` 逻辑对齐：provider 精确匹配 → email 关联 → 新建
- [x] 启动时 `db.AutoMigrate(...)`

**验收**：用一小段 `_test.go` 或临时 main 插入用户并查出。

---

### Phase 3 — JWT + Keys + StoreAPI（1–1.5 天）

- [x] JWT 签发 / 解析中间件（Cookie + Bearer）
- [x] API Key 中间件
- [x] `/keys`：创建（返回 store_url）、列表、删除、`/{id}/status`
- [x] `/storeapi`：一次性消费，返回 plaintext env 格式
- [x] 后台 goroutine：每小时清理过期 store tokens

**验收**：手工造 JWT / 直接插 DB 用户后，完整走通「创建 key → 拉 storeapi → 用 API Key 访问下一阶段」。

---

### Phase 4 — Profile HTTP API（0.5–1 天）

- [x] `/profile`（API Key）：GET / POST / history / rollback
- [x] `/web/profile`（JWT）：GET / POST
- [x] content 最长 5000 校验
- [x] **不实现** `/cos/sts`

**验收**：用 API Key 写入 Markdown，文件系统出现 `data/profiles/<uid>/current.md` 与 history。

---

### Phase 5 — OAuth（按 provider 增量，2–3 天）

建议顺序（由易到难）：

1. **GitHub**（最常用，流程标准）✅
2. **Google**
3. **Microsoft**
4. **Twitter + PKCE**
5. **Apple**（form_post + client_secret JWT，最后做）

每完成一个 provider：

- [x] `/auth/github/login` + `/callback`（其余 provider 待做）
- [x] state / redirect 白名单
- [x] 写 Cookie 并跳转 `FRONTEND_URL/auth/callback`
- [x] `/auth/me`、`/auth/logout`

**验收**：本地起前端或用 curl 跟 callback（可用 OAuth 测试应用）；至少 GitHub 全链路登录成功。

---

### Phase 6 — 横切能力打磨（0.5–1 天）

- [x] CORS、全局限流、路由级限流
- [x] panic recover、统一错误 JSON（`{"detail":"..."}` 尽量对齐 FastAPI）
- [x] graceful shutdown（`signal.Notify` + `Server.Shutdown`）
- [x] `root` 路径策略：Python 有 `root_path=/api`；Go 可用反代 strip，或路由统一挂在 `/api`

**验收**：错误状态码与限流行为与 Python 大致一致。

---

### Phase 7 — 测试与练习加分项（持续）

- [x] handler 层用 `httptest` 做表驱动测试
- [x] service 层 mock `ProfileStore` / repos（接口收益在这里）
- [ ] 可选：用 `air` / `docker compose` 一键起 Go 后端
- [ ] 可选：抽 oauth provider interface；若查询变复杂再考虑 Raw SQL
- [ ] **停下来体验基础版**（校验手感、有无文档痛点）→ 再决定是否进入 Phase 8

---

### Phase 8 — 可选：接入 Huma（仅评估通过后）

> 详见 §4.5。已采用 **混合**：JSON CRUD 走 Huma，OAuth / storeapi 留 Gin。

- [x] 试点并扩展 JSON API（`/auth/me`、`/keys`、`/profile`、`/web/profile`）
- [x] 评估中间件与 OAuth / plaintext 路由兼容性 → 特殊路由留 Gin
- [x] 启用 OpenAPI UI（dev：`/api/docs`；prod 关闭）
- [x] 决定：混合采用 Huma（校验失败多为 **422**，与旧 Gin 400 略有差异）

---

## 8. 建议实施顺序总览

```text
Phase 0  脚手架 / health          ← 纯 Gin，不上 Huma
    ↓
Phase 1  localfs ProfileStore + 单测     ← 去 COS 的核心
    ↓
Phase 2  SQLite + GORM models / repos
    ↓
Phase 3  JWT + API Keys + storeapi
    ↓
Phase 4  /profile + /web/profile
    ↓
Phase 5  OAuth（GitHub → … → Apple）
    ↓
Phase 6  CORS / 限流 / shutdown
    ↓
Phase 7  测试 + 你体验基础版
    ↓
        ┌─ 觉得文档/校验不够 → Phase 8 评估接入 Huma
        └─ 觉得 Gin 已够用   → 结束，保持基础版
```

---

## 9. 运行与联调约定

| 项 | 建议 |
|----|------|
| 默认端口 | `8080`（Python 可继续用原端口，避免冲突） |
| DB | `./data/whoami.db`（与 Python 的 `whoami.db` **分开**，避免锁/方言问题） |
| Profile | `./data/profiles/` |
| 环境文件 | `backend-go/.env`；`WHOAMI_ENV=dev` 时可选 `.env.dev` |
| 前端 | 将 API base 临时指到 Go 端口做对比测试 |
| 生产密钥 | 练习可用假 OAuth client；真登录再填 GitHub 等凭据 |

### 最小 `.env.example`（Go 版）

```env
ENV=dev
HTTP_ADDR=:8080
DATABASE_URL=./data/whoami.db
PROFILE_DATA_DIR=./data/profiles
JWT_SECRET=change-me-to-a-random-secret
JWT_EXPIRE_DAYS=7
FRONTEND_URL=http://localhost:3000
API_KEY_PREFIX=wai_

GITHUB_CLIENT_ID=
GITHUB_CLIENT_SECRET=
GITHUB_REDIRECT_URI=http://localhost:8080/api/auth/github/callback

# Google / Microsoft / Apple / Twitter 按 Phase 5 再补
```

---

## 10. 风险与刻意简化

| 点 | 处理 |
|----|------|
| 多实例 OAuth state | 练习单进程内存即可；多实例需 Redis（不做） |
| SQLite 并发写 | 开 `WAL`；profile 文件另用 per-user lock |
| 与 Python Cookie Secure | 本地 HTTP 时 Secure Cookie 可能种不上——dev 可加 `COOKIE_SECURE=false` |
| API 细微差异 | 以联调日志为准逐步对齐，不追求第一周 100% 字节级兼容 |
| Apple / Twitter | 可标为 optional，不影响主练习路径 |

---

## 11. 完成定义（Definition of Done）

**基础版**可认为「主路径完成」当且仅当：

1. 不依赖任何云 COS / STS。
2. GitHub 登录 → 创建 API Key → storeapi 取密钥 → `/profile` 读写 → 历史与回滚在本地文件可见。
3. `/web/profile` 在 JWT 下可读写同一份本地档案。
4. Python `backend/` 目录无功能改动。
5. 关键业务有基础测试（至少 `localfs` + 1～2 个 HTTP 冒烟测试）。
6. **未强制引入 Huma**；是否进入 Phase 8 由体验评估决定（见 §4.5）。

---

## 12. 下一步

从 **Phase 0 + Phase 1** 开始落地：**纯 Gin 基础版**，先可运行 + 本地档案存储带测试。  
Huma 留到主路径跑通、你检查过校验与文档体验之后再议。
