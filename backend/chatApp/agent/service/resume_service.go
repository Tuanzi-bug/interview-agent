package service

import (
	"ai-eino-interview-agent/chatApp/agent/resume"
	"ai-eino-interview-agent/chatApp/agent/session"
	"ai-eino-interview-agent/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

// ResumeParseResult 简历解析结果
type ResumeParseResult struct {
	BasicInfo struct {
		Name      string `json:"name"`
		WorkYears string `json:"work_years"`
		Contact   string `json:"contact"`
	} `json:"basic_info"`
	Education []struct {
		School         string `json:"school"`
		Major          string `json:"major"`
		Degree         string `json:"degree"`
		GraduationYear string `json:"graduation_year"`
	} `json:"education"`
	WorkExperience []struct {
		Company          string `json:"company"`
		Position         string `json:"position"`
		Duration         string `json:"duration"`
		Responsibilities string `json:"responsibilities"`
	} `json:"work_experience"`
	TechStack                   []string      `json:"tech_stack"`
	Projects                    []interface{} `json:"projects"`
	Skills                      []string      `json:"skills"`
	Certifications              []string      `json:"certifications"`
	Strengths                   string        `json:"strengths"`
	PotentialWeaknesses         string        `json:"potential_weaknesses"`
	RecommendedDifficulty       string        `json:"recommended_difficulty"`
	InterviewFocusAreas         []string      `json:"interview_focus_areas"`
	SuggestedQuestionDirections []string      `json:"suggested_questions_directions"`
}

// ParseResumeAndSave 调用简历解析智能体解析简历，并将结果保存到数据库
// 从 session 中获取 userID 和 resumeFilePath，无需通过参数传递
// 返回简历ID和解析结果
// 同时将解析结果存储到会话上下文中，供后续智能体使用
func ParseResumeAndSave(ctx context.Context) (uint64, *ResumeParseResult, error) {
	// 添加 120 秒超时
	timeoutCtx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	// 初始化会话上下文管理器
	sessionMgr := session.NewSessionContextManager(timeoutCtx)

	// 从 session 中获取参数
	userId, resumeFilePath, _, _, fileSize, _, err := sessionMgr.GetResumeUploadParams()
	if err != nil {
		log.Printf("[ParseResumeAndSave] 从 session 获取参数失败: %v", err)
		return 0, nil, fmt.Errorf("failed to get resume params from session: %w", err)
	}

	if userId == 0 || resumeFilePath == "" {
		return 0, nil, fmt.Errorf("invalid resume params: userID=%d, resumeFilePath=%s", userId, resumeFilePath)
	}

	// 创建简历解析智能体
	agent := resume.NewResumeParserAgent(userId)

	// 创建 runner
	runner := adk.NewRunner(timeoutCtx, adk.RunnerConfig{
		Agent: agent,
	})

	// 构建查询消息，包含简历文件路径
	query := fmt.Sprintf(`请解析以下简历文件并提取关键信息：

简历文件路径：%s

请按照以下步骤进行：
1. 使用 pdf_to_text 工具解析简历内容
2. 从简历中提取所有关键信息
3. 分析候选人的背景特点
4. 生成面试建议

请返回完整的 JSON 格式结果。`, resumeFilePath)

	// 创建用户消息
	userMsg := &schema.Message{
		Role:    schema.User,
		Content: query,
	}

	messages := []adk.Message{
		userMsg,
	}

	// 运行智能体
	iter := runner.Run(timeoutCtx, messages)

	var lastMessage string
	for {
		select {
		case <-timeoutCtx.Done():
			log.Printf("[ParseResumeAndSave] 超时：等待智能体响应超过 120 秒")
			return 0, nil, fmt.Errorf("timeout waiting for resume parsing (120s)")
		default:
		}

		event, ok := iter.Next()
		if !ok {
			break
		}

		if event.Err != nil {
			log.Printf("[ParseResumeAndSave] 错误: %v", event.Err)
			return 0, nil, fmt.Errorf("error during resume parsing: %w", event.Err)
		}

		// 收集最后一条消息
		if event.Output != nil && event.Output.MessageOutput != nil {
			lastMessage = event.Output.MessageOutput.Message.Content
		}
	}

	// 解析智能体响应
	parseResult := parseResumeResponse(lastMessage)
	if parseResult == nil {
		log.Printf("[ParseResumeAndSave] 无法解析简历响应")
		return 0, nil, fmt.Errorf("failed to parse resume response")
	}
	// 将解析结果保存到数据库
	resumeID, err := saveResumeToDatabase(ctx, userId, fileSize, parseResult)
	if err != nil {
		log.Printf("[ParseResumeAndSave] 保存简历失败: %v", err)
		return 0, nil, fmt.Errorf("failed to save resume: %w", err)
	}

	log.Printf("[ParseResumeAndSave] 简历解析成功，简历ID: %d", resumeID)
	return resumeID, parseResult, nil
}

// parseResumeResponse 从智能体响应解析简历数据
func parseResumeResponse(agentResponse string) *ResumeParseResult {
	result := &ResumeParseResult{}

	// 尝试直接解析 JSON
	if err := json.Unmarshal([]byte(agentResponse), result); err != nil {
		// 尝试从文本中提取 JSON
		jsonStr := extractJSONFromResponse(agentResponse)
		if jsonStr == "" {
			log.Printf("[parseResumeResponse] 无法提取 JSON")
			return nil
		}

		// 尝试解析提取的 JSON
		if err := json.Unmarshal([]byte(jsonStr), result); err != nil {
			log.Printf("[parseResumeResponse] 解析提取的 JSON 失败: %v", err)
			return nil
		}
	}

	return result
}

// convertResumeToMap 将 ResumeParseResult 转换为 map[string]interface{}
// 用于存储到会话上下文
func convertResumeToMap(parseResult *ResumeParseResult) map[string]interface{} {
	if parseResult == nil {
		return make(map[string]interface{})
	}

	// 序列化为 JSON 再反序列化为 map，确保完整的数据转换
	data, _ := json.Marshal(parseResult)
	var resultMap map[string]interface{}
	_ = json.Unmarshal(data, &resultMap)
	return resultMap
}

// saveResumeToDatabase 将简历解析结果保存到数据库
func saveResumeToDatabase(ctx context.Context, userId uint, fileSize int64, parseResult *ResumeParseResult) (uint64, error) {
	// 将解析结果转换为 JSON 字符串存储
	contentJSON, err := json.Marshal(parseResult)
	if err != nil {
		log.Printf("[saveResumeToDatabase] 序列化简历数据失败: %v", err)
		return 0, fmt.Errorf("failed to marshal resume data: %w", err)
	}

	// 创建简历记录
	resumeRecord := &model.Resume{
		UserID:    userId,
		Content:   string(contentJSON),
		FileName:  fmt.Sprintf("resume_%s.json", parseResult.BasicInfo.Name),
		FileSize:  fileSize,
		FileType:  "json",
		IsDefault: 1,
		Deleted:   0,
	}

	// 调用 DAO 方法保存到数据库
	resumeID, err := model.ResumeDao.CreateResume(resumeRecord)
	if err != nil {
		log.Printf("[saveResumeToDatabase] 创建简历记录失败: %v", err)
		return 0, fmt.Errorf("failed to create resume record: %w", err)
	}

	log.Printf("[saveResumeToDatabase] 简历记录已保存，ID: %d, 用户ID: %d", resumeID, userId)
	return resumeID, nil
}
