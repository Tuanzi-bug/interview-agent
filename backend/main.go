package main

// @title AI-Eino 智能面试系统 API文档
// @version 1.0.0
// @description 智能面试助手系统，提供用户管理、简历管理、面试管理等功能
// @contact.name 技术支持
// @contact.url http://example.com/support
// @contact.email support@example.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8000
// @BasePath /api/v1
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

import (
	"ai-eino-interview-agent/internal/config"
	"ai-eino-interview-agent/internal/eino/milvus"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 1. 加载 .env 文件（如果存在）
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
		log.Println("Application will use system environment variables or config.yaml defaults")
	} else {
		log.Println("Successfully loaded .env file")
	}

	// 2. 加载配置文件
	cfg, err := config.LoadConfig("/Users/hongtairen/Documents/ai_project/chat/go-eino-interview-agent/backend/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 7. 初始化 Milvus Manager（向量数据库、Embedding、检索等服务）
	log.Println("Initializing Milvus Manager...")
	ctx := context.Background()
	milvusManager, err := milvus.InitMilvusManager(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Milvus Manager: %v", err)
	}
	// 进行健康检查
	if err := milvusManager.HealthCheck(ctx); err != nil {
		log.Printf("Warning: Milvus health check failed: %v", err)
	}
	log.Println("Milvus Manager initialized successfully")

	tool := milvusManager.RetrieverService
	q1 := "GMP模型是什么"
	opts2 := &milvus.RetrieveOptions{
		Language: milvus.LanguageGolang,
		TopK:     5,
	}
	data, err := tool.RetrieveWithOptions(ctx, q1, opts2)
	if err != nil {
		log.Println("检索失败", err)
	}
	for i, doc := range data {
		fmt.Printf("Document %d:\n", i+1)
		fmt.Printf("  ID: %s\n", doc.ID)
		fmt.Printf("  Content: %s\n", doc.Content)
		// 如果需要输出元数据
		if doc.MetaData != nil {
			fmt.Printf("  MetaData: %+v\n", doc.MetaData)
		}
		fmt.Println("---")
	}

	// 创建一个通道来监听中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 等待中断信号
	<-quit
	log.Println("Shutting down server...")
	// 关闭 Milvus Manager
	if milvusManager != nil {
		if err := milvusManager.Close(); err != nil {
			log.Printf("Warning: Failed to close Milvus Manager: %v", err)
		}
	}

	log.Println("Server exiting")
}
