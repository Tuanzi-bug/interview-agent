package session

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/eino/adk"
)

// SessionContextManager 管理 Eino 会话值的工具类
// 用于在 Agent 执行期间存储和检索面试相关的上下文信息
type SessionContextManager struct {
	ctx context.Context
}

// NewSessionContextManager 创建新的会话上下文管理器
func NewSessionContextManager(ctx context.Context) *SessionContextManager {
	return &SessionContextManager{ctx: ctx}
}

// ===== 会话值键定义 =====
const (
	// 简历相关
	KeyResumeContent  = "resume-content"   // 用于问题生成器初始化
	KeyResumeFilePath = "resume-file-path" // 用于简历上传参数

	// 面试配置
	KeyInterviewType       = "interview-type"
	KeyInterviewDomain     = "interview-domain"
	KeyInterviewDifficulty = "interview-difficulty"

	// 对话历史
	KeyConversationHistory  = "conversation-history"
	KeyCurrentQuestion      = "current-question"
	KeyCurrentQuestionIndex = "current-question-index"
	KeyCurrentDimension     = "current-dimension"

	// 用户信息
	KeyUserID       = "user-id"
	KeyRecordID     = "record-id"
	KeyResumeUserID = "resume-user-id"

	// 评估相关
	KeyEvaluationContext = "evaluation-context"
	KeyParsedResume      = "parsed-resume"

	// 简历上传相关
	KeyResumeFileName      = "resume-file-name"
	KeyResumeFileType      = "resume-file-type"
	KeyResumeFileSize      = "resume-file-size"
	KeyResumeUploadContent = "resume-upload-content"
)

// ===== 简历相关操作 =====

// SetResumeContent 存储简历内容（用于问题生成器初始化）
func (m *SessionContextManager) SetResumeContent(content string) error {
	adk.AddSessionValue(m.ctx, KeyResumeContent, content)
	return nil
}

// GetResumeContent 获取简历内容
func (m *SessionContextManager) GetResumeContent() (string, error) {
	val, ok := adk.GetSessionValue(m.ctx, KeyResumeContent)
	if !ok {
		return "", nil
	}
	if val == nil {
		return "", nil
	}
	content, isString := val.(string)
	if !isString {
		return "", fmt.Errorf("resume content is not a string")
	}
	return content, nil
}

// ===== 面试配置操作 =====

// SetInterviewConfig 存储面试配置
func (m *SessionContextManager) SetInterviewConfig(interviewType, domain, difficulty string) error {
	adk.AddSessionValue(m.ctx, KeyInterviewType, interviewType)
	adk.AddSessionValue(m.ctx, KeyInterviewDomain, domain)
	adk.AddSessionValue(m.ctx, KeyInterviewDifficulty, difficulty)
	return nil
}

// GetInterviewConfig 获取面试配置
func (m *SessionContextManager) GetInterviewConfig() (interviewType, domain, difficulty string, err error) {
	val, ok := adk.GetSessionValue(m.ctx, KeyInterviewType)
	if ok && val != nil {
		interviewType, _ = val.(string)
	}
	val, ok = adk.GetSessionValue(m.ctx, KeyInterviewDomain)
	if ok && val != nil {
		domain, _ = val.(string)
	}
	val, ok = adk.GetSessionValue(m.ctx, KeyInterviewDifficulty)
	if ok && val != nil {
		difficulty, _ = val.(string)
	}
	return
}

// ===== 对话历史操作 =====

// Message 代表一条对话消息
type Message struct {
	Role    string `json:"role"` // "user" 或 "assistant"
	Content string `json:"content"`
}

// AddConversationMessage 添加对话消息到历史
func (m *SessionContextManager) AddConversationMessage(role, content string) error {
	// 获取现有历史
	history, _ := m.GetConversationHistory()

	// 添加新消息
	history = append(history, Message{Role: role, Content: content})

	// 存储更新后的历史
	adk.AddSessionValue(m.ctx, KeyConversationHistory, history)
	return nil
}

