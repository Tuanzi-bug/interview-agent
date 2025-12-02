package interview

import (
	"ai-eino-interview-agent/chatApp/agent/bearAgent"
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

// PersistentSession 持久化会话
type PersistentSession struct {
	SessionID      string
	UserID         uint
	MessageHistory []adk.Message       // 对话历史
	MessageChan    chan UserMessageReq // 用户消息通道
	ResponseChan   chan SSEResponse    // 响应通道
	LastActiveTime time.Time           // 最后活跃时间
	IsEnded        bool                // 是否已结束
	Ctx            context.Context
	Cancel         context.CancelFunc
	mu             sync.Mutex
}

// UserMessageReq 用户消息请求
type UserMessageReq struct {
	Message   string
	UseResume bool
	ResumeID  uint64
}

// SSEResponse SSE响应
type SSEResponse struct {
	EventType string // message, error, done, end_interview
	Data      string
}

// PersistentSessionManager 持久会话管理器
type PersistentSessionManager struct {
	sessions map[string]*PersistentSession
	mu       sync.RWMutex
}

// 全局会话管理器实例
var persistentSessionMgr = &PersistentSessionManager{
	sessions: make(map[string]*PersistentSession),
}

// GetPersistentSessionManager 获取持久会话管理器
func GetPersistentSessionManager() *PersistentSessionManager {
	return persistentSessionMgr
}

// CreateSession 创建新会话
func (m *PersistentSessionManager) CreateSession(userID uint) *PersistentSession {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 使用 userID + 时间戳 作为会话ID
	sessionID := generatePersistentSessionID(userID)

	ctx, cancel := context.WithCancel(context.Background())

	session := &PersistentSession{
		SessionID:      sessionID,
		UserID:         userID,
		MessageHistory: make([]adk.Message, 0),
		MessageChan:    make(chan UserMessageReq, 10),
		ResponseChan:   make(chan SSEResponse, 100),
		LastActiveTime: time.Now(),
		IsEnded:        false,
		Ctx:            ctx,
		Cancel:         cancel,
	}

	m.sessions[sessionID] = session

	// 启动会话超时检查
	go m.watchSessionTimeout(sessionID)

	// 启动消息处理协程
	go session.processMessages()

	log.Printf("[PersistentSession] 创建新会话: sessionID=%s, userID=%d", sessionID, userID)
	return session
}

// GetSession 获取会话
func (m *PersistentSessionManager) GetSession(sessionID string) *PersistentSession {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sessions[sessionID]
}

// GetSessionByUserID 通过用户ID获取活跃会话（返回第一个找到的）
func (m *PersistentSessionManager) GetSessionByUserID(userID uint) *PersistentSession {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, session := range m.sessions {
		if session.UserID == userID && !session.IsEnded {
			return session
		}
	}
	return nil
}

// GetAllSessionsByUserID 获取用户的所有活跃会话
func (m *PersistentSessionManager) GetAllSessionsByUserID(userID uint) []*PersistentSession {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var sessions []*PersistentSession
	for _, session := range m.sessions {
		if session.UserID == userID && !session.IsEnded {
			sessions = append(sessions, session)
		}
	}
	return sessions
}

// EndSession 结束会话
func (m *PersistentSessionManager) EndSession(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if session, exists := m.sessions[sessionID]; exists {
		session.mu.Lock()
		session.IsEnded = true
		session.Cancel()
		session.mu.Unlock()

		// 发送结束事件
		select {
		case session.ResponseChan <- SSEResponse{EventType: "end_interview", Data: "面试已结束"}:
		default:
		}

		close(session.MessageChan)
		log.Printf("[PersistentSession] 会话已结束: sessionID=%s", sessionID)
	}
}

// RemoveSession 移除会话
func (m *PersistentSessionManager) RemoveSession(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, sessionID)
	log.Printf("[PersistentSession] 会话已移除: sessionID=%s", sessionID)
}

// watchSessionTimeout 监控会话超时
func (m *PersistentSessionManager) watchSessionTimeout(sessionID string) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			session := m.GetSession(sessionID)
			if session == nil {
				return
			}

			session.mu.Lock()
			timeSinceLastActive := time.Since(session.LastActiveTime)
			isEnded := session.IsEnded
			session.mu.Unlock()

			if isEnded {
				return
			}

			// 超过30分钟无活动
			if timeSinceLastActive > 30*time.Minute {
				log.Printf("[PersistentSession] 会话超时: sessionID=%s, 最后活跃时间: %v",
					sessionID, session.LastActiveTime)

				// 发送超时事件
				select {
				case session.ResponseChan <- SSEResponse{EventType: "timeout", Data: "会话已超时（30分钟无活动），面试已结束"}:
				default:
				}

				m.EndSession(sessionID)
				return
			}
		}
	}
}

