package feishu

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkdocx "github.com/larksuite/oapi-sdk-go/v3/service/docx/v1"
)

// SDK 使用文档：https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/server-side-sdk/golang-sdk-guide/preparations
// 复制该 Demo 后, 需要将 "YOUR_APP_ID", "YOUR_APP_SECRET" 替换为自己应用的 APP_ID, APP_SECRET.
// 以下示例代码默认根据文档示例值填充，如果存在代码问题，请在 API 调试台填上相关必要参数后再复制代码使用

// 飞书文档块的顶层结构（API 响应的 data.items 数组元素）
type FeishuBlock struct {
	BlockID   string        `json:"block_id"`           // 块ID
	BlockType int           `json:"block_type"`         // 块类型
	ParentID  string        `json:"parent_id"`          // 父块ID（用于处理嵌套，如列表）
	Page      *PageBlock    `json:"page,omitempty"`     // block_type=1：文档标题块
	Text      *TextBlock    `json:"text,omitempty"`     // block_type=2：普通文本块
	Heading1  *HeadingBlock `json:"heading1,omitempty"` // block_type=3：一级标题
	Heading2  *HeadingBlock `json:"heading2,omitempty"` // block_type=4：二级标题
	Heading3  *HeadingBlock `json:"heading3,omitempty"` // block_type=5：三级标题
	Heading4  *HeadingBlock `json:"heading4,omitempty"` // block_type=6：四级标题
	Heading5  *HeadingBlock `json:"heading5,omitempty"` // block_type=7：五级标题
	Heading6  *HeadingBlock `json:"heading6,omitempty"` // block_type=8：六级标题
	Heading7  *HeadingBlock `json:"heading7,omitempty"` // block_type=9：七级标题
	Heading8  *HeadingBlock `json:"heading8,omitempty"` // block_type=10：八级标题
	Heading9  *HeadingBlock `json:"heading9,omitempty"` // block_type=11：九级标题
	Bullet    *BulletBlock  `json:"bullet,omitempty"`   // block_type=12：无序列表
	Ordered   *OrderedBlock `json:"ordered,omitempty"`  // block_type=13：有序列表
	Code      *CodeBlock    `json:"code,omitempty"`     // block_type=14：代码块
	Quote     *QuoteBlock   `json:"quote,omitempty"`    // block_type=15：引用块
	Callout   *CalloutBlock `json:"callout,omitempty"`  // block_type=16：标注块
	Divider   *DividerBlock `json:"divider,omitempty"`  // block_type=17：分割线
	Image     *ImageBlock   `json:"image,omitempty"`    // block_type=18：图片
	Table     *TableBlock   `json:"table,omitempty"`    // block_type=19：表格
}

// PageBlock：文档标题块（block_type=1）
type PageBlock struct {
	Elements []TextElement `json:"elements"`
	Style    struct{}      `json:"style"` // 暂无需处理样式
}

// TextBlock：普通文本块（block_type=2）
type TextBlock struct {
	Elements []TextElement `json:"elements"`
	Style    struct{}      `json:"style"`
}

// HeadingBlock：标题块（heading1-heading9 通用）
type HeadingBlock struct {
	Elements []TextElement `json:"elements"`
	Style    struct{}      `json:"style"`
}

// BulletBlock：无序列表块
type BulletBlock struct {
	Elements []TextElement `json:"elements"`
	Style    struct{}      `json:"style"`
}

// OrderedBlock：有序列表块
type OrderedBlock struct {
	Elements []TextElement `json:"elements"`
	Style    struct{}      `json:"style"`
}

// CodeBlock：代码块
type CodeBlock struct {
	Elements []TextElement `json:"elements"`
	Language int           `json:"language"` // 语言类型
	Wrap     bool          `json:"wrap"`
	Style    struct{}      `json:"style"`
}

// QuoteBlock：引用块
type QuoteBlock struct {
	Elements []TextElement `json:"elements"`
	Style    struct{}      `json:"style"`
}

