package response

import (
	"backend/pkg/errors"
)

// Response 统一返回结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Success 成功响应
func Success(data interface{}) *Response {
	return &Response{
		Code:    errors.Success.Code,
		Message: errors.Success.Message,
		Data:    data,
	}
}

// Error 错误响应
func Error(err *errors.AppError) *Response {
	return &Response{
		Code:    err.Code,
		Message: err.Message,
		Data:    nil,
	}
}

// PageData 分页数据结构
type PageData struct {
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	List     interface{} `json:"list"`
}
