package tool

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

type LearnUrl struct {
	Classify string `json:"classify" jsonschema:"required,description=编程语言分类"`
	Url      string `json:"url" jsonschema:"description=学习网站url"`
}

func GetUrl(_ context.Context, url *LearnUrl) (string, error) {
	UrlList := []LearnUrl{
		{Classify: "go", Url: "www.wangzhongyang,com"},
		{Classify: "java", Url: "www.java,com"},
		{Classify: "python", Url: "www.python,com"},
	}
	for _, urlList := range UrlList {
		if urlList.Classify == url.Classify {
			return urlList.Url, nil
		}
	}
	return "", errors.New("里面没有这个分类")
}

func CreateTool() tool.InvokableTool {
	getUrlTool := utils.NewTool(&schema.ToolInfo{
		Name: "get_Url",
		Desc: "根据编程语言分类获取学习网站url,例如：get_learnUrl(classify='go')",
		ParamsOneOf: schema.NewParamsOneOfByParams(
			map[string]*schema.ParameterInfo{
				"classify": &schema.ParameterInfo{
					Type:     schema.String,
					Required: true,
					Desc:     "编程语言分类（只能是go、java、python中的一个）",
				},
			}),
	}, GetUrl)
	fmt.Printf("使用get_learnUrl工具获取学习网站url\n")
	return getUrlTool
}
