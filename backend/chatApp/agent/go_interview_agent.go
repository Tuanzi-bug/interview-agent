package agent

import (
	"context"
	"fmt"
	"time"

	"ai-eino-interview-agent/chatApp/agent/service"
	"ai-eino-interview-agent/chatApp/chat"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
)

// 结构体定义已在specialized_agent_interface.go中定义，此处复用

// GoInterviewAgent Go专项面试智能体
// 实现SpecializedAgent接口
// 负责根据Go语言面试需求生成问题并评估回答
// 支持不同难度和细分领域的面试需求
type GoInterviewAgent struct {
	BaseSpecializedAgent
	openaiChatModel model.ToolCallingChatModel
	supervisorName  string
	userId          uint
}

// Description 返回智能体描述信息
func (g *GoInterviewAgent) Description(ctx context.Context) string {
	return g.config.Description
}

// Name 返回智能体名称
func (g *GoInterviewAgent) Name(ctx context.Context) string {
	return g.config.AgentName
}

// NewGoInterviewAgent 创建Go专项面试智能体
func NewGoInterviewAgent(ctx context.Context, difficulty string, domain string) SpecializedAgent {
	// 这里userId临时设置为0，实际应该从上下文或参数中获取
	userId := uint(0)

	config := &SpecializedAgentConfig{
		AgentName:      "GoInterviewAgent",
		Domain:         domain,
		Description:    "Go语言专项面试智能体，专注于Go语言深度技术面试",
		ExpertiseLevel: LevelExpert,
		MaxIterations:  20,
		Timeout:        30 * time.Second,
		EnableCache:    true,
		CacheTTL:       1 * time.Hour,
		RetryAttempts:  3,
		FallbackAgent:  "GeneralInterviewAgent",
	}

	// 创建BaseSpecializedAgent实例
	baseAgent := NewBaseSpecializedAgent(config)

	// 构建Go知识库
	knowledgeBase := buildGoKnowledgeBase()
	baseAgent.knowledgeBase = knowledgeBase

	// 正确初始化GoInterviewAgent，包含BaseSpecializedAgent
	agent := &GoInterviewAgent{
		BaseSpecializedAgent: *baseAgent,
		openaiChatModel:      chat.CreatOpenAiChatModel(ctx, userId),
		supervisorName:       "GoInterviewSupervisor",
		userId:               userId,
	}

	return agent
}

// ProcessInterview 处理面试请求
func (g *GoInterviewAgent) ProcessInterview(ctx context.Context, req InterviewRequest) (*InterviewResponse, error) {
	startTime := time.Now()
	defer func() {
		responseTime := time.Since(startTime)
		g.UpdateMetrics(true, responseTime)
	}()

	// 根据难度和类别选择问题模板
	questionTemplate := g.selectQuestionTemplate(req.Difficulty, req.Domain)
	if questionTemplate == nil {
		return nil, fmt.Errorf("no suitable question template found for difficulty: %s", req.Difficulty)
	}

	// 个性化问题生成
	personalizedQuestion := g.personalizeQuestion(questionTemplate, req.ResumeContent, req.PreviousContext)

	response := &InterviewResponse{
		RequestID:  req.SessionID,
		Question:   personalizedQuestion,
		Category:   questionTemplate.Category,
		Difficulty: questionTemplate.Difficulty,
		Hints:      questionTemplate.ExpectedPoints[:2], // 提供前两个关键点作为提示
		TimeLimit:  questionTemplate.TimeLimit,
		NextAction: "wait_for_answer",
		Context: map[string]interface{}{
			"question_template_id": questionTemplate.ID,
			"category":             questionTemplate.Category,
			"difficulty":           questionTemplate.Difficulty,
			"expected_points":      questionTemplate.ExpectedPoints,
		},
	}

	return response, nil
}

