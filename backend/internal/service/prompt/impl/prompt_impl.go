package impl

import (
	"ai-eino-interview-agent/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// PromptService 提示词服务实现
type PromptService struct {
}

// NewPromptService 创建提示词服务实例
func NewPromptService() *PromptService {
	return &PromptService{}
}

// GetPromptTemplate 获取提示词模板
func (s *PromptService) GetPromptTemplate(ctx context.Context, agentName string) (string, error) {
	key := PromptKeyPrefix + agentName
	client := repository.GetRedis()
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", fmt.Errorf("提示词不存在: %s", agentName)
		}
		return "", err
	}

	var promptData PromptData
	if err := json.Unmarshal([]byte(val), &promptData); err != nil {
		return "", fmt.Errorf("解析提示词失败: %w", err)
	}
	return promptData.Instruction, nil
}
