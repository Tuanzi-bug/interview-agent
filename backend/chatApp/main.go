package main

import (
	"ai-eino-agent/chatApp/chat"
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()

	//使用message模版
	fmt.Printf("===create messages===\n")
	message := chat.MessagesTemplate()
	//fmt.Printf("messages: %+v\n\n", chat.MessagesTemplate())

	//创建llm
	fmt.Printf("===create llm===\n")
	model := chat.CreatOpenAiChatModel(ctx)
	//log.Printf("create llm success\n\n")

	//使用llm生成Steam流式回复
	fmt.Printf("===llm stream ===\n")
	streamResult := chat.Stream(ctx, model, message)
	chat.ReportSteam(streamResult)
}
