package domain

// POI (Point of Interest) 兴趣地点
type POI struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Location string `json:"location"` // "经度,纬度"
	Address  string `json:"address"`
	Distance string `json:"distance,omitempty"` // 距离中心点距离
}

// RoutePlan 路线规划结果
type RoutePlan struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	Distance    int    `json:"distance"`
	Duration    int    `json:"duration"`
	Paths       []Path `json:"paths"`
}

// Path 具体路线方案
type Path struct {
	Distance int      `json:"distance"`
	Duration int      `json:"duration"`
	Steps    []string `json:"steps"` // 路线指示步骤
}

// AmapResponse 高德API基础响应格式
type AmapResponse struct {
	Status   string `json:"status"` // 1:成功, 0:失败
	Info     string `json:"info"`
	Infocode string `json:"infocode"`
}
