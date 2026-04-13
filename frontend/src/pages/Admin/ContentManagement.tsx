import { useState, useEffect } from 'react'
import { Table, Space, Button, Modal, Form, Input, InputNumber, message, Popconfirm } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import axios from '../../utils/request'

const ContentManagement = () => {
  const [attractions, setAttractions] = useState([])
  const [loading, setLoading] = useState(false)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)

  const [modalVisible, setModalVisible] = useState(false)
  const [form] = Form.useForm()
  const [editingId, setEditingId] = useState<number | null>(null)

  const fetchAttractions = async (p = 1) => {
    setLoading(true)
    try {
      const res = await axios.get('/api/v1/tour/attractions', {
        params: { page: p, size: 10 }
      })
      if (res.data.code === 0) {
        setAttractions(res.data.data.list || [])
        setTotal(res.data.data.total)
      }
    } catch (error) {
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchAttractions(page)
  }, [page])

  const handleDelete = async (id: number) => {
    try {
      const res = await axios.delete(`/api/v1/admin/attractions/${id}`)
      if (res.data.code === 0) {
        message.success('删除成功')
        fetchAttractions(page)
      }
    } catch (e) {}
  }

  const openModal = (record?: any) => {
    setEditingId(record ? record.id : null)
    if (record) {
      form.setFieldsValue({
        name: record.name,
        description: record.description,
        address: record.address,
        location_lng: record.location_lng,
        location_lat: record.location_lat,
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
        const res = await axios.put(`/api/v1/admin/attractions/${editingId}`, values)
        if (res.data.code === 0) message.success('更新成功')
      } else {
        const res = await axios.post('/api/v1/admin/attractions', values)
        if (res.data.code === 0) message.success('创建成功')
      }
      setModalVisible(false)
      fetchAttractions(page)
    } catch (e) {}
  }

  const columns = [
    { title: 'ID', dataIndex: 'id' },
    { title: '名称', dataIndex: 'name' },
    { title: '地址', dataIndex: 'address' },
    { title: '热度', dataIndex: 'heat_level' },
    {
      title: '操作',
      render: (_: any, record: any) => (
        <Space>
          <Button type="link" onClick={() => openModal(record)}>编辑</Button>
          <Popconfirm title="确定删除该景点吗？" onConfirm={() => handleDelete(record.id)}>
            <Button type="link" danger>删除</Button>
          </Popconfirm>
        </Space>
      )
    }
  ]

  return (
    <div>
      <div className="flex justify-between mb-4">
        <h2 className="text-xl font-bold">景点管理</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => openModal()}>新增景点</Button>
      </div>

      <Table 
        columns={columns} 
        dataSource={attractions} 
        rowKey="id"
        loading={loading}
        pagination={{
          current: page,
          total,
          onChange: (p) => setPage(p)
        }}
      />

      <Modal
        title={editingId ? '编辑景点' : '新增景点'}
        open={modalVisible}
        onOk={handleSave}
        onCancel={() => setModalVisible(false)}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="景点名称" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="description" label="简介" rules={[{ required: true }]}>
            <Input.TextArea rows={4} />
          </Form.Item>
          <Form.Item name="address" label="详细地址" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <div className="flex gap-4">
            <Form.Item name="location_lng" label="经度" rules={[{ required: true }]} className="flex-1">
              <InputNumber className="w-full" step={0.000001} />
            </Form.Item>
            <Form.Item name="location_lat" label="纬度" rules={[{ required: true }]} className="flex-1">
              <InputNumber className="w-full" step={0.000001} />
            </Form.Item>
          </div>
        </Form>
      </Modal>
    </div>
  )
}

export default ContentManagement
