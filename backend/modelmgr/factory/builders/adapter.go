package builders

import (
	"ai-eino-interview-agent/modelmgr/factory"
	"context"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

// einoChatModelAdapter 是通用的 eino ChatModel 到 factory.ChatModel 的适配器。
// 可以被 OpenAI、Ark 等所有基于 eino 的构建器共享使用。
type einoChatModelAdapter struct {
	model model.ToolCallingChatModel // eino 的 ChatModel 实例
}

// newEinoChatModelAdapter 创建一个新的 eino ChatModel 适配器。
func newEinoChatModelAdapter(model model.ToolCallingChatModel) factory.ChatModel {
	return &einoChatModelAdapter{model: model}
}

// Generate 生成完整响应（非流式）。
func (a *einoChatModelAdapter) Generate(ctx context.Context, input interface{}) (interface{}, error) {
	// 类型断言：将 input 转换为 []*schema.Message
	messages, ok := input.([]*schema.Message)
	if !ok {
		return nil, factory.ErrInvalidInput
	}

	// 调用底层模型的 Generate 方法
	return a.model.Generate(ctx, messages)
}

// Stream 流式生成响应。
func (a *einoChatModelAdapter) Stream(ctx context.Context, input interface{}) (interface{}, error) {
	// 类型断言：将 input 转换为 []*schema.Message
	messages, ok := input.([]*schema.Message)
	if !ok {
		return nil, factory.ErrInvalidInput
	}

	// 调用底层模型的 Stream 方法
	return a.model.Stream(ctx, messages)
}
