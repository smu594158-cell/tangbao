import { useState, useEffect } from 'react'
import { Table, Tag, Space, Button, Input, Modal, Form, Select, message, Popconfirm } from 'antd'
import { PlusOutlined, ExclamationCircleOutlined } from '@ant-design/icons'
import axios from '../../utils/request'

const { Option } = Select

const UserManagement = () => {
  const [users, setUsers] = useState([])
  const [loading, setLoading] = useState(false)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [keyword, setKeyword] = useState('')
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([])

  const [modalVisible, setModalVisible] = useState(false)
  const [form] = Form.useForm()
  const [editingId, setEditingId] = useState<number | null>(null)

  const fetchUsers = async (p = 1, k = '') => {
    setLoading(true)
    try {
      const res = await axios.get('/api/v1/admin/users', {
        params: { page: p, size: 10, keyword: k }
      })
      if (res.data.code === 0) {
        setUsers(res.data.data.list || [])
        setTotal(res.data.data.total)
      } else {
        message.error(res.data.message)
      }
    } catch (error) {
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchUsers(page, keyword)
  }, [page, keyword])

  const handleSearch = (value: string) => {
    setKeyword(value)
    setPage(1)
  }

  const handleDelete = async (id: number) => {
    try {
      const res = await axios.delete(`/api/v1/admin/users/${id}`)
      if (res.data.code === 0) {
        message.success('删除成功')
        fetchUsers(page, keyword)
      }
    } catch (e) {}
  }

  const handleBatchDelete = async () => {
    if (selectedRowKeys.length === 0) return message.warning('请选择用户')
    try {
      const res = await axios.post('/api/v1/admin/users/batch-delete', { ids: selectedRowKeys })
      if (res.data.code === 0) {
        message.success('批量删除成功')
        setSelectedRowKeys([])
        fetchUsers(page, keyword)
      }
    } catch (e) {}
  }

  const handleBatchRole = async (role: number) => {
    if (selectedRowKeys.length === 0) return message.warning('请选择用户')
    try {
      const res = await axios.post('/api/v1/admin/users/batch-role', { ids: selectedRowKeys, role })
      if (res.data.code === 0) {
        message.success('批量修改权限成功')
        setSelectedRowKeys([])
        fetchUsers(page, keyword)
      }
    } catch (e) {}
  }

  const openModal = (record?: any) => {
    setEditingId(record ? record.id : null)
    if (record) {
      form.setFieldsValue({
        nickname: record.nickname,
        role: record.role,
        status: record.status
      })
    } else {
      form.resetFields()
    }
    setModalVisible(true)
  }

  const handleSave = async () => {
    try {
      const values = await form.validateFields()
      if (editingId) {
        // Update
        const res = await axios.put(`/api/v1/admin/users/${editingId}`, values)
        if (res.data.code === 0) message.success('更新成功')
      } else {
        // Create
        const res = await axios.post('/api/v1/admin/users', values)
        if (res.data.code === 0) message.success('创建成功')
      }
      setModalVisible(false)
      fetchUsers(page, keyword)
    } catch (e) {}
  }

  const columns = [
    { title: 'ID', dataIndex: 'id' },
    { title: '用户名', dataIndex: 'username' },
    { title: '昵称', dataIndex: 'nickname' },
    { 
      title: '角色', 
      dataIndex: 'role',
      render: (role: number) => <Tag color={role === 9 ? 'red' : 'blue'}>{role === 9 ? '管理员' : '普通用户'}</Tag>
    },
    { 
      title: '状态', 
      dataIndex: 'status',
      render: (status: number) => <Tag color={status === 1 ? 'success' : 'default'}>{status === 1 ? '正常' : '禁用'}</Tag>
    },
    {
      title: '操作',
      render: (_: any, record: any) => (
        <Space>
          <Button type="link" onClick={() => openModal(record)}>编辑</Button>
          <Popconfirm title="确定删除吗？" onConfirm={() => handleDelete(record.id)}>
            <Button type="link" danger>删除</Button>
          </Popconfirm>
        </Space>
      )
    }
  ]

  return (
    <div>
      <div className="flex justify-between mb-4">
        <Space>
          <Input.Search 
            placeholder="搜索用户名/昵称" 
            onSearch={handleSearch} 
            enterButton 
            allowClear
          />
          <Button type="primary" icon={<PlusOutlined />} onClick={() => openModal()}>新建用户</Button>
        </Space>
        <Space>
          {selectedRowKeys.length > 0 && (
            <>
              <Popconfirm title="确定批量删除选中的用户吗？" onConfirm={handleBatchDelete} icon={<ExclamationCircleOutlined style={{ color: 'red' }} />}>
                <Button danger>批量删除</Button>
              </Popconfirm>
              <Button onClick={() => handleBatchRole(9)}>设为管理员</Button>
              <Button onClick={() => handleBatchRole(1)}>设为普通用户</Button>
            </>
          )}
        </Space>
      </div>

      <Table 
        rowSelection={{
          selectedRowKeys,
          onChange: setSelectedRowKeys,
        }}
        columns={columns} 
        dataSource={users} 
        rowKey="id"
        loading={loading}
        pagination={{
          current: page,
          total,
          onChange: (p) => setPage(p)
        }}
      />

      <Modal
        title={editingId ? '编辑用户' : '新建用户'}
        open={modalVisible}
        onOk={handleSave}
        onCancel={() => setModalVisible(false)}
      >
        <Form form={form} layout="vertical">
          {!editingId && (
            <>
              <Form.Item name="username" label="用户名" rules={[{ required: true }]}>
                <Input />
              </Form.Item>
              <Form.Item name="password" label="密码" rules={[{ required: true }]}>
                <Input.Password />
              </Form.Item>
            </>
          )}
          <Form.Item name="nickname" label="昵称" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="role" label="角色" rules={[{ required: true }]}>
            <Select>
              <Option value={1}>普通用户</Option>
              <Option value={9}>管理员</Option>
            </Select>
          </Form.Item>
          <Form.Item name="status" label="状态" initialValue={1}>
            <Select>
              <Option value={1}>正常</Option>
              <Option value={0}>禁用</Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default UserManagement
