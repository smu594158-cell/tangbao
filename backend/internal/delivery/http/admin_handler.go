package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"backend/internal/usecase"
	"backend/pkg/errors"
	"backend/pkg/response"
)

type AdminHandler struct {
	userUseCase usecase.UserUseCase
	tourUseCase usecase.TourUsecase
}

func NewAdminHandler(r *gin.RouterGroup, uc usecase.UserUseCase, tc usecase.TourUsecase) {
	handler := &AdminHandler{
		userUseCase: uc,
		tourUseCase: tc,
	}

	{
		// 用户管理相关
		r.GET("/users", handler.ListUsers)
		r.POST("/users", handler.CreateUser)
		r.PUT("/users/:id", handler.UpdateUser)
		r.DELETE("/users/:id", handler.DeleteUser)
		r.POST("/users/batch-delete", handler.BatchDeleteUsers)
		r.POST("/users/batch-role", handler.BatchUpdateRole)

		// 内容管理相关 (景点)
		r.POST("/attractions", handler.CreateAttraction)
		r.PUT("/attractions/:id", handler.UpdateAttraction)
		r.DELETE("/attractions/:id", handler.DeleteAttraction)
	}
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	keyword := c.Query("keyword")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	size, err := strconv.Atoi(c.DefaultQuery("size", "10"))
	if err != nil || size < 1 {
		size = 10
	}

	users, total, appErr := h.userUseCase.ListUsers(c.Request.Context(), page, size, keyword)
	if appErr != nil {
		c.JSON(http.StatusOK, response.Error(appErr))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"list":  users,
		"total": total,
	}))
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Nickname string `json:"nickname" binding:"required,max=64"`
	Role     int8   `json:"role" binding:"required"`
}

func (h *AdminHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.Error(errors.ErrInvalidParam))
		return
	}

	user, appErr := h.userUseCase.CreateUser(c.Request.Context(), req.Username, req.Password, req.Nickname, req.Role)
	if appErr != nil {
		c.JSON(http.StatusOK, response.Error(appErr))
		return
	}

	c.JSON(http.StatusOK, response.Success(user))
}

type UpdateUserRequest struct {
	Nickname string `json:"nickname" binding:"required,max=64"`
	Role     int8   `json:"role" binding:"required"`
	Status   int8   `json:"status"`
}

func (h *AdminHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.Error(errors.ErrInvalidParam))
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.Error(errors.ErrInvalidParam))
		return
	}

	user, appErr := h.userUseCase.UpdateUser(c.Request.Context(), id, req.Nickname, req.Role, req.Status)
	if appErr != nil {
		c.JSON(http.StatusOK, response.Error(appErr))
		return
	}

	c.JSON(http.StatusOK, response.Success(user))
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.Error(errors.ErrInvalidParam))
		return
	}

	appErr := h.userUseCase.DeleteUser(c.Request.Context(), id)
	if appErr != nil {
		c.JSON(http.StatusOK, response.Error(appErr))
		return
	}

	c.JSON(http.StatusOK, response.Success(nil))
}

type CreateAttractionRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Address     string  `json:"address" binding:"required"`
	LocationLng float64 `json:"location_lng" binding:"required"`
	LocationLat float64 `json:"location_lat" binding:"required"`
}

func (h *AdminHandler) CreateAttraction(c *gin.Context) {
	var req CreateAttractionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.Error(errors.ErrInvalidParam))
		return
	}

	attraction, err := h.tourUseCase.CreateAttraction(c.Request.Context(), req.Name, req.Description, req.Address, req.LocationLng, req.LocationLat)
	if err != nil {
		c.JSON(http.StatusOK, response.Error(errors.ErrInternalServer))
		return
	}
	c.JSON(http.StatusOK, response.Success(attraction))
}

func (h *AdminHandler) UpdateAttraction(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.Error(errors.ErrInvalidParam))
		return
	}

	type UpdateAttractionRequest struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description" binding:"required"`
		Address     string  `json:"address" binding:"required"`
		LocationLng float64 `json:"location_lng" binding:"required"`
		LocationLat float64 `json:"location_lat" binding:"required"`
	}

	var req UpdateAttractionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.Error(errors.ErrInvalidParam))
		return
	}

	attraction, err := h.tourUseCase.UpdateAttraction(c.Request.Context(), id, req.Name, req.Description, req.Address, req.LocationLng, req.LocationLat)
	if err != nil {
		c.JSON(http.StatusOK, response.Error(errors.ErrInternalServer))
		return
	}
	c.JSON(http.StatusOK, response.Success(attraction))
}

func (h *AdminHandler) DeleteAttraction(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.Error(errors.ErrInvalidParam))
		return
	}

	err = h.tourUseCase.DeleteAttraction(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusOK, response.Error(errors.ErrInternalServer))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil))
}

type BatchIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"required,min=1"`
}

func (h *AdminHandler) BatchDeleteUsers(c *gin.Context) {
	var req BatchIDsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.Error(errors.ErrInvalidParam))
		return
	}

	appErr := h.userUseCase.BatchDeleteUsers(c.Request.Context(), req.IDs)
	if appErr != nil {
		c.JSON(http.StatusOK, response.Error(appErr))
		return
	}

	c.JSON(http.StatusOK, response.Success(nil))
}

type BatchUpdateRoleRequest struct {
	IDs  []uint64 `json:"ids" binding:"required,min=1"`
	Role int8     `json:"role" binding:"required"`
}

func (h *AdminHandler) BatchUpdateRole(c *gin.Context) {
	var req BatchUpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.Error(errors.ErrInvalidParam))
		return
	}

	appErr := h.userUseCase.BatchUpdateRole(c.Request.Context(), req.IDs, req.Role)
	if appErr != nil {
		c.JSON(http.StatusOK, response.Error(appErr))
		return
	}

	c.JSON(http.StatusOK, response.Success(nil))
}