// EvaluateAnswer 评估回答
func (g *GoInterviewAgent) EvaluateAnswer(ctx context.Context, req EvaluationRequest) (*EvaluationResponse, error) {
	startTime := time.Now()
	defer func() {
		responseTime := time.Since(startTime)
		g.UpdateMetrics(true, responseTime)
	}()

	// 获取评估标准
	evaluationCriteria := g.getEvaluationCriteria(req.Category, req.Difficulty)

	// 多维度评估
	detailedScores := make(map[string]float32)
	var totalScore float32

	for category, weight := range evaluationCriteria.Weights {
		score := g.evaluateCategory(req.Answer, category, evaluationCriteria.Criteria)
		weightedScore := score * weight
		detailedScores[category] = weightedScore
		totalScore += weightedScore
	}

	// 生成改进建议
	improvements := g.generateImprovements(req.Answer, detailedScores, evaluationCriteria)

	// 确定等级
	grade := g.determineGrade(totalScore)

	response := &EvaluationResponse{
		RequestID:      req.RequestID,
		OverallScore:   totalScore,
		MaxScore:       100.0,
		Grade:          grade,
		Strengths:      g.identifyStrengths(req.Answer, detailedScores),
		Weaknesses:     g.identifyWeaknesses(req.Answer, detailedScores),
		Improvements:   improvements,
		DetailedScores: detailedScores,
		NextTopic:      g.suggestNextTopic(req.Category, detailedScores),
		Context: map[string]interface{}{
			"evaluation_criteria": evaluationCriteria,
			"time_spent":          req.TimeSpent,
			"answer_length":       len(req.Answer),
		},
	}

	return response, nil
}

// selectQuestionTemplate 选择问题模板
func (g *GoInterviewAgent) selectQuestionTemplate(difficulty, category string) *QuestionTemplate {
	// 从知识库中筛选合适的问题模板
	var candidates []QuestionTemplate

	for _, template := range g.knowledgeBase.Questions {
		if template.Difficulty == difficulty && (category == "" || template.Category == category) {
			candidates = append(candidates, template)
		}
	}

	if len(candidates) == 0 {
		return nil
	}

	// 根据上下文和历史记录选择最合适的问题
	// 这里可以实现更复杂的算法，如基于用户历史表现的自适应选择
	return &candidates[0]
}

// personalizeQuestion 个性化问题
func (g *GoInterviewAgent) personalizeQuestion(template *QuestionTemplate, resumeContent string, context map[string]interface{}) string {
	// 基于简历内容进行个性化
	if resumeContent != "" {
		// 提取简历中的关键信息，如项目经验、技能等
		// 将相关信息融入问题中
		return fmt.Sprintf("基于您在简历中提到的项目经验，%s", template.Question)
	}

	return template.Question
}

// getEvaluationCriteria 获取评估标准
func (g *GoInterviewAgent) getEvaluationCriteria(category, difficulty string) *EvaluationCriteria {
	for _, criteria := range g.knowledgeBase.Evaluations {
		if criteria.Category == category {
			return &criteria
		}
	}

	// 默认评估标准
	return &EvaluationCriteria{
		Category: "general",
		Weights: map[string]float32{
			"technical_depth":      0.4,
			"practical_experience": 0.3,
			"problem_solving":      0.2,
			"communication":        0.1,
		},
		Criteria: []EvaluationItem{
			{Name: "technical_depth", Description: "技术深度理解", MaxScore: 40},
			{Name: "practical_experience", Description: "实践经验", MaxScore: 30},
			{Name: "problem_solving", Description: "问题解决能力", MaxScore: 20},
			{Name: "communication", Description: "沟通表达", MaxScore: 10},
		},
	}
}

// evaluateCategory 评估具体类别
func (g *GoInterviewAgent) evaluateCategory(answer, category string, criteria []EvaluationItem) float32 {
	// 这里实现具体的评估逻辑
	// 可以结合NLP技术进行语义分析、关键词匹配等

	switch category {
	case "technical_depth":
		return g.evaluateTechnicalDepth(answer)
	case "practical_experience":
		return g.evaluatePracticalExperience(answer)
	case "problem_solving":
		return g.evaluateProblemSolving(answer)
	case "communication":
		return g.evaluateCommunication(answer)
	default:
		return 70.0 // 默认中等分数
	}
}

