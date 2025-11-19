package question

import (
	"ai-eino-interview-agent/chatApp/chat"
	tool2 "ai-eino-interview-agent/chatApp/tool"
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/adk"
	componenttool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

// NewQuestionAgent 基于现有模板构建的模拟面试智能体
func NewQuestionAgent() adk.Agent {
	ctx := context.Background()
	baseAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "QuestionAgent",
		Description: "一个专业面试提问的智能体",
		Instruction: `你是一个专业的面试官。根据提供的简历或信息，快速生成一个面试问题和对话。

任务：
1. 如果提供了PDF路径，使用pdf_to_text工具解析简历
2. 根据背景生成1个问题和2-3条对话
3. 只返回JSON，不要返回其他文本

必须返回的JSON格式：
{
  "questions": [{
    "question_text": "问题内容",
    "eval_dimension": "纬度",
    "order": 1
  }],
  "dialogues": [
    {"speaker_type": "interviewer", "content": "提问", "display_order": 1},
    {"speaker_type": "candidate", "content": "回答", "display_order": 2},
    {"speaker_type": "interviewer", "content": "追问", "display_order": 3}
  ]
}`,

		Model: chat.CreatOpenAiChatModel(ctx),
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []componenttool.BaseTool{
					tool2.CreatePDFToTextTool(),
				},
			},
		},
		MaxIterations: 20,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create question generator agent: %w", err))
	}
	return baseAgent
}
