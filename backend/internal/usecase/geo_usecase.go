package usecase

import (
	"context"
	"fmt"

	"backend/internal/domain"
	"backend/pkg/errors"
)

// GeoProvider 地理信息服务接口（用于解耦 mcp.Server）
type GeoProvider interface {
	SearchPOI(keywords, city string) ([]map[string]interface{}, error)
	RoutePlanning(origin, destination, mode string) (map[string]interface{}, error)
}

// GeoUseCase 定义地理信息相关用例
type GeoUseCase interface {
	SearchPOIs(ctx context.Context, keywords, city string) ([]*domain.POI, *errors.AppError)
	PlanRoute(ctx context.Context, origin, destination, mode string) (*domain.RoutePlan, *errors.AppError)
}

type geoUseCase struct {
	geoProvider GeoProvider
}

// NewGeoUseCase 实例一个GeoUseCase
func NewGeoUseCase(geoProvider GeoProvider) GeoUseCase {
	return &geoUseCase{
		geoProvider: geoProvider,
	}
}

func (u *geoUseCase) SearchPOIs(ctx context.Context, keywords, city string) ([]*domain.POI, *errors.AppError) {
	if city == "" {
		city = "杭州" // 默认杭州
	}

	rawPois, err := u.geoProvider.SearchPOI(keywords, city)
	if err != nil {
		return nil, errors.ErrMCPServiceFailed
	}

	var pois []*domain.POI
	for _, raw := range rawPois {
		poi := &domain.POI{
			ID:       getString(raw, "id"),
			Name:     getString(raw, "name"),
			Type:     getString(raw, "type"),
			Location: getString(raw, "location"),
			Address:  getString(raw, "address"),
		}
		pois = append(pois, poi)
	}

	return pois, nil
}

func (u *geoUseCase) PlanRoute(ctx context.Context, origin, destination, mode string) (*domain.RoutePlan, *errors.AppError) {
	rawRoute, err := u.geoProvider.RoutePlanning(origin, destination, mode)
	if err != nil && rawRoute == nil {
		return nil, errors.ErrMCPServiceFailed
	}

	// 简单解析部分数据作为演示，实际需根据高德返回格式精细解析
	routeData, ok := rawRoute["route"].(map[string]interface{})
	if !ok {
		return nil, errors.ErrMCPServiceFailed
	}

	plan := &domain.RoutePlan{
		Origin:      getString(routeData, "origin"),
		Destination: getString(routeData, "destination"),
	}

	// 解析 paths (仅以步行/驾车为例)
	if pathsData, ok := routeData["paths"].([]interface{}); ok && len(pathsData) > 0 {
		if firstPath, ok := pathsData[0].(map[string]interface{}); ok {
			plan.Distance = getInt(firstPath, "distance")
			plan.Duration = getInt(firstPath, "duration")
			// 具体 steps 解析可在此扩充
		}
	}

	return plan, nil
}

// 辅助函数: 安全地从 map 中提取string
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return ""
}

// 辅助函数: 安全地提取可能为 string 格式的int (高德API经常把数字返回为字符串）
func getInt(m map[string]interface{}, key string) int {
	if val, ok := m[key]; ok {
		if s, ok := val.(string); ok {
			var i int
			fmt.Sscanf(s, "%d", &i)
			return i
		}
		if f, ok := val.(float64); ok {
			return int(f)
		}
	}
	return 0
}
