package ext

import (
	"ai-eino-interview-agent/chatApp/agent/question"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

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

// GenerateInterviewQuestions 调用智能体生成面试问题
// 返回生成的问题列表和对话列表
func GenerateInterviewQuestions(ctx context.Context, resumeContent string) (*QuestionGeneratorResult, error) {
	// 创建问题生成智能体
	agent := question.NewQuestionAgent()

	// 创建 runner
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent: agent,
	})

	// 构建查询消息
	query := fmt.Sprintf(`如果简历存在则请根据以下简历内容生成适当数量的面试问题。

简历内容：
%s

请按照以下JSON格式返回结果：
{
  "questions": [
    {
      "question_text": "问题文本",
      "eval_dimension": "professional_field|project_experience|technical_depth|technical_foundation|team_collaboration|system_architecture_design",
      "order": 1
    }
  ],
  "dialogues": [
    {
      "speaker_type": "interviewer",
      "content": "问题内容",
      "display_order": 1
    },
    {
      "speaker_type": "candidate",
      "content": "回答内容",
      "display_order": 2
    }
  ]
}`, resumeContent)

	// 创建用户消息
	userMsg := &schema.Message{
		Role:    schema.User,
		Content: query,
	}

	messages := []adk.Message{
		userMsg,
	}

	// 运行智能体
	iter := runner.Run(ctx, messages)

	var lastMessage string
	for {
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
			log.Printf("Failed to extract JSON from response. Raw message: %s", lastMessage)
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
				return result, nil
			}
			// 尝试解析为对话数组
			var dialogues []DialogueData
			if err := json.Unmarshal([]byte(jsonStr), &dialogues); err == nil {
				result.Dialogues = dialogues
				log.Printf("[DEBUG] 解析为对话数组，共 %d 条对话", len(dialogues))
				// 从对话中提取问题：找到第一个 speaker_type="interviewer" 的对话作为问题
				for i, d := range dialogues {
					log.Printf("[DEBUG] 对话 %d: speaker_type='%s', content='%.50s'", i, d.SpeakerType, d.Content)
					if d.SpeakerType == "interviewer" {
						result.Questions = append(result.Questions, QuestionData{
							QuestionText:  d.Content,
							EvalDimension: "professional_field", // 默认维度
							Order:         1,
						})
						log.Printf("[DEBUG] 从对话中提取问题成功")
						break // 只取第一个问题
					}
				}
				if len(result.Questions) == 0 {
					log.Printf("[DEBUG] 警告：从对话数组中未找到 interviewer 类型的对话")
				}
				return result, nil
			}
		}

		// 尝试解析为对象
		if err := json.Unmarshal([]byte(jsonStr), result); err != nil {
			log.Printf("Failed to parse extracted JSON: %s, Error: %v", jsonStr, err)
			return nil, fmt.Errorf("failed to parse extracted JSON: %w", err)
		}
	}

	return result, nil
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

// ParseQuestionResponse 解析问题生成响应中的 JSON
// 这是一个辅助函数，用于从智能体的文本响应中提取结构化数据
func ParseQuestionResponse(responseText string) (*QuestionGeneratorResult, error) {
	// 尝试直接解析
	result := &QuestionGeneratorResult{}
	if err := json.Unmarshal([]byte(responseText), result); err != nil {
		// 尝试从文本中提取 JSON
		jsonStr := extractJSON(responseText)
		if jsonStr == "" {
			log.Printf("Failed to extract JSON from response: %s", responseText)
			return nil, fmt.Errorf("no JSON found in response")
		}
		if err := json.Unmarshal([]byte(jsonStr), result); err != nil {
			return nil, fmt.Errorf("failed to parse extracted JSON: %w", err)
		}
	}
	return result, nil
}
