package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/internal/usecase"
	"backend/pkg/errors"
	"backend/pkg/response"
)

type AuthHandler struct {
	authUseCase usecase.AuthUseCase
}

func NewAuthHandler(r *gin.Engine, uc usecase.AuthUseCase) {
	handler := &AuthHandler{
		authUseCase: uc,
	}

	authGroup := r.Group("/api/v1/auth")
	{
		authGroup.POST("/register", handler.Register)
		authGroup.POST("/login", handler.Login)
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Nickname string `json:"nickname" binding:"required,max=64"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.Error(errors.ErrInvalidParam))
		return
	}

	user, appErr := h.authUseCase.Register(c.Request.Context(), req.Username, req.Password, req.Nickname)
	if appErr != nil {
		c.JSON(http.StatusOK, response.Error(appErr))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"id":       user.ID,
		"username": user.Username,
	}))
}

type LoginRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	CaptchaId string `json:"captcha_id"`
	Captcha   string `json:"captcha"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.Error(errors.ErrInvalidParam))
		return
	}

	// 验证验证码 (如果有传的话，可以配置为管理员必须验证)
	if req.CaptchaId != "" || req.Captcha != "" {
		if !VerifyCaptcha(req.CaptchaId, req.Captcha) {
			c.JSON(http.StatusOK, response.Error(errors.New(40010, "验证码错误")))
			return
		}
	}

	token, user, appErr := h.authUseCase.Login(c.Request.Context(), req.Username, req.Password)
	if appErr != nil {
		c.JSON(http.StatusOK, response.Error(appErr))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"nickname": user.Nickname,
			"role":     user.Role,
		},
	}))
}
