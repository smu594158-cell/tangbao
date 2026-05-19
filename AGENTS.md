# HZTour (杭州 AI 旅游助手) — Agent 指南

## 概述
基于 Go (Gin + Eino/Ollama) + React SPA 的杭州 AI 旅游助手（毕业设计）。

## 快速命令

```bash
# 后端（在项目根目录或 backend/ 下执行）
make run-backend                    # go run cmd/server/main.go
make test-backend                   # go test -v -cover ./...
make build-backend                  # go build -o bin/server cmd/server/main.go

# 前端
make run-frontend                   # cd frontend && npm run dev
make build-frontend                 # cd frontend && npm run build  (tsc -b && vite build)

# 基础设施
make docker-up                      # 启动 MySQL 8.0、Ollama、Loki+Promtail+Grafana
make init                           # docker-up + 等待就绪 + 提示拉取 qwen:4b 模型
```

**启动后端前：** `cp backend/.env.example backend/.env`，然后填入 `AMAP_WEB_KEY`。

## 架构

| 层 | 后端 | 前端 |
|-------|---------|----------|
| 入口 | `backend/cmd/server/main.go` | `frontend/src/main.tsx` |
| HTTP | `internal/delivery/http/` (Gin handlers) | `src/pages/` (React 组件) |
| 业务 | `internal/usecase/` | `src/store/` (空目录 — Zustand 未接入) |
| 数据 | `internal/repository/` (GORM + MySQL) | `src/api/` + `src/utils/request.ts` (Axios) |
| 领域 | `internal/domain/` (实体 + 接口) | `src/services/` (空目录) |
| 工具包 | `pkg/{errors,response,middleware,mcp,utils}` | `src/components/{Layout,Auth}` |

### 主要依赖
- **后端：** Gin、GORM、Eino+Ollama (qwen:4b)、golang-jwt、godotenv、base64Captcha、bcrypt
- **前端：** React 18、Ant Design 5、Tailwind CSS 3、Vite 5、Zustand 5、react-router-dom 6、Axios、framer-motion

## API 约定

- **统一响应格式：** `{"code": int, "message": string, "data": any}`
- **错误码：** 0=成功，100xx=基础错误，400xx=认证错误，500xx=业务错误
- **鉴权：** `Authorization` 头传 Bearer token；JWT 荷载包含 `userID`、`username`、`role`
- **管理后台：** 路由挂载在 `/api/v1/admin` 下（JWT + role=9 + 操作日志中间件）
- **地理接口：** `/api/v1/geo` 的 handler **不鉴权**（方便测试）

## 注意事项

1. **`.env` 是必需的**且被 gitignore — 从 `.env.example` 复制
2. **MySQL + Ollama** 必须在启动后端前运行（docker-compose 或本地安装）
3. **启动时自动迁移** — GORM `AutoMigrate` 每次启动都会执行，数据库为空时会初始化 4 个景点种子数据
4. **聊天记录清理**（保留 30 天）是领域概念，但 **暂无定时任务** 实现
5. **前端构建必须通过 `tsc -b`**，`tsconfig.json` 启用了 `strict: true`、`noUnusedLocals`、`noUnusedParameters`
6. **前端开发服务器代理** `/api` → `http://localhost:8080`（端口 3000）
7. **路径别名：** `@/` → `src/`（tsconfig 和 vite.config 均已配置）
8. **空目录（脚手架）：** `src/store/`、`src/services/`、`src/api/`、`src/router/` 均为空，仅 `src/utils/request.ts` 存在
9. **前端同时依赖 Google Maps (`@react-google-maps/api`) 和高德地图 JS API v2.0** — 添加地图功能前确认实际使用的是哪个
10. **仅 1 个测试文件：** `backend/pkg/middleware/auth_test.go`
11. **`main.go` 中硬编码了 JWT 密钥和高德 Key 的默认值** — 生产环境务必修改
12. **限流：** 每个 IP 每秒 5 个请求，突发 10 个（内存 map 实现，重启后重置）

## Docker 服务栈

- MySQL 8.0 (hztour-mysql, 端口 3306)、Ollama (hztour-ollama, 端口 11434)
- Loki (3100)、Promtail、Grafana (3000, admin/admin)
- 后端服务在 docker-compose 中**被注释掉了**