// evaluateTechnicalDepth 评估技术深度
func (g *GoInterviewAgent) evaluateTechnicalDepth(answer string) float32 {
	score := float32(50.0) // 基础分

	// 关键词匹配
	deepKeywords := []string{"底层原理", "源码分析", "实现机制", "设计思想", "性能优化"}
	for _, keyword := range deepKeywords {
		if contains(answer, keyword) {
			score += 10.0
		}
	}

	// 技术细节程度
	if len(answer) > 200 {
		score += 10.0
	}

	return min(score, 100.0)
}

// evaluatePracticalExperience 评估实践经验
func (g *GoInterviewAgent) evaluatePracticalExperience(answer string) float32 {
	score := float32(40.0) // 基础分

	// 项目经验关键词
	experienceKeywords := []string{"项目", "线上", "优化", "调优", "踩坑", "最佳实践"}
	for _, keyword := range experienceKeywords {
		if contains(answer, keyword) {
			score += 12.0
		}
	}

	return min(score, 100.0)
}

// evaluateProblemSolving 评估问题解决能力
func (g *GoInterviewAgent) evaluateProblemSolving(answer string) float32 {
	score := float32(45.0) // 基础分

	// 问题解决关键词
	solvingKeywords := []string{"分析", "定位", "解决", "方案", "思路"}
	for _, keyword := range solvingKeywords {
		if contains(answer, keyword) {
			score += 11.0
		}
	}

	return min(score, 100.0)
}

// evaluateCommunication 评估沟通能力
func (g *GoInterviewAgent) evaluateCommunication(answer string) float32 {
	score := float32(60.0) // 基础分

	// 表达清晰度
	if len(answer) > 100 && hasClearStructure(answer) {
		score += 20.0
	}

	// 逻辑性
	if hasLogicalFlow(answer) {
		score += 20.0
	}

	return min(score, 100.0)
}

// generateImprovements 生成改进建议
func (g *GoInterviewAgent) generateImprovements(answer string, scores map[string]float32, criteria *EvaluationCriteria) []string {
	var improvements []string

	for category, score := range scores {
		if score < 70.0 {
			switch category {
			case "technical_depth":
				improvements = append(improvements, "建议深入学习Go语言的底层实现原理")
				improvements = append(improvements, "多阅读Go标准库源码，理解设计思想")
			case "practical_experience":
				improvements = append(improvements, "多参与实际项目，积累项目经验")
				improvements = append(improvements, "总结项目中遇到的问题和解决方案")
			case "problem_solving":
				improvements = append(improvements, "培养系统化的问题分析思维")
				improvements = append(improvements, "学习常见的设计模式和解决方案")
			case "communication":
				improvements = append(improvements, "提高技术表达的条理性和逻辑性")
				improvements = append(improvements, "多进行技术分享和演讲练习")
			}
		}
	}

	return improvements
}

// identifyStrengths 识别优势
func (g *GoInterviewAgent) identifyStrengths(answer string, scores map[string]float32) []string {
	var strengths []string

	for category, score := range scores {
		if score >= 85.0 {
			switch category {
			case "technical_depth":
				strengths = append(strengths, "对Go语言有深入的技术理解")
			case "practical_experience":
				strengths = append(strengths, "具备丰富的项目实践经验")
			case "problem_solving":
				strengths = append(strengths, "展现出优秀的问题解决能力")
			case "communication":
				strengths = append(strengths, "技术表达清晰，沟通能力强")
			}
		}
	}

	return strengths
}

// identifyWeaknesses 识别弱点
func (g *GoInterviewAgent) identifyWeaknesses(answer string, scores map[string]float32) []string {
	var weaknesses []string

	for category, score := range scores {
		if score < 60.0 {
			switch category {
			case "technical_depth":
				weaknesses = append(weaknesses, "Go语言技术深度需要加强")
			case "practical_experience":
				weaknesses = append(weaknesses, "缺乏足够的项目实践验证")
			case "problem_solving":
				weaknesses = append(weaknesses, "问题解决思路不够清晰")
			case "communication":
				weaknesses = append(weaknesses, "技术表达能力有待提高")
			}
		}
	}

	return weaknesses
}

// suggestNextTopic 建议下一个话题
func (g *GoInterviewAgent) suggestNextTopic(currentCategory string, scores map[string]float32) string {
	// 找到得分最低的类别作为下一个话题
	minScore := float32(100.0)
	nextCategory := "fundamentals"

	for category, score := range scores {
		if score < minScore {
			minScore = score
			nextCategory = category
		}
	}

	return nextCategory
}

