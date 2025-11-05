package chat

import (
	myTool "ai-eino-agent/chatApp/tool"
	"context"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	einoTool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"log"
)

func CreatOpenAiChatModel(ctx context.Context) model.ToolCallingChatModel {
	//你的api Key
	key := "c1c8f7ce-266f-4af5-a832-9a8457f36e74"
	//模型名称
	modelName := "doubao-seed-1-6-251015"
	//api url
	url := "https://ark.cn-beijing.volces.com/api/v3"

	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey:  key,
		Model:   modelName,
		BaseURL: url,
	})
	if err != nil {
		log.Fatalf("create openai chat model failed: %v", err)
	}
	var toolList []einoTool.InvokableTool            // 用字节tool包的InvokableTool
	toolList = append(toolList, myTool.CreateTool()) // 加入get_Url工具

	// 3. 收集工具元信息（和之前逻辑一样，只是遍历的切片类型变了）
	toolInfos := make([]*schema.ToolInfo, 0, len(toolList))
	for _, toolItem := range toolList {
		// 调用工具的Info()方法（tool.InvokableTool接口自带Info()）
		info, err := toolItem.Info(ctx)
		if err != nil {
			log.Fatalf("获取工具说明书失败：%v", err)
		}
		toolInfos = append(toolInfos, info)
	}

	return chatModel
}
