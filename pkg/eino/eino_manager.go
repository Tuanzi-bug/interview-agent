package eino

import (
	"context"
	"log"

	"ai-eino-interview-agent/internal/config"

	"github.com/cloudwego/eino"
	"github.com/cloudwego/eino-ext/components/model/openai"
)

// EinoAgent 全局Eino代理实例
var EinoAgent *eino.Agent

// InitEino 初始化Eino框架
func InitEino() error {
	cfg := config.Global.Eino

	// 创建OpenAI模型配置
	modelConfig := &openai.ModelConfig{
		APIKey:      cfg.APIKey,
		BaseURL:     config.Global.OpenAI.BaseURL,
		ModelName:   config.Global.OpenAI.ModelName,
		Temperature: cfg.Temperature,
		MaxTokens:   cfg.MaxTokens,
	}

	// 初始化OpenAI模型
	model, err := openai.NewModel(modelConfig)
	if err != nil {
		return err
	}

	// 创建Eino代理
	EinoAgent = eino.NewAgent(
		eino.WithModel(model),
		eino.WithTemperature(cfg.Temperature),
		eino.WithMaxTokens(cfg.MaxTokens),
	)

	log.Println("Eino框架初始化成功")
	return nil
}

// GetEinoAgent 获取Eino代理实例
func GetEinoAgent() *eino.Agent {
	return EinoAgent
}

// GenerateInterviewQuestions 生成面试问题
func GenerateInterviewQuestions(ctx context.Context, jobTitle, difficulty, interviewType, resumeContent string) ([]string, error) {
	prompt := `
你是一名专业的技术面试官，请根据以下信息生成5-8个高质量的面试问题：

职位：` + jobTitle + `
难度：` + difficulty + `
面试类型：` + interviewType + `

简历内容：
` + resumeContent + `

请确保问题涵盖相关技术栈、项目经验和技能，并且符合该职位的要求。问题应该既有理论知识，也有实际应用场景。
`

	response, err := EinoAgent.Chat(ctx, prompt)
	if err != nil {
		return nil, err
	}

	// 解析问题列表（这里需要根据实际返回格式进行解析）
	questions := []string{response}
	return questions, nil
}

// AnalyzeResume 分析简历
func AnalyzeResume(ctx context.Context, resumeContent string) (string, error) {
	prompt := `
请对以下简历进行专业分析，找出优势和可以改进的地方，并提供具体建议：

` + resumeContent + `

分析应包括：
1. 简历整体结构评估
2. 技术栈展示评估
3. 项目经验描述评估
4. 优势和亮点
5. 需要改进的地方
6. 具体改进建议
`

	return EinoAgent.Chat(ctx, prompt)
}

// EvaluateAnswer 评估面试回答
func EvaluateAnswer(ctx context.Context, question, answer, jobTitle string) (string, float64, error) {
	prompt := `
请作为技术面试官，评估以下面试回答：

问题：` + question + `
回答：` + answer + `
职位：` + jobTitle + `

请提供：
1. 详细的反馈意见
2. 评分（0-100分）
3. 改进建议

评分标准：
- 技术准确性
- 回答完整性
- 表达清晰度
- 专业深度
- 问题相关性

请以JSON格式返回：{"feedback": "...", "score": 85.5}
`

	response, err := EinoAgent.Chat(ctx, prompt)
	if err != nil {
		return "", 0, err
	}

	// 这里需要解析JSON响应获取feedback和score
	// 简化实现，实际需要使用json.Unmarshal
	return response, 85.0, nil
}