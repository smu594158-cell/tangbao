package http

import (
	"net/http"
	"strconv"

	"backend/internal/domain"
	"backend/pkg/errors"
	"backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// TourHandler 处理旅游景点及内容生成的 HTTP 请求
type TourHandler struct {
	usecase domain.TourUsecase
}

// NewTourHandler 注册路由
func NewTourHandler(r *gin.Engine, uc domain.TourUsecase, authMiddleware gin.HandlerFunc) {
	handler := &TourHandler{usecase: uc}

	api := r.Group("/api/v1/tour")

	// 公开接口：获取景点列表和详情
	api.GET("/attractions", handler.ListAttractions)
	api.GET("/attractions/:id", handler.GetAttraction)

	// 需要鉴权的接口：生成文案
	authApi := api.Group("")
	authApi.Use(authMiddleware)
	{
		authApi.POST("/content/generate", handler.GenerateContent)
	}
}

// ListAttractions 获取景点列表 (分页)
func (h *TourHandler) ListAttractions(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")

	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)

	list, total, err := h.usecase.ListAttractions(c.Request.Context(), page, size)
	if err != nil {
		c.JSON(http.StatusOK, response.Error(errors.ErrInternalServer))
		return
	}

	c.JSON(http.StatusOK, response.Success(response.PageData{
		Total:    total,
		Page:     page,
		PageSize: size,
		List:     list,
	}))
}

// GetAttraction 获取单个景点详情
func (h *TourHandler) GetAttraction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.Error(errors.ErrInvalidParam))
		return
	}

	attraction, err := h.usecase.GetAttractionInfo(c.Request.Context(), id)
	if err != nil {
		if err == errors.ErrNotFound {
			c.JSON(http.StatusOK, response.Error(errors.ErrNotFound))
			return
		}
		c.JSON(http.StatusOK, response.Error(errors.ErrInternalServer))
		return
	}

	c.JSON(http.StatusOK, response.Success(attraction))
}

// GenerateContent 触发 AI 文本生成
func (h *TourHandler) GenerateContent(c *gin.Context) {
	var req domain.GenerateTextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.Error(errors.ErrInvalidParam))
		return
	}

	genText, err := h.usecase.GenerateAttractionText(c.Request.Context(), &req)
	if err != nil {
		appErr, ok := err.(*errors.AppError)
		if !ok {
			appErr = errors.ErrInternalServer
		}
		c.JSON(http.StatusOK, response.Error(appErr))
		return
	}

	c.JSON(http.StatusOK, response.Success(genText))
}
