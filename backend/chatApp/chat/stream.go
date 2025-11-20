package chat

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

func ReportSteam(sr *schema.StreamReader[*schema.Message]) {
	defer sr.Close()

	//i := 0

	for {
		message, err := sr.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatalf("recv message failed: %v", err)
		}
		content := message.Content

		fmt.Printf(content)
		//i++
	}
}

func Stream(ctx context.Context, llm model.ToolCallingChatModel, in []*schema.Message) *schema.StreamReader[*schema.Message] {

	// 执行链
	result, err := llm.Stream(ctx, in)
	if err != nil {
		log.Fatalf("llm generate failed: %v", err)
	}
	return result
}
