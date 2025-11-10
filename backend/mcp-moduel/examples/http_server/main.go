package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ai-eino-interview-agent/mcp-moduel/internal/server"
	"ai-eino-interview-agent/mcp-moduel/internal/server/tools"
)

func main() {
	// 创建 HTTP 服务器配置
	config := &server.ServerConfig{
		Address:       ":8080",
		ServerName:    "mcp-http-server",
		ServerVersion: "1.0.0",
		EnableCORS:    true,
		AllowedOrigins: []string{
			"*",
		},
	}
	if err := config.Validate(); err != nil {
		log.Fatalf("Invalid server config: %v", err)
	}

	// 创建 HTTP 服务器实例
	srv := server.NewServer(config)

	// 注册默认工具
	defaultTools := tools.GetDefaultTools()
	for _, tool := range defaultTools {
		if t, ok := tool.(server.Tool); ok {
			if err := srv.RegisterTool(t); err != nil {
				log.Fatalf("Failed to register tool: %v", err)
			}
			log.Printf("Registered tool: %s", t.Name())
		}
	}

	// 设置关闭
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// 在 goroutine 中启动服务器
	go func() {
		log.Printf("MCP HTTP Server starting on %s", config.Address)
		log.Printf("Health check endpoint: http://localhost%s/health", config.Address)
		log.Printf("MCP endpoint: http://localhost%s/mcp", config.Address)
		if err := srv.Start(); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// 等待中断信号
	<-sigChan
	log.Println("Shutting down server...")

	// 关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Stop(ctx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}

	log.Println("Server stopped")
}
