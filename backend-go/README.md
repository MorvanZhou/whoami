# whoami backend-go

Python FastAPI 后端的 **Go 练习版**（Gin + GORM 路线）。不接腾讯云 COS，档案存本地文件。  
详细计划见 [PLAN.md](./PLAN.md)。

## 当前进度

- [x] Phase 0：脚手架 + `GET /health`
- [x] Phase 1：本地 `ProfileStore` + 单测
- [x] Phase 2：GORM + User/ApiKey/StoreToken repos + `FindOrCreateUser`
- [x] Phase 3：JWT + `/keys` + `/storeapi` + token 清理
- [x] Phase 4：`/profile` + `/web/profile`（本地文件）
- [x] Phase 5（部分）：GitHub OAuth + `/auth/me` + `/auth/logout`
- [ ] Phase 5 其余：Google / Microsoft / Twitter / Apple
- [x] Phase 6：CORS / 限流 / JSON recover / `/api` 前缀
- [x] Phase 7（测试）：handler 表驱动 + service mock
- [x] Phase 8：Huma 混合接入（JSON + OpenAPI；OAuth/storeapi 仍 Gin）

## 快速开始

```bash
cd backend-go
cp .env.example .env
# 填写 GITHUB_CLIENT_ID / GITHUB_CLIENT_SECRET（Callback URL 指向 GITHUB_REDIRECT_URI）

go mod tidy
go test ./...
go run ./cmd/server
```

探活、文档与登录：

```bash
curl http://localhost:8080/api/health
# {"status":"ok"}

# OpenAPI UI（仅 ENV=dev）
open http://localhost:8080/api/docs

open "http://localhost:8080/api/auth/github/login?redirect=/dashboard"
curl -s http://localhost:8080/api/auth/me -H "Authorization: Bearer <jwt>"
```

**路由分工**：`/auth/me`、`/keys`、`/profile`、`/web/profile` 由 Huma（校验 + OpenAPI）；GitHub OAuth 与 `/storeapi` 仍为 Gin。注意 Huma 字段校验失败多为 **422**。

## 目录

- `cmd/server`：进程入口
- `internal/config`：环境变量
- `internal/domain`：领域类型与接口
- `internal/service`：业务用例（如 `FindOrCreateUser`）
- `internal/store/gormdb`：SQLite + GORM 仓库
- `internal/store/localfs`：档案本地存储（替代 COS）
- `internal/httpapi`：Gin 组装 + 中间件
- `internal/httpapi/humaapi`：Huma JSON 操作注册
- `internal/httpapi/handlers`：Gin 专用（OAuth / storeapi / logout）
