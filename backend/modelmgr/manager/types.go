package manager

import (
	"ai-eino-interview-agent/modelmgr/factory"
	"context"
)

// Manager 是内部模型管理器接口。
type Manager interface {
	List(ctx context.Context, opts *ListOptions) ([]*InternalModel, error)
	Get(ctx context.Context, id int64) (*InternalModel, error)
	MGet(ctx context.Context, ids []int64) ([]*InternalModel, error)
}

// ListOptions 是查询模型列表的内部选项。
type ListOptions struct {
	Status     []int
	Limit      int
	Cursor     string
	NameFilter string
}

// InternalModel 是模型的内部完整表示。
type InternalModel struct {
	ID                int64             `yaml:"id"`
	Name              string            `yaml:"name"`
	DisplayName       string            `yaml:"display_name"`
	IconURI           string            `yaml:"icon_uri"`
	IconURL           string            `yaml:"icon_url"`
	Description       *MultilingualText `yaml:"description"`
	DefaultParameters []*Parameter      `yaml:"default_parameters"`
	Meta              ModelMeta         `yaml:",inline"`
}

// ModelMeta 包含模型的元信息和配置。
type ModelMeta struct {
	Protocol   factory.Protocol       `yaml:"protocol"`
	Capability *Capability            `yaml:"capability"`
	ConnConfig *factory.Config        `yaml:"config"`
	Status     int                    `yaml:"status"`
	Extra      map[string]interface{} `yaml:"extra,omitempty"`
}

// MultilingualText 表示多语言文本。
type MultilingualText struct {
	ZH string `json:"zh,omitempty" yaml:"zh,omitempty"`
	EN string `json:"en,omitempty" yaml:"en,omitempty"`
}

// Parameter 定义模型的一个参数。
type Parameter struct {
	Name       string            `yaml:"name"`
	Type       string            `yaml:"type"`
	Min        string            `yaml:"min"`
	Max        string            `yaml:"max"`
	DefaultVal map[string]string `yaml:"default_val"`
	Precision  int               `yaml:"precision"`
}

// Capability 描述模型支持的功能和限制。
type Capability struct {
	// FunctionCall 是否支持函数调用（Function Calling）
	FunctionCall bool `yaml:"support_function_call"`

	// Streaming 是否支持流式输出
	Streaming bool `yaml:"support_streaming"`

	// Vision 是否支持视觉输入（图片理解）
	Vision bool `yaml:"support_vision"`

	// InputModal 支持的输入模态（如 ["text", "image"]）
	InputModal []string `yaml:"input_modal,omitempty"`

	// InputTokens 最大输入 token 数
	InputTokens int `yaml:"max_input_tokens"`

	// JSONMode 是否支持 JSON 模式输出
	JSONMode bool `yaml:"support_json_mode,omitempty"`

	// MaxTokens 最大总 token 数（输入+输出）
	MaxTokens int `yaml:"max_tokens,omitempty"`

	// OutputModal 支持的输出模态（如 ["text"]）
	OutputModal []string `yaml:"output_modal,omitempty"`

	// OutputTokens 最大输出 token 数
	OutputTokens int `yaml:"max_output_tokens"`

	// PrefixCaching 是否支持前缀缓存（提高重复请求的性能）
	PrefixCaching bool `yaml:"support_prefix_caching,omitempty"`

	// Reasoning 是否支持推理模式（如 OpenAI o1）
	Reasoning bool `yaml:"support_reasoning,omitempty"`

	// PrefillResponse 是否支持预填充响应（如 Claude 的 prefill）
	PrefillResponse bool `yaml:"support_prefill_response,omitempty"`
}
