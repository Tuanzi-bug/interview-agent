package service

import (
	"ai-eino-interview-agent/chatApp/agent/question"
	"ai-eino-interview-agent/chatApp/agent/session"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

// QuestionData 问题数据结构
type QuestionData struct {
	QuestionText  string `json:"question_text"`
	EvalDimension string `json:"eval_dimension"`
	Order         uint32 `json:"order"`
}

// DialogueData 对话数据结构
type DialogueData struct {
	SpeakerType  string `json:"speaker_type"` // "interviewer" 或 "candidate"
	Content      string `json:"content"`
	DisplayOrder uint32 `json:"display_order"`
}

// QuestionGeneratorResult 问题生成结果
type QuestionGeneratorResult struct {
	Questions []QuestionData `json:"questions"`
	Dialogues []DialogueData `json:"dialogues"`
}

// BuildInterviewPrompt 构建面试问题生成的提示词
// interviewType: "综合面试" 或 "专项面试"
// domain: 面试领域（综合面试：校招/社招；专项面试：java/golang等）
// difficulty: 难度级别（简单/中等/困难）
func BuildInterviewPrompt(questionIndex int, query string, resumeID int64, hasResume bool, dimension string, followUpCount int, interviewType string, domain string, difficulty string) string {
	// 根据面试类型选择维度
	var dimensionMap map[string]string
	if interviewType == "综合面试" {
		dimensionMap = map[string]string{
			"professional_field":         "专业领域",
			"project_experience":         "项目经历",
			"technical_depth":            "技术深度",
			"technical_foundation":       "技术基础",
			"team_collaboration":         "团队协作",
			"system_architecture_design": "系统架构设计",
		}
	} else {
		// 专项面试
		dimensionMap = map[string]string{
			"basic_knowledge_mastery":                "基础知识掌握",
			"working_principle_practical_experience": "工作原理与实践经验",
			"advanced_features_application":          "高级特性应用",
			"problem_troubleshooting_skills":         "问题排查能力",
			"architecture_design_thinking":           "架构设计思维",
			"performance_optimization_ability":       "性能优化能力",
		}
	}
	dimensionCN := dimensionMap[dimension]

	// 构建面试类型和难度的描述
	difficultyDesc := ""
	switch difficulty {
	case "简单":
		difficultyDesc = "初级难度"
	case "中等":
		difficultyDesc = "中级难度"
	case "困难":
		difficultyDesc = "高级难度"
	default:
		difficultyDesc = "中级难度"
	}

	interviewTypeDesc := interviewType
	domainDesc := domain

	if followUpCount == 0 {
		// 主问题
		if questionIndex == 1 {
			// 第一个问题
			if hasResume && resumeID != 0 {
				return fmt.Sprintf(`
面试类型：%s
面试领域：%s
难度级别：%s
评估维度：%s

简历ID：%d

用户补充信息：%s

要求：
1. 只返回JSON格式
2. 只生成面试官的提问，不要生成用户回答
3. dialogues数组中只包含speaker_type为"interviewer"的提问
4. 问题围绕评估维度"%s"进行，难度为%s
5. 生成一个主问题（不是追问）

JSON格式：
{
  "questions": [{"question_text": "问题内容", "eval_dimension": "%s", "order": 1}],
  "dialogues": [{"speaker_type": "interviewer", "content": "提问内容", "display_order": 1}]
}`, interviewTypeDesc, domainDesc, difficultyDesc, dimensionCN, resumeID, query, dimensionCN, difficultyDesc, dimension)
			}
			if hasResume {
				return fmt.Sprintf(`根据以下信息生成一个面试问题。

面试类型：%s
面试领域：%s
难度级别：%s
评估维度：%s

%s

要求：
1. 只返回JSON格式
2. 只生成面试官的提问，不要生成用户回答
3. dialogues数组中只包含speaker_type为"interviewer"的提问
4. 问题围绕评估维度"%s"进行，难度为%s
5. 生成一个主问题（不是追问）

JSON格式：
{
  "questions": [{"question_text": "问题内容", "eval_dimension": "%s", "order": 1}],
  "dialogues": [{"speaker_type": "interviewer", "content": "提问内容", "display_order": 1}]
}`, interviewTypeDesc, domainDesc, difficultyDesc, dimensionCN, query, dimensionCN, difficultyDesc, dimension)
			}
			return fmt.Sprintf(`生成一个面试问题。

面试类型：%s
面试领域：%s
难度级别：%s
评估维度：%s

要求：
1. 只返回JSON格式
2. 只生成面试官的提问，不要生成用户回答
3. dialogues数组中只包含speaker_type为"interviewer"的提问
4. 问题围绕评估维度"%s"进行，难度为%s
5. 生成一个主问题（不是追问）

JSON格式：
{
  "questions": [{"question_text": "问题内容", "eval_dimension": "%s", "order": 1}],
  "dialogues": [{"speaker_type": "interviewer", "content": "提问内容", "display_order": 1}]
}`, interviewTypeDesc, domainDesc, difficultyDesc, dimensionCN, dimensionCN, difficultyDesc, dimension)
		}

		// 后续主问题
		return fmt.Sprintf(`根据以下简历和用户的回答，生成下一个面试问题。

面试类型：%s
面试领域：%s
难度级别：%s
评估维度：%s

%s

要求：
1. 只返回JSON格式
2. 只生成面试官的提问，不要生成用户回答
3. dialogues数组中只包含speaker_type为"interviewer"的提问
4. 问题围绕评估维度"%s"进行，难度为%s，且与之前的问题不同
5. 生成一个主问题（不是追问）

JSON格式：
{
  "questions": [{"question_text": "问题内容", "eval_dimension": "%s", "order": %d}],
  "dialogues": [{"speaker_type": "interviewer", "content": "提问内容", "display_order": 1}]
}`, interviewTypeDesc, domainDesc, difficultyDesc, dimensionCN, query, dimensionCN, difficultyDesc, dimension, questionIndex)
	}

	// 追问
	return fmt.Sprintf(`根据用户对上一个问题的回答，生成一个追问问题。

面试类型：%s
面试领域：%s
难度级别：%s
评估维度：%s

用户的回答：
%s

要求：
1. 只返回JSON格式
2. 只生成面试官的追问，不要生成用户回答
3. dialogues数组中只包含speaker_type为"interviewer"的追问
4. 追问围绕评估维度"%s"进行，难度为%s
5. 追问基于用户回答内容，深入探讨相关话题
6. 这是第%d个追问

JSON格式：
{
  "questions": [{"question_text": "追问内容", "eval_dimension": "%s", "order": %d}],
  "dialogues": [{"speaker_type": "interviewer", "content": "追问内容", "display_order": 1}]
}`, interviewTypeDesc, domainDesc, difficultyDesc, dimensionCN, query, dimensionCN, difficultyDesc, followUpCount, dimension, questionIndex)
}

// GenerateInterviewQuestions 调用智能体生成面试问题
// 返回生成的问题列表和对话列表
// isFirstQuestion: 是否是第一个问题（决定是否需要简历解析工具）
func GenerateInterviewQuestions(ctx context.Context, prompt string, userId uint, isFirstQuestion bool, interview_type string) (*QuestionGeneratorResult, error) {
	// 添加 30分钟 超时，防止无限等待（API 响应可能需要较长时间）
	timeoutCtx, cancel := context.WithTimeout(ctx, 1800*time.Second)
	defer cancel()

	//根据面试类型，使用不同的智能体
	agent := question.NewQuestionAgent(userId, isFirstQuestion)
	if interview_type == "专项面试" {
		agent = question.NewSpecialQuestionAgent(userId)
	}

	// 创建 runner
	runner := adk.NewRunner(timeoutCtx, adk.RunnerConfig{
		Agent: agent,
	})

	// 构建查询消息
	query := fmt.Sprintf(`根据以下提示词内容生成面试问题。

提示词内容：
%s

要求：
1. 只返回JSON格式
2. 只生成面试官的提问，不要生成用户回答
3. dialogues数组中只包含speaker_type为"interviewer"的提问

JSON格式：
{
  "questions": [{"question_text": "问题内容", "eval_dimension": "professional_field|project_experience|technical_depth|technical_foundation|team_collaboration|system_architecture_design", "order": 1}],
  "dialogues": [{"speaker_type": "interviewer", "content": "提问内容", "display_order": 1}]
}`, prompt)

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
			return nil, fmt.Errorf("timeout waiting for question generation")
		default:
		}

		event, ok := iter.Next()
		if !ok {
			break
		}

		if event.Err != nil {
			return nil, fmt.Errorf("error during question generation: %w", event.Err)
		}

		// 收集最后一条消息
		if event.Output != nil && event.Output.MessageOutput != nil {
			lastMessage = event.Output.MessageOutput.Message.Content
		}
	}

	// 解析 JSON 结果
	result := &QuestionGeneratorResult{}

	if err := json.Unmarshal([]byte(lastMessage), result); err != nil {
		// 尝试从文本中提取 JSON
		jsonStr := extractJSON(lastMessage)
		if jsonStr == "" {
			return nil, fmt.Errorf("failed to parse question generation result: %w", err)
		}

		// 清理 JSON 字符串中的非法字符
		jsonStr = cleanJSON(jsonStr)

		// 如果提取的是数组格式，需要包装成对象
		trimmedJSON := strings.TrimSpace(jsonStr)
		if len(trimmedJSON) > 0 && trimmedJSON[0] == '[' {
			// 尝试解析为问题数组
			var questions []QuestionData
			if err := json.Unmarshal([]byte(jsonStr), &questions); err == nil {
				result.Questions = questions
				logQuestionResult(result)
				return result, nil
			}
			// 尝试解析为对话数组
			var dialogues []DialogueData
			if err := json.Unmarshal([]byte(jsonStr), &dialogues); err == nil {
				result.Dialogues = dialogues
				// 从对话中提取问题：找到第一个 speaker_type="interviewer" 的对话作为问题
				for _, d := range dialogues {
					if d.SpeakerType == "interviewer" {
						result.Questions = append(result.Questions, QuestionData{
							QuestionText:  d.Content,
							EvalDimension: "professional_field",
							Order:         1,
						})
						break
					}
				}
				logQuestionResult(result)
				return result, nil
			}
		}

		// 尝试解析为对象
		if err := json.Unmarshal([]byte(jsonStr), result); err != nil {
			return nil, fmt.Errorf("failed to parse extracted JSON: %w", err)
		}
	}

	logQuestionResult(result)
	return result, nil
}

