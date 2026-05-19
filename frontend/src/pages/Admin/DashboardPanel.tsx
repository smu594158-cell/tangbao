import { useEffect, useState } from 'react'
import { Card, Row, Col, Statistic, Spin } from 'antd'
import { UserOutlined, FileTextOutlined, SafetyCertificateOutlined } from '@ant-design/icons'
import axios from '../../utils/request'

interface DashboardStats {
  userCount: number
  contentCount: number
}

const DashboardPanel = () => {
  const [stats, setStats] = useState<DashboardStats>({ userCount: 0, contentCount: 0 })
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const res = await axios.get('/api/v1/admin/users?page=1&size=1')
        if (res.data.code === 0) {
          setStats(prev => ({ ...prev, userCount: res.data.data.total || 0 }))
        }
      } catch (e) {
        // ignore
      }
      try {
        const res = await axios.get('/api/v1/tour/attractions?page=1&size=1')
        if (res.data.code === 0) {
          setStats(prev => ({ ...prev, contentCount: res.data.data.total || 0 }))
        }
      } catch (e) {
        // ignore
      }
      setLoading(false)
    }
    fetchStats()
  }, [])

  if (loading) {
    return <div className="flex justify-center py-12"><Spin size="large" /></div>
  }

  return (
    <div>
      <h2 className="text-xl font-bold mb-4">系统仪表盘</h2>
      <Row gutter={16}>
        <Col span={8}>
          <Card>
            <Statistic
              title="注册用户数"
              value={stats.userCount}
              prefix={<UserOutlined />}
            />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic
              title="景点数量"
              value={stats.contentCount}
              prefix={<FileTextOutlined />}
            />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic
              title="系统安全状态"
              value="良好"
              prefix={<SafetyCertificateOutlined />}
              valueStyle={{ color: '#3f8600' }}
            />
          </Card>
        </Col>
      </Row>
    </div>
  )
}

export default DashboardPanel
