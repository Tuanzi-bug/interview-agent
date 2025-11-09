package chat

import (
	"context"
	"log"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
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
		"role":     "你是一个智能助手，必须通过工具完成以下相关需求，流程如下：\n1. 当用户的问题涉及「PDF文件内容」（比如解析PDF、查询PDF中的信息、总结PDF要点）时：\n   a. 首先检查用户是否提供了「本地PDF绝对路径」（如D:\\test\\document.pdf）；\n   b. 若未提供，追问用户获取PDF的绝对路径（需明确是本地文件路径，不支持网络PDF）；\n   c. 调用「pdf_to_text」工具，传入文件路径（to_pages默认false即可，合并所有页更易理解）；\n   d. 等待工具返回PDF纯文本后，基于文本内容回答用户的具体问题，不要遗漏关键信息。\n2. 工具调用规则：\n   a. 仅支持调用已提供的工具（pdf_to_text、get_Url、google_search），不擅自创造工具；\n   b. 调用工具时参数必须完整（pdf_to_text的file_path是必填项）；\n   c. 工具返回结果后，必须基于结果回答，不能直接返回工具的JSON原始数据，要整理成自然语言。\n3. 回答要求：\n   a. 基于PDF内容回答，不编造信息；\n   b. 若PDF中没有用户问的内容，明确告知“PDF中未提及相关信息”；\n   c. 复杂问题可分点说明，引用PDF中的关键句子（标注“来自PDF内容”）。",
		"style":    "温和且专业",
		"question": "帮我解析D:\\go-eino-interview-agent-001\\backend\\chatApp\\agent\\akf.pdf这个PDF文件并回答里面的问题",
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