// logQuestionResult 输出 AI 生成的问题内容到终端控制台
func logQuestionResult(result *QuestionGeneratorResult) {
	log.Println("============ AI Generated Interview Question ============")
	if len(result.Questions) > 0 {
		for i, q := range result.Questions {
			log.Printf("[Question %d] Dimension: %s, Order: %d", i+1, q.EvalDimension, q.Order)
			log.Printf("[Question %d] Content: %s", i+1, q.QuestionText)
		}
	}
	if len(result.Dialogues) > 0 {
		for i, d := range result.Dialogues {
			log.Printf("[Dialogue %d] Speaker: %s, Order: %d", i+1, d.SpeakerType, d.DisplayOrder)
			log.Printf("[Dialogue %d] Content: %s", i+1, d.Content)
		}
	}
	log.Println("==========================================================")
}

// extractJSON 从文本中提取 JSON 字符串
func extractJSON(text string) string {
	// 先清理文本中的非法字符
	text = cleanJSON(text)

	// 优先查找对象格式 {...}
	start := strings.Index(text, "{")
	if start != -1 {
		// 从后往前找最后一个 }
		end := strings.LastIndex(text, "}")
		if end != -1 && end > start {
			return text[start : end+1]
		}
	}

	// 如果没有找到对象，尝试查找数组格式 [...]
	start = strings.Index(text, "[")
	if start != -1 {
		// 从后往前找最后一个 ]
		end := strings.LastIndex(text, "]")
		if end != -1 && end > start {
			return text[start : end+1]
		}
	}

	return ""
}

