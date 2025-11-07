package modelmgr

import (
	"context"
	"time"
)

// ModelManager 是模型管理器的核心接口，负责管理模型元数据和创建模型实例。
type ModelManager interface {
	ListModels(ctx context.Context, opts ...*ListOptions) ([]*Model, error)
	GetModel(ctx context.Context, id int64) (*Model, error)
	GetModels(ctx context.Context, ids []int64) ([]*Model, error)
	CreateChatModel(ctx context.Context, modelID int64, opts ...*ModelOptions) (ChatModel, error)
}

// ChatModel 是聊天模型的统一接口，屏蔽了不同 AI 服务提供商的实现差异。
type ChatModel interface {
	Generate(ctx context.Context, messages []Message) (*Response, error)
	Stream(ctx context.Context, messages []Message) (<-chan *StreamChunk, error)
}

// Builder 是模型构建器接口，负责根据配置创建具体的 ChatModel 实例。
type Builder interface {
	Build(ctx context.Context, config *BuildConfig) (ChatModel, error)
}

// Logger 是日志接口，用于记录模型管理器的运行日志。
type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
}

// Config 是 ModelManager 的初始化配置。
type Config struct {
	// ConfigPath 配置文件路径，指向包含模型 YAML 配置文件的目录。
	// 支持以下几种方式：
	// 1. 绝对路径："/path/to/config"
	// 2. 相对路径：相对于当前工作目录的路径
	// 3. 空字符串：自动按以下顺序查找配置目录：
	//    - 环境变量 MODELMGR_CONFIG_PATH
	//    - ./config
	//    - ./modelmgr/config
	//    - 项目根目录下的 backend/modelmgr/config
	ConfigPath string
	// TemplatePath 模板文件路径，用于环境变量模式。
	TemplatePath string
	// LoadFromEnv 是否从环境变量加载敏感配置。
	LoadFromEnv bool
	// CustomBuilders 自定义构建器映射。
	CustomBuilders map[Protocol]Builder
	// Logger 日志记录器，用于记录运行时日志。
	Logger Logger
	// DefaultTimeout 默认超时时间，用于模型请求。
	DefaultTimeout time.Duration
}

// ListOptions 是查询模型列表时的过滤和分页选项。
type ListOptions struct {
	// Status 状态过滤，只返回指定状态的模型。
	// 例如：[]ModelStatus{StatusActive} 只返回激活的模型。
	Status []ModelStatus
	// Limit 返回结果的最大数量。
	// 如果不设置或设置为 0，使用默认值（通常为 100）。
	Limit int
	// Cursor 分页游标，用于获取下一页结果。
	// 从上一次查询的响应中获取。
	Cursor string
	// NameFilter 名称过滤，支持模糊匹配。
	// 例如："gpt" 会匹配 "gpt-4", "gpt-3.5-turbo" 等。
	NameFilter string
}

// ModelOptions 是创建聊天模型实例时的运行时选项。
type ModelOptions struct {
	// Temperature 温度参数，控制输出的随机性。
	// 范围：0.0-2.0
	// - 0.0: 确定性输出（每次相同）
	// - 1.0: 平衡的随机性（推荐）
	// - 2.0: 高度随机
	Temperature *float32
	// MaxTokens 最大生成 token 数。
	// 限制模型生成的最大长度，防止过长的输出。
	MaxTokens *int
	// TopP 核采样参数（nucleus sampling）。
	// 范围：0.0-1.0
	// 与 Temperature 类似，但使用不同的采样策略。
	TopP *float32
	// CustomConfig 自定义配置，用于传递特定模型的额外参数。
	// 例如：frequency_penalty, presence_penalty 等。
	CustomConfig map[string]interface{}
}

// BuildConfig 是构建器创建模型实例时使用的配置。
type BuildConfig struct {
	// Protocol 协议类型（openai、claude 等）
	Protocol Protocol
	// Model 模型名称（如 "gpt-4", "claude-3-opus" 等）
	Model string
	// APIKey API 密钥，用于身份验证
	APIKey string
	// BaseURL API 基础 URL，用于自定义端点或代理
	BaseURL string
	// Temperature 温度参数
	Temperature *float32
	// MaxTokens 最大 token 数
	MaxTokens *int
	// Timeout 请求超时时间
	Timeout time.Duration
	// Extra 额外的配置参数，由具体的构建器解释
	Extra map[string]interface{}
}

// Model 表示一个 AI 模型的元数据信息（对外暴露的简化结构）。
type Model struct {
	// ID 模型的唯一标识符
	ID int64 `json:"id"`

	// Name 模型的内部名称（如 "gpt-4", "claude-3-opus"）
	Name string `json:"name"`

	// DisplayName 模型的显示名称（用于 UI 展示）
	DisplayName string `json:"display_name"`

	// Description 模型的描述信息
	Description string `json:"description"`

	// Protocol 模型使用的协议（openai、claude 等）
	Protocol Protocol `json:"protocol"`

	// Status 模型的当前状态（激活、待定、已删除等）
	Status ModelStatus `json:"status"`

	// Capability 模型的能力信息（可选）
	Capability *ModelCapability `json:"capability"`

	// IconURL 模型的图标 URL（用于 UI 展示）
	IconURL string `json:"icon_url"`
}