// determineGrade 确定等级
func (g *GoInterviewAgent) determineGrade(score float32) string {
	if score >= 90 {
		return "优秀"
	} else if score >= 80 {
		return "良好"
	} else if score >= 70 {
		return "中等"
	} else if score >= 60 {
		return "及格"
	}
	return "需改进"
}

// buildGoKnowledgeBase 构建Go知识库
func buildGoKnowledgeBase() *KnowledgeBase {
	return &KnowledgeBase{
		Domain:     "go",
		Categories: []string{"fundamentals", "concurrency", "memory", "stdlib", "performance", "testing", "microservices"},
		Questions: []QuestionTemplate{
			// 基础知识类别
			{
				ID:             "go_fundamentals_001",
				Category:       "fundamentals",
				Difficulty:     "beginner",
				Question:       "请解释Go语言中值类型和引用类型的区别，并举例说明",
				Keywords:       []string{"值类型", "引用类型", "slice", "map", "channel"},
				ExpectedPoints: []string{"基本类型是值类型", "slice/map/channel是引用类型", "内存分配差异", "函数参数传递差异"},
				TimeLimit:      180,
			},
			{
				ID:             "go_fundamentals_002",
				Category:       "fundamentals",
				Difficulty:     "intermediate",
				Question:       "Go语言中的接口有什么特点？如何实现和使用接口？",
				Keywords:       []string{"接口", "隐式实现", "空接口", "interface{}"},
				ExpectedPoints: []string{"接口是方法签名的集合", "隐式实现机制", "空接口可以存储任何类型", "接口类型断言"},
				TimeLimit:      240,
			},
			{
				ID:             "go_fundamentals_003",
				Category:       "fundamentals",
				Difficulty:     "intermediate",
				Question:       "Go语言的错误处理机制有什么特点？与传统的try/catch有什么不同？",
				Keywords:       []string{"错误处理", "error接口", "panic", "recover"},
				ExpectedPoints: []string{"error接口设计", "显式错误检查", "panic/recover机制", "错误包装和上下文"},
				TimeLimit:      240,
			},

			// 并发类别
			{
				ID:             "go_concurrency_001",
				Category:       "concurrency",
				Difficulty:     "intermediate",
				Question:       "详细解释Goroutine的调度机制，包括GMP模型的原理",
				Keywords:       []string{"Goroutine", "GMP模型", "调度器", "M:P关系"},
				ExpectedPoints: []string{"GMP模型组成", "工作窃取算法", "调度策略", "系统调用处理"},
				TimeLimit:      300,
			},
			{
				ID:             "go_concurrency_002",
				Category:       "concurrency",
				Difficulty:     "intermediate",
				Question:       "Go语言中的channel有什么作用？如何使用channel进行goroutine间通信？",
				Keywords:       []string{"channel", "goroutine通信", "select", "缓冲channel"},
				ExpectedPoints: []string{"channel基本概念和使用", "无缓冲与有缓冲channel区别", "select多路复用", "channel的关闭和range遍历"},
				TimeLimit:      270,
			},
			{
				ID:             "go_concurrency_003",
				Category:       "concurrency",
				Difficulty:     "advanced",
				Question:       "如何在Go中实现安全的并发访问共享数据？请举例说明sync包的使用场景。",
				Keywords:       []string{"sync包", "互斥锁", "读写锁", "WaitGroup", "原子操作"},
				ExpectedPoints: []string{"Mutex和RWMutex使用", "WaitGroup同步多个goroutine", "原子操作实现无锁并发", "sync.Map适用于读多写少场景"},
				TimeLimit:      300,
			},

			// 内存管理类别
			{
				ID:             "go_memory_001",
				Category:       "memory",
				Difficulty:     "advanced",
				Question:       "深入分析Go语言的垃圾回收机制，包括三色标记法和写屏障原理",
				Keywords:       []string{"垃圾回收", "三色标记", "写屏障", "内存管理"},
				ExpectedPoints: []string{"三色标记算法", "写屏障机制", "GC触发条件", "性能优化策略"},
				TimeLimit:      360,
			},
			{
				ID:             "go_memory_002",
				Category:       "memory",
				Difficulty:     "intermediate",
				Question:       "Go语言中的内存分配策略是什么？简述mspan和内存池的概念。",
				Keywords:       []string{"内存分配", "mspan", "内存池", "arena", "mcache"},
				ExpectedPoints: []string{"内存分配器设计", "mspan结构和大小分类", "mcache和mcentral角色", "内存碎片处理"},
				TimeLimit:      330,
			},

			// 标准库类别
			{
				ID:             "go_stdlib_001",
				Category:       "stdlib",
				Difficulty:     "intermediate",
				Question:       "Go语言的context包有什么作用？如何使用context进行请求生命周期管理？",
				Keywords:       []string{"context", "超时控制", "取消信号", "请求范围值"},
				ExpectedPoints: []string{"context基本概念", "WithCancel/WithTimeout使用", "context在请求链路中的传递", "避免context滥用"},
				TimeLimit:      270,
			},
			{
				ID:             "go_stdlib_002",
				Category:       "stdlib",
				Difficulty:     "beginner",
				Question:       "Go语言中json包的使用场景有哪些？如何处理复杂的JSON数据结构？",
				Keywords:       []string{"json", "序列化", "反序列化", "tag"},
				ExpectedPoints: []string{"json.Marshal/Unmarshal基本用法", "结构体tag使用", "自定义编解码", "处理嵌套和复杂结构"},
				TimeLimit:      240,
			},

			// 性能优化类别
			{
				ID:             "go_performance_001",
				Category:       "performance",
				Difficulty:     "advanced",
				Question:       "如何优化Go程序的性能？请分享一些实用的性能调优技巧。",
				Keywords:       []string{"性能优化", "pprof", "内存分配", "并发优化"},
				ExpectedPoints: []string{"性能分析工具使用", "内存分配优化", "goroutine数量控制", "CPU密集型与IO密集型优化策略"},
				TimeLimit:      330,
			},
			{
				ID:             "go_performance_002",
				Category:       "performance",
				Difficulty:     "intermediate",
				Question:       "Go语言中的内存逃逸是什么？如何避免不必要的内存分配？",
				Keywords:       []string{"内存逃逸", "栈分配", "堆分配", "编译器优化"},
				ExpectedPoints: []string{"内存逃逸概念和原因", "使用go build -gcflags=-m分析逃逸", "避免大对象频繁创建", "预分配和对象池使用"},
				TimeLimit:      300,
			},

			// 测试类别
			{
				ID:             "go_testing_001",
				Category:       "testing",
				Difficulty:     "intermediate",
				Question:       "Go语言的testing包提供了哪些功能？如何编写有效的单元测试和基准测试？",
				Keywords:       []string{"testing", "单元测试", "基准测试", "表驱动测试"},
				ExpectedPoints: []string{"testing包基本用法", "表驱动测试模式", "基准测试编写和分析", "测试覆盖率计算"},
				TimeLimit:      300,
			},
			{
				ID:             "go_testing_002",
				Category:       "testing",
				Difficulty:     "intermediate",
				Question:       "如何在Go中进行模拟测试（Mock）？mock和stub的区别是什么？",
				Keywords:       []string{"mock", "stub", "接口模拟", "测试隔离"},
				ExpectedPoints: []string{"接口设计便于测试", "mock库使用（如testify/mock）", "依赖注入方式", "mock与stub概念区分"},
				TimeLimit:      270,
			},
		},
		Evaluations: []EvaluationCriteria{
			{
				Category: "general",
				Weights: map[string]float32{
					"technical_depth":      0.4,
					"practical_experience": 0.3,
					"problem_solving":      0.2,
					"communication":        0.1,
				},
				Criteria: []EvaluationItem{
					{Name: "technical_depth", Description: "技术深度理解", MaxScore: 40},
					{Name: "practical_experience", Description: "实践经验", MaxScore: 30},
					{Name: "problem_solving", Description: "问题解决能力", MaxScore: 20},
					{Name: "communication", Description: "沟通表达", MaxScore: 10},
				},
			},
		},
		Version:   "1.0.0",
		UpdatedAt: time.Now(),
	}
}