// cleanJSON 清理 JSON 字符串中的非法字符
func cleanJSON(jsonStr string) string {
	// 移除 BOM 标记
	if len(jsonStr) >= 3 && jsonStr[0] == 0xEF && jsonStr[1] == 0xBB && jsonStr[2] == 0xBF {
		jsonStr = jsonStr[3:]
	}

	// 使用 strings.Builder 来构建清理后的字符串
	var builder strings.Builder
	for i := 0; i < len(jsonStr); i++ {
		b := jsonStr[i]

		// 保留 ASCII 可打印字符和常见的 JSON 字符
		if b >= 0x20 && b <= 0x7E {
			// ASCII 可打印字符
			builder.WriteByte(b)
		} else if b == '\t' || b == '\n' || b == '\r' {
			// 保留制表符、换行符、回车符
			builder.WriteByte(b)
		} else if b >= 0x80 {
			// 处理 UTF-8 多字节字符
			// 检查是否是有效的 UTF-8 序列开始
			if (b & 0xE0) == 0xC0 {
				// 2 字节字符
				if i+1 < len(jsonStr) {
					builder.WriteByte(b)
					i++
					builder.WriteByte(jsonStr[i])
				}
			} else if (b & 0xF0) == 0xE0 {
				// 3 字节字符
				if i+2 < len(jsonStr) {
					builder.WriteByte(b)
					i++
					builder.WriteByte(jsonStr[i])
					i++
					builder.WriteByte(jsonStr[i])
				}
			} else if (b & 0xF8) == 0xF0 {
				// 4 字节字符
				if i+3 < len(jsonStr) {
					builder.WriteByte(b)
					i++
					builder.WriteByte(jsonStr[i])
					i++
					builder.WriteByte(jsonStr[i])
					i++
					builder.WriteByte(jsonStr[i])
				}
			}
		}
		// 其他字符（控制字符等）都被跳过
	}

	return builder.String()
}

