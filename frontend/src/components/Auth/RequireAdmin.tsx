import { Navigate, useLocation } from 'react-router-dom'
import { message } from 'antd'

interface RequireAdminProps {
  children: JSX.Element
}

const RequireAdmin = ({ children }: RequireAdminProps) => {
  const location = useLocation()
  const userStr = localStorage.getItem('user')
  
  if (!userStr) {
    message.warning('请先登录')
    return <Navigate to="/login" state={{ from: location }} replace />
  }

  try {
    const user = JSON.parse(userStr)
    // role 9 is Admin
    if (user.role !== 9) {
      message.error('无权限访问该页面：需要管理员权限')
      // Redirect to a forbidden page or home
      return <Navigate to="/" replace />
    }
  } catch (e) {
    return <Navigate to="/login" replace />
  }

  return children
}

export default RequireAdmin
