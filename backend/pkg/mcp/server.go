package mcp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Server 代表一个本地的 MCP (Model Context Protocol) Server 代理
// 它的职责是封装所有外部的 API 调用（如高德地图），将其转化为统一的工具接口
type Server struct {
	AmapKey string // 高德 Web API Key
}

func NewMCPServer(amapKey string) *Server {
	return &Server{
		AmapKey: amapKey,
	}
}

// SearchPOI 工具: 周边兴趣点搜索
func (s *Server) SearchPOI(keywords, city string) ([]map[string]interface{}, error) {
	baseURL := "https://restapi.amap.com/v3/place/text"
	params := url.Values{}
	params.Add("key", s.AmapKey)
	params.Add("keywords", keywords)
	params.Add("city", city)
	params.Add("offset", "10") // 默认返回10条
	params.Add("page", "1")

	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Status string                   `json:"status"`
		Pois   []map[string]interface{} `json:"pois"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Status != "1" {
		return nil, fmt.Errorf("amap api error")
	}

	return result.Pois, nil
}

// RoutePlanning 工具: 路径规划 (支持 driving, transit, walking, bicycling)
func (s *Server) RoutePlanning(origin, destination, mode string) (map[string]interface{}, error) {
	var baseURL string
	switch mode {
	case "driving":
		baseURL = "https://restapi.amap.com/v3/direction/driving"
	case "transit":
		baseURL = "https://restapi.amap.com/v3/direction/transit/integrated"
	case "walking":
		baseURL = "https://restapi.amap.com/v3/direction/walking"
	case "bicycling":
		baseURL = "https://restapi.amap.com/v4/direction/bicycling"
	default:
		return nil, fmt.Errorf("unsupported route mode: %s", mode)
	}

	params := url.Values{}
	params.Add("key", s.AmapKey)
	params.Add("origin", origin)
	params.Add("destination", destination)
	if mode == "transit" {
		params.Add("city", "杭州") // 公交规划需要城市参数
	}

	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	resp, err := http.Get(reqURL)
	if err != nil {
		fmt.Printf("[MCP Error] HTTP GET failed: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[MCP Error] Read response body failed: %v\n", err)
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("[MCP Error] Unmarshal response failed: %v, Body: %s\n", err, string(body))
		return nil, err
	}

	// 检查高德API的返回状态
	if status, ok := result["status"].(string); !ok || status != "1" {
		fmt.Printf("[MCP Error] Amap API returned error: %s\n", string(body))
		// 降级策略: 返回默认路线或本地缓存
		return s.fallbackRoutePlan(origin, destination, mode), fmt.Errorf("amap api error: %s", string(body))
	}

	return result, nil
}

// fallbackRoutePlan 提供降级路线数据
func (s *Server) fallbackRoutePlan(origin, destination, mode string) map[string]interface{} {
	fmt.Printf("[MCP Info] Using fallback route plan for %s -> %s (%s)\n", origin, destination, mode)
	return map[string]interface{}{
		"status": "1",
		"info":   "OK (Fallback)",
		"route": map[string]interface{}{
			"origin":      origin,
			"destination": destination,
			"paths": []interface{}{
				map[string]interface{}{
					"distance": "1000",
					"duration": "600",
					"steps": []interface{}{
						map[string]interface{}{
							"instruction": "系统降级：请使用本地导航或第三方地图App获取精确路线",
						},
					},
				},
			},
		},
	}
}

