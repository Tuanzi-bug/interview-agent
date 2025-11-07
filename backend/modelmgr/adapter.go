package modelmgr

import (
	"ai-eino-interview-agent/modelmgr/factory"
	"context"
	"errors"
	"time"

	"github.com/cloudwego/eino/schema"
)

// builderAdapter 将 modelmgr.Builder 适配为 factory.Builder
type builderAdapter struct {
	builder Builder
}

func (a *builderAdapter) Build(ctx context.Context, config *factory.Config) (factory.ChatModel, error) {
	buildConfig := &BuildConfig{
		Model:       config.Model,
		APIKey:      config.APIKey,
		BaseURL:     config.BaseURL,
		Temperature: config.Temperature,
		MaxTokens:   config.MaxTokens,
		Timeout:     time.Duration(config.Timeout) * time.Millisecond,
		Extra:       config.Extra,
	}

	chatModel, err := a.builder.Build(ctx, buildConfig)
	if err != nil {
		return nil, err
	}

	return &factoryChatModelAdapter{model: chatModel}, nil
}

// factoryChatModelAdapter 将 modelmgr.ChatModel 适配为 factory.ChatModel
type factoryChatModelAdapter struct {
	model ChatModel
}

func (a *factoryChatModelAdapter) Generate(ctx context.Context, input interface{}) (interface{}, error) {
	messages, ok := input.([]Message)
	if !ok {
		return nil, errors.New("invalid input type, expected []Message")
	}
	return a.model.Generate(ctx, messages)
}

func (a *factoryChatModelAdapter) Stream(ctx context.Context, input interface{}) (interface{}, error) {
	messages, ok := input.([]Message)
	if !ok {
		return nil, errors.New("invalid input type, expected []Message")
	}
	return a.model.Stream(ctx, messages)
}

// chatModelAdapter 将 factory.ChatModel 适配为 modelmgr.ChatModel
type chatModelAdapter struct {
	model factory.ChatModel
}

func (a *chatModelAdapter) Generate(ctx context.Context, messages []Message) (*Response, error) {
	// 将 modelmgr.Message 转换为 eino schema.Message
	schemaMessages := convertToSchemaMessages(messages)

	resp, err := a.model.Generate(ctx, schemaMessages)
	if err != nil {
		return nil, err
	}

	// 转换响应
	return convertFromSchemaResponse(resp)
}

func (a *chatModelAdapter) Stream(ctx context.Context, messages []Message) (<-chan *StreamChunk, error) {
	// 将 modelmgr.Message 转换为 eino schema.Message
	schemaMessages := convertToSchemaMessages(messages)

	stream, err := a.model.Stream(ctx, schemaMessages)
	if err != nil {
		return nil, err
	}

	// 转换流式响应
	return convertFromSchemaStream(stream)
}

// convertToSchemaMessages 将 modelmgr.Message 转换为 eino schema.Message
func convertToSchemaMessages(messages []Message) []*schema.Message {
	schemaMessages := make([]*schema.Message, len(messages))
	for i, msg := range messages {
		schemaMessages[i] = &schema.Message{
			Role:    schema.RoleType(msg.Role),
			Content: msg.Content,
		}
	}
	return schemaMessages
}

// convertFromSchemaResponse 将 eino 响应转换为 modelmgr.Response
func convertFromSchemaResponse(resp interface{}) (*Response, error) {
	schemaMsg, ok := resp.(*schema.Message)
	if !ok {
		return nil, errors.New("invalid response type from eino model")
	}

	return &Response{
		Content:      schemaMsg.Content,
		Role:         string(schemaMsg.Role),
		FinishReason: "stop",
	}, nil
}

// convertFromSchemaStream 将 eino 流式响应转换为 modelmgr StreamChunk 通道
func convertFromSchemaStream(stream interface{}) (<-chan *StreamChunk, error) {
	schemaStream, ok := stream.(*schema.StreamReader[*schema.Message])
	if !ok {
		return nil, errors.New("invalid stream type from eino model")
	}

	chunks := make(chan *StreamChunk, 10)
	go func() {
		defer close(chunks)
		for {
			chunk, err := schemaStream.Recv()
			if err != nil {
				// 流结束或错误
				if chunk == nil {
					break
				}
				chunks <- &StreamChunk{
					Content: "",
					Delta:   "",
					Done:    true,
					Error:   err,
				}
				break
			}

			chunks <- &StreamChunk{
				Content: chunk.Content,
				Delta:   chunk.Content, // 简化处理，实际应该计算增量
				Done:    false,
			}
		}
	}()

	return chunks, nil
}
