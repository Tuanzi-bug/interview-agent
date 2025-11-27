package chat

import (
	"ai-eino-interview-agent/internal/errors"
	usermodel "ai-eino-interview-agent/internal/model"
	"ai-eino-interview-agent/internal/service/common"
	"context"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
)

func CreatOpenAiChatModel(ctx context.Context, userId uint) (model.ToolCallingChatModel, error) {
	result, err := usermodel.UserModelDao.GetDefaultUserModel(int64(userId))
	if err != nil {

		return nil, errors.NewDBError("Failed to get user model", err)
	}
	apiKey, err := common.DecryptAPIKey(result.APIKeyEncrypted)
	if err != nil {
		return nil, errors.NewInternalError("Failed to decrypt API key", err)
	}
	key := apiKey
	//模型名称
	modelName := result.ModelKey
	//api url
	url := result.BaseURL

	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey:  key,
		Model:   modelName,
		BaseURL: url,
	})
	if err != nil {
		return nil, errors.NewOpenAIError("Failed to create OpenAI chat model", err)
	}

	return chatModel, nil
}
