package interview

import (
	"ai-eino-interview-agent/internal/middleware"
	"ai-eino-interview-agent/internal/model"
	"io"
	"log"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"golang.org/x/net/context"
)

// TestAgentRequest 测试请求结构
type TestAgentRequest struct {
	UserID    uint          `json:"user_id,omitempty"`          // 可选用户ID（优先使用JWT中的用户ID）
	Message   string        `json:"message" binding:"required"` // 用户输入消息
	History   []ChatMessage `json:"history,omitempty"`          // 可选的对话历史
	UseResume bool          `json:"use_resume,omitempty"`       // 是否使用简历
	ResumeID  uint64        `json:"resume_id,omitempty"`        // 简历ID（前端传入）
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

// ========== 新版持久会话接口 ==========

// StartPersistentInterviewRequest 开始持久面试请求
type StartPersistentInterviewRequest struct {
	UserID    uint   `json:"user_id,omitempty"`    // 可选用户ID
	UseResume bool   `json:"use_resume,omitempty"` // 是否使用简历
	ResumeID  uint64 `json:"resume_id,omitempty"`  // 简历ID
}

// SendMessageRequest 发送消息请求
type SendMessageRequest struct {
	SessionID string `json:"session_id" binding:"required"` // 会话ID
	Message   string `json:"message" binding:"required"`    // 用户消息
}

// EndInterviewRequest 结束面试请求
type EndInterviewRequest struct {
	SessionID string `json:"session_id" binding:"required"` // 会话ID
}

// StartPersistentInterview 开始持久面试会话
// 这个接口会建立一个持久的SSE连接，直到面试结束或超时
// @router /api/interview/persistent/start [POST]
func StartPersistentInterview(ctx context.Context, c *app.RequestContext) {
	// 1. 解析请求
	var req StartPersistentInterviewRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, TestAgentResponse{
			Success: false,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 2. 获取用户ID
	jwtUserID := middleware.GetUserID(c)
	userID := jwtUserID
	if userID == 0 {
		userID = req.UserID
	}
	if userID == 0 {
		c.JSON(consts.StatusUnauthorized, TestAgentResponse{
			Success: false,
			Message: "未登录或未提供user_id",
		})
		return
	}

	log.Printf("[PersistentInterview] 开始持久面试: userID=%d", userID)

	// 3. 创建新会话（支持多会话，每个会话独立上下文）
	mgr := GetPersistentSessionManager()

	// 注意：不再自动结束旧会话，允许用户同时拥有多个独立会话
	session := mgr.CreateSession(userID)

	// 5. 设置SSE响应头
	c.SetStatusCode(consts.StatusOK)
	c.Response.Header.Set("Content-Type", "text/event-stream")
	c.Response.Header.Set("Cache-Control", "no-cache")
	c.Response.Header.Set("Connection", "keep-alive")
	c.Response.Header.Set("Access-Control-Allow-Origin", "*")
	c.Response.Header.Set("X-Session-ID", session.SessionID)

	// 6. 创建管道
	pipeReader, pipeWriter := io.Pipe()
	c.SetBodyStream(pipeReader, -1)

	go func() {
		defer pipeWriter.Close()
		defer mgr.RemoveSession(session.SessionID)

		// 发送会话ID
		writeSSEEvent(pipeWriter, "session_created", session.SessionID)

		// 如果使用简历，自动发送第一条消息
		if req.UseResume && req.ResumeID > 0 {
			resume, err := model.ResumeDao.GetResumeByID(req.ResumeID)
			if err != nil {
				writeSSEEvent(pipeWriter, "error", "获取简历失败")
				return
			}
			if resume.Content != "" {
				session.SendMessage(UserMessageReq{
					Message:   "请分析以下简历并开始面试：\n" + resume.Content,
					UseResume: true,
					ResumeID:  req.ResumeID,
				})
			}
		}

		// 持续监听并发送响应
		session.StreamResponses(pipeWriter)

		log.Printf("[PersistentInterview] SSE连接关闭: sessionID=%s", session.SessionID)
	}()
}

// SendInterviewMessage 发送面试消息
// @router /api/interview/persistent/message [POST]
func SendInterviewMessage(ctx context.Context, c *app.RequestContext) {
	var req SendMessageRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	mgr := GetPersistentSessionManager()
	session := mgr.GetSession(req.SessionID)
	if session == nil {
		c.JSON(consts.StatusNotFound, map[string]interface{}{
			"success": false,
			"message": "会话不存在或已结束",
		})
		return
	}

	if session.IsEnded {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "会话已结束",
		})
		return
	}

	// 发送消息
	if !session.SendMessage(UserMessageReq{Message: req.Message}) {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "发送消息失败",
		})
		return
	}

	log.Printf("[PersistentInterview] 收到用户消息: sessionID=%s, message=%s",
		req.SessionID, truncateString(req.Message, 50))

	c.JSON(consts.StatusOK, map[string]interface{}{
		"success": true,
		"message": "消息已发送",
	})
}