// 实现adk.Agent接口的Run方法 - 暂时移除，因为相关类型未定义
// 建议直接使用ProcessInterview和EvaluateAnswer方法

// isFirstInteraction 判断是否是首次交互 - 暂时注释，因为adk.Message类型可能不存在
/*
func isFirstInteraction(messages []adk.Message) bool {
	for _, msg := range messages {
		if msg.Role == adk.RoleAssistant && msg.Content != "" {
			return false
		}
	}
	return true
}
*/

// 辅助函数
func contains(text, keyword string) bool {
	return len(text) > 0 && len(keyword) > 0 && (text == keyword ||
		(len(text) > len(keyword) && containsSubstring(text, keyword)))
}

func containsSubstring(text, substring string) bool {
	for i := 0; i <= len(text)-len(substring); i++ {
		if text[i:i+len(substring)] == substring {
			return true
		}
	}
	return false
}

func hasClearStructure(answer string) bool {
	// 简单的结构检测逻辑
	structureKeywords := []string{"首先", "其次", "然后", "最后", "第一", "第二", "总结"}
	for _, keyword := range structureKeywords {
		if contains(answer, keyword) {
			return true
		}
	}
	return false
}

func hasLogicalFlow(answer string) bool {
	// 简单的逻辑流检测
	logicKeywords := []string{"因为", "所以", "因此", "但是", "然而", "而且"}
	count := 0
	for _, keyword := range logicKeywords {
		if contains(answer, keyword) {
			count++
		}
	}
	return count >= 2
}

