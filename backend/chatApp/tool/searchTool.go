package tool

import (
	"ai-eino-interview-agent/chatApp/config"
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/tool/googlesearch"
	"github.com/cloudwego/eino/components/tool"
	"log"
)

func CreatGoogleSearchTool(ctx context.Context) (tool.InvokableTool, error) {
	loadConfig, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("加载配置失败：%v", err)
	}

	googleAPIKey := loadConfig.Google.APIKey
	googleSearchEngineID := loadConfig.Google.SearchEngineID

	searchConfig := &googlesearch.Config{
		APIKey:         googleAPIKey,
		SearchEngineID: googleSearchEngineID,
		Lang:           "zh-CN",
		Num:            5,
	}

	searchTool, err := googlesearch.NewTool(ctx, searchConfig)
	if err != nil {
		return nil, fmt.Errorf("创建谷歌搜索工具失败：%v", err)
	}

	log.Println("谷歌搜索工具（google_search）创建成功！")
	return searchTool, nil

}
