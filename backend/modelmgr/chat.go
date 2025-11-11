package modelmgr

import (
	"context"
	"fmt"
	"log"
	"time"
)

type ChatModelMgr struct {
}

func NewChatModelMg() *ChatModelMgr {
	return &ChatModelMgr{}
}
func (c *ChatModelMgr) Init(name string) (ChatModel, error) {
	// 1. 初始化模型管理器
	cfg := &Config{
		ConfigPath:     "",
		LoadFromEnv:    false,
		DefaultTimeout: 30 * time.Second,
	}
	mgr, err := New(cfg)
	if err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	ctx := context.Background()

	// 2. 查找可用的豆包模型
	models, err := mgr.ListModels(ctx, &ListOptions{
		Status:     []ModelStatus{StatusActive},
		NameFilter: name,
	})
	if err != nil {
		log.Fatalf("查询模型失败: %v", err)
	}

	if len(models) == 0 {
		log.Fatal("未找到豆包模型，请检查配置文件")
	}

	// 3. 获取模型详细信息
	arkModel := models[0]
	modelDetail, err := mgr.GetModel(ctx, arkModel.ID)
	if err != nil {
		log.Fatalf("获取模型详情失败: %v", err)
	}
	fmt.Println(modelDetail)
	// 4. 创建聊天模型实例
	temperature := float32(0.7)
	maxTokens := 1000
	chatModel, err := mgr.CreateChatModel(ctx, arkModel.ID, &ModelOptions{
		Temperature: &temperature,
		MaxTokens:   &maxTokens,
	})
	if err != nil {
		return nil, err
	}
	return chatModel, nil
}
