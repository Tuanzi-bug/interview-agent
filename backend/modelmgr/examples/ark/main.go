package main

import (
	"ai-eino-interview-agent/modelmgr"
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	// 1. 初始化模型管理器
	cfg := &modelmgr.Config{
		ConfigPath:     "",
		LoadFromEnv:    false,
		DefaultTimeout: 30 * time.Second,
	}
	mgr, err := modelmgr.New(cfg)
	if err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	ctx := context.Background()

	// 2. 查找可用的豆包模型
	models, err := mgr.ListModels(ctx, &modelmgr.ListOptions{
		Status:     []modelmgr.ModelStatus{modelmgr.StatusActive},
		NameFilter: "doubao",
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
	chatModel, err := mgr.CreateChatModel(ctx, arkModel.ID, &modelmgr.ModelOptions{
		Temperature: &temperature,
		MaxTokens:   &maxTokens,
	})
	if err != nil {
		return
	}

	// 5. 发送单次对话请求
	messages := []modelmgr.Message{
		{
			Role:    modelmgr.RoleSystem,
			Content: "你是豆包，字节跳动开发的 AI 助手",
		},
		{
			Role:    modelmgr.RoleUser,
			Content: "你好！请用一句话介绍你自己。",
		},
	}

	response, err := chatModel.Generate(ctx, messages)
	if err != nil {
		log.Fatalf("生成响应失败: %v", err)
	}
	fmt.Println(response.Content)
	// 6. 发起流式对话请求
	streamMessages := []modelmgr.Message{
		{
			Role:    modelmgr.RoleUser,
			Content: "请用三个要点介绍火山引擎 Ark 的特点",
		},
	}

	streamChan, err := chatModel.Stream(ctx, streamMessages)
	if err != nil {
		log.Fatalf("启动流式生成失败: %v", err)
	}

	// 接收流式响应数据
	for chunk := range streamChan {
		if chunk.Error != nil {
			log.Fatalf("流式生成错误: %v", chunk.Error)
		}
		fmt.Print(chunk.Content)
	}

	// 7. 进行多轮对话交互
	conversation := []modelmgr.Message{
		{
			Role:    modelmgr.RoleUser,
			Content: "什么是 Go 语言？",
		},
	}

	// 第一轮对话
	resp1, err := chatModel.Generate(ctx, conversation)
	if err != nil {
		log.Fatalf("第一轮对话失败: %v", err)
	}

	// 添加响应到对话历史，继续第二轮对话
	conversation = append(conversation, modelmgr.Message{
		Role:    modelmgr.RoleAssistant,
		Content: resp1.Content,
	})

	conversation = append(conversation, modelmgr.Message{
		Role:    modelmgr.RoleUser,
		Content: "它的主要优势是什么？",
	})

	resp2, err := chatModel.Generate(ctx, conversation)
	if err != nil {
		log.Fatalf("第二轮对话失败: %v", err)
	}
	fmt.Print(resp2.Content)
}
