import { Layout, Menu, Button } from 'antd'
import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import { MessageOutlined, CompassOutlined, GlobalOutlined, HomeOutlined, LogoutOutlined, SafetyCertificateOutlined } from '@ant-design/icons'
import { useEffect, useState } from 'react'

const { Header, Content, Sider } = Layout

const checkAuth = (): { valid: boolean; isAdmin: boolean } => {
  const token = localStorage.getItem('token')
  const expiry = localStorage.getItem('token_expiry')
  if (expiry && new Date().getTime() > parseInt(expiry, 10)) {
    localStorage.removeItem('token')
    localStorage.removeItem('token_expiry')
    localStorage.removeItem('user')
    return { valid: false, isAdmin: false }
  }
  if (!token) return { valid: false, isAdmin: false }
  try {
    const user = JSON.parse(localStorage.getItem('user') || '{}')
    return { valid: true, isAdmin: user.role === 9 }
  } catch {
    return { valid: !!token, isAdmin: false }
  }
}

const MainLayout = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const [isAdmin, setIsAdmin] = useState(false)

  useEffect(() => {
    const { valid, isAdmin: admin } = checkAuth()
    setIsAdmin(admin)
    if (!valid && location.pathname !== '/login') {
      navigate('/login', { replace: true })
    }
  }, []) // only run on mount

  const handleLogout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('token_expiry')
    localStorage.removeItem('user')
    navigate('/login')
  }

  // 动态生成菜单
  const menuItems = [
    {
      key: '/',
      icon: <HomeOutlined />,
      label: '首页',
    },
    {
      key: '/chat',
      icon: <MessageOutlined />,
      label: 'AI 导游问答',
    },
    {
      key: '/attractions',
      icon: <CompassOutlined />,
      label: '景点资讯',
    },
    {
      key: '/map',
      icon: <GlobalOutlined />,
      label: '地图与路线',
    }
  ]

  if (isAdmin) {
    menuItems.push({
      key: '/admin',
      icon: <SafetyCertificateOutlined />,
      label: '管理员控制台',
    })
  }

  return (
    <Layout className="min-h-screen">
      <Header className="flex items-center justify-between bg-white px-6 shadow-sm z-10">
        <div className="text-xl font-bold text-blue-600">AI 杭州旅游助手</div>
        <Button type="text" icon={<LogoutOutlined />} onClick={handleLogout}>退出登录</Button>
      </Header>
      <Layout>
        <Sider width={200} className="bg-white border-r border-gray-100">
          <Menu
            mode="inline"
            selectedKeys={[location.pathname]}
            style={{ height: '100%', borderRight: 0 }}
            items={menuItems}
            onClick={({ key }) => navigate(key)}
          />
        </Sider>
        <Layout className="p-6 bg-gray-50">
          <Content className="bg-white p-6 rounded-lg shadow-sm min-h-[280px]">
            <Outlet />
          </Content>
        </Layout>
      </Layout>
    </Layout>
  )
}

export default MainLayout