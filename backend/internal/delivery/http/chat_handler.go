package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/internal/usecase"
	"backend/pkg/errors"
	"backend/pkg/response"
)

// ChatHandler 处理智能问答相关的HTTP 请求
type ChatHandler struct {
	chatUseCase usecase.ChatUseCase
}

// NewChatHandler 实例化并注册路由
func NewChatHandler(r *gin.Engine, uc usecase.ChatUseCase, authMiddleware gin.HandlerFunc) {
	handler := &ChatHandler{
		chatUseCase: uc,
	}

	// 注册路由组，需前置 Auth 中间件
	chatGroup := r.Group("/api/v1/chat")
	chatGroup.Use(authMiddleware) // 保护聊天接口
	{
		chatGroup.POST("/message", handler.SendMessage)
	}
}

// ChatMessageRequest 定义请求结构
type ChatMessageRequest struct {
	SessionID string `json:"session_id" binding:"required,min=8,max=64"`
	Content   string `json:"content" binding:"required,min=1,max=1000"`
}

// ChatMessageResponse 定义返回结构
type ChatMessageResponse struct {
	Reply string `json:"reply"`
}

// SendMessage 发送消息并获取 AI 回复
func (h *ChatHandler) SendMessage(c *gin.Context) {
	var req ChatMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.Error(errors.ErrInvalidParam))
		return
	}

	// JWT Token 解析获取 UserID
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusOK, response.Error(errors.ErrUnauthorized))
		return
	}
	userID := userIDVal.(uint64)

	// 调用应用层获取回答
	reply, appErr := h.chatUseCase.SendMessage(c.Request.Context(), userID, req.SessionID, req.Content)
	if appErr != nil {
		c.JSON(http.StatusOK, response.Error(appErr))
		return
	}

	// 成功返回
	c.JSON(http.StatusOK, response.Success(ChatMessageResponse{
		Reply: reply,
	}))
}
