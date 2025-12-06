package service

import (
	"ai-eino-interview-agent/chatApp/agent/resume"
	"ai-eino-interview-agent/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
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
// 参数说明：
//   - ctx: 上下文
//   - userId: 用户ID
//   - resumeFilePath: 上传的简历文件路径（已保存到 backend/uploads/resumes）
//   - fileSize: 文件大小
//
// 返回简历ID和解析结果
func ParseResumeAndSave(ctx context.Context, userId uint, resumeFilePath string, fileSize int64) (uint64, *ResumeParseResult, error) {
	// 添加 120 秒超时
	timeoutCtx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	// 创建简历解析智能体
	agent, err := resume.NewResumeParserAgent(userId)
	if err != nil {
		log.Printf("[ParseResumeAndSave] 创建简历解析智能体失败: %v", err)
		return 0, nil, err
	}

	// 创建 runner
	runner := adk.NewRunner(timeoutCtx, adk.RunnerConfig{
		Agent: agent,
	})

	// 构建查询消息，包含简历文件路径
	query := fmt.Sprintf(`【重要】请立即解析以下简历文件并提取关键信息：

简历文件路径：%s

【必须执行的步骤】：
1. 【第一步】立即使用 pdf_to_text 工具解析简历文件，获取完整的简历文本内容
2. 【第二步】从解析的简历文本中提取所有关键信息（姓名、工作年限、联系方式、教育背景、工作经历、技术栈、项目经验、技能、证书等）
3. 【第三步】分析候选人的背景特点和核心竞争力
4. 【第四步】生成面试建议和推荐难度

【重要提示】：
- 不要跳过 pdf_to_text 工具调用
- 必须从简历内容中提取真实的信息，不要返回空数据
- 所有JSON字段都必须填充实际内容
- 只返回JSON格式，不要返回其他文本

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
	if lastMessage == "" {
		log.Printf("[ParseResumeAndSave] 智能体未返回任何响应")
		return 0, nil, fmt.Errorf("agent returned empty response")
	}

	log.Printf("[ParseResumeAndSave] 智能体响应内容: %s", lastMessage)
	parseResult := parseResumeResponse(lastMessage)
	if parseResult == nil {
		log.Printf("[ParseResumeAndSave] 无法解析简历响应")
		return 0, nil, fmt.Errorf("failed to parse resume response")
	}

	// 验证解析结果是否有效（不能全是空数据）
	if !isValidResumeResult(parseResult) {
		log.Printf("[ParseResumeAndSave] 解析结果无效（全是空数据），请检查简历文件是否正确")
		return 0, nil, fmt.Errorf("resume parsing result is empty or invalid")
	}

	// 将解析结果保存到数据库
	resumeID, err := saveResumeToDatabase(ctx, userId, resumeFilePath, fileSize, parseResult)
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
		log.Printf("[parseResumeResponse] 直接解析 JSON 失败: %v，尝试提取 JSON", err)
		// 尝试从文本中提取 JSON
		jsonStr := ExtractJSONFromResponse(agentResponse)
		if jsonStr == "" {
			log.Printf("[parseResumeResponse] 无法提取 JSON，原始响应: %s", agentResponse)
			return nil
		}

		log.Printf("[parseResumeResponse] 提取的 JSON: %s", jsonStr)
		// 尝试解析提取的 JSON
		if err := json.Unmarshal([]byte(jsonStr), result); err != nil {
			log.Printf("[parseResumeResponse] 解析提取的 JSON 失败: %v", err)
			return nil
		}
	}

	return result
}

// isValidResumeResult 检查解析结果是否有效（不能全是空数据）
func isValidResumeResult(result *ResumeParseResult) bool {
	if result == nil {
		return false
	}

	// 检查基本信息是否有内容
	if result.BasicInfo.Name != "" || result.BasicInfo.WorkYears != "" || result.BasicInfo.Contact != "" {
		return true
	}

	// 检查教育背景
	if len(result.Education) > 0 {
		return true
	}

	// 检查工作经历
	if len(result.WorkExperience) > 0 {
		return true
	}

	// 检查技术栈
	if len(result.TechStack) > 0 {
		return true
	}

	// 检查项目经验
	if len(result.Projects) > 0 {
		return true
	}

	// 检查技能
	if len(result.Skills) > 0 {
		return true
	}

	// 检查证书
	if len(result.Certifications) > 0 {
		return true
	}

	// 检查其他字段
	if result.Strengths != "" || result.PotentialWeaknesses != "" || result.RecommendedDifficulty != "" {
		return true
	}

	// 检查面试关注领域
	if len(result.InterviewFocusAreas) > 0 {
		return true
	}

	// 检查建议的提问方向
	if len(result.SuggestedQuestionDirections) > 0 {
		return true
	}

	// 如果所有字段都是空的，返回 false
	return false
}

// saveResumeToDatabase 将简历解析结果保存到数据库
// 参数说明：
//   - ctx: 上下文
//   - userId: 用户ID
//   - resumeFilePath: 原始 PDF 文件路径（已保存到 backend/uploads/resumes）
//   - fileSize: 文件大小
//   - parseResult: 解析后的简历数据
func saveResumeToDatabase(ctx context.Context, userId uint, resumeFilePath string, fileSize int64, parseResult *ResumeParseResult) (uint64, error) {
	// 将解析结果转换为 JSON 字符串存储
	contentJSON, err := json.Marshal(parseResult)
	if err != nil {
		log.Printf("[saveResumeToDatabase] 序列化简历数据失败: %v", err)
		return 0, fmt.Errorf("failed to marshal resume data: %w", err)
	}

	// 从文件路径中提取文件名
	fileName := filepath.Base(resumeFilePath)
	log.Printf("[saveResumeToDatabase] 原始文件路径: %s, 提取的文件名: %s", resumeFilePath, fileName)

	// 创建简历记录
	// Content 字段存储解析后的 JSON 数据，用于快速查询
	// FileName 字段存储原始 PDF 文件名
	// FileType 字段标记为 "pdf"，表示这是 PDF 简历
	resumeRecord := &model.Resume{
		UserID:    userId,
		Content:   string(contentJSON),
		FileName:  fileName,
		FileSize:  fileSize,
		FileType:  "pdf",
		IsDefault: 1,
		Deleted:   0,
	}

	// 调用 DAO 方法保存到数据库
	resumeID, err := model.ResumeDao.CreateResume(resumeRecord)
	if err != nil {
		log.Printf("[saveResumeToDatabase] 创建简历记录失败: %v", err)
		return 0, fmt.Errorf("failed to create resume record: %w", err)
	}

	log.Printf("[saveResumeToDatabase] 简历记录已保存，ID: %d, 用户ID: %d, 文件路径: %s", resumeID, userId, resumeFilePath)
	return resumeID, nil
}
