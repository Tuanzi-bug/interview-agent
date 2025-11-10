package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"ai-eino-interview-agent/mcp-moduel/internal/client"
)

func main() {
	serverExe := findServerExe()
	mcpConfigJSON := fmt.Sprintf(`{
		"transport_type": "stdio",
		"command": "%s",
		"command_args": [],
		"command_env": [],
		"timeout": 30,
		"retry_times": 3
	}`, serverExe)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	//工具初始化
	toolFetcher := client.NewEinoToolFetcher()
	tools, err := toolFetcher.GetTools(ctx, mcpConfigJSON)
	if err != nil {
		log.Fatalf("获取工具列表失败: %v", err)
	}
	if len(tools) == 0 {
		log.Fatalf("没有可用的工具")
	}
	fmt.Printf("获取到 %d 个工具\n", len(tools))
	for _, tool := range tools {
		fmt.Printf("  - %s: %s\n", tool.Name, tool.Description)
	}

	selectedTool := findTool(tools, "echo")
	if selectedTool == nil {
		selectedTool = tools[0]
	}

	invoker := client.NewMcpCallImpl()
	args := &client.InvocationArgs{
		Tool: selectedTool,
		Body: map[string]any{
			"message": "Hello, MCP Stdio Server!",
		},
		McpConfig: mcpConfigJSON,
	}

	_, response, err := invoker.Do(ctx, args)
	if err != nil {
		log.Fatalf("工具调用失败: %v", err)
	}

	fmt.Printf("\n工具 '%s' 调用结果:\n%s\n", selectedTool.Name, response)

	if calcTool := findTool(tools, "calculate"); calcTool != nil {
		calcArgs := &client.InvocationArgs{
			Tool: calcTool,
			Body: map[string]any{
				"operation": "add",
				"a":         10,
				"b":         20,
			},
			McpConfig: mcpConfigJSON,
		}

		if _, calcResponse, err := invoker.Do(ctx, calcArgs); err == nil {
			fmt.Printf("\n计算器工具调用结果:\n%s\n", calcResponse)
		}
	}
}

// 获取服务器可执行文件路径
func findServerExe() string {
	if envPath := os.Getenv("STDIO_SERVER_PATH"); envPath != "" {
		if _, err := os.Stat(envPath); err == nil {
			return envPath
		}
		log.Fatalf("STDIO_SERVER_PATH is set but file not found: %s", envPath)
	}

	var possiblePaths []string

	if _, currentFile, _, ok := runtime.Caller(0); ok {
		currentDir := filepath.Dir(currentFile)
		possiblePaths = append(possiblePaths,
			filepath.Join(currentDir, "..", "stdio_server", "stdio_server"),
			filepath.Join(filepath.Dir(currentDir), "stdio_server", "stdio_server"),
		)
	}

	if wd, err := os.Getwd(); err == nil {
		possiblePaths = append(possiblePaths,
			filepath.Join(wd, "..", "stdio_server", "stdio_server"),
			filepath.Join(wd, "stdio_server", "stdio_server"),
			filepath.Join(wd, "examples", "stdio_server", "stdio_server"),
			filepath.Join(wd, "backend", "mcp-moduel", "examples", "stdio_server", "stdio_server"),
			filepath.Join(wd, "mcp-moduel", "examples", "stdio_server", "stdio_server"),
		)
	}

	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		possiblePaths = append(possiblePaths,
			filepath.Join(exeDir, "stdio_server"),
			filepath.Join(exeDir, "..", "stdio_server", "stdio_server"),
			filepath.Join(exeDir, "..", "examples", "stdio_server", "stdio_server"),
		)
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			if absPath, err := filepath.Abs(path); err == nil {
				return absPath
			}
		}
	}

	log.Fatalf("找不到 stdio_server，请设置 STDIO_SERVER_PATH 环境变量或先构建服务器")
	return ""
}

// 获取工具
func findTool(tools []*client.ToolInfo, name string) *client.ToolInfo {
	for _, tool := range tools {
		if tool.Name == name {
			return tool
		}
	}
	return nil
}
