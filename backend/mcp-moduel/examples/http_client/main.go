package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"ai-eino-interview-agent/mcp-moduel/internal/client"
)

func main() {
	// HTTP 服务器地址
	serverURL := "http://localhost:8080/mcp"

	// 配置 MCP 客户端（使用 HTTP 传输）
	mcpConfigJSON := fmt.Sprintf(`{
		"transport_type": "http",
		"server_url": "%s",
		"timeout": 30,
		"retry_times": 3
	}`, serverURL)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 工具初始化 - 获取工具列表
	toolFetcher := client.NewEinoToolFetcher()
	tools, err := toolFetcher.GetTools(ctx, mcpConfigJSON)
	if err != nil {
		log.Fatalf("获取工具列表失败: %v", err)
	}

	if len(tools) == 0 {
		log.Fatalf("没有可用的工具")
	}

	fmt.Printf("获取到 %d 个工具:\n", len(tools))
	for _, tool := range tools {
		fmt.Printf("  - %s: %s\n", tool.Name, tool.Description)
	}

	// 查找文本处理工具
	textProcessorTool := findTool(tools, "text_processor")
	if textProcessorTool == nil {
		log.Fatalf("未找到 text_processor 工具")
	}

	// 创建工具调用器
	invoker := client.NewMcpCallImpl()

	// 示例 1: 反转字符串
	fmt.Println("\n=== 示例 1: 反转字符串 ===")
	reverseArgs := &client.InvocationArgs{
		Tool: textProcessorTool,
		Body: map[string]any{
			"text":      "Hello, World!",
			"operation": "reverse",
		},
		McpConfig: mcpConfigJSON,
	}

	_, reverseResponse, err := invoker.Do(ctx, reverseArgs)
	if err != nil {
		log.Printf("工具调用失败: %v", err)
	} else {
		fmt.Printf("输入: Hello, World!\n")
		fmt.Printf("操作: reverse\n")
		fmt.Printf("结果: %s\n", reverseResponse)
	}

	// 示例 2: 转换为大写
	fmt.Println("\n=== 示例 2: 转换为大写 ===")
	upperArgs := &client.InvocationArgs{
		Tool: textProcessorTool,
		Body: map[string]any{
			"text":      "hello world",
			"operation": "uppercase",
		},
		McpConfig: mcpConfigJSON,
	}

	_, upperResponse, err := invoker.Do(ctx, upperArgs)
	if err != nil {
		log.Printf("工具调用失败: %v", err)
	} else {
		fmt.Printf("输入: hello world\n")
		fmt.Printf("操作: uppercase\n")
		fmt.Printf("结果: %s\n", upperResponse)
	}

	// 示例 3: 计算字符数
	fmt.Println("\n=== 示例 3: 计算字符数 ===")
	countArgs := &client.InvocationArgs{
		Tool: textProcessorTool,
		Body: map[string]any{
			"text":      "Hello, 世界!",
			"operation": "count",
		},
		McpConfig: mcpConfigJSON,
	}

	_, countResponse, err := invoker.Do(ctx, countArgs)
	if err != nil {
		log.Printf("工具调用失败: %v", err)
	} else {
		fmt.Printf("输入: Hello, 世界!\n")
		fmt.Printf("操作: count\n")
		fmt.Printf("结果: %s\n", countResponse)
	}

	// 示例 4: 统计单词数
	fmt.Println("\n=== 示例 4: 统计单词数 ===")
	wordCountArgs := &client.InvocationArgs{
		Tool: textProcessorTool,
		Body: map[string]any{
			"text":      "This is a sample text for word counting",
			"operation": "word_count",
		},
		McpConfig: mcpConfigJSON,
	}

	_, wordCountResponse, err := invoker.Do(ctx, wordCountArgs)
	if err != nil {
		log.Printf("工具调用失败: %v", err)
	} else {
		fmt.Printf("输入: This is a sample text for word counting\n")
		fmt.Printf("操作: word_count\n")
		fmt.Printf("结果: %s\n", wordCountResponse)
	}

	// 示例 5: 去除空格
	fmt.Println("\n=== 示例 5: 去除空格 ===")
	removeSpacesArgs := &client.InvocationArgs{
		Tool: textProcessorTool,
		Body: map[string]any{
			"text":      "Hello   World   with   spaces",
			"operation": "remove_spaces",
		},
		McpConfig: mcpConfigJSON,
	}

	_, removeSpacesResponse, err := invoker.Do(ctx, removeSpacesArgs)
	if err != nil {
		log.Printf("工具调用失败: %v", err)
	} else {
		fmt.Printf("输入: Hello   World   with   spaces\n")
		fmt.Printf("操作: remove_spaces\n")
		fmt.Printf("结果: %s\n", removeSpacesResponse)
	}

	fmt.Println("\n所有示例执行完成!")
}

// findTool 查找指定名称的工具
func findTool(tools []*client.ToolInfo, name string) *client.ToolInfo {
	for _, tool := range tools {
		if tool.Name == name {
			return tool
		}
	}
	return nil
}
