package chat

import (
	"context"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"log"
)

func creatTemplate() prompt.ChatTemplate {
	//创建模版，使用Fstring
	return prompt.FromMessages(schema.FString,

		//系统提示词，定义角色和语气
		schema.SystemMessage("你是一个{role},你需要用{style}的语气回答问题，你的目标是解答程序员的面试问题，专注帮助程序员提升面试表现"),
		//插入需要历史消息
		schema.MessagesPlaceholder("chat_history", true),
		//用户消息模版
		schema.UserMessage("问题：{question}"),
	)
}

func MessagesTemplate() []*schema.Message {

	template := creatTemplate()

	messages, err := template.Format(context.Background(), map[string]any{
		"role":     "经验丰富的大厂开发面试专家，专注于帮助程序员解答面试问题",
		"style":    "温和且专业",
		"question": "你好",
		"chat_history": []*schema.Message{
			schema.UserMessage("你好"),
			schema.AssistantMessage("嘿！我是你的程序员面试！记住，每个优秀的程序员都是从 Debug 中成长起来的。有什么我可以帮你的吗？", nil),
			schema.UserMessage("我觉得自己写的代码太烂了"),
			schema.AssistantMessage("每个程序员都经历过这个阶段！重要的是你在不断学习和进步。让我们一起看看代码，我相信通过重构和优化，它会变得更好。记住，Rome wasn't built in a day，代码质量是通过持续改进来提升的。", nil),
		},
	})
	if err != nil {
		log.Fatalf("format template failed: %v", err)
	}
	return messages
}