// CalloutBlock：标注块
type CalloutBlock struct {
	Elements []TextElement `json:"elements"`
	Style    struct{}      `json:"style"`
}

// DividerBlock：分割线块
type DividerBlock struct {
	Style struct{} `json:"style"`
}

// ImageBlock：图片块
type ImageBlock struct {
	Token  string   `json:"token"`
	Width  int      `json:"width"`
	Height int      `json:"height"`
	Style  struct{} `json:"style"`
}

// TableBlock：表格块
type TableBlock struct {
	Rows    int      `json:"rows"`
	Columns int      `json:"columns"`
	Style   struct{} `json:"style"`
}

// TextElement：文本元素（text_run/mention_user 等）
type TextElement struct {
	TextRun     *TextRun     `json:"text_run,omitempty"`     // 普通文本
	MentionUser *MentionUser `json:"mention_user,omitempty"` // 提及用户
}

// TextRun：普通文本内容
type TextRun struct {
	Content          string            `json:"content"` // 文本内容
	TextElementStyle *TextElementStyle `json:"text_element_style,omitempty"`
}

// MentionUser：提及用户
type MentionUser struct {
	UserID           string            `json:"user_id"` // 用户ID（可通过飞书API查询用户名，这里先占位）
	TextElementStyle *TextElementStyle `json:"text_element_style,omitempty"`
}

// TextElementStyle：文本样式（粗体/斜体等，暂无需处理，保留结构）
type TextElementStyle struct {
	Bold          bool `json:"bold"`
	InlineCode    bool `json:"inline_code"`
	Italic        bool `json:"italic"`
	Strikethrough bool `json:"strikethrough"`
	Underline     bool `json:"underline"`
}

