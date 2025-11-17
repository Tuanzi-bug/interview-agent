package response

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// Response 统一响应结构体
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// NewResponse 创建新的响应
func NewResponse(code int, message string, data interface{}) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// Success 成功响应
func Success(ctx context.Context, c *app.RequestContext, data interface{}) {
	resp := NewResponse(200, "Success", data)
	c.JSON(consts.StatusOK, resp)
}

// SuccessWithMessage 带自定义消息的成功响应
func SuccessWithMessage(ctx context.Context, c *app.RequestContext, message string, data interface{}) {
	resp := NewResponse(200, message, data)
	c.JSON(consts.StatusOK, resp)
}

// Error 错误响应
func Error(ctx context.Context, c *app.RequestContext, code int, message string) {
	resp := NewResponse(code, message, nil)
	httpCode := getHTTPStatusCode(code)
	c.JSON(httpCode, resp)
}

// ErrorWithData 带数据的错误响应
func ErrorWithData(ctx context.Context, c *app.RequestContext, code int, message string, data interface{}) {
	resp := NewResponse(code, message, data)
	httpCode := getHTTPStatusCode(code)
	c.JSON(httpCode, resp)
}

// BadRequest 400 错误
func BadRequest(ctx context.Context, c *app.RequestContext, message string) {
	Error(ctx, c, 400, message)
}

// Unauthorized 401 未授权
func Unauthorized(ctx context.Context, c *app.RequestContext, message string) {
	Error(ctx, c, 401, message)
}

// Forbidden 403 禁止访问
func Forbidden(ctx context.Context, c *app.RequestContext, message string) {
	Error(ctx, c, 403, message)
}

// NotFound 404 未找到
func NotFound(ctx context.Context, c *app.RequestContext, message string) {
	Error(ctx, c, 404, message)
}

// InternalServerError 500 服务器错误
func InternalServerError(ctx context.Context, c *app.RequestContext, message string) {
	Error(ctx, c, 500, message)
}

// getHTTPStatusCode 根据业务码获取 HTTP 状态码
func getHTTPStatusCode(code int) int {
	switch code {
	case 200:
		return consts.StatusOK
	case 400:
		return consts.StatusBadRequest
	case 401:
		return consts.StatusUnauthorized
	case 403:
		return consts.StatusForbidden
	case 404:
		return consts.StatusNotFound
	case 500:
		return consts.StatusInternalServerError
	default:
		return consts.StatusInternalServerError
	}
}
