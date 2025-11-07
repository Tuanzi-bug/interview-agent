package factory

import (
	"context"
	"errors"
)

// Factory 是模型工厂接口，负责根据协议创建模型实例。
type Factory interface {
	// Create 根据协议和配置创建聊天模型实例。
	Create(ctx context.Context, protocol Protocol, config *Config) (ChatModel, error)
	// Support 检查是否支持指定的协议。
	Support(protocol Protocol) bool
}

// ChatModel 是内部聊天模型接口，对接底层 SDK。
type ChatModel interface {
	Generate(ctx context.Context, input interface{}) (interface{}, error)
	Stream(ctx context.Context, input interface{}) (interface{}, error)
}

// Builder 是构建器接口，负责创建特定协议的模型实例。
type Builder interface {
	Build(ctx context.Context, config *Config) (ChatModel, error)
}

// BuilderFunc 是构建器函数类型，可以直接作为 Builder 使用。
type BuilderFunc func(ctx context.Context, config *Config) (ChatModel, error)

// Build 实现 Builder 接口，调用函数本身。
func (f BuilderFunc) Build(ctx context.Context, config *Config) (ChatModel, error) {
	return f(ctx, config)
}

// Protocol 表示协议类型。
type Protocol string

// Config 是工厂创建模型时使用的配置。
type Config struct {
	// Model 模型名称（如 "gpt-4", "claude-3-opus"）
	Model string `yaml:"model"`
	// APIKey API 密钥
	APIKey string `yaml:"api_key"`
	// BaseURL API 基础 URL
	BaseURL string `yaml:"base_url"`
	// Temperature 温度参数（0.0-2.0）
	Temperature *float32 `yaml:"temperature"`
	// MaxTokens 最大生成 token 数
	MaxTokens *int `yaml:"max_tokens"`
	// TopP 核采样参数（0.0-1.0）
	TopP *float32 `yaml:"top_p"`
	// TopK Top-K 采样参数
	TopK *int `yaml:"top_k"`
	// FrequencyPenalty 频率惩罚（-2.0 到 2.0）
	FrequencyPenalty *float32 `yaml:"frequency_penalty"`
	// PresencePenalty 存在惩罚（-2.0 到 2.0）
	PresencePenalty *float32 `yaml:"presence_penalty"`
	// Stop 停止序列（遇到这些字符串时停止生成）
	Stop []string `yaml:"stop"`
	// Timeout 请求超时时间（毫秒）
	Timeout int64 `yaml:"timeout"`
	// Extra 额外的配置参数（用于特定模型的扩展）
	Extra map[string]interface{} `yaml:"extra"`
}

// 预定义的错误变量
var (
	ErrInvalidInput         = errors.New("invalid input")
	ErrProtocolNotSupported = errors.New("protocol not supported")
)
