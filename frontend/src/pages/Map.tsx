import { useState, useEffect, useRef } from 'react'
import { Form, Input, Button, Card, Typography, message, Space } from 'antd'
import { EnvironmentOutlined, FlagOutlined, CompassOutlined } from '@ant-design/icons'
import axios from 'axios'

const { Title, Text } = Typography

declare global {
  interface Window {
    AMap: any;
  }
}

interface PathStep {
  distance: number
  duration: number
  steps: string[]
}

const Map = () => {
  const [loading, setLoading] = useState(false)
  const [pathResult, setPathResult] = useState<PathStep | null>(null)
  const [originCoord, setOriginCoord] = useState<string>('')
  const [destCoord, setDestCoord] = useState<string>('')
  const [currentMode, setCurrentMode] = useState<string>('walking')
  const [destName, setDestName] = useState<string>('')
  
  const mapContainer = useRef<HTMLDivElement>(null)
  const mapRef = useRef<any>(null)
  const routePluginRef = useRef<any>(null)

  useEffect(() => {
    if (!window.AMap || !mapContainer.current) return
    mapRef.current = new window.AMap.Map(mapContainer.current, {
      zoom: 11,
      center: [120.15507, 30.274084], // 杭州市中心
    })

    return () => {
      if (mapRef.current) {
        mapRef.current.destroy()
      }
    }
  }, [])

  const drawRouteOnMap = (origin: string, destination: string, mode: string) => {
    if (!window.AMap || !mapRef.current) return

    // 清除旧的路线
    if (routePluginRef.current) {
      routePluginRef.current.clear()
    }

    const [oLng, oLat] = origin.split(',')
    const [dLng, dLat] = destination.split(',')
    const startLngLat = new window.AMap.LngLat(Number(oLng), Number(oLat))
    const endLngLat = new window.AMap.LngLat(Number(dLng), Number(dLat))

    let pluginName = 'AMap.Walking'
    if (mode === 'driving') pluginName = 'AMap.Driving'
    if (mode === 'transit') pluginName = 'AMap.Transfer'
    if (mode === 'bicycling') pluginName = 'AMap.Bicycling'

    window.AMap.plugin(pluginName, () => {
      let routeObj: any
      // panel 属性如果不指定，高德地图默认只会在地图上画线，不会生成详细面板，这通常没问题。
      const commonOpts = {
        map: mapRef.current,
        hideMarkers: false,
        isOutline: true,
        outlineColor: '#ffeeee',
        autoFitView: true,
      }

      if (mode === 'walking') routeObj = new window.AMap.Walking(commonOpts)
      else if (mode === 'driving') routeObj = new window.AMap.Driving(commonOpts)
      else if (mode === 'transit') routeObj = new window.AMap.Transfer({ ...commonOpts, city: '全国' }) // 设为全国支持跨城公交
      else if (mode === 'bicycling') routeObj = new window.AMap.Bicycling(commonOpts)

      routePluginRef.current = routeObj

      routeObj.search(startLngLat, endLngLat, (status: string, result: any) => {
        if (status === 'complete') {
          console.log('绘制路线成功', result)
          // 自动缩放地图以显示完整的路线
          if (routeObj.setFitView) {
              routeObj.setFitView()
          } else if (mapRef.current.setFitView) {
              mapRef.current.setFitView()
          }
        } else {
          message.warning('地图路线绘制失败：' + (result.info || result))
          console.error('地图路线绘制失败', result)
        }
      })
    })
  }

  const onFinish = async (values: any) => {
    setLoading(true)
    try {
      // 1. 获取起点坐标
      const originRes = await axios.get('/api/v1/geo/poi/search', {
        params: { keywords: values.origin }
      })
      if (originRes.data.code !== 0 || !originRes.data.data || originRes.data.data.length === 0) {
        message.error('无法找到起点位置')
        setLoading(false)
        return
      }
      const originLocation = originRes.data.data[0].location

      // 2. 获取终点坐标
      const destRes = await axios.get('/api/v1/geo/poi/search', {
        params: { keywords: values.destination }
      })
      if (destRes.data.code !== 0 || !destRes.data.data || destRes.data.data.length === 0) {
        message.error('无法找到终点位置')
        setLoading(false)
        return
      }
      const destLocation = destRes.data.data[0].location
      
      setOriginCoord(originLocation)
      setDestCoord(destLocation)
      setCurrentMode(values.type || 'walking')
      setDestName(values.destination)

      // 3. 调用 Web 服务 API 规划路径获取详情（展示用）
      const res = await axios.get('/api/v1/geo/route/plan', {
        params: {
          origin: originLocation,
          destination: destLocation,
          mode: values.type || 'walking'
        }
      })

      if (res.data.code === 0) {
        const planData = res.data.data
        if (!planData.steps) {
          planData.steps = ['从起点出发', '沿规划路线前往', '到达终点']
        }
        setPathResult(planData)
        message.success('路径规划成功')
        
        // 4. 在前端渲染地图路线
        drawRouteOnMap(originLocation, destLocation, values.type || 'walking')
      } else {
        message.error(res.data.message || '获取路径失败')
      }
    } catch (error) {
      message.error('网络请求失败')
    } finally {
      setLoading(false)
    }
  }

  const invokeAmapApp = () => {
    if (!originCoord || !destCoord) {
      message.warning('请先完成路径规划')
      return
    }

    const [oLng, oLat] = originCoord.split(',')
    const [dLng, dLat] = destCoord.split(',')
    
    // t=0:驾车, t=1:公交, t=2:步行, t=3:骑行
    let t = '2'
    if (currentMode === 'driving') t = '0'
    else if (currentMode === 'transit') t = '1'
    else if (currentMode === 'bicycling') t = '3'

    const scheme = `amapuri://route/plan/?slat=${oLat}&slon=${oLng}&sname=起点&dlat=${dLat}&dlon=${dLng}&dname=${encodeURIComponent(destName)}&dev=0&t=${t}`
    const h5Url = `https://uri.amap.com/navigation?from=${oLng},${oLat},起点&to=${dLng},${dLat},${encodeURIComponent(destName)}&mode=${currentMode === 'driving' ? 'car' : currentMode === 'transit' ? 'bus' : currentMode === 'bicycling' ? 'ride' : 'walk'}&policy=1&src=mypage&coordinate=gaode&callnative=1`

    // 尝试唤起 App
    window.location.href = scheme

    // 容错降级：如果 2 秒后页面没有被切换（即未安装 App），则跳转到 H5 导航页
    setTimeout(() => {
      if (document.visibilityState !== 'hidden') {
        window.open(h5Url, '_blank')
      }
    }, 2000)
  }

  return (
    <div className="h-full flex flex-col md:flex-row gap-6">
      <div className="w-full md:w-1/3 flex flex-col gap-6">
        <Card title="路径规划 (MCP 高德地图)" className="shadow-sm">
          <Form layout="vertical" onFinish={onFinish} initialValues={{ type: 'walking' }}>
            <Form.Item 
              label="起点 (例如: 西湖)" 
              name="origin" 
              rules={[{ required: true, message: '请输入起点' }]}
            >
              <Input prefix={<EnvironmentOutlined />} placeholder="输入起点名称" />
            </Form.Item>
            <Form.Item 
              label="终点 (例如: 灵隐寺)" 
              name="destination" 
              rules={[{ required: true, message: '请输入终点' }]}
            >
              <Input prefix={<FlagOutlined />} placeholder="输入终点名称" />
            </Form.Item>
            <Form.Item label="出行方式" name="type">
              <select className="w-full p-2 border border-gray-300 rounded-md outline-none focus:border-blue-500">
                <option value="walking">步行</option>
                <option value="driving">驾车</option>
                <option value="transit">公交</option>
                <option value="bicycling">骑行</option>
              </select>
            </Form.Item>
            <Button type="primary" htmlType="submit" className="w-full bg-blue-600 mb-2" loading={loading}>
              开始规划
            </Button>
            <Button 
              type="default" 
              className="w-full border-blue-600 text-blue-600 hover:bg-blue-50" 
              icon={<CompassOutlined />}
              onClick={invokeAmapApp}
              disabled={!pathResult}
            >
              一键导航 (支持唤起高德App)
            </Button>
          </Form>
        </Card>

        {pathResult && (
          <Card title="规划结果" className="shadow-sm flex-1 overflow-y-auto">
            <Space direction="vertical" className="w-full">
              <div className="bg-blue-50 p-4 rounded-lg flex justify-between items-center">
                <div>
                  <Text type="secondary">总距离</Text>
                  <div className="text-xl font-bold text-blue-600">{(pathResult.distance / 1000).toFixed(2)} km</div>
                </div>
                <div className="text-right">
                  <Text type="secondary">预计耗时</Text>
                  <div className="text-xl font-bold text-green-600">{Math.ceil(pathResult.duration / 60)} 分钟</div>
                </div>
              </div>
              
              <div className="mt-4">
                <Title level={5}>详细步骤</Title>
                <ul className="list-disc pl-5 space-y-2 text-gray-600 text-sm">
                  {pathResult.steps.map((step, index) => (
                    <li key={index}>{step}</li>
                  ))}
                </ul>
              </div>
            </Space>
          </Card>
        )}
      </div>

      <div className="w-full md:w-2/3 bg-gray-100 rounded-lg border border-gray-300 overflow-hidden min-h-[400px]">
        <div ref={mapContainer} className="w-full h-full" />
      </div>
    </div>
  )
}

export default Map