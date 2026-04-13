# 杭州 AI 旅游助手 (HZTour)

基于 Go (Gin + Eino + Ollama) 和 React (Vite + Ant Design) 构建的智能旅游推荐系统。

## 项目架构

本项目采用前后端分离的架构：
- **后端 (`/backend`)**：提供 RESTful API，处理业务逻辑，与 MySQL 数据库交互，并通过 CloudWeGo/Eino 框架集成 Ollama 本地大模型进行旅游推文的智能生成。同时集成了高德地图 Web 服务 API (MCP) 进行地理信息处理。
- **前端 (`/frontend`)**：采用 React 18 和 Vite 构建，使用 Ant Design 提供高质量 UI，Tailwind CSS 负责响应式布局，并深度集成了高德地图 JS API (v2.0) 实现交互式路线规划与展示。

> ⚠️ **注意**：本项目移除了所有静态页面生成服务（SSG）逻辑，前端采用纯 SPA 模式，所有交互均通过 React 客户端渲染完成。

## 开发指南

### 环境要求
- **Go**: 1.22+ (用于后端编译)
- **Node.js**: 18+ (用于前端构建)
- **MySQL**: 8.0+
- **Ollama**: 本地安装并拉取 `qwen:4b` 模型

### 后端运行
```bash
cd backend
# 复制环境变量模板并配置数据库及API Key
cp .env.example .env

# 安装依赖并启动服务 (默认端口 8080)
go mod tidy
go run cmd/server/main.go
```

### 前端运行
```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器 (Vite)
npm run dev
```

## 生产构建与部署

前端构建静态资源供 Nginx 或其他静态服务器托管：
```bash
cd frontend
npm run build
# 产物将生成在 frontend/dist 目录下
```

后端编译可执行文件：
```bash
cd backend
go build -o server.exe ./cmd/server/main.go
# 运行编译后的文件
./server.exe
```
