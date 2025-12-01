package agent

import (
	"fmt"
	"time"

	"github.com/cloudwego/eino/adk"
)

// SpecializedAgent 专项智能体统一接口
type SpecializedAgent interface {
	adk.Agent

	// 领域识别
	CanHandle(domain string) bool

	// 专项能力评估
	GetExpertiseLevel() ExpertiseLevel

	// 知识库管理
	GetKnowledgeBase() *KnowledgeBase

	// 性能指标
	GetPerformanceMetrics() *AgentMetrics

	// 健康检查
	HealthCheck() error
}

// ExpertiseLevel 专业能力等级
type ExpertiseLevel int

const (
	LevelBeginner ExpertiseLevel = iota
	LevelIntermediate
	LevelAdvanced
	LevelExpert
)

// KnowledgeBase 知识库结构
type KnowledgeBase struct {
	Domain      string               `json:"domain"`
	Categories  []string             `json:"categories"`
	Questions   []QuestionTemplate   `json:"questions"`
	Evaluations []EvaluationCriteria `json:"evaluations"`
	Version     string               `json:"version"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

// QuestionTemplate 问题模板
type QuestionTemplate struct {
	ID             string   `json:"id"`
	Category       string   `json:"category"`
	Difficulty     string   `json:"difficulty"`
	Question       string   `json:"question"`
	Keywords       []string `json:"keywords"`
	ExpectedPoints []string `json:"expected_points"`
	TimeLimit      int      `json:"time_limit"`
}

// EvaluationCriteria 评估标准
type EvaluationCriteria struct {
	Category string             `json:"category"`
	Weights  map[string]float32 `json:"weights"`
	Criteria []EvaluationItem   `json:"criteria"`
}

// EvaluationItem 评估项
type EvaluationItem struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	MaxScore    float32  `json:"max_score"`
	KeyPoints   []string `json:"key_points"`
}

// AgentMetrics Agent性能指标
type AgentMetrics struct {
	TotalRequests   int64         `json:"total_requests"`
	SuccessRequests int64         `json:"success_requests"`
	AvgResponseTime time.Duration `json:"avg_response_time"`
	ErrorRate       float32       `json:"error_rate"`
	LastActiveTime  time.Time     `json:"last_active_time"`
	CacheHitRate    float32       `json:"cache_hit_rate"`
}

// BaseSpecializedAgent 专项智能体基类
type BaseSpecializedAgent struct {
	agentType     string
	domain        string
	knowledgeBase *KnowledgeBase
	metrics       *AgentMetrics
	config        *SpecializedAgentConfig
}

// SpecializedAgentConfig 专项智能体配置
type SpecializedAgentConfig struct {
	AgentName      string
	Domain         string
	Description    string
	ExpertiseLevel ExpertiseLevel
	MaxIterations  int
	Timeout        time.Duration
	EnableCache    bool
	CacheTTL       time.Duration
	RetryAttempts  int
	FallbackAgent  string
}

// NewBaseSpecializedAgent 创建基类实例
func NewBaseSpecializedAgent(config *SpecializedAgentConfig) *BaseSpecializedAgent {
	return &BaseSpecializedAgent{
		agentType: config.AgentName,
		domain:    config.Domain,
		knowledgeBase: &KnowledgeBase{
			Domain:    config.Domain,
			Version:   "1.0.0",
			UpdatedAt: time.Now(),
		},
		metrics: &AgentMetrics{
			LastActiveTime: time.Now(),
		},
		config: config,
	}
}

// CanHandle 检查是否能处理指定领域
func (b *BaseSpecializedAgent) CanHandle(domain string) bool {
	return b.domain == domain
}

// GetExpertiseLevel 获取专业能力等级
func (b *BaseSpecializedAgent) GetExpertiseLevel() ExpertiseLevel {
	return b.config.ExpertiseLevel
}

// GetKnowledgeBase 获取知识库
func (b *BaseSpecializedAgent) GetKnowledgeBase() *KnowledgeBase {
	return b.knowledgeBase
}

// GetPerformanceMetrics 获取性能指标
func (b *BaseSpecializedAgent) GetPerformanceMetrics() *AgentMetrics {
	return b.metrics
}

// HealthCheck 健康检查
func (b *BaseSpecializedAgent) HealthCheck() error {
	if time.Since(b.metrics.LastActiveTime) > b.config.Timeout*3 {
		return fmt.Errorf("agent %s health check failed: timeout", b.agentType)
	}
	return nil
}

// UpdateMetrics 更新性能指标
func (b *BaseSpecializedAgent) UpdateMetrics(success bool, responseTime time.Duration) {
	b.metrics.TotalRequests++
	if success {
		b.metrics.SuccessRequests++
	}

	// 更新平均响应时间
	if b.metrics.AvgResponseTime == 0 {
		b.metrics.AvgResponseTime = responseTime
	} else {
		b.metrics.AvgResponseTime = (b.metrics.AvgResponseTime + responseTime) / 2
	}

	// 更新错误率
	b.metrics.ErrorRate = float32(b.metrics.TotalRequests-b.metrics.SuccessRequests) / float32(b.metrics.TotalRequests)
	b.metrics.LastActiveTime = time.Now()
}

// SpecializedAgentRegistry 专项智能体注册中心
type SpecializedAgentRegistry struct {
	agents map[string]SpecializedAgent
}

// NewSpecializedAgentRegistry 创建注册中心
func NewSpecializedAgentRegistry() *SpecializedAgentRegistry {
	return &SpecializedAgentRegistry{
		agents: make(map[string]SpecializedAgent),
	}
}

// Register 注册专项智能体
func (r *SpecializedAgentRegistry) Register(agent SpecializedAgent) error {
	agentType := agent.GetKnowledgeBase().Domain
	if _, exists := r.agents[agentType]; exists {
		return fmt.Errorf("agent type %s already registered", agentType)
	}
	r.agents[agentType] = agent
	return nil
}

// Get 获取专项智能体
func (r *SpecializedAgentRegistry) Get(agentType string) (SpecializedAgent, bool) {
	agent, exists := r.agents[agentType]
	return agent, exists
}

// FindByDomain 根据领域查找智能体
func (r *SpecializedAgentRegistry) FindByDomain(domain string) []SpecializedAgent {
	var result []SpecializedAgent
	for _, agent := range r.agents {
		if agent.CanHandle(domain) {
			result = append(result, agent)
		}
	}
	return result
}

// GetAll 获取所有智能体
func (r *SpecializedAgentRegistry) GetAll() []SpecializedAgent {
	var result []SpecializedAgent
	for _, agent := range r.agents {
		result = append(result, agent)
	}
	return result
}

// InterviewRequest 面试请求
type InterviewRequest struct {
	UserID          string                 `json:"user_id"`
	SessionID       string                 `json:"session_id"`
	InterviewType   string                 `json:"interview_type"`
	Domain          string                 `json:"domain"`
	Difficulty      string                 `json:"difficulty"`
	ResumeContent   string                 `json:"resume_content"`
	PreviousContext map[string]interface{} `json:"previous_context"`
}

// InterviewResponse 面试响应
type InterviewResponse struct {
	RequestID  string                 `json:"request_id"`
	Question   string                 `json:"question"`
	Category   string                 `json:"category"`
	Difficulty string                 `json:"difficulty"`
	Hints      []string               `json:"hints"`
	TimeLimit  int                    `json:"time_limit"`
	NextAction string                 `json:"next_action"`
	Context    map[string]interface{} `json:"context"`
}

// EvaluationRequest 评估请求
type EvaluationRequest struct {
	RequestID  string                 `json:"request_id"`
	Question   string                 `json:"question"`
	Answer     string                 `json:"answer"`
	Category   string                 `json:"category"`
	Difficulty string                 `json:"difficulty"`
	TimeSpent  int                    `json:"time_spent"`
	Context    map[string]interface{} `json:"context"`
}

// EvaluationResponse 评估响应
type EvaluationResponse struct {
	RequestID      string                 `json:"request_id"`
	OverallScore   float32                `json:"overall_score"`
	MaxScore       float32                `json:"max_score"`
	Grade          string                 `json:"grade"`
	Strengths      []string               `json:"strengths"`
	Weaknesses     []string               `json:"weaknesses"`
	Improvements   []string               `json:"improvements"`
	DetailedScores map[string]float32     `json:"detailed_scores"`
	NextTopic      string                 `json:"next_topic"`
	Context        map[string]interface{} `json:"context"`
}
