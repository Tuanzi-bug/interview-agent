package bearAgent

import (
	"ai-eino-interview-agent/internal/config"
	"ai-eino-interview-agent/internal/repository"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
	"golang.org/x/net/context"
)

// 模拟的面试问答记录
const testInterviewRecord = `
# 面试问答记录

## 候选人信息
- 姓名：张三
- 应聘岗位：Go后端开发工程师
- 工作年限：3年

## 问答环节

### Q1: 请介绍一下你在电商订单系统重构项目中的具体工作？
回答：这个项目是我负责的核心项目，将单体订单系统拆分为微服务架构。我负责了订单服务、库存服务、支付服务的设计和实现，使用Redis分布式锁解决库存超卖问题，基于Kafka实现异步通知机制。项目上线后QPS从500提升到3000。

### Q2: 你提到使用Redis分布式锁解决库存超卖，能详细说说具体实现吗？
回答：使用SET key value NX PX命令加锁，用UUID+线程ID作为value防止误删，用Lua脚本保证解锁原子性，用看门狗机制实现锁续期。

### Q3: 说说你对微服务架构的理解？
回答：微服务是将应用构建为小型服务的方法。拆分原则包括单一职责、高内聚低耦合、按业务能力拆分、考虑团队组织、渐进式拆分。

### Q4: 你在团队中扮演什么角色？
回答：担任技术负责人，负责方案设计评审；每周做Code Review；每月组织技术分享；带过2个校招新人。

### Q5: 你对Go语言的GMP调度模型了解多少？
回答：G是协程，M是系统线程，P是处理器。P持有本地队列存放G，M绑定P执行G，还有Work Stealing和Hand Off优化机制。

### Q6: 你平时是如何学习新技术的？
回答：看官方文档、关注技术博客、阅读开源项目源码、在side project中实践。最近关注AI、eBPF、Rust等技术。
`

// initTestEnv 初始化测试环境
func initTestEnv() {
	envPaths := []string{"../../../.env", "../../../../.env", "../../.env"}
	for _, p := range envPaths {
		if err := godotenv.Load(p); err == nil {
			break
		}
	}

	configPath := findTestConfigFile()
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	cfg.ExpandEnv()

	if err = repository.InitDatabase(cfg.Database); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
}

// findTestConfigFile 查找配置文件
func findTestConfigFile() string {
	paths := []string{"../../../config.yaml", "../../../../config.yaml", "../../config.yaml"}
	for _, p := range paths {
		absPath, _ := filepath.Abs(p)
		if _, err := os.Stat(absPath); err == nil {
			return absPath
		}
	}
	return "../../../config.yaml"
}

// TestReportGeneration 测试社招面试报告生成智能体
func TestReportGeneration(t *testing.T) {
	log.Println("========================================")
	log.Println("[启动] 社招面试报告生成智能体测试")
	log.Println("========================================")

	initTestEnv()

	ctx := context.Background()
	var userId uint = 2

	// 1. 创建智能体
	log.Println("[Agent] 正在创建智能体...")
	startTime := time.Now()
	agent := ReportGeneration(userId)
	log.Printf("[Agent] 创建完成，耗时: %v", time.Since(startTime))

	// 2. 创建 Runner
	runner := adk.NewRunner(ctx, adk.RunnerConfig{Agent: agent})

	// 3. 构建消息并运行
	userMsg := schema.UserMessage("请根据以下面试问答记录生成面试评估报告：\n" + testInterviewRecord)
	messages := []adk.Message{userMsg}

	log.Println("[Agent] 开始生成面试报告...")
	startTime = time.Now()
	iter := runner.Run(ctx, messages)

	fmt.Println("\n========== 生成的面试报告 ==========")

	for {
		event, ok := iter.Next()
		if !ok {
			break
		}

		if event.Err != nil {
			t.Fatalf("Agent执行出错: %v", event.Err)
		}

		if event.Output != nil && event.Output.MessageOutput != nil {
			content := event.Output.MessageOutput.Message.Content
			if content != "" {
				fmt.Print(content)
			}
		}
	}

	fmt.Println("\n\n========== 报告生成完成 ==========")
	log.Printf("[Agent] 耗时: %v", time.Since(startTime))
}
