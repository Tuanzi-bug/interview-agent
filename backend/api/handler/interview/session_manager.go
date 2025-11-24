package interview

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"
)

// 错误定义
var (
	ErrSessionNotFound = errors.New("session not found")
	ErrInvalidRequest  = errors.New("invalid request")
)

// InterviewSession 面试会话
type InterviewSession struct {
	SessionID       string        // 会话ID
	UserID          uint          // 用户ID
	RecordID        uint64        // 面试记录ID
	CurrentQuestion int           // 当前问题索引
	AllQuestions    []interface{} // 所有问题
	AllDialogues    []interface{} // 所有对话
	UserAnswer      string        // 用户的答案
	AnswerReceived  bool          // 是否收到答案
	CreatedAt       time.Time     // 创建时间
	StartTime       time.Time     // 面试开始时间
	LastActivity    time.Time     // 最后活动时间
	ResumeFilePath  string        // 简历文件路径
	HasResume       bool          // 是否有简历
	Query           string        // 用户查询
	AnswerChan      chan string   // 答案通道
	Done            bool          // 面试是否完成
	Type            string        // 面试类型（综合面试/专项面试）
	Domain          string        // 面试领域
	Difficulty      string        // 难度级别
	mu              sync.Mutex    // 锁
}

// SessionManager 会话管理器
type SessionManager struct {
	sessions map[string]*InterviewSession
	mu       sync.RWMutex
	ticker   *time.Ticker
}

var sessionManager = &SessionManager{
	sessions: make(map[string]*InterviewSession),
}

// init 初始化会话管理器
func init() {
	// 启动清理过期会话的 goroutine（30分钟无活动则清理）
	sessionManager.ticker = time.NewTicker(5 * time.Minute)
	go func() {
		for range sessionManager.ticker.C {
			sessionManager.cleanupExpiredSessions()
		}
	}()
}

// CreateSession 创建新会话
func (sm *SessionManager) CreateSession(userID uint, recordID uint64, resumeFilePath string, hasResume bool, query string, interviewType string, domain string, difficulty string) *InterviewSession {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sessionID := generateSessionID()
	now := time.Now()
	session := &InterviewSession{
		SessionID:       sessionID,
		UserID:          userID,
		RecordID:        recordID,
		CurrentQuestion: 0,
		AllQuestions:    []interface{}{},
		AllDialogues:    []interface{}{},
		CreatedAt:       now,
		StartTime:       now,
		LastActivity:    now,
		ResumeFilePath:  resumeFilePath,
		HasResume:       hasResume,
		Query:           query,
		Type:            interviewType,
		Domain:          domain,
		Difficulty:      difficulty,
		AnswerChan:      make(chan string, 1),
		Done:            false,
	}

	sm.sessions[sessionID] = session
	return session
}

// GetSession 获取会话
func (sm *SessionManager) GetSession(sessionID string) *InterviewSession {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.sessions[sessionID]
}

// UpdateSession 更新会话
func (sm *SessionManager) UpdateSession(session *InterviewSession) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	session.LastActivity = time.Now()
	sm.sessions[session.SessionID] = session
}

// DeleteSession 删除会话
func (sm *SessionManager) DeleteSession(sessionID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if session, ok := sm.sessions[sessionID]; ok {
		close(session.AnswerChan)
		delete(sm.sessions, sessionID)
	}
}

// SubmitAnswer 提交答案
func (sm *SessionManager) SubmitAnswer(sessionID, answer string) error {
	sm.mu.RLock()
	session, ok := sm.sessions[sessionID]
	sm.mu.RUnlock()

	if !ok {
		return ErrSessionNotFound
	}

	session.mu.Lock()
	defer session.mu.Unlock()

	session.UserAnswer = answer
	session.AnswerReceived = true
	session.LastActivity = time.Now()

	// 非阻塞发送答案到通道
	select {
	case session.AnswerChan <- answer:
	default:
	}

	return nil
}

// GetAnswer 获取答案（阻塞等待，带超时）
func (sm *SessionManager) GetAnswer(sessionID string, timeout time.Duration) (string, bool) {
	sm.mu.RLock()
	session, ok := sm.sessions[sessionID]
	sm.mu.RUnlock()

	if !ok {
		return "", false
	}

	// 如果已经有答案，直接返回
	session.mu.Lock()
	if session.AnswerReceived {
		answer := session.UserAnswer
		session.AnswerReceived = false
		session.UserAnswer = ""
		session.mu.Unlock()
		return answer, true
	}
	session.mu.Unlock()

	// 等待答案通道
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case answer := <-session.AnswerChan:
		session.mu.Lock()
		session.AnswerReceived = false
		session.UserAnswer = ""
		session.mu.Unlock()
		return answer, true
	case <-timer.C:
		return "", false
	}
}

// ClearAnswer 清空答案标志
func (sm *SessionManager) ClearAnswer(sessionID string) {
	sm.mu.RLock()
	session, ok := sm.sessions[sessionID]
	sm.mu.RUnlock()

	if !ok {
		return
	}

	session.mu.Lock()
	defer session.mu.Unlock()

	session.AnswerReceived = false
	session.UserAnswer = ""
}

// cleanupExpiredSessions 清理过期会话（30分钟无活动）
func (sm *SessionManager) cleanupExpiredSessions() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	now := time.Now()
	for sessionID, session := range sm.sessions {
		if now.Sub(session.LastActivity) > 120*time.Minute {
			close(session.AnswerChan)
			delete(sm.sessions, sessionID)
		}
	}
}

// GetSessionManager 获取会话管理器实例
func GetSessionManager() *SessionManager {
	return sessionManager
}

// generateSessionID 生成唯一的会话ID
func generateSessionID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("session_%s_%d", hex.EncodeToString(b), time.Now().UnixNano())
}