// 转换飞书块列表为 Markdown 字符串
func BlocksToMarkdown(blocks []FeishuBlock) string {
	var mdBuilder strings.Builder

	// 递归解析单个块
	var parseBlock func(block FeishuBlock)
	parseBlock = func(block FeishuBlock) {
		switch block.BlockType {
		case 1: // 文档标题（page 块）
			if block.Page != nil {
				content := parseTextElements(block.Page.Elements)
				if content != "" {
					mdBuilder.WriteString(fmt.Sprintf("# %s\n\n", content))
				}
			}
		case 2: // 普通文本块（text 块）
			if block.Text != nil {
				content := parseTextElements(block.Text.Elements)
				if content != "" {
					mdBuilder.WriteString(fmt.Sprintf("%s\n\n", content))
				}
			}
		case 3: // 一级标题（heading1 块）
			if block.Heading1 != nil {
				content := parseTextElements(block.Heading1.Elements)
				if content != "" {
					mdBuilder.WriteString(fmt.Sprintf("# %s\n\n", content))
				}
			}
		case 4: // 二级标题（heading2 块）
			if block.Heading2 != nil {
				content := parseTextElements(block.Heading2.Elements)
				if content != "" {
					mdBuilder.WriteString(fmt.Sprintf("## %s\n\n", content))
				}
			}
		case 5: // 三级标题（heading3 块）
			if block.Heading3 != nil {
				content := parseTextElements(block.Heading3.Elements)
				if content != "" {
					mdBuilder.WriteString(fmt.Sprintf("### %s\n\n", content))
				}
			}
		case 6: // 四级标题（heading4 块）
			if block.Heading4 != nil {
				content := parseTextElements(block.Heading4.Elements)
				if content != "" {
					mdBuilder.WriteString(fmt.Sprintf("#### %s\n\n", content))
				}
			}
		case 7: // 五级标题（heading5 块）
			if block.Heading5 != nil {
				content := parseTextElements(block.Heading5.Elements)
				if content != "" {
					mdBuilder.WriteString(fmt.Sprintf("##### %s\n\n", content))
				}
			}
		case 8: // 六级标题（heading6 块）
			if block.Heading6 != nil {
				content := parseTextElements(block.Heading6.Elements)
				if content != "" {
					mdBuilder.WriteString(fmt.Sprintf("###### %s\n\n", content))
				}
			}
		case 9: // 七级标题（heading7 块）
			if block.Heading7 != nil {
				content := parseTextElements(block.Heading7.Elements)
				if content != "" {
					mdBuilder.WriteString(fmt.Sprintf("####### %s\n\n", content))
				}
			}
		case 10: // 八级标题（heading8 块）
			if block.Heading8 != nil {
				content := parseTextElements(block.Heading8.Elements)
				if content != "" {
					mdBuilder.WriteString(fmt.Sprintf("######## %s\n\n", content))
				}
			}
		case 11: // 九级标题（heading9 块）
			if block.Heading9 != nil {
				content := parseTextElements(block.Heading9.Elements)
				if content != "" {
					mdBuilder.WriteString(fmt.Sprintf("######### %s\n\n", content))
				}
			}
		case 12: // 无序列表（bullet 块）
			if block.Bullet != nil {
				content := parseTextElements(block.Bullet.Elements)
				if content != "" {
					mdBuilder.WriteString(fmt.Sprintf("- %s\n", content))
				}
			}
		case 13: // 有序列表（ordered 块）
			if block.Ordered != nil {
				content := parseTextElements(block.Ordered.Elements)
				if content != "" {
					mdBuilder.WriteString(fmt.Sprintf("1. %s\n", content))
				}
			}
		case 14: // 代码块（code 块）
			if block.Code != nil {
				content := parseTextElements(block.Code.Elements)
				language := getCodeLanguage(block.Code.Language)
				if content != "" {
					mdBuilder.WriteString(fmt.Sprintf("```%s\n%s\n```\n\n", language, content))
				}
			}
		case 15: // 引用块（quote 块）
			if block.Quote != nil {
				content := parseTextElements(block.Quote.Elements)
				if content != "" {
					mdBuilder.WriteString(fmt.Sprintf("> %s\n\n", strings.ReplaceAll(content, "\n", "\n> ")))
				}
			}
		case 16: // 标注块（callout 块）
			if block.Callout != nil {
				content := parseTextElements(block.Callout.Elements)
				if content != "" {
					mdBuilder.WriteString(fmt.Sprintf("> 💡 %s\n\n", content))
				}
			}
		case 17: // 分割线（divider 块）
			mdBuilder.WriteString("---\n\n")
		case 18: // 图片（image 块）
			if block.Image != nil {
				mdBuilder.WriteString(fmt.Sprintf("![image](%s)\n\n", block.Image.Token))
			}
		case 19: // 表格（table 块）
			if block.Table != nil {
				mdBuilder.WriteString(fmt.Sprintf("<!-- 表格: %dx%d -->\n\n", block.Table.Rows, block.Table.Columns))
			}
		default:
			// 未知块类型，尝试从各种可能的字段提取文本内容
			content := extractTextFromBlock(block)
			if content != "" {
				mdBuilder.WriteString(fmt.Sprintf("%s\n\n", content))
			} else {
				// 如果都没有，至少输出块类型信息
				mdBuilder.WriteString(fmt.Sprintf("<!-- 未处理的块类型: %d, block_id: %s -->\n\n", block.BlockType, block.BlockID))
			}
		}
	}

	// 遍历所有块解析
	for _, block := range blocks {
		parseBlock(block)
	}

	return mdBuilder.String()
}

