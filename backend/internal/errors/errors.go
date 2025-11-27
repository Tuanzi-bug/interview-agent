package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode 错误码类型
type ErrorCode string

const (
	ErrCodeInternal     ErrorCode = "INTERNAL_ERROR"
	ErrCodeInvalidParam ErrorCode = "INVALID_PARAM"
	ErrCodeNotFound     ErrorCode = "NOT_FOUND"
	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden    ErrorCode = "FORBIDDEN"
	ErrCodeDBError      ErrorCode = "DATABASE_ERROR"
	ErrCodeRedisError   ErrorCode = "REDIS_ERROR"
	ErrCodeMilvusError  ErrorCode = "MILVUS_ERROR"
	ErrCodeModelError   ErrorCode = "MODEL_ERROR"
	ErrCodeValidation   ErrorCode = "VALIDATION_ERROR"
	ErrCodeFeishuError  ErrorCode = "FEISHU_ERROR"
	ErrCodeOpenAIError  ErrorCode = "OPENAI_ERROR"
)

// AppError 应用错误结构
// 简化设计：只保留必要的字段
type AppError struct {
	Code       ErrorCode `json:"code"`    // 错误码，用于程序识别
	Message    string    `json:"message"` // 用户友好的错误消息
	HTTPStatus int       `json:"-"`       // HTTP状态码，不返回给客户端
	Err        error     `json:"-"`       // 底层错误，用于错误链追踪
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap 返回底层错误（用于错误链追踪）
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError 创建新的应用错误（不带底层错误）
func NewAppError(code ErrorCode, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

// WrapError 包装现有错误
// 作用：保留原始错误信息，同时添加业务层面的错误码和消息
// 例如：数据库错误 -> 包装成 DBError，保留原始错误用于调试
func WrapError(err error, code ErrorCode, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Err:        err, // 保留原始错误，可以通过 Unwrap() 获取
	}
}

// 预定义错误构造函数
func NewInternalError(message string, err error) *AppError {
	return WrapError(err, ErrCodeInternal, message, http.StatusInternalServerError)
}

func NewInvalidParamError(message string) *AppError {
	return NewAppError(ErrCodeInvalidParam, message, http.StatusBadRequest)
}

func NewNotFoundError(resource string) *AppError {
	msg := fmt.Sprintf("%s not found", resource)
	return NewAppError(ErrCodeNotFound, msg, http.StatusNotFound)
}

func NewDBError(message string, err error) *AppError {
	return WrapError(err, ErrCodeDBError, message, http.StatusInternalServerError)
}

func NewValidationError(message string) *AppError {
	return NewAppError(ErrCodeValidation, message, http.StatusBadRequest)
}

func NewMilvusError(message string, err error) *AppError {
	return WrapError(err, ErrCodeMilvusError, message, http.StatusInternalServerError)
}

func NewFeishuError(message string, err error) *AppError {
	return WrapError(err, ErrCodeFeishuError, message, http.StatusBadGateway)
}

func NewOpenAIError(message string, err error) *AppError {
	return WrapError(err, ErrCodeOpenAIError, message, http.StatusBadGateway)
}

func NewModelError(message string, err error) *AppError {
	return WrapError(err, ErrCodeModelError, message, http.StatusInternalServerError)
}

// As 检查错误链中是否存在指定类型的错误
func As(err error, target interface{}) bool {
	if err == nil {
		return false
	}
	if appErr, ok := err.(*AppError); ok {
		if t, ok := target.(**AppError); ok {
			*t = appErr
			return true
		}
	}
	// 检查错误链
	if unwrapped := Unwrap(err); unwrapped != nil {
		return As(unwrapped, target)
	}
	return false
}

// Unwrap 返回底层错误
func Unwrap(err error) error {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Unwrap()
	}
	return nil
}
