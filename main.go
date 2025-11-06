package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ai-eino-interview-agent/internal/config"
	"ai-eino-interview-agent/internal/repository"
	"ai-eino-interview-agent/pkg/eino"
	"ai-eino-interview-agent/pkg/hertz"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/cors"
)

func main() {
	// 1. 加载配置
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. 初始化数据库
	log.Println("Initializing database connection...")
	err = repository.InitDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database initialized successfully")

	// 3. 初始化Redis
	log.Println("Initializing Redis connection...")
	err = repository.InitRedis(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	log.Println("Redis initialized successfully")

	// 4. 初始化Eino框架
	log.Println("Initializing Eino framework...")
	err = eino.InitEino(cfg.Eino)
	if err != nil {
		log.Fatalf("Failed to initialize Eino: %v", err)
	}
	log.Println("Eino initialized successfully")

	// 5. 初始化Hertz服务器
	log.Println("Initializing Hertz server...")
	err = hertz.InitHertz(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Hertz: %v", err)
	}
	log.Println("Hertz server initialized successfully")

	// 6. 启动服务器
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}