// GetConversationHistory 获取对话历史
func (m *SessionContextManager) GetConversationHistory() ([]Message, error) {
	val, ok := adk.GetSessionValue(m.ctx, KeyConversationHistory)
	if !ok {
		return []Message{}, nil
	}
	if val == nil {
		return []Message{}, nil
	}

	// 处理两种可能的类型：[]Message 或 []interface{}
	switch v := val.(type) {
	case []Message:
		return v, nil
	case []interface{}:
		var history []Message
		// 转换 []interface{} 为 []Message
		for _, item := range v {
			if msgMap, ok := item.(map[string]interface{}); ok {
				msg := Message{
					Role:    fmt.Sprintf("%v", msgMap["role"]),
					Content: fmt.Sprintf("%v", msgMap["content"]),
				}
				history = append(history, msg)
			}
		}
		return history, nil
	default:
		return []Message{}, fmt.Errorf("conversation history has unexpected type: %T", val)
	}
}

// ClearConversationHistory 清空对话历史
func (m *SessionContextManager) ClearConversationHistory() error {
	adk.AddSessionValue(m.ctx, KeyConversationHistory, []Message{})
	return nil
}

// ===== 当前问题操作 =====

// SetCurrentQuestion 存储当前问题
func (m *SessionContextManager) SetCurrentQuestion(question string) error {
	adk.AddSessionValue(m.ctx, KeyCurrentQuestion, question)
	return nil
}

// GetCurrentQuestion 获取当前问题
func (m *SessionContextManager) GetCurrentQuestion() (string, error) {
	val, ok := adk.GetSessionValue(m.ctx, KeyCurrentQuestion)
	if !ok {
		return "", nil
	}
	if val == nil {
		return "", nil
	}
	question, isString := val.(string)
	if !isString {
		return "", fmt.Errorf("current question is not a string")
	}
	return question, nil
}

// SetCurrentQuestionIndex 存储当前问题索引
func (m *SessionContextManager) SetCurrentQuestionIndex(index int) error {
	adk.AddSessionValue(m.ctx, KeyCurrentQuestionIndex, index)
	return nil
}

// GetCurrentQuestionIndex 获取当前问题索引
func (m *SessionContextManager) GetCurrentQuestionIndex() (int, error) {
	val, ok := adk.GetSessionValue(m.ctx, KeyCurrentQuestionIndex)
	if !ok {
		return 0, nil
	}
	if val == nil {
		return 0, nil
	}
	index, isInt := val.(int)
	if !isInt {
		return 0, fmt.Errorf("current question index is not an int")
	}
	return index, nil
}

// SetCurrentDimension 存储当前评估维度
func (m *SessionContextManager) SetCurrentDimension(dimension string) error {
	adk.AddSessionValue(m.ctx, KeyCurrentDimension, dimension)
	return nil
}

// GetCurrentDimension 获取当前评估维度
func (m *SessionContextManager) GetCurrentDimension() (string, error) {
	val, ok := adk.GetSessionValue(m.ctx, KeyCurrentDimension)
	if !ok {
		return "", nil
	}
	if val == nil {
		return "", nil
	}
	dimension, isString := val.(string)
	if !isString {
		return "", fmt.Errorf("current dimension is not a string")
	}
	return dimension, nil
}

// ===== 用户信息操作 =====

// SetUserInfo 存储用户信息
func (m *SessionContextManager) SetUserInfo(userID uint, recordID uint64) error {
	adk.AddSessionValue(m.ctx, KeyUserID, userID)
	adk.AddSessionValue(m.ctx, KeyRecordID, recordID)
	return nil
}

// GetUserInfo 获取用户信息
func (m *SessionContextManager) GetUserInfo() (userID uint, recordID uint64, err error) {
	val, ok := adk.GetSessionValue(m.ctx, KeyUserID)
	if ok && val != nil {
		if uid, isUint := val.(uint); isUint {
			userID = uid
		}
	}
	val, ok = adk.GetSessionValue(m.ctx, KeyRecordID)
	if ok && val != nil {
		if rid, isUint64 := val.(uint64); isUint64 {
			recordID = rid
		}
	}
	return
}

// ===== 评估相关操作 =====