// processMessages 处理用户消息
func (s *PersistentSession) processMessages() {
	for {
		select {
		case <-s.Ctx.Done():
			log.Printf("[PersistentSession] 消息处理协程退出: sessionID=%s", s.SessionID)
			return

		case msg, ok := <-s.MessageChan:
			if !ok {
				log.Printf("[PersistentSession] 消息通道已关闭: sessionID=%s", s.SessionID)
				return
			}

			s.mu.Lock()
			s.LastActiveTime = time.Now()
			s.mu.Unlock()

			log.Printf("[PersistentSession] 处理用户消息: sessionID=%s, message=%s",
				s.SessionID, truncateString(msg.Message, 50))

			// 处理消息并生成回复
			s.handleUserMessage(msg)
		}
	}
}

// handleUserMessage 处理用户消息
func (s *PersistentSession) handleUserMessage(msg UserMessageReq) {
	userMessage := msg.Message

	// 处理简历相关逻辑（如需要可扩展）
	if msg.UseResume && msg.ResumeID > 0 {
		// TODO: 从数据库获取简历内容
		// 这里简化处理，实际可以调用 model.ResumeDao.GetResumeByID
	}

	// 检查消息是否为空
	if strings.TrimSpace(userMessage) == "" {
		s.ResponseChan <- SSEResponse{EventType: "error", Data: "消息内容不能为空"}
		return
	}

	// 添加用户消息到历史
	s.mu.Lock()
	s.MessageHistory = append(s.MessageHistory, schema.UserMessage(userMessage))
	currentHistory := make([]adk.Message, len(s.MessageHistory))
	copy(currentHistory, s.MessageHistory)
	s.mu.Unlock()

	// 创建面试智能体
	agent := bearAgent.QuestionGeneratorAgent(s.UserID)

	// 创建 Runner
	runner := adk.NewRunner(s.Ctx, adk.RunnerConfig{
		Agent: agent,
	})

	// 运行 Agent
	iter := runner.Run(s.Ctx, currentHistory)

	var assistantResponse strings.Builder

	for {
		event, ok := iter.Next()
		if !ok {
			break
		}

		// 检查会话是否已结束
		select {
		case <-s.Ctx.Done():
			return
		default:
		}

		if event.Err != nil {
			s.ResponseChan <- SSEResponse{EventType: "error", Data: event.Err.Error()}
			return
		}

		if event.Output != nil && event.Output.MessageOutput != nil {
			content := event.Output.MessageOutput.Message.Content
			if content != "" {
				assistantResponse.WriteString(content)
				s.ResponseChan <- SSEResponse{EventType: "message", Data: content}
			}
		}
	}

	// 将助手回复添加到历史
	if assistantResponse.Len() > 0 {
		s.mu.Lock()
		s.MessageHistory = append(s.MessageHistory, schema.AssistantMessage(assistantResponse.String(), nil))
		s.mu.Unlock()
	}

	// 发送本轮完成事件
	s.ResponseChan <- SSEResponse{EventType: "turn_complete", Data: ""}
}

// SendMessage 发送用户消息
func (s *PersistentSession) SendMessage(msg UserMessageReq) bool {
	s.mu.Lock()
	if s.IsEnded {
		s.mu.Unlock()
		return false
	}
	s.LastActiveTime = time.Now()
	s.mu.Unlock()

	select {
	case s.MessageChan <- msg:
		return true
	default:
		return false
	}
}

// StreamResponses 流式输出响应到writer
func (s *PersistentSession) StreamResponses(w io.Writer) {
	for {
		select {
		case <-s.Ctx.Done():
			writeSSEEvent(w, "end_interview", "会话已结束")
			return

		case resp, ok := <-s.ResponseChan:
			if !ok {
				writeSSEEvent(w, "end_interview", "会话已关闭")
				return
			}
			writeSSEEvent(w, resp.EventType, resp.Data)

			// 如果是结束事件，退出循环
			if resp.EventType == "end_interview" || resp.EventType == "timeout" {
				return
			}
		}
	}
}

// generatePersistentSessionID 生成持久会话ID
func generatePersistentSessionID(userID uint) string {
	// 使用纳秒级时间戳确保唯一性
	return fmt.Sprintf("ps_%d_%d", userID, time.Now().UnixNano())
}
