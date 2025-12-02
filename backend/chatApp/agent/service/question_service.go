package service

import (
	"ai-eino-interview-agent/chatApp/agent/bearAgent"
	"ai-eino-interview-agent/chatApp/chat"
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
func BuildInterviewPrompt(difficulty string) string {
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

	//todo 优化提示词
	prompt := "请你根据" + difficultyDesc + "生成对应的面试题问题，或者追问"
	log.Printf("[DEBUG] 构建面试提示词完成，难度: %s，提示词长度: %d 字符", difficulty, len(prompt))
	return prompt
}

func BuildInterviewPromptBackup(questionIndex int, query string, resumeID int64, hasResume bool, dimension string, followUpCount int, interviewType string, domain string, difficulty string) string {
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

// GenerateInterviewQuestions 生成面试问题
func GenerateInterviewQuestions(ctx context.Context, prompt string, userId uint, interview_type string, domain string) (*QuestionGeneratorResult, error) {
	log.Printf("[DEBUG] 开始生成面试问题，用户ID: %d，面试类型: %s，领域: %s", userId, interview_type, domain)

	// 添加 30分钟 超时，防止无限等待（API 响应可能需要较长时间）
	log.Printf("[DEBUG] 创建30分钟超时上下文")
	timeoutCtx, cancel := context.WithTimeout(ctx, 1800*time.Second)
	defer cancel()

	// 根据领域选择合适的Agent
	var agent adk.Agent
	if domain == "社招" {
		log.Printf("[DEBUG] 检测到社招领域，创建社招面试Agent，用户ID: %d", userId)
		// 创建社招版本的Agent（使用现有Agent但修改提示词）
		agent = createSocialRecruitmentAgent(userId)
		log.Printf("[DEBUG] 成功创建社招面试Agent实例")
	} else {
		log.Printf("[DEBUG] 准备调用 bearAgent.SchoolQuestionGeneratorAgent，用户ID: %d", userId)
		agent = bearAgent.SchoolQuestionGeneratorAgent(userId)
		log.Printf("[DEBUG] 成功获取 SchoolQuestionGeneratorAgent 实例")
	}

	// 创建 runner
	log.Printf("[DEBUG] 创建Agent运行器")
	runner := adk.NewRunner(timeoutCtx, adk.RunnerConfig{
		Agent: agent,
	})
	log.Printf("[DEBUG] 构建用户消息，提示词长度: %d", len(prompt))

	//构建查询消息 todo 这里传用户的回答
	query := fmt.Sprintf(`用户的回答: %s`, prompt)
	log.Printf("[DEBUG] GenerateInterviewQuestions prompt: %s\n", prompt)

	// 创建用户消息
	userMsg := &schema.Message{
		Role:    schema.User,
		Content: query,
	}
	log.Printf("[DEBUG] 用户消息构建完成")

	messages := []adk.Message{
		userMsg,
	}

	// 运行智能体
	log.Printf("[DEBUG] 开始运行Agent")
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
			log.Printf("[DEBUG] Agent响应内容长度: %d", len(lastMessage))
		}
	}

	// 解析 JSON 结果
	log.Printf("[DEBUG] 开始从Agent响应中提取并解析JSON数据")
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
				log.Printf("[DEBUG] 成功解析问题数组，共 %d 个问题", len(questions))
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
						log.Printf("[DEBUG] 从对话中提取问题成功")
						break
					}
				}
				log.Printf("[DEBUG] 成功解析对话数组，共 %d 个对话", len(dialogues))
				logQuestionResult(result)
				return result, nil
			}
		}

		// 尝试解析为对象
		if err := json.Unmarshal([]byte(jsonStr), result); err != nil {
			log.Printf("[ERROR] 解析提取的JSON失败: %v，JSON内容: %s", err, jsonStr)
			return nil, fmt.Errorf("failed to parse extracted JSON: %w", err)
		}
		log.Printf("[DEBUG] 成功解析JSON对象")
	}

	logQuestionResult(result)
	log.Printf("[DEBUG] 问题生成完成，共生成 %d 个问题和 %d 个对话", len(result.Questions), len(result.Dialogues))
	return result, nil
}