// SetEvaluationContext 存储评估上下文
func (m *SessionContextManager) SetEvaluationContext(contextData map[string]interface{}) error {
	data, _ := json.Marshal(contextData)
	adk.AddSessionValue(m.ctx, KeyEvaluationContext, string(data))
	return nil
}

// GetEvaluationContext 获取评估上下文
func (m *SessionContextManager) GetEvaluationContext() (map[string]interface{}, error) {
	val, ok := adk.GetSessionValue(m.ctx, KeyEvaluationContext)
	if !ok {
		return make(map[string]interface{}), nil
	}
	if val == nil {
		return make(map[string]interface{}), nil
	}

	contextStr, isString := val.(string)
	if !isString {
		return nil, fmt.Errorf("evaluation context is not a string")
	}

	var contextData map[string]interface{}
	if err := json.Unmarshal([]byte(contextStr), &contextData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal evaluation context: %w", err)
	}
	return contextData, nil
}

// SetParsedResume 存储解析后的简历
func (m *SessionContextManager) SetParsedResume(resume map[string]interface{}) error {
	data, _ := json.Marshal(resume)
	adk.AddSessionValue(m.ctx, KeyParsedResume, string(data))
	return nil
}

// GetParsedResume 获取解析后的简历
func (m *SessionContextManager) GetParsedResume() (map[string]interface{}, error) {
	val, ok := adk.GetSessionValue(m.ctx, KeyParsedResume)
	if !ok {
		return make(map[string]interface{}), nil
	}
	if val == nil {
		return make(map[string]interface{}), nil
	}

	resumeStr, isString := val.(string)
	if !isString {
		return nil, fmt.Errorf("parsed resume is not a string")
	}

	var resume map[string]interface{}
	if err := json.Unmarshal([]byte(resumeStr), &resume); err != nil {
		return nil, fmt.Errorf("failed to unmarshal parsed resume: %w", err)
	}
	return resume, nil
}

// ===== 简历上传参数操作 =====

// SetResumeUploadParams 存储简历上传参数（包括 userID 和 resumeFilePath）
func (m *SessionContextManager) SetResumeUploadParams(userID uint, resumeFilePath, fileName, fileType string, fileSize int64, content string) error {
	adk.AddSessionValue(m.ctx, KeyResumeUserID, userID)
	adk.AddSessionValue(m.ctx, KeyResumeFilePath, resumeFilePath)
	adk.AddSessionValue(m.ctx, KeyResumeFileName, fileName)
	adk.AddSessionValue(m.ctx, KeyResumeFileType, fileType)
	adk.AddSessionValue(m.ctx, KeyResumeFileSize, fileSize)
	adk.AddSessionValue(m.ctx, KeyResumeUploadContent, content)
	return nil
}

// GetResumeUploadParams 获取简历上传参数
func (m *SessionContextManager) GetResumeUploadParams() (userID uint, resumeFilePath, fileName, fileType string, fileSize int64, content string, err error) {
	// 获取 userID
	val, ok := adk.GetSessionValue(m.ctx, KeyResumeUserID)
	if ok && val != nil {
		userID, _ = val.(uint)
	}

	// 获取 resumeFilePath
	val, ok = adk.GetSessionValue(m.ctx, KeyResumeFilePath)
	if ok && val != nil {
		resumeFilePath, _ = val.(string)
	}

	// 获取文件名
	val, ok = adk.GetSessionValue(m.ctx, KeyResumeFileName)
	if ok && val != nil {
		fileName, _ = val.(string)
	}

	// 获取文件类型
	val, ok = adk.GetSessionValue(m.ctx, KeyResumeFileType)
	if ok && val != nil {
		fileType, _ = val.(string)
	}

	// 获取文件大小
	val, ok = adk.GetSessionValue(m.ctx, KeyResumeFileSize)
	if ok && val != nil {
		fileSize, _ = val.(int64)
	}

	// 获取内容
	val, ok = adk.GetSessionValue(m.ctx, KeyResumeUploadContent)
	if ok && val != nil {
		content, _ = val.(string)
	}

	return
}
