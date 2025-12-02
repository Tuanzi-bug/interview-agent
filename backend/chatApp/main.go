package main

import (
	"ai-eino-interview-agent/chatApp/agent/bearAgent"
	"ai-eino-interview-agent/internal/config"
	"ai-eino-interview-agent/internal/repository"
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
	"golang.org/x/net/context"
)

// 测试用简历
const testResume = `
# 个人简历

## 基本信息
- 姓名：张三
- 工作年限：3年
- 求职意向：Go后端开发工程师

## 教育背景
- 2018.09 - 2022.06  XX大学  计算机科学与技术  本科

## 工作经历

### 某科技有限公司 | 后端开发工程师 | 2022.07 - 至今

**项目一：电商订单系统重构**
- 项目背景：负责公司核心订单系统的技术重构，将单体应用拆分为微服务架构
- 技术栈：Go、Gin、gRPC、MySQL、Redis、Kafka、Docker、K8s
- 主要职责：
  - 设计并实现订单服务、库存服务、支付服务的微服务拆分方案
  - 使用Redis实现分布式锁解决库存超卖问题
  - 基于Kafka实现订单状态变更的异步通知机制
- 项目成果：系统QPS从500提升至3000，订单处理延迟从200ms降低至50ms

**项目二：实时数据分析平台**
- 技术栈：Go、ClickHouse、Elasticsearch、Grafana
- 主要职责：设计高性能数据采集服务，支持每秒10万条数据写入

## 技术技能
- 编程语言：Go（精通）、Python（熟练）
- 数据库：MySQL、Redis、ClickHouse
- 中间件：Kafka、Elasticsearch
- 云原生：Docker、Kubernetes
`

func main() {
	// 初始化配置和数据库
	log.Println("========================================")
	log.Println("[启动] 面试智能体测试程序")
	log.Println("========================================")

	initApp()

	ctx := context.Background()
	var userId uint = 2

	// 1. 创建面试调度智能体
	log.Println("[Agent] 正在创建面试智能体...")
	startTime := time.Now()
	agent := bearAgent.QuestionGeneratorAgent(userId)
	log.Printf("[Agent] 智能体创建完成，耗时: %v", time.Since(startTime))

	// 2. 创建 Runner
	log.Println("[Runner] 正在创建 Agent Runner...")
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent: agent,
	})
	log.Println("[Runner] Runner 创建完成")

	// 3. 打印测试说明
	fmt.Println("\n====== 面试智能体交互测试 ======")
	fmt.Println("命令说明:")
	fmt.Println("  - 输入 'resume' 发送测试简历开始面试")
	fmt.Println("  - 输入 'exit' 或 'quit' 退出程序")
	fmt.Println("  - 输入其他内容作为面试回答")
	fmt.Println("==================================\n")

	reader := bufio.NewReader(os.Stdin)
	messageHistory := []adk.Message{} // 保存对话历史

	for {
		fmt.Print("👤 你: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("[错误] 读取输入失败: %v", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// 退出命令
		if input == "exit" || input == "quit" {
			log.Println("[退出] 用户主动退出程序")
			fmt.Println("感谢使用，再见！")
			break
		}

		// 发送测试简历
		if input == "resume" {
			input = "请分析以下简历并开始面试：\n" + testResume
			log.Println("[测试] 发送测试简历...")
			fmt.Println("📄 已发送测试简历，等待分析...")
		}

		// 记录用户输入
		log.Printf("[用户输入] %s", truncateString(input, 100))

		// 构建消息（包含历史记录）
		userMsg := schema.UserMessage(input)
		messageHistory = append(messageHistory, userMsg)

		// 运行 Agent
		log.Println("[Agent] 开始处理请求...")
		startTime := time.Now()
		iter := runner.Run(ctx, messageHistory)

		fmt.Print("\n🤖 面试官: ")
		var fullResponse strings.Builder
		eventCount := 0

		for {
			event, ok := iter.Next()
			if !ok {
				break
			}
			eventCount++

			if event.Err != nil {
				log.Printf("[错误] Agent 执行出错: %v", event.Err)
				fmt.Printf("\n❌ 出错了: %v\n", event.Err)
				break
			}

			// 打印事件详情（调试用）
			if event.Output != nil {
				// 打印消息输出
				if event.Output.MessageOutput != nil {
					content := event.Output.MessageOutput.Message.Content
					if content != "" {
						fmt.Print(content)
						fullResponse.WriteString(content)
					}
				}
			}
		}

		// 记录响应信息
		elapsed := time.Since(startTime)
		log.Printf("[Agent] 处理完成，耗时: %v，事件数: %d", elapsed, eventCount)
		log.Printf("[Agent响应] %s", truncateString(fullResponse.String(), 200))

		// 保存 AI 响应到历史
		if fullResponse.Len() > 0 {
			aiMsg := schema.AssistantMessage(fullResponse.String(), nil)
			messageHistory = append(messageHistory, aiMsg)
		}

		fmt.Println("\n")
	}
}

// truncateString 截断字符串用于日志显示
func truncateString(s string, maxLen int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) > maxLen {
		return s[:maxLen] + "..."
	}
	return s
}

// initApp 初始化应用（配置、数据库）
func initApp() {
	// 加载 .env 文件
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// 加载配置文件
	configPath := findConfigFile()
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	cfg.ExpandEnv()

	// 初始化数据库
	err = repository.InitDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("数据库初始化成功")

	// // 初始化 Milvus
	// ctx := context.Background()
	// _, err = milvus.InitMilvusManager(ctx, cfg)
	// if err != nil {
	// 	log.Fatalf("Failed to initialize Milvus: %v", err)
	// }
	log.Println("Milvus 初始化成功")
}

// findConfigFile 查找配置文件
func findConfigFile() string {
	// 尝试多个可能的路径
	paths := []string{
		"../config.yaml",
		"../../config.yaml",
		"config.yaml",
	}

	for _, p := range paths {
		absPath, _ := filepath.Abs(p)
		if _, err := os.Stat(absPath); err == nil {
			return absPath
		}
	}

	return "../config.yaml"
}