// 从块中提取文本内容（尝试所有可能的字段）
func extractTextFromBlock(block FeishuBlock) string {
	// 尝试从 Text 字段提取
	if block.Text != nil {
		return parseTextElements(block.Text.Elements)
	}
	// 尝试从各种标题字段提取
	if block.Heading1 != nil {
		return parseTextElements(block.Heading1.Elements)
	}
	if block.Heading2 != nil {
		return parseTextElements(block.Heading2.Elements)
	}
	if block.Heading3 != nil {
		return parseTextElements(block.Heading3.Elements)
	}
	if block.Heading4 != nil {
		return parseTextElements(block.Heading4.Elements)
	}
	if block.Heading5 != nil {
		return parseTextElements(block.Heading5.Elements)
	}
	if block.Heading6 != nil {
		return parseTextElements(block.Heading6.Elements)
	}
	if block.Heading7 != nil {
		return parseTextElements(block.Heading7.Elements)
	}
	if block.Heading8 != nil {
		return parseTextElements(block.Heading8.Elements)
	}
	if block.Heading9 != nil {
		return parseTextElements(block.Heading9.Elements)
	}
	// 尝试从列表字段提取
	if block.Bullet != nil {
		return parseTextElements(block.Bullet.Elements)
	}
	if block.Ordered != nil {
		return parseTextElements(block.Ordered.Elements)
	}
	// 尝试从代码块提取
	if block.Code != nil {
		return parseTextElements(block.Code.Elements)
	}
	// 尝试从引用块提取
	if block.Quote != nil {
		return parseTextElements(block.Quote.Elements)
	}
	// 尝试从标注块提取
	if block.Callout != nil {
		return parseTextElements(block.Callout.Elements)
	}
	// 尝试从 Page 块提取
	if block.Page != nil {
		return parseTextElements(block.Page.Elements)
	}
	return ""
}

// 获取代码语言名称
func getCodeLanguage(langCode int) string {
	langMap := map[int]string{
		1:  "plaintext",
		2:  "abap",
		3:  "ada",
		4:  "arduino",
		5:  "autoit",
		6:  "c",
		7:  "clojure",
		8:  "coffeescript",
		9:  "cpp",
		10: "csharp",
		11: "css",
		12: "dart",
		13: "delphi",
		14: "dockerfile",
		15: "erlang",
		16: "fortran",
		17: "foxpro",
		18: "go",
		19: "groovy",
		20: "haskell",
		21: "html",
		22: "java",
		23: "javascript",
		24: "json",
		25: "julia",
		26: "kotlin",
		27: "latex",
		28: "lisp",
		29: "lua",
		30: "matlab",
		31: "nginx",
		32: "objectivec",
		33: "perl",
		34: "php",
		35: "powershell",
		36: "prolog",
		37: "protobuf",
		38: "python",
		39: "r",
		40: "ruby",
		41: "rust",
		42: "scala",
		43: "scheme",
		44: "shell",
		45: "sql",
		46: "swift",
		47: "thrift",
		48: "typescript",
		49: "vb",
		50: "verilog",
		51: "vhdl",
		52: "xml",
		53: "yaml",
	}
	if lang, ok := langMap[langCode]; ok {
		return lang
	}
	return "plaintext"
}

// 解析文本元素（TextElement 数组 → 字符串）
// 处理普通文本（text_run）和提及用户（mention_user）
func parseTextElements(elements []TextElement) string {
	var contentBuilder strings.Builder
	for _, elem := range elements {
		if elem.TextRun != nil {
			// 普通文本：直接拼接内容
			content := elem.TextRun.Content
			// 处理样式（粗体、斜体等）
			if elem.TextRun.TextElementStyle != nil {
				style := elem.TextRun.TextElementStyle
				if style.Bold {
					content = fmt.Sprintf("**%s**", content)
				}
				if style.Italic {
					content = fmt.Sprintf("*%s*", content)
				}
				if style.InlineCode {
					content = fmt.Sprintf("`%s`", content)
				}
				if style.Strikethrough {
					content = fmt.Sprintf("~~%s~~", content)
				}
			}
			contentBuilder.WriteString(content)
		} else if elem.MentionUser != nil {
			// 提及用户：转为 @用户名（这里 user_id 可替换为真实用户名，需调用飞书用户API）
			// 临时方案：用 user_id 后缀占位，实际可通过飞书 API 查询用户名
			userID := elem.MentionUser.UserID
			if len(userID) > 6 {
				userName := fmt.Sprintf("用户_%s", userID[len(userID)-6:])
				contentBuilder.WriteString(fmt.Sprintf("@%s", userName))
			} else {
				contentBuilder.WriteString(fmt.Sprintf("@%s", userID))
			}
		}
	}
	return contentBuilder.String()
}

