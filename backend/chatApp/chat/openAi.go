package chat

import (
	"context"
	"log"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/compose"
)

var GlobalToolsNode *compose.ToolsNode

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

	//var toolList []einoTool.BaseTool
	//
	//// 加入谷歌搜索工具
	////googleTool, err := myTool.CreatGoogleSearchTool(ctx)
	///*if err != nil {
	//	log.Fatalf("初始化谷歌搜索工具失败：%v", err)
	//}*/
	//
	////toolList = append(toolList, googleTool)
	//toolList = append(toolList, myTool.CreatePDFToTextTool())
	//toolList = append(toolList, myTool.CreateTool())
	//toolList = append(toolList, myTool.NewResumeAnalysisTool()) // 加入get_Url工具
	//
	//// 3. 收集工具元信息（和之前逻辑一样，只是遍历的切片类型变了）
	//toolInfos := make([]*schema.ToolInfo, 0, len(toolList))
	//for _, toolItem := range toolList {
	//	// 调用工具的Info()方法（tool.InvokableTool接口自带Info()）
	//	info, err := toolItem.Info(ctx)
	//	if err != nil {
	//		log.Fatalf("获取工具说明书失败：%v", err)
	//	}
	//	// 手动强化工具描述（让大模型更清楚工具用途）
	//	switch info.Name {
	//	case "google_search":
	//		info.Desc = "用于搜索网页信息，比如官方文档、最新教程、新闻等，必须传入搜索关键词query，可选指定结果数量（1-10）"
	//	case "get_Url":
	//		info.Desc = "用于获取指定编程语言的学习网站URL，支持的分类：go、java、python，必须传入classify参数"
	//	}
	//	toolInfos = append(toolInfos, info)
	//}
	//
	//// 5. 关键修正：处理 WithTools 绑定结果（必须使用绑定后的新模型）
	//toolmodel, err := chatModel.WithTools(toolInfos)
	//if err != nil {
	//	log.Fatalf("大模型绑定工具失败：%v", err)
	//}
	//
	//// 6. 修正：完善 ToolsNode 错误处理（打印日志，避免静默失败）
	//GlobalToolsNode, err = compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
	//	Tools: toolList,
	//})
	//if err != nil {
	//	log.Fatalf("创建 ToolsNode 失败：%v", err)
	//}
	//
	//log.Println("大模型绑定工具成功，ToolsNode 初始化完成！")
	return chatModel
}
