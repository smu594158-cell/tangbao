# 管理后台系统 - 操作手册与 API 文档

## 1. 核心特性概述
- **安全登录与认证**: 管理员登录需校验账号密码及图片验证码；采用 JWT Token 进行会话管理，自动在前端进行过期检测与登出。
- **RBAC 权限控制**: 后端所有 `/api/v1/admin/*` 接口强制经过 `AdminAuth` 中间件验证；前端路由由 `<RequireAdmin>` 高阶组件守护。
- **用户管理**: 支持对注册用户进行增、删、改、查，提供关键词搜索及基于勾选的批量删除、批量角色修改功能。
- **内容管理**: 支持对“景点 (Attraction)”进行 CRUD 操作，数据实时入库。
- **操作日志与防刷**: 管理员执行修改操作(POST, PUT, DELETE)将被 `AdminOperationLog` 中间件拦截并异步落库 `admin_logs` 表。集成 IP 维度的令牌桶限流算法。

---

## 2. API 文档

### 2.1 登录与验证码
#### 获取验证码
- **GET** `/api/v1/auth/captcha`
- **Response**: `{ "code": 0, "data": { "captcha_id": "xxx", "image_url": "base64..." } }`

#### 登录 (带验证码)
- **POST** `/api/v1/auth/login`
- **Body**: `{ "username": "admin", "password": "123", "captcha_id": "xxx", "captcha": "1234" }`

### 2.2 用户管理 (需 Admin 权限)
#### 获取用户列表 (带搜索与分页)
- **GET** `/api/v1/admin/users?page=1&size=10&keyword=test`

#### 创建用户
- **POST** `/api/v1/admin/users`
- **Body**: `{ "username": "user1", "password": "pwd", "nickname": "nick", "role": 1 }`

#### 更新用户
- **PUT** `/api/v1/admin/users/:id`
- **Body**: `{ "nickname": "new", "role": 9, "status": 1 }`

#### 批量删除与批量角色
- **POST** `/api/v1/admin/users/batch-delete`
- **Body**: `{ "ids": [1, 2, 3] }`
- **POST** `/api/v1/admin/users/batch-role`
- **Body**: `{ "ids": [1, 2], "role": 9 }`

### 2.3 内容管理 - 景点 (需 Admin 权限)
- **POST** `/api/v1/admin/attractions` - 创建景点
- **PUT** `/api/v1/admin/attractions/:id` - 更新景点
- **DELETE** `/api/v1/admin/attractions/:id` - 删除景点

---

## 3. 操作手册

### 3.1 登录管理系统
1. 打开浏览器访问 `/login`。
2. 只有 `role = 9` (管理员) 的账号在登录后，侧边栏才会显示**“管理员控制台”**入口。
3. 若无权限，强行访问 `/admin` 将被重定向并提示“权限不足”。

### 3.2 使用仪表盘
点击“管理员控制台”，默认进入“仪表盘”。在此可以查看系统当前的活跃用户数、生成推文数及系统安全状态总览。

### 3.3 用户管理
1. **搜索**：在顶部搜索框输入用户名或昵称，点击搜索即可过滤列表。
2. **新增/编辑**：点击右上角“新建用户”或表格行内的“编辑”按钮，在弹窗中填写表单后保存。
3. **批量操作**：勾选表格左侧的复选框，顶部操作栏将出现“批量删除”、“设为管理员”、“设为普通用户”的快捷按钮。点击前系统会弹出二次确认框，防止误操作。

### 3.4 内容管理 (景点)
切换至“内容管理”标签页，您可以在此维护系统中预设的景点信息，包括经纬度。所有修改将实时影响 C 端用户的“景点资讯”页面数据。