// EndInterview 结束面试
// @router /api/interview/persistent/end [POST]
func EndInterview(ctx context.Context, c *app.RequestContext) {
	var req EndInterviewRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	mgr := GetPersistentSessionManager()
	session := mgr.GetSession(req.SessionID)
	if session == nil {
		c.JSON(consts.StatusNotFound, map[string]interface{}{
			"success": false,
			"message": "会话不存在",
		})
		return
	}

	mgr.EndSession(req.SessionID)

	log.Printf("[PersistentInterview] 面试已结束: sessionID=%s", req.SessionID)

	c.JSON(consts.StatusOK, map[string]interface{}{
		"success": true,
		"message": "面试已结束",
	})
}

// GetInterviewSession 获取当前用户的所有活跃会话
// @router /api/interview/persistent/session [GET]
func GetInterviewSession(ctx context.Context, c *app.RequestContext) {
	jwtUserID := middleware.GetUserID(c)
	if jwtUserID == 0 {
		c.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "未登录",
		})
		return
	}

	mgr := GetPersistentSessionManager()
	sessions := mgr.GetAllSessionsByUserID(jwtUserID)

	if len(sessions) == 0 {
		c.JSON(consts.StatusOK, map[string]interface{}{
			"success":    true,
			"has_active": false,
			"sessions":   []interface{}{},
			"message":    "无活跃会话",
		})
		return
	}

	// 构建会话列表
	sessionList := make([]map[string]interface{}, 0, len(sessions))
	for _, session := range sessions {
		sessionList = append(sessionList, map[string]interface{}{
			"session_id":       session.SessionID,
			"last_active_time": session.LastActiveTime.Format(time.RFC3339),
			"message_count":    len(session.MessageHistory),
		})
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"success":    true,
		"has_active": true,
		"count":      len(sessions),
		"sessions":   sessionList,
	})
}

// 测试用简历（已弃用，改为从数据库获取）
// const testResume = `
// # 个人简历
//
// ## 基本信息
// - 姓名：张三
// - 工作年限：3年
// - 求职意向：Go后端开发工程师
//
// ## 教育背景
// - 2018.09 - 2022.06  XX大学  计算机科学与技术  本科
//
// ## 工作经历
//
// ### 某科技有限公司 | 后端开发工程师 | 2022.07 - 至今
//
// **项目一：电商订单系统重构**
// - 项目背景：负责公司核心订单系统的技术重构，将单体应用拆分为微服务架构
// - 技术栈：Go、Gin、gRPC、MySQL、Redis、Kafka、Docker、K8s
// - 主要职责：
//   - 设计并实现订单服务、库存服务、支付服务的微服务拆分方案
//   - 使用Redis实现分布式锁解决库存超卖问题
//   - 基于Kafka实现订单状态变更的异步通知机制
// - 项目成果：系统QPS从500提升至3000，订单处理延迟从200ms降低至50ms
//
// **项目二：实时数据分析平台**
// - 技术栈：Go、ClickHouse、Elasticsearch、Grafana
// - 主要职责：设计高性能数据采集服务，支持每秒10万条数据写入
//
// ## 技术技能
// - 编程语言：Go（精通）、Python（熟练）
// - 数据库：MySQL、Redis、ClickHouse
// - 中间件：Kafka、Elasticsearch
// - 云原生：Docker、Kubernetes
// `

