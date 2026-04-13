# 权限控制与安全隔离测试

## 权限矩阵

系统包含两种角色：普通用户 (Role=1) 和管理员 (Role=9)。

| 模块 | 功能/接口 | 路由路径 | 普通用户 (User) | 管理员 (Admin) | 备注 |
|---|---|---|---|---|---|
| IAM | 注册、登录 | `/api/v1/auth/*` | ✅ 允许 | ✅ 允许 | 公开接口 |
| Chat | 智能问答 | `/api/v1/chat/message` | ✅ 允许 | ✅ 允许 | 需要 JWT 鉴权 |
| Tour | 景点列表、详情 | `/api/v1/tour/attractions*` | ✅ 允许 | ✅ 允许 | 公开接口 |
| Tour | 生成推文 | `/api/v1/tour/content/generate` | ✅ 允许 | ✅ 允许 | 需要 JWT 鉴权 |
| Geo | POI搜索、路线规划 | `/api/v1/geo/*` | ✅ 允许 | ✅ 允许 | 公开接口 (未来可加鉴权) |
| Admin | 用户管理列表 | `/api/v1/admin/users` | ❌ 拦截 (403) | ✅ 允许 | 需要 JWT + AdminAuth 鉴权 |
| Frontend | 管理员控制台菜单 | `/admin` | ❌ 隐藏且拦截 | ✅ 显示并可访问 | 纯前端 `RequireAdmin` 路由拦截 |

## 单元测试与集成测试验证

### 1. 后端中间件单元测试

在 `backend/pkg/middleware/auth_test.go` 中（假设实现），测试 `AdminAuth`：
- **用例 1**：普通用户 Token 请求 Admin 接口，返回 `40007 禁止访问：权限不足`。
- **用例 2**：管理员 Token 请求 Admin 接口，返回 `200 成功` 并下发数据。

### 2. 前端组件测试 (E2E 伪代码)

**用例 1：普通用户无法访问管理员菜单**
```javascript
test('普通用户登录后不应看到“管理员控制台”菜单', async ({ page }) => {
  await loginAs(page, 'user', '123456'); // Role = 1
  const menuVisible = await page.isVisible('text="管理员控制台"');
  expect(menuVisible).toBeFalsy();
  
  // 强行通过 URL 访问
  await page.goto('/admin');
  await page.waitForURL('/'); // 预期被 RequireAdmin 拦截回首页
});
```

**用例 2：管理员可正常访问并操作控制台**
```javascript
test('管理员登录后可以访问“管理员控制台”并拉取数据', async ({ page }) => {
  await loginAs(page, 'admin', '123456'); // Role = 9
  const menuVisible = await page.isVisible('text="管理员控制台"');
  expect(menuVisible).toBeTruthy();
  
  await page.click('text="管理员控制台"');
  await page.waitForURL('/admin');
  
  // 验证用户列表数据是否渲染
  const tableVisible = await page.isVisible('.ant-table');
  expect(tableVisible).toBeTruthy();
});
```
