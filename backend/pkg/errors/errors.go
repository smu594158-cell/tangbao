package errors

import (
	"fmt"
)

// AppError 定义了业务应用中的标准错误结果
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

// New 创建一个新的AppError
func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// 预定义全局统一错误码
var (
	// 基础错误 100xx
	Success           = New(0, "success")
	ErrInvalidParam   = New(10001, "请求参数无效")
	ErrInternalServer = New(10002, "服务器内部错误")
	ErrDatabase       = New(10003, "数据库操作异常")

	// 用户权限相关 400xx
	ErrUserNotFound    = New(40001, "用户不存在")
	ErrPasswordInvalid = New(40002, "密码错误")
	ErrTokenInvalid    = New(40003, "Token无效或已过期")
	ErrUnauthorized    = New(40004, "未登录或登录已过期")
	ErrUserExists      = New(40005, "用户名已存在")
	ErrNotFound        = New(40006, "请求的资源不存在")
	ErrForbidden       = New(40007, "禁止访问：权限不足")

	// 业务逻辑相关 500xx
	ErrChatModelFailed  = New(50001, "AI大模型响应失败")
	ErrMCPServiceFailed = New(50002, "MCP高德服务调用失败")
)