func Test() string {
	// 创建 Client
	client := lark.NewClient("cli_a9afad5abfb85bc0", "RDIAVuYOukhGNdZcn1zO9dLJS8up7rYL")
	// 创建请求对象
	req := larkdocx.NewListDocumentBlockReqBuilder().
		DocumentId(`SQbFdHo6Wo9fOixcLVecp1yTnLh`).
		PageSize(500).
		DocumentRevisionId(-1).
		Build()

	// 发起请求
	resp, err := client.Docx.V1.DocumentBlock.List(context.Background(), req, larkcore.WithUserAccessToken("u-fSJtoUlf9cuaVCduZCHUt_1g5awBggiNigEamRuw08au"))

	// 处理错误
	if err != nil {
		fmt.Printf("请求失败：%v\n", err)
		return ""
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Printf("logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
		return ""
	}

	// 将响应转换为 JSON 以便解析
	respJSON, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("序列化响应失败：%v\n", err)
		return ""
	}

	// 解析 JSON 响应
	var apiResponse struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Items []FeishuBlock `json:"items"`
		} `json:"data"`
	}

	// 先解析为通用的 map 结构，以便调试
	var rawResponse map[string]interface{}
	if err := json.Unmarshal(respJSON, &rawResponse); err == nil {
		// 检查是否有 data.items
		if data, ok := rawResponse["data"].(map[string]interface{}); ok {
			if items, ok := data["items"].([]interface{}); ok {
				fmt.Printf("调试：发现 %d 个块，开始分析块类型...\n", len(items))
				blockTypeCount := make(map[int]int)
				for i, item := range items {
					if itemMap, ok := item.(map[string]interface{}); ok {
						if blockType, ok := itemMap["block_type"].(float64); ok {
							bt := int(blockType)
							blockTypeCount[bt]++
							// 只显示前10个块的详细信息
							if i < 10 {
								fmt.Printf("  块 %d: block_type=%d, block_id=%v\n", i+1, bt, itemMap["block_id"])
							}
						}
					}
				}
				fmt.Println("块类型统计：")
				for bt, count := range blockTypeCount {
					fmt.Printf("  block_type %d: %d 个\n", bt, count)
				}
			}
		}
	}

	err = json.Unmarshal(respJSON, &apiResponse)
	if err != nil {
		fmt.Printf("解析响应失败：%v\n", err)
		fmt.Println("尝试查看原始响应 JSON（前2000字符）：")
		jsonStr := string(respJSON)
		if len(jsonStr) > 2000 {
			fmt.Println(jsonStr[:2000] + "...")
		} else {
			fmt.Println(jsonStr)
		}
		return ""
	}

	if apiResponse.Code != 0 {
		fmt.Printf("飞书API返回错误：code=%d, msg=%s\n", apiResponse.Code, apiResponse.Msg)
		return ""
	}

	if len(apiResponse.Data.Items) == 0 {
		fmt.Println("警告：未获取到任何文档块，请检查文档ID和权限")
		return ""
	}

	fmt.Printf("成功获取 %d 个文档块，开始转换为 Markdown...\n", len(apiResponse.Data.Items))

	// 转换为 Markdown
	mdContent := BlocksToMarkdown(apiResponse.Data.Items)

	// // 输出为 Markdown 文件
	// outputPath := "飞书文档转换结果.md"
	// err = os.WriteFile(outputPath, []byte(mdContent), 0644)
	// if err != nil {
	// 	fmt.Printf("保存 Markdown 文件失败：%v\n", err)
	// 	return ""
	// }

	// fmt.Printf("转换成功！Markdown 文件已保存至：%s\n", outputPath)
	fmt.Println("转换结果预览：")
	fmt.Println("----------------------------------------")
	fmt.Println(mdContent)
	return mdContent

}
