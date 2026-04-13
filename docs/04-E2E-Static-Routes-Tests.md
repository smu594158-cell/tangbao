# 静态页面路由安全与 E2E 测试

## 背景
为满足 PRD 要求，所有通过大模型生成的纯静态页面（SSR/SSG产物）及其控制台 `/generator` 仅用于 SEO 或后台预览，不得向普通 C 端用户暴露。

## 修复措施
1. **导航栏剔除**: 在 `MainLayout.tsx` 中，仅对管理员角色或特定环境显示 "静态页面生成" 入口。
2. **构建白名单过滤**: 在 `App.tsx` 中通过环境变量或域名检测，移除了普通环境下的 `/generator` React 路由。
3. **SEO 隔离**: 在 `public/robots.txt` 和 `public/sitemap.xml` 中剔除了相关的爬虫索引路径。

## E2E 测试用例设计 (Playwright / Cypress 伪代码)

### 用例1：普通用户无法看到生成器菜单
```javascript
test('普通用户登录后不应看到“静态页面生成”菜单', async ({ page }) => {
  await page.goto('/login');
  await page.fill('input[placeholder="用户名 (默认: admin)"]', 'user');
  await page.fill('input[placeholder="密码 (默认: 123456)"]', '123456');
  await page.click('button:has-text("登 录")');
  
  // 验证侧边栏不存在该菜单
  const menuVisible = await page.isVisible('text="静态页面生成"');
  expect(menuVisible).toBeFalsy();
});
```

### 用例2：强行访问白名单外的路由重定向/404
```javascript
test('普通环境下强行访问 /generator 路由不予渲染', async ({ page }) => {
  // 设置环境变量模拟普通C端生产环境
  await page.goto('/generator');
  
  // 验证页面不包含生成器核心按钮
  const btnVisible = await page.isVisible('text="开始生成代码"');
  expect(btnVisible).toBeFalsy();
});
```

### 用例3：SEO 文件验证
```javascript
test('robots.txt 和 sitemap.xml 配置正确', async ({ request }) => {
  const robots = await request.get('/robots.txt');
  const robotsText = await robots.text();
  expect(robotsText).toContain('Disallow: /generator');

  const sitemap = await request.get('/sitemap.xml');
  const sitemapText = await sitemap.text();
  expect(sitemapText).not.toContain('/generator');
});
```