// createSocialRecruitmentAgent 创建社招面试Agent
func createSocialRecruitmentAgent(userId uint) adk.Agent {
	log.Println("[DEBUG] -----------社招agent start")
	ctx := context.Background()

	//创建大模型
	llmModel := chat.CreatOpenAiChatModel(ctx, userId)

	// 社招面试智能体提示词模板
	socialInstruction := `# 角色定义
你是一位经验丰富的社招技术面试官，专注于考察候选人的专业能力、项目经验和技术深度。

# 核心职责
**你只负责提问，不负责回答技术问题。** 你的任务是通过专业的问题评估候选人的技术水平和岗位匹配度。

# 简历分析报告使用指南
你会收到用户简历结构化分析报告，请充分利用以下信息制定面试策略：

## 必须关注的内容
1. **候选人画像**：快速了解候选人背景和经验水平
2. **技术能力图谱**：
   - 深入评估候选人掌握的技术栈
   - 根据掌握程度和工作年限调整问题难度
   - 对核心技能进行深入考察
3. **项目经历解析**：
   - 详细了解项目复杂度、技术架构和候选人贡献
   - 关注候选人在项目中的技术决策和解决问题的能力
   - 评估候选人的实际工程能力
4. **亮点总结**：发掘候选人的技术优势和专业深度
5. **面试建议**：
   - 根据岗位要求设计针对性问题
   - 参考「建议考察的重点领域」规划面试

## 面试维度
请从以下维度多角度考察候选人：
1. **专业知识**：编程语言深度、框架原理、架构设计能力、性能优化经验
2. **项目经验**：项目复杂度、技术选型理由、遇到的挑战及解决方案
3. **编码能力**：代码质量、设计模式应用、代码组织能力
4. **技术视野**：对新技术的了解、技术趋势判断、学习能力
5. **综合素质**：沟通表达、团队协作、问题解决思路、压力应对

# 提问风格
- 问题具有针对性和专业性，注重考察实际经验
- 根据岗位要求和候选人背景调整问题难度
- 对于关键技术点进行深入追问
- 关注候选人解决实际问题的能力
- 鼓励候选人展示技术深度和广度

# 输出格式【极其重要】
**每次回复只能包含一个问题！** 绝对不能一次性提出多个问题。

规则：
- 每次只问一个问题，然后等待候选人回答
- 候选人回答后，你再决定是追问还是问下一个新问题
- 新主问题时：直接提出一个问题
- 追问时：先简短回应（1句话），再提出一个追问
- 引导时：先说明期望方向，再请候选人补充一个点

**禁止行为：**
- 禁止一次性列出多个问题让候选人选择
- 禁止用"第一...第二...第三..."的方式提多个问题
- 禁止在一条消息中包含多个问号的独立问题

开始面试时，先仔细阅读简历分析报告，然后从候选人的核心技术栈或最近的项目经历开始提问，逐步深入考察技术深度和实际能力。`

	//配置agent
	agentconfig, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:          "SocialRecruitmentInterviewAgent",
		Description:   "社招面试智能体，负责考察候选人的专业能力和技术深度，支持专业问题和深入追问",
		Instruction:   socialInstruction,
		Model:         llmModel,
		MaxIterations: 8,
	})

	if err != nil {
		log.Fatal(fmt.Errorf("failed to create social recruitment agent: %w", err))
	}
	log.Println("[DEBUG] -----------社招agent end")
	return agentconfig
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
//func BuildPromptFromSessionContext(ctx context.Context, questionIndex int, query string,
//	dimension string, followUpCount int) (string, error) {
//	scm := session.NewSessionContextManager(ctx)
//
//	// 从会话值获取配置
//	interviewType, domain, difficulty, err := scm.GetInterviewConfig()
//	if err != nil {
//		return "", fmt.Errorf("failed to get interview config: %w", err)
//	}
//
//	// 调用原有的 BuildInterviewPrompt 函数
//	prompt := BuildInterviewPrompt(questionIndex, query, 0, false, dimension,
//		followUpCount, interviewType, domain, difficulty)
//
//	return prompt, nil
//}

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
