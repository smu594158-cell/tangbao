import { useState } from 'react'
import { Layout, Menu } from 'antd'
import { UserOutlined, DashboardOutlined, EnvironmentOutlined } from '@ant-design/icons'
import UserManagement from './UserManagement'
import ContentManagement from './ContentManagement'
import DashboardPanel from './DashboardPanel'

const { Sider, Content } = Layout

const AdminLayout = () => {
  const [selectedKey, setSelectedKey] = useState('dashboard')

  const menuItems = [
    { key: 'dashboard', icon: <DashboardOutlined />, label: '仪表盘' },
    { key: 'users', icon: <UserOutlined />, label: '用户管理' },
    { key: 'content', icon: <EnvironmentOutlined />, label: '内容管理 (景点)' },
  ]

  const renderContent = () => {
    switch (selectedKey) {
      case 'dashboard':
        return <DashboardPanel />
      case 'users':
        return <UserManagement />
      case 'content':
        return <ContentManagement />
      default:
        return <DashboardPanel />
    }
  }

  return (
    <Layout className="h-full bg-white">
      <Sider width={200} className="bg-white border-r">
        <Menu
          mode="inline"
          selectedKeys={[selectedKey]}
          items={menuItems}
          onClick={({ key }) => setSelectedKey(key)}
          style={{ height: '100%', borderRight: 0 }}
        />
      </Sider>
      <Layout className="px-6 py-4 bg-gray-50">
        <Content className="bg-white p-6 rounded-lg shadow-sm min-h-[500px]">
          {renderContent()}
        </Content>
      </Layout>
    </Layout>
  )
}

export default AdminLayout