func min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

// Run 实现adk.Agent接口的Run方法
func (g *GoInterviewAgent) Run(ctx context.Context, input *adk.AgentInput, options ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	// 简单实现，返回nil
	return nil
}

// GenerateFirstQuestion 生成第一个面试问题
func (g *GoInterviewAgent) GenerateFirstQuestion(ctx context.Context, difficulty string, domain string) (*service.QuestionData, error) {
	// 创建面试请求
	req := InterviewRequest{
		Difficulty:      difficulty,
		Domain:          domain,
		ResumeContent:   "",
		PreviousContext: make(map[string]interface{}),
	}

	// 处理面试请求
	response, err := g.ProcessInterview(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("生成第一个问题失败: %w", err)
	}

	// 构建QuestionData返回
	questionData := &service.QuestionData{
		QuestionText:  response.Question,
		EvalDimension: response.Category,
		Order:         1,
	}

	return questionData, nil
}

// GenerateNextQuestion 基于前一个问题和答案生成下一个问题
func (g *GoInterviewAgent) GenerateNextQuestion(ctx context.Context, previousQuestion string, userAnswer string, difficulty string, domain string) (*service.QuestionData, error) {
	// 首先评估前一个答案
	evalReq := EvaluationRequest{
		Question:   previousQuestion,
		Answer:     userAnswer,
		Category:   domain,
		Difficulty: difficulty,
		TimeSpent:  300, // 默认5分钟
		Context:    make(map[string]interface{}),
	}

	evalResponse, err := g.EvaluateAnswer(ctx, evalReq)
	if err != nil {
		return nil, fmt.Errorf("评估前一个答案失败: %w", err)
	}

	// 使用评估结果中的nextTopic作为下一个问题的领域
	nextDomain := evalResponse.NextTopic
	if nextDomain == "" {
		nextDomain = domain // 如果没有建议的下一个话题，则使用当前领域
	}

	// 创建新的面试请求
	req := InterviewRequest{
		Difficulty: difficulty,
		Domain:     nextDomain,
		PreviousContext: map[string]interface{}{
			"previous_question": previousQuestion,
			"previous_answer":   userAnswer,
			"evaluation":        evalResponse,
		},
	}

	// 处理面试请求获取下一个问题
	response, err := g.ProcessInterview(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("生成下一个问题失败: %w", err)
	}

	// 构建QuestionData返回
	questionData := &service.QuestionData{
		QuestionText:  response.Question,
		EvalDimension: response.Category,
		Order:         2, // 这里可以根据上下文动态设置
	}

	return questionData, nil
}
