package chat

import (
	usermodel "ai-eino-interview-agent/internal/model"
	"ai-eino-interview-agent/internal/service/common"
	"context"
	"github.com/cloudwego/eino/components/model"

	"log"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/compose"
)

var GlobalToolsNode *compose.ToolsNode

func CreatOpenAiChatModel(ctx context.Context, userId uint) model.ToolCallingChatModel {
	result, err := usermodel.UserModelDao.GetDefaultUserModel(int64(userId))
	if err != nil {
		log.Println(err)
	}
	apiKey, err := common.DecryptAPIKey(result.APIKeyEncrypted)
	if err != nil {
		log.Println(err)
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
		log.Fatalf("create openai chat model failed: %v", err)
	}

	return chatModel
}
