import { Routes, Route } from 'react-router-dom'
import { ConfigProvider } from 'antd'
import zhCN from 'antd/locale/zh_CN'

// Components
import MainLayout from './components/Layout/MainLayout'

// Pages
import Home from './pages/Home'
import Chat from './pages/Chat'
import Attractions from './pages/Attractions'
import Map from './pages/Map'
import Login from './pages/Login'
import AdminDashboard from './pages/Admin/Dashboard'
import RequireAdmin from './components/Auth/RequireAdmin'

function App() {
  return (
    <ConfigProvider locale={zhCN}>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/" element={<MainLayout />}>
          <Route index element={<Home />} />
          <Route path="chat" element={<Chat />} />
          <Route path="attractions" element={<Attractions />} />
          <Route path="map" element={<Map />} />
          <Route 
            path="admin" 
            element={
              <RequireAdmin>
                <AdminDashboard />
              </RequireAdmin>
            } 
          />
        </Route>
      </Routes>
    </ConfigProvider>
  )
}

export default App