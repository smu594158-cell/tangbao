package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/internal/usecase"
	"backend/pkg/errors"
	"backend/pkg/response"
)

type GeoHandler struct {
	geoUseCase usecase.GeoUseCase
}

func NewGeoHandler(r *gin.Engine, uc usecase.GeoUseCase) {
	handler := &GeoHandler{
		geoUseCase: uc,
	}

	geoGroup := r.Group("/api/v1/geo")
	{
		geoGroup.GET("/poi/search", handler.SearchPOIs)
		geoGroup.GET("/route/plan", handler.PlanRoute)
	}
}

// SearchPOIsRequest 搜索请求
type SearchPOIsRequest struct {
	Keywords string `form:"keywords" binding:"required"`
	City     string `form:"city"` // 选填，默认杭州
}

func (h *GeoHandler) SearchPOIs(c *gin.Context) {
	var req SearchPOIsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, response.Error(errors.ErrInvalidParam))
		return
	}

	pois, appErr := h.geoUseCase.SearchPOIs(c.Request.Context(), req.Keywords, req.City)
	if appErr != nil {
		c.JSON(http.StatusOK, response.Error(appErr))
		return
	}

	c.JSON(http.StatusOK, response.Success(pois))
}

// PlanRouteRequest 路线规划请求
type PlanRouteRequest struct {
	Origin      string `form:"origin" binding:"required"`      // 经度,纬度
	Destination string `form:"destination" binding:"required"` // 经度,纬度
	Mode        string `form:"mode" binding:"required,oneof=driving transit walking bicycling"`
}

func (h *GeoHandler) PlanRoute(c *gin.Context) {
	var req PlanRouteRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, response.Error(errors.ErrInvalidParam))
		return
	}

	plan, appErr := h.geoUseCase.PlanRoute(c.Request.Context(), req.Origin, req.Destination, req.Mode)
	if appErr != nil {
		c.JSON(http.StatusOK, response.Error(appErr))
		return
	}

	c.JSON(http.StatusOK, response.Success(plan))
}

