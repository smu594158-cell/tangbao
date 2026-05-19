import { useState, useRef, useEffect } from 'react'
import { Input, Button, List, Avatar, Spin, message } from 'antd'
import { SendOutlined, UserOutlined, RobotOutlined } from '@ant-design/icons'
import axios from '../utils/request'

interface Message {
  role: 'user' | 'assistant'
  content: string
}

const Chat = () => {
  const [messages, setMessages] = useState<Message[]>([
    { role: 'assistant', content: '你好！我是你的专属杭州导游。有什么我可以帮你的吗？' }
  ])
  const [input, setInput] = useState('')
  const [loading, setLoading] = useState(false)
  const messagesEndRef = useRef<HTMLDivElement>(null)
  
  // 生成一个随机的 session_id，在实际应用中应该持久化或从服务端获取，后端要求 min=8
  const [sessionId] = useState(() => Math.random().toString(36).substring(2, 10))

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }

  useEffect(() => {
    scrollToBottom()
  }, [messages])

  const handleSend = async () => {
    if (!input.trim()) return

    const userMsg = input.trim()
    setInput('')
    setMessages(prev => [...prev, { role: 'user', content: userMsg }])
    setLoading(true)

    try {
      const res = await axios.post('/api/v1/chat/message', 
        {
          session_id: sessionId,
          content: userMsg
        }
      )

      if (res.data.code === 0) {
        setMessages(prev => [...prev, { role: 'assistant', content: res.data.data.reply }])
      } else {
        message.error(res.data.message)
      }
    } catch (error) {
      message.error('发送失败，请重试')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex flex-col h-full bg-white rounded-lg overflow-hidden">
      <div className="p-4 border-b border-gray-100">
        <h2 className="text-xl font-semibold text-gray-800">AI 导游问答</h2>
        <p className="text-sm text-gray-500">向我提问关于杭州的任何旅游问题</p>
      </div>

      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        <List
          itemLayout="horizontal"
          dataSource={messages}
          renderItem={(msg) => (
            <List.Item className={`border-b-0 ${msg.role === 'user' ? 'flex-row-reverse' : ''}`}>
              <List.Item.Meta
                className={`${msg.role === 'user' ? 'text-right' : ''} flex items-start max-w-[80%] ${msg.role === 'user' ? 'ml-auto' : ''}`}
                avatar={
                  <Avatar 
                    icon={msg.role === 'user' ? <UserOutlined /> : <RobotOutlined />} 
                    className={msg.role === 'user' ? 'bg-blue-500 ml-4' : 'bg-green-500 mr-4'}
                  />
                }
                title={<span className="text-gray-500 text-xs">{msg.role === 'user' ? '你' : '杭州导游'}</span>}
                description={
                  <div className={`mt-1 inline-block p-3 rounded-lg ${
                    msg.role === 'user' 
                      ? 'bg-blue-500 text-white rounded-tr-none' 
                      : 'bg-gray-100 text-gray-800 rounded-tl-none'
                  }`}>
                    {msg.content}
                  </div>
                }
              />
            </List.Item>
          )}
        />
        {loading && (
          <div className="flex items-center space-x-2 text-gray-400">
            <Spin size="small" />
            <span>导游正在思考...</span>
          </div>
        )}
        <div ref={messagesEndRef} />
      </div>

      <div className="p-4 border-t border-gray-100 bg-gray-50">
        <div className="flex space-x-2">
          <Input.TextArea
            value={input}
            onChange={e => setInput(e.target.value)}
            onPressEnter={(e) => {
              if (!e.shiftKey) {
                e.preventDefault()
                handleSend()
              }
            }}
            placeholder="输入你的问题... (Shift + Enter 换行)"
            autoSize={{ minRows: 1, maxRows: 4 }}
            className="flex-1"
          />
          <Button 
            type="primary" 
            icon={<SendOutlined />} 
            onClick={handleSend}
            loading={loading}
            className="bg-blue-600 h-auto"
          >
            发送
          </Button>
        </div>
      </div>
    </div>
  )
}

export default Chat