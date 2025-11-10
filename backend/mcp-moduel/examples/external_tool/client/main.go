package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"ai-eino-interview-agent/mcp-moduel/internal/client"
)

func main() {
	// HTTP 服务器地址（使用外部工具服务器的端口）
	serverURL := "http://localhost:8081/mcp"

	// 配置 MCP 客户端（使用 HTTP 传输）
	// 注意：这里的 api_key 字段仅用于示例，实际 API Key 需要在服务器端通过环境变量设置
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

	// 查找天气工具
	weatherTool := findTool(tools, "get_weather")
	if weatherTool == nil {
		log.Fatalf("未找到 get_weather 工具")
	}

	// 创建工具调用器
	invoker := client.NewMcpCallImpl()

	// 示例 1: 查询北京的天气
	fmt.Println("\n=== 示例 1: 查询北京的天气 ===")
	beijingArgs := &client.InvocationArgs{
		Tool: weatherTool,
		Body: map[string]any{
			"city":  "Beijing",
			"units": "metric",
		},
		McpConfig: mcpConfigJSON,
	}

	_, beijingResponse, err := invoker.Do(ctx, beijingArgs)
	if err != nil {
		log.Printf("工具调用失败: %v", err)
	} else {
		fmt.Printf("查询城市: Beijing\n")
		fmt.Printf("结果: %s\n", beijingResponse)
	}

	// 示例 2: 查询纽约的天气（使用华氏度）
	fmt.Println("\n=== 示例 2: 查询纽约的天气（华氏度） ===")
	newYorkArgs := &client.InvocationArgs{
		Tool: weatherTool,
		Body: map[string]any{
			"city":  "New York",
			"units": "imperial",
		},
		McpConfig: mcpConfigJSON,
	}

	_, newYorkResponse, err := invoker.Do(ctx, newYorkArgs)
	if err != nil {
		log.Printf("工具调用失败: %v", err)
	} else {
		fmt.Printf("查询城市: New York\n")
		fmt.Printf("结果: %s\n", newYorkResponse)
	}

	// 示例 3: 查询伦敦的天气
	fmt.Println("\n=== 示例 3: 查询伦敦的天气 ===")
	londonArgs := &client.InvocationArgs{
		Tool: weatherTool,
		Body: map[string]any{
			"city":  "London",
			"units": "metric",
		},
		McpConfig: mcpConfigJSON,
	}

	_, londonResponse, err := invoker.Do(ctx, londonArgs)
	if err != nil {
		log.Printf("工具调用失败: %v", err)
	} else {
		fmt.Printf("查询城市: London\n")
		fmt.Printf("结果: %s\n", londonResponse)
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
