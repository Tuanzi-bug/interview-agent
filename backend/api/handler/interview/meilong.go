package interview

import (
	"ai-eino-interview-agent/chatApp/agent/bearAgent"
	"ai-eino-interview-agent/internal/middleware"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"golang.org/x/net/context"
	"io"
	"log"
	"strings"
	"time"
)

// TestAgentRequest 测试请求结构
type TestAgentRequest struct {
	UserID    uint          `json:"user_id,omitempty"`          // 可选用户ID（优先使用JWT中的用户ID）
	Message   string        `json:"message" binding:"required"` // 用户输入消息
	History   []ChatMessage `json:"history,omitempty"`          // 可选的对话历史
	UseResume bool          `json:"use_resume,omitempty"`       // 是否使用测试简历
}

// ChatMessage 对话消息结构
type ChatMessage struct {
	Role    string `json:"role"`    // user 或 assistant
	Content string `json:"content"` // 消息内容
}

// TestAgentResponse 测试响应结构
type TestAgentResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Response string `json:"response,omitempty"`
	Duration string `json:"duration,omitempty"`
}

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

// StartInterviewMeilong .
// @router /api/interview/meilong/stream [POST]
func StartInterviewMeilong(ctx context.Context, c *app.RequestContext) {
	// 1. 解析请求
	var req TestAgentRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, TestAgentResponse{
			Success: false,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 2. 获取用户ID（优先从JWT获取，否则使用请求体中的user_id）
	jwtUserID := middleware.GetUserID(c)
	log.Printf("[TestAgentStream] JWT中的UserID: %d, 请求体中的UserID: %d", jwtUserID, req.UserID)

	userID := jwtUserID
	if userID == 0 {
		userID = req.UserID
		log.Printf("[TestAgentStream] JWT未获取到用户ID，使用请求体中的UserID: %d", userID)
	}
	if userID == 0 {
		c.JSON(consts.StatusUnauthorized, TestAgentResponse{
			Success: false,
			Message: "未登录或未提供user_id，请在Header中添加 Authorization: Bearer {token}",
		})
		return
	}

	log.Printf("[TestAgentStream] 最终使用的UserID: %d, Message: %s", userID, truncateString(req.Message, 50))

	// 3. 设置SSE响应头
	c.SetStatusCode(consts.StatusOK)
	c.Response.Header.Set("Content-Type", "text/event-stream")
	c.Response.Header.Set("Cache-Control", "no-cache")
	c.Response.Header.Set("Connection", "keep-alive")
	c.Response.Header.Set("Access-Control-Allow-Origin", "*")

	// 4. 创建管道用于流式传输
	pipeReader, pipeWriter := io.Pipe()
	c.SetBodyStream(pipeReader, -1)

	go func() {
		defer pipeWriter.Close()

		// 创建面试智能体
		agent := bearAgent.QuestionGeneratorAgent(userID)

		// 创建 Runner
		runner := adk.NewRunner(ctx, adk.RunnerConfig{
			Agent: agent,
		})

		// 构建消息历史
		// 注意：LLM API 要求每条消息的 content 字段不能为空，否则会返回 400 错误
		var messageHistory []adk.Message
		for _, msg := range req.History {
			// 跳过空内容的消息，防止 "missing messages.content" 错误
			if strings.TrimSpace(msg.Content) == "" {
				log.Printf("[TestAgentStream] 跳过空内容的历史消息, Role: %s", msg.Role)
				continue
			}
			if msg.Role == "user" {
				messageHistory = append(messageHistory, schema.UserMessage(msg.Content))
			} else if msg.Role == "assistant" {
				messageHistory = append(messageHistory, schema.AssistantMessage(msg.Content, nil))
			}
		}

		// 处理当前消息
		userMessage := req.Message
		if req.UseResume {
			userMessage = "请分析以下简历并开始面试：\n" + testResume
		}

		// 确保用户消息不为空，避免 "missing messages.content" 错误
		if strings.TrimSpace(userMessage) == "" {
			writeSSEEvent(pipeWriter, "error", "消息内容不能为空")
			return
		}
		messageHistory = append(messageHistory, schema.UserMessage(userMessage))
		log.Printf("messageHistory:%s", messageHistory)
		// 运行 Agent
		iter := runner.Run(ctx, messageHistory)

		for {
			event, ok := iter.Next()
			if !ok {
				break
			}

			if event.Err != nil {
				writeSSEEvent(pipeWriter, "error", event.Err.Error())
				return
			}

			if event.Output != nil && event.Output.MessageOutput != nil {
				content := event.Output.MessageOutput.Message.Content
				if content != "" {
					writeSSEEvent(pipeWriter, "message", content)
				}
			}
		}

		// 发送完成事件
		writeSSEEvent(pipeWriter, "done", "")
	}()
}

// writeSSEEvent 写入SSE事件
func writeSSEEvent(w io.Writer, eventType string, data string) {
	w.Write([]byte("event: " + eventType + "\n"))
	w.Write([]byte("data: " + data + "\n\n"))
}

// truncateString 截断字符串用于日志显示
func truncateString(s string, maxLen int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) > maxLen {
		return s[:maxLen] + "..."
	}
	return s
}

// HealthCheck 健康检查接口
// @router /api/test/health [GET]
func HealthCheck(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, map[string]interface{}{
		"success": true,
		"message": "服务正常运行",
		"time":    time.Now().Format(time.RFC3339),
	})
}
