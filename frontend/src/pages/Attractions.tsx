import { useState, useEffect } from 'react'
import { Card, Row, Col, Typography, Tag, Spin, message, Button, Modal, Form, Input, Space, Divider } from 'antd'
import { CompassOutlined, FireOutlined, EnvironmentOutlined, SendOutlined } from '@ant-design/icons'
import axios from '../utils/request'

const { Title, Paragraph, Text } = Typography
const { Meta } = Card

interface Attraction {
  id: number
  name: string
  description: string
  address: string
  heat_level: number
  location_lng?: number
  location_lat?: number
}

const Attractions = () => {
  const [attractions, setAttractions] = useState<Attraction[]>([])
  const [loading, setLoading] = useState(false)
  const [genLoading, setGenLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [generatedText, setGeneratedText] = useState('')

  const [form] = Form.useForm()

  useEffect(() => {
    fetchAttractions()
  }, [])

  const fetchAttractions = async () => {
    setLoading(true)
    try {
      const res = await axios.get('/api/v1/tour/attractions?page=1&size=20')
      if (res.data.code === 0) {
        setAttractions(res.data.data.list || [])
      } else {
        message.error(res.data.message)
      }
    } catch (error) {
      message.error('获取景点列表失败')
    } finally {
      setLoading(false)
    }
  }

  const handleGenerateByID = async (id: number) => {
    setGenLoading(true)
    try {
      const res = await axios.post('/api/v1/tour/content/generate', 
        { attraction_id: id, word_count: 300 }
      )
      
      if (res.data.code === 0) {
        setGeneratedText(res.data.data.generated_content)
        setModalVisible(true)
      } else {
        message.error(res.data.message)
      }
    } catch (error) {
      message.error('生成失败')
    } finally {
      setGenLoading(false)
    }
  }

  const handleGenerateByName = async (values: any) => {
    if (!values.locationName) {
      message.warning('请输入地点名称')
      return
    }

    setGenLoading(true)
    try {
      const res = await axios.post('/api/v1/tour/content/generate', 
        { location_name: values.locationName, word_count: 300 }
      )
      
      if (res.data.code === 0) {
        setGeneratedText(res.data.data.generated_content)
        setModalVisible(true)
      } else {
        message.error(res.data.message)
      }
    } catch (error) {
      message.error('生成失败')
    } finally {
      setGenLoading(false)
    }
  }

  // 辅助方法：返回带有默认图片的封面
  const getAttractionImage = (_name: string, index: number) => {
    // 使用Unsplash占位图，每次根据index不同展示不同的图
    return `https://source.unsplash.com/400x300/?scenery,china,${index}`
  }

  // 预设热门示例地点
  const presetLocations = [
    { name: '西湖', type: '自然风光', color: 'green' },
    { name: '灵隐寺', type: '历史文化', color: 'orange' },
    { name: '钱江新城', type: '城市地标', color: 'blue' },
    { name: '宋城', type: '历史文化', color: 'orange' },
    { name: '千岛湖', type: '自然风光', color: 'green' },
  ]

  return (
    <div className="h-full overflow-y-auto pr-4">
      <div className="mb-6">
        <Title level={2} className="!mb-1">AI 旅游推文生成</Title>
        <Paragraph type="secondary">输入任意目的地或选择下方热门景点，一键生成包含特色、游览时间、交通提示的高质量旅游推文</Paragraph>
      </div>

      {/* 自定义推文生成区域 */}
      <Card className="mb-8 shadow-sm border-blue-100 bg-blue-50/30">
        <Form form={form} layout="vertical" onFinish={handleGenerateByName}>
          <Form.Item 
            label={<span className="text-base font-medium">输入想去的目的地</span>}
            name="locationName" 
            rules={[{ required: true, message: '请输入或选择一个地点' }]}
          >
            <Input 
              size="large" 
              prefix={<EnvironmentOutlined className="text-gray-400" />} 
              placeholder="例如：北京故宫、三亚、杭州西湖..." 
              allowClear
            />
          </Form.Item>
          
          <div className="mb-6">
            <Text type="secondary" className="mr-3">热门示例快速选择：</Text>
            <Space wrap>
              {presetLocations.map(loc => (
                <Tag 
                  key={loc.name} 
                  color={loc.color} 
                  className="cursor-pointer px-3 py-1 text-sm hover:opacity-80 transition-opacity"
                  onClick={() => form.setFieldsValue({ locationName: loc.name })}
                >
                  {loc.name} ({loc.type})
                </Tag>
              ))}
            </Space>
          </div>

          <Form.Item className="mb-0">
            <Button 
              type="primary" 
              htmlType="submit" 
              size="large" 
              loading={genLoading}
              icon={<SendOutlined />}
              className="w-full md:w-auto md:px-12 bg-gradient-to-r from-blue-500 to-indigo-600 border-0"
            >
              立即生成高质量推文
            </Button>
          </Form.Item>
        </Form>
      </Card>

      <Divider className="my-8" />

      <div className="mb-6">
        <Title level={3}>杭州知名景点推荐</Title>
        <Paragraph type="secondary">为您精选杭州热门景点，点击直接生成专属推文</Paragraph>
      </div>

      {loading ? (
        <div className="flex justify-center py-12"><Spin size="large" /></div>
      ) : (
        <Row gutter={[16, 16]}>
          {attractions.map((item, index) => (
            <Col xs={24} sm={12} lg={8} xl={6} key={item.id}>
              <Card
                hoverable
                className="h-full flex flex-col"
                cover={
                  <div className="h-48 overflow-hidden">
                    <img
                      alt={item.name}
                      src={getAttractionImage(item.name, index)}
                      className="w-full h-full object-cover hover:scale-110 transition-transform duration-300"
                    />
                  </div>
                }
                actions={[
                  <Button 
                    type="link" 
                    icon={<SendOutlined />}
                    onClick={() => handleGenerateByID(item.id)}
                    loading={genLoading}
                  >
                    生成推文
                  </Button>
                ]}
              >
                <Meta
                  title={
                    <div className="flex justify-between items-center">
                      <span className="truncate">{item.name}</span>
                      <Tag color="volcano" icon={<FireOutlined />}>
                        {item.heat_level}
                      </Tag>
                    </div>
                  }
                  description={
                    <div className="mt-2">
                      <p className="text-gray-500 text-xs mb-2 flex items-start">
                        <CompassOutlined className="mt-1 mr-1" />
                        <span className="truncate">{item.address}</span>
                      </p>
                      <Paragraph ellipsis={{ rows: 2 }} className="text-sm text-gray-600 h-10">
                        {item.description}
                      </Paragraph>
                    </div>
                  }
                />
              </Card>
            </Col>
          ))}
        </Row>
      )}

      <Modal
        title="AI 旅游推文"
        open={modalVisible}
        onOk={() => setModalVisible(false)}
        onCancel={() => setModalVisible(false)}
        width={700}
        footer={[
          <Button key="close" type="primary" onClick={() => setModalVisible(false)}>
            完成
          </Button>
        ]}
      >
        <div className="bg-gray-50 p-6 rounded-lg mt-4 whitespace-pre-wrap text-gray-700 text-base leading-relaxed border border-gray-100 shadow-inner">
          {generatedText}
        </div>
      </Modal>
    </div>
  )
}

export default Attractions