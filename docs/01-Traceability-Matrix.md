# 需求-代码追溯表 (Requirements to Code Traceability Matrix)

| 需求编号 | 需求描述 (PRD) | 对应领域/模块 (DDD Context) | 核心代码路径 (后端) | 核心代码路径 (前端) | 状态 |
|---|---|---|---|---|---|
| REQ-01 | 用户注册与登录，密码哈希加密，脱敏存储 | IAM (身份与访问管理) | `backend/internal/domain/user.go` <br> `backend/internal/usecase/auth_usecase.go` | `frontend/src/pages/Login.tsx` | 已完成 |
| REQ-02 | AI 智能问答（杭州导游角色），上下文支持 | Chat (对话域) | `backend/internal/domain/tour.go` <br> `backend/internal/usecase/chat_usecase.go` | `frontend/src/pages/Chat.tsx` | 已完成 |
| REQ-03 | 景点资讯与推文智能生成（集成 Ollama 大模型） | Tour & Content (旅游与内容域) | `backend/internal/domain/tour.go` <br> `backend/internal/usecase/tour_usecase.go` | `frontend/src/pages/Attractions.tsx` | 已完成 |
| REQ-04 | 交互式地图渲染与出行路径规划（高德地图 JSAPI 与 MCP 服务集成） | Geo (地理域) | `backend/pkg/mcp/server.go` <br> `backend/internal/usecase/geo_usecase.go` | `frontend/src/pages/Map.tsx` | 已完成 |
| REQ-N01 | 并发能力与资源限制 (500并发, <2GB) | Infrastructure (基础设施) | `deploy/docker-compose.yml` <br> `backend/cmd/server/main.go` | N/A | 已完成 |
| REQ-N02 | 安全性 (防SQL注入, XSS, CSRF, JWT鉴权, 密钥隔离) | IAM / API Gateway | `backend/pkg/middleware/auth.go` <br> `backend/pkg/utils/jwt.go` | `frontend/src/pages/Login.tsx` | 已完成 |
| REQ-N03 | 日志与监控 (Loki+Promtail+Grafana) | Infrastructure | `deploy/docker-compose.yml` | N/A | 基础设施搭建完成 |

## 变更记录
| 版本 | 变更日期 | 责任人 | 变更说明 |
|---|---|---|---|
| V1.0 | 2026-03-24 | 汤诚浩 | 初始版本创建，完成各模块核心代码路径映射。 |
| V1.1 | 2026-03-29 | 汤诚浩 | 架构升级为纯 SPA：彻底移除原 `REQ-04` (静态页面生成服务) 的前后端追溯记录；更新 `REQ-03` 和新的 `REQ-04` 描述以匹配最新的功能实现（旅游推文生成与地图深度集成）；补充 `REQ-N02` 密钥隔离要求。 |