// BuildPromptFromSessionContext 从会话值构建提示词
// 优化版本：直接从会话值获取所需信息，无需参数传递
func BuildPromptFromSessionContext(ctx context.Context, questionIndex int, query string,
	dimension string, followUpCount int) (string, error) {
	scm := session.NewSessionContextManager(ctx)

	// 从会话值获取配置
	interviewType, domain, difficulty, err := scm.GetInterviewConfig()
	if err != nil {
		return "", fmt.Errorf("failed to get interview config: %w", err)
	}

	// 调用原有的 BuildInterviewPrompt 函数
	prompt := BuildInterviewPrompt(questionIndex, query, 0, false, dimension,
		followUpCount, interviewType, domain, difficulty)

	return prompt, nil
}

// GenerateInterviewQuestionsWithSessionContext 使用会话值生成面试问题
// 优化版本：从会话值获取参数，减少参数传递
//func GenerateInterviewQuestionsWithSessionContext(ctx context.Context, questionIndex int,
//	query string, dimension string, followUpCount int, userId uint) (*QuestionGeneratorResult, error) {
//
//	// 从会话值构建提示词
//	prompt, err := BuildPromptFromSessionContext(ctx, questionIndex, query, dimension, followUpCount)
//	if err != nil {
//		return nil, fmt.Errorf("failed to build prompt from session context: %w", err)
//	}
//
//	// 只有第一个问题才需要简历解析工具
//	isFirstQuestion := questionIndex == 1 && followUpCount == 0
//	return GenerateInterviewQuestions(ctx, prompt, userId, isFirstQuestion, "")
//}
