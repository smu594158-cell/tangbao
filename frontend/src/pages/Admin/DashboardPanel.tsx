import { Card, Row, Col, Statistic } from 'antd'
import { UserOutlined, FileTextOutlined, SafetyCertificateOutlined } from '@ant-design/icons'

const DashboardPanel = () => {
  return (
    <div>
      <h2 className="text-xl font-bold mb-4">系统仪表盘</h2>
      <Row gutter={16}>
        <Col span={8}>
          <Card>
            <Statistic
              title="活跃用户数"
              value={120}
              prefix={<UserOutlined />}
            />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic
              title="生成推文数"
              value={345}
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