// ModelCapability 描述模型支持的功能和限制。
type ModelCapability struct {
	// SupportFunctionCall 是否支持函数调用（Function Calling）
	SupportFunctionCall bool `json:"support_function_call"`
	// SupportStreaming 是否支持流式输出
	SupportStreaming bool `json:"support_streaming"`
	// SupportVision 是否支持视觉输入（图片）
	SupportVision bool `json:"support_vision"`
	// MaxInputTokens 最大输入 token 数
	MaxInputTokens int `json:"max_input_tokens"`
	// MaxOutputTokens 最大输出 token 数
	MaxOutputTokens int `json:"max_output_tokens"`
	// SupportedModalities 支持的模态类型（如 "text", "image", "audio"）
	SupportedModalities []string `json:"supported_modalities"`
}

// Message 表示对话中的一条消息。
type Message struct {
	// Role 消息角色（"system", "user", "assistant"）
	Role string `json:"role"`
	// Content 消息内容（文本）
	Content string `json:"content"`
	// Name 消息发送者的名称（可选，用于多用户场景）
	Name string `json:"name,omitempty"`
	// Extra 额外的自定义字段（如图片 URL、函数调用等）
	Extra map[string]interface{} `json:"extra,omitempty"`
}

// Response 表示模型生成的完整响应（非流式）。
type Response struct {
	// Content 生成的文本内容
	Content string `json:"content"`
	// Role 响应的角色（通常是 "assistant"）
	Role string `json:"role"`
	// FinishReason 完成原因
	// 可能的值：
	//   - "stop": 正常完成
	//   - "length": 达到最大 token 限制
	//   - "content_filter": 内容被过滤
	FinishReason string `json:"finish_reason"`
	// Usage token 使用统计（可选）
	Usage *Usage `json:"usage,omitempty"`
	// Extra 额外的响应字段（如函数调用结果等）
	Extra map[string]interface{} `json:"extra,omitempty"`
}

// StreamChunk 表示流式响应中的一个数据块。
type StreamChunk struct {
	// Content 累积的完整内容（从开始到当前）
	Content string `json:"content"`
	// Delta 本次新增的内容（增量）
	Delta string `json:"delta"`
	// Done 是否已完成生成
	Done bool `json:"done"`
	// Error 如果发生错误，包含错误信息
	Error error `json:"error,omitempty"`
	// Extra 额外的响应字段
	Extra map[string]interface{} `json:"extra,omitempty"`
}

// Usage 记录 API 调用的 token 使用情况。
type Usage struct {
	// PromptTokens 输入消息使用的 token 数
	PromptTokens int `json:"prompt_tokens"`
	// CompletionTokens 生成内容使用的 token 数
	CompletionTokens int `json:"completion_tokens"`
	// TotalTokens 总 token 数（PromptTokens + CompletionTokens）
	TotalTokens int `json:"total_tokens"`
}

// Protocol 表示 AI 模型使用的协议类型。
type Protocol string

// 协议类型常量
const (
	// ProtocolOpenAI 表示 OpenAI 兼容的 API 协议
	ProtocolOpenAI Protocol = "openai"
	// ProtocolClaude 表示 Anthropic Claude 的 API 协议
	ProtocolClaude Protocol = "claude"
	// ProtocolDeepseek 表示 Deepseek 的 API 协议
	ProtocolDeepseek Protocol = "deepseek"
	// ProtocolGemini 表示 Google Gemini 的 API 协议
	ProtocolGemini Protocol = "gemini"
	// ProtocolArk 表示火山引擎 Ark 的 API 协议
	ProtocolArk Protocol = "ark"
	// ProtocolOllama 表示 Ollama 本地模型的 API 协议
	ProtocolOllama Protocol = "ollama"
	// ProtocolQwen 表示阿里云通义千问的 API 协议
	ProtocolQwen Protocol = "qwen"
)

// ModelStatus 表示模型的当前状态。
type ModelStatus int

// 模型状态常量
const (
	// StatusUnknown 表示未知状态（默认值）
	StatusUnknown ModelStatus = 0
	// StatusActive 表示模型已激活，可以正常使用
	StatusActive ModelStatus = 1
	// StatusPending 表示模型待定，可能正在配置或测试中
	StatusPending ModelStatus = 2
	// StatusDeleted 表示模型已删除，不应再使用
	StatusDeleted ModelStatus = 3
)

// 消息角色常量
const (
	RoleSystem    = "system"    // 系统消息，用于设置对话的行为和背景
	RoleUser      = "user"      // 用户消息
	RoleAssistant = "assistant" // 助手消息（模型的回复）
)
