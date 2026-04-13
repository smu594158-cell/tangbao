import { useState, useEffect } from 'react'
import { Form, Input, Button, message, Tabs } from 'antd'
import { UserOutlined, LockOutlined, IdcardOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import axios from 'axios'

const Login = () => {
  const [loading, setLoading] = useState(false)
  const [activeTab, setActiveTab] = useState('login')
  const [captchaUrl, setCaptchaUrl] = useState('')
  const [captchaId, setCaptchaId] = useState('')
  const navigate = useNavigate()

  const fetchCaptcha = async () => {
    try {
      const res = await axios.get('/api/v1/auth/captcha')
      if (res.data.code === 0) {
        setCaptchaId(res.data.data.captcha_id)
        setCaptchaUrl(res.data.data.image_url)
      }
    } catch (error) {
      console.error('Failed to fetch captcha')
    }
  }

  useEffect(() => {
    fetchCaptcha()
  }, [])

  const onLogin = async (values: any) => {
    setLoading(true)
    try {
      const payload = {
        ...values,
        captcha_id: captchaId
      }
      const res = await axios.post('/api/v1/auth/login', payload)
      if (res.data.code === 0) {
        message.success('登录成功')
        // 设置 token 和过期时间 (假设 7 天)
        const expiry = new Date().getTime() + 7 * 24 * 60 * 60 * 1000
        localStorage.setItem('token', res.data.data.token)
        localStorage.setItem('token_expiry', expiry.toString())
        localStorage.setItem('user', JSON.stringify(res.data.data.user))
        navigate('/')
      } else {
        message.error(res.data.message || '登录失败')
        fetchCaptcha() // 刷新验证码
      }
    } catch (error) {
      message.error('网络错误')
      fetchCaptcha()
    } finally {
      setLoading(false)
    }
  }

  const onRegister = async (values: any) => {
    setLoading(true)
    try {
      const res = await axios.post('/api/v1/auth/register', values)
      if (res.data.code === 0) {
        message.success('注册成功，请登录')
        setActiveTab('login')
        fetchCaptcha()
      } else {
        message.error(res.data.message || '注册失败')
      }
    } catch (error) {
      message.error('网络错误')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="bg-white p-8 rounded-lg shadow-md w-96">
        <div className="text-center mb-8">
          <h2 className="text-2xl font-bold text-gray-800">AI 杭州旅游助手</h2>
          <p className="text-gray-500 mt-2">欢迎使用</p>
        </div>
        
        <Tabs activeKey={activeTab} onChange={setActiveTab} centered items={[
          {
            key: 'login',
            label: '登录',
            children: (
              <Form name="login" onFinish={onLogin} size="large">
                <Form.Item name="username" rules={[{ required: true, message: '请输入用户名' }]}>
                  <Input prefix={<UserOutlined />} placeholder="用户名 (默认: admin)" />
                </Form.Item>
                <Form.Item name="password" rules={[{ required: true, message: '请输入密码' }]}>
                  <Input.Password prefix={<LockOutlined />} placeholder="密码 (默认: 123456)" size="large" />
                </Form.Item>
                <Form.Item>
                  <div className="flex gap-4">
                    <Form.Item
                      name="captcha"
                      noStyle
                      rules={[{ required: true, message: '请输入验证码' }]}
                    >
                      <Input placeholder="验证码" className="flex-1" />
                    </Form.Item>
                    <img 
                      src={captchaUrl} 
                      alt="验证码" 
                      className="h-10 cursor-pointer rounded border"
                      onClick={fetchCaptcha}
                      title="点击刷新验证码"
                    />
                  </div>
                </Form.Item>
                <Form.Item>
                  <Button type="primary" htmlType="submit" className="w-full bg-blue-600" loading={loading}>
                    登录
                  </Button>
                </Form.Item>
              </Form>
            )
          },
          {
            key: 'register',
            label: '注册',
            children: (
              <Form name="register" onFinish={onRegister} size="large">
                <Form.Item name="username" rules={[{ required: true, message: '请输入用户名' }, { min: 3, message: '至少3个字符' }]}>
                  <Input prefix={<UserOutlined />} placeholder="用户名" />
                </Form.Item>
                <Form.Item name="password" rules={[{ required: true, message: '请输入密码' }, { min: 6, message: '至少6个字符' }]}>
                  <Input.Password prefix={<LockOutlined />} placeholder="密码" />
                </Form.Item>
                <Form.Item name="nickname" rules={[{ required: true, message: '请输入昵称' }]}>
                  <Input prefix={<IdcardOutlined />} placeholder="昵称" />
                </Form.Item>
                <Form.Item>
                  <Button type="primary" htmlType="submit" className="w-full bg-blue-600" loading={loading}>
                    注册
                  </Button>
                </Form.Item>
              </Form>
            )
          }
        ]} />
      </div>
    </div>
  )
}

export default Login