// ========== 旧版接口（保留方便回滚）==========
// StartInterviewMeilong 旧版接口 - 每次请求独立，不保持长连接
// 如需使用持久连接版本，请使用 /api/interview/persistent/* 系列接口
// @router /api/interview/meilong/stream [POST]
func StartInterviewMeilong(ctx context.Context, c *app.RequestContext) {
	// 直接调用新版持久会话接口
	// 如需回滚到原有逻辑，请取消下面注释的代码并注释掉这行调用
	StartPersistentInterview(ctx, c)

	// ========== 原有逻辑（如需回滚请取消下面注释）==========
	// // 1. 解析请求
	// var req TestAgentRequest
	// if err := c.BindAndValidate(&req); err != nil {
	// 	c.JSON(consts.StatusBadRequest, TestAgentResponse{
	// 		Success: false,
	// 		Message: "请求参数错误: " + err.Error(),
	// 	})
	// 	return
	// }
	//
	// // 2. 获取用户ID（优先从JWT获取，否则使用请求体中的user_id）
	// jwtUserID := middleware.GetUserID(c)
	// log.Printf("[TestAgentStream] JWT中的UserID: %d, 请求体中的UserID: %d", jwtUserID, req.UserID)
	//
	// userID := jwtUserID
	// if userID == 0 {
	// 	userID = req.UserID
	// 	log.Printf("[TestAgentStream] JWT未获取到用户ID，使用请求体中的UserID: %d", userID)
	// }
	// if userID == 0 {
	// 	c.JSON(consts.StatusUnauthorized, TestAgentResponse{
	// 		Success: false,
	// 		Message: "未登录或未提供user_id，请在Header中添加 Authorization: Bearer {token}",
	// 	})
	// 	return
	// }
	//
	// log.Printf("[TestAgentStream] 最终使用的UserID: %d, Message: %s", userID, truncateString(req.Message, 50))
	//
	// // 3. 设置SSE响应头
	// c.SetStatusCode(consts.StatusOK)
	// c.Response.Header.Set("Content-Type", "text/event-stream")
	// c.Response.Header.Set("Cache-Control", "no-cache")
	// c.Response.Header.Set("Connection", "keep-alive")
	// c.Response.Header.Set("Access-Control-Allow-Origin", "*")
	//
	// // 4. 创建管道用于流式传输
	// pipeReader, pipeWriter := io.Pipe()
	// c.SetBodyStream(pipeReader, -1)
	//
	// go func() {
	// 	defer pipeWriter.Close()
	//
	// 	// 创建面试智能体
	// 	agent := bearAgent.QuestionGeneratorAgent(userID)
	//
	// 	// 创建 Runner
	// 	runner := adk.NewRunner(ctx, adk.RunnerConfig{
	// 		Agent: agent,
	// 	})
	//
	// 	// 构建消息历史
	// 	// 注意：LLM API 要求每条消息的 content 字段不能为空，否则会返回 400 错误
	// 	var messageHistory []adk.Message
	// 	for _, msg := range req.History {
	// 		// 跳过空内容的消息，防止 "missing messages.content" 错误
	// 		if strings.TrimSpace(msg.Content) == "" {
	// 			log.Printf("[TestAgentStream] 跳过空内容的历史消息, Role: %s", msg.Role)
	// 			continue
	// 		}
	// 		if msg.Role == "user" {
	// 			messageHistory = append(messageHistory, schema.UserMessage(msg.Content))
	// 		} else if msg.Role == "assistant" {
	// 			messageHistory = append(messageHistory, schema.AssistantMessage(msg.Content, nil))
	// 		}
	// 	}
	//
	// 	// 处理当前消息
	// 	userMessage := req.Message
	// 	if req.UseResume {
	// 		// 原先使用硬编码简历的代码（已注释）
	// 		// userMessage = "请分析以下简历并开始面试：\n" + testResume
	//
	// 		// 通过前端传的简历ID从数据库获取简历
	// 		// TODO: 目前写死 resumeID=4 用于测试，后续改为使用 req.ResumeID
	// 		resumeID := uint64(4) // 写死用于测试
	// 		// resumeID := req.ResumeID // 正式使用前端传参时启用这行
	// 		log.Printf("[TestAgentStream] 请求的ResumeID: %d (写死测试值: 4)", req.ResumeID)
	//
	// 		resume, err := model.ResumeDao.GetResumeByID(resumeID)
	// 		if err != nil {
	// 			log.Printf("[TestAgentStream] 获取简历失败: resumeID=%d, err=%v", resumeID, err)
	// 			writeSSEEvent(pipeWriter, "error", "获取简历失败，简历不存在")
	// 			return
	// 		}
	// 		if resume.Content == "" {
	// 			log.Printf("[TestAgentStream] 简历内容为空: resumeID=%d", resumeID)
	// 			writeSSEEvent(pipeWriter, "error", "简历内容为空，请重新上传简历")
	// 			return
	// 		}
	// 		log.Printf("[TestAgentStream] 获取简历成功: resumeID=%d, userID=%d", resume.ID, resume.UserID)
	// 		userMessage = "请分析以下简历并开始面试：\n" + resume.Content
	// 	}
	//
	// 	// 确保用户消息不为空，避免 "missing messages.content" 错误
	// 	if strings.TrimSpace(userMessage) == "" {
	// 		writeSSEEvent(pipeWriter, "error", "消息内容不能为空")
	// 		return
	// 	}
	// 	messageHistory = append(messageHistory, schema.UserMessage(userMessage))
	// 	log.Printf("messageHistory:%s", messageHistory)
	// 	// 运行 Agent
	// 	iter := runner.Run(ctx, messageHistory)
	//
	// 	for {
	// 		event, ok := iter.Next()
	// 		if !ok {
	// 			break
	// 		}
	//
	// 		if event.Err != nil {
	// 			writeSSEEvent(pipeWriter, "error", event.Err.Error())
	// 			return
	// 		}
	//
	// 		if event.Output != nil && event.Output.MessageOutput != nil {
	// 			content := event.Output.MessageOutput.Message.Content
	// 			if content != "" {
	// 				writeSSEEvent(pipeWriter, "message", content)
	// 			}
	// 		}
	// 	}
	//
	// 	// 发送完成事件
	// 	writeSSEEvent(pipeWriter, "done", "")
	// }()
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
