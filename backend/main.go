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
	"ai-eino-interview-agent/internal/repository"
	"ai-eino-interview-agent/pkg/eino"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ai-eino-interview-agent/internal/config"
	"ai-eino-interview-agent/internal/eino/milvus"
	"ai-eino-interview-agent/pkg/hertz"

	"github.com/joho/godotenv"
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
	cfg, err := config.LoadConfig("backend/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 3. 展开配置中的环境变量引用（${VAR_NAME}）
	cfg.ExpandEnv()
	log.Println("Environment variables expanded in configuration")

	// 4. 初始化数据库
	log.Println("Initializing database connection...")
	err = repository.InitDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database initialized successfully")

	// 5. 初始化Redis
	log.Println("Initializing Redis connection...")
	err = repository.InitRedis(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	log.Println("Redis initialized successfully")

	// 6. 初始化Eino框架
	log.Println("Initializing Eino framework...")
	err = eino.InitEino(cfg.Eino)
	if err != nil {
		log.Fatalf("Failed to initialize Eino: %v", err)
	}
	log.Println("Eino initialized successfully")

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

	// 8. 初始化Hertz服务器
	log.Println("Initializing Hertz server...")
	err = hertz.InitHertz()
	if err != nil {
		log.Fatalf("Failed to initialize Hertz: %v", err)
	}
	log.Println("Hertz server initialized successfully")

	// 9. 启动服务器
	srv := hertz.GetHertzServer()

	// 创建一个通道来监听中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 在单独的goroutine中启动服务器
	go func() {
		log.Printf("Server is running on %s:%d", cfg.Host, cfg.Port)
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号
	<-quit
	log.Println("Shutting down server...")

	// 创建一个带有超时的上下文，用于优雅关闭
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// 关闭 Milvus Manager
	if milvusManager != nil {
		if err := milvusManager.Close(); err != nil {
			log.Printf("Warning: Failed to close Milvus Manager: %v", err)
		}
	}

	log.Println("Server exiting")
}
