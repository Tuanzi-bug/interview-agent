package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"ai-eino-interview-agent/mcp-moduel/examples/external_tool"
	"ai-eino-interview-agent/mcp-moduel/internal/server"
)

func main() {
	// 获取当前工作目录
	wd, _ := os.Getwd()

	// 获取源代码文件的真实路径（不受执行目录影响）
	_, currentSourceFile, _, ok := runtime.Caller(0)
	if !ok {
		log.Printf("Warning: Could not get source file path")
	}

	// 基于源代码文件路径查找 .env 文件
	// main.go 在 server/ 目录，.env 在 external_tool/ 目录（上一级）
	var envPaths []string
	if currentSourceFile != "" {
		sourceDir := filepath.Dir(currentSourceFile) // server/ 目录
		externalToolDir := filepath.Dir(sourceDir)   // external_tool/ 目录
		envPaths = append(envPaths, filepath.Join(externalToolDir, ".env"))
	}

	// 也尝试从其他常见位置加载
	envPaths = append(envPaths,
		".env",
		filepath.Join(wd, ".env"),
		"../.env",
		filepath.Join(wd, "..", ".env"),
		filepath.Join(wd, "..", "..", ".env"),
	)

	loaded := false
	var loadedPath string
	for _, envPath := range envPaths {
		absPath, _ := filepath.Abs(envPath)
		// 使用 Overload 而不是 Load，确保 .env 文件中的值会覆盖已存在的环境变量
		if err := godotenv.Overload(envPath); err == nil {
			log.Printf("Loaded .env file from: %s (absolute: %s)", envPath, absPath)
			loadedPath = absPath
			loaded = true
			break
		}
	}
	if !loaded {
		log.Printf("Warning: Could not load .env file from any of: %v (continuing with system environment variables)", envPaths)
		log.Printf("Current working directory: %s", wd)
		if currentSourceFile != "" {
			log.Printf("Source file directory: %s", filepath.Dir(currentSourceFile))
		}
	} else {
		log.Printf("Successfully loaded .env file from: %s", loadedPath)
	}

	// 创建 HTTP 服务器配置
	config := &server.ServerConfig{
		Address:       ":8081",
		ServerName:    "mcp-weather-server",
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

	// 从环境变量读取 API Key（优先从 .env 文件加载）
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		// 调试信息：显示调试信息
		log.Printf("DEBUG: Current working directory: %s", wd)
		if currentSourceFile != "" {
			log.Printf("DEBUG: Source file directory: %s", filepath.Dir(currentSourceFile))
		}
		if loaded {
			log.Printf("DEBUG: .env file was loaded from: %s", loadedPath)
			log.Printf("DEBUG: But WEATHER_API_KEY is still empty. This might indicate a problem with the .env file format.")
		}
		log.Fatalf(`
ERROR: WEATHER_API_KEY environment variable is not set!

To fix this:
1. Get a free API key from https://openweathermap.org/api
2. Sign up for a free account
3. Copy your API key
4. Create a .env file in the external_tool directory with:
   WEATHER_API_KEY=your_api_key_here
5. Restart the server

Or you can run the server with:
   WEATHER_API_KEY=your_api_key_here ./server
`)
	}

	log.Printf("Weather API Key loaded (length: %d)", len(apiKey))

	// 注册外部工具 - 天气查询工具
	weatherTool := external_tool.NewWeatherTool(apiKey)
	if err := srv.RegisterTool(weatherTool); err != nil {
		log.Fatalf("Failed to register weather tool: %v", err)
	}
	log.Printf("Registered tool: %s", weatherTool.Name())

	// 设置关闭信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// 在 goroutine 中启动服务器
	go func() {
		log.Printf("MCP Weather Server starting on %s", config.Address)
		log.Printf("Health check endpoint: http://localhost%s/health", config.Address)
		log.Printf("MCP endpoint: http://localhost%s/mcp", config.Address)
		if err := srv.Start(); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// 等待中断信号
	<-sigChan
	log.Println("Shutting down server...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Stop(ctx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}

	log.Println("Server stopped")
}
