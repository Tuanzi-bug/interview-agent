package modelmgr

import "errors"

// 预定义的错误变量
var (
	// ErrModelNotFound 表示请求的模型不存在
	ErrModelNotFound = errors.New("model not found")
	// ErrProtocolNotSupported 表示不支持的协议类型
	ErrProtocolNotSupported = errors.New("protocol not supported")
	// ErrInvalidConfig 表示配置无效
	ErrInvalidConfig = errors.New("invalid config")
	// ErrBuilderNotFound 表示找不到对应协议的构建器
	ErrBuilderNotFound = errors.New("builder not found")
	// ErrModelCreationFailed 表示模型实例创建失败
	ErrModelCreationFailed = errors.New("model creation failed")
	// ErrInvalidModelID 表示模型 ID 无效
	ErrInvalidModelID = errors.New("invalid model id")
)

// Error 是自定义错误类型。
type Error struct {
	// Code 错误代码，用于程序化识别错误类型
	Code string
	// Message 错误消息，用于展示给用户
	Message string
	// Err 底层错误，保留错误链
	Err error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

// NewError 创建一个新的自定义错误。
func NewError(code, message string, err error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
