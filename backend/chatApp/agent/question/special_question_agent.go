package question

import (
	"ai-eino-interview-agent/chatApp/chat"
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/adk"
)

// NewSpecialQuestionAgent 专项面试智能体：仅根据提示词中的主题生成问题，不依赖简历
func NewSpecialQuestionAgent(userId uint) adk.Agent {
	ctx := context.Background()

	baseAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "SpecialQuestionAgent",
		Description: "一个只围绕单一技术主题提问的专项面试智能体",
		Instruction: `你是一位专项面试官。所有提问必须严格围绕用户消息中的“面试主题”和“难度”。

要求：
1. 不使用也不请求任何简历或候选人背景信息。
2. 只生成与单一面试主题紧密相关的主问题，可附带1-2个追问，整体由浅入深。
3. 根据用户给出的难度词（如初级/中级/高级）控制问题深度。
4. 输出必须是 JSON，且只包含面试官视角内容。

JSON格式：
{
  "questions": [{"question_text": "问题内容", "eval_dimension": "topic", "order": 1}],
  "dialogues": [{"speaker_type": "interviewer", "content": "提问内容", "display_order": 1}]
}`,

		Model:         chat.CreatOpenAiChatModel(ctx, userId),
		MaxIterations: 4,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create question generator agent: %w", err))
	}
	return baseAgent
}
