package tool

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"ai-eino-interview-agent/internal/service"

	pdfParser "github.com/cloudwego/eino-ext/components/document/parser/pdf"
	"github.com/cloudwego/eino/components/document/parser"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

// PDFToTextRequest 大模型调用工具的入参结构体（明确参数要求）
type PDFToTextRequest struct {
	FilePath string `json:"file_path" jsonschema:"required,description=本地PDF文件的绝对路径（例如：D:\\test\\document.pdf 或 /home/user/document.pdf）"`
	ToPages  bool   `json:"to_pages" jsonschema:"default=false,description=是否按页面分割文本（true=分页输出，false=合并所有页为一个文本，默认false）"`
}

// PDFToTextResult 工具返回的结构化结果（大模型可直接解析）
type PDFToTextResult struct {
	Success    bool                   `json:"success" jsonschema:"description=解析是否成功"`
	Content    string                 `json:"content,omitempty" jsonschema:"description=合并后的纯文本（ToPages=false时返回）"`
	Pages      []PDFPageText          `json:"pages,omitempty" jsonschema:"description=分页文本（ToPages=true时返回）"`
	TotalPages int                    `json:"total_pages" jsonschema:"description=PDF总页数"`
	ErrorMsg   string                 `json:"error_msg,omitempty" jsonschema:"description=错误信息（失败时返回）"`
	Meta       map[string]interface{} `json:"meta,omitempty" jsonschema:"description=元数据（方便追溯）"`
}

// PDFPageText 单页文本结构（分页模式下使用）
type PDFPageText struct {
	PageNum int    `json:"page_num" jsonschema:"description=页码（从1开始）"`
	Content string `json:"content" jsonschema:"description=单页纯文本"`
}

// ConvertPDFToText 核心逻辑：PDF转纯文本（工具执行入口）
func ConvertPDFToText(ctx context.Context, req *PDFToTextRequest) (*PDFToTextResult, error) {
	result := PDFToTextResult{
		Meta: map[string]interface{}{
			"file_path":  req.FilePath,
			"to_pages":   req.ToPages,
			"parse_time": time.Now().Format("2006-01-02 15:04:05"),
			"cache_hit":  false,
		},
	}

	if req.FilePath == "" {
		result.Success = false
		result.ErrorMsg = "参数错误：必须传入 file_path（本地PDF文件的绝对路径）"
		return &result, errors.New(result.ErrorMsg)
	}

	cache := service.NewResumeTextCache()
	cachedText, hit, err := cache.Get(ctx, req.FilePath)
	if hit && err == nil {
		log.Printf("[PDFParser] Cache HIT for %s", req.FilePath)
		result.Meta["cache_hit"] = true
		return reconstructResultFromCache(cachedText, req.ToPages, &result), nil
	}

	if err != nil {
		log.Printf("[PDFParser] Cache error: %v", err)
	} else {
		log.Printf("[PDFParser] Cache MISS for %s", req.FilePath)
	}

	file, err := os.Open(req.FilePath)
	if err != nil {
		result.Success = false
		result.ErrorMsg = fmt.Sprintf("打开PDF文件失败：%v（请检查路径是否正确、文件是否存在）", err)
		return &result, errors.New(result.ErrorMsg)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("close file failed: %v", err)
		}
	}(file)

	pdfIns, err := pdfParser.NewPDFParser(ctx, &pdfParser.Config{
		ToPages: req.ToPages,
	})
	if err != nil {
		result.Success = false
		result.ErrorMsg = fmt.Sprintf("初始化PDF解析器失败：%v", err)
		return &result, errors.New(result.ErrorMsg)
	}

	docs, err := pdfIns.Parse(ctx, file,
		parser.WithURI(req.FilePath),
		parser.WithExtraMeta(result.Meta),
	)
	if err != nil {
		result.Success = false
		result.ErrorMsg = fmt.Sprintf("PDF解析失败：%v（仅支持文本型PDF，不支持扫描件/加密PDF）", err)
		return &result, errors.New(result.ErrorMsg)
	}

	result.Success = true
	result.TotalPages = len(docs)

	var textToCache string
	if req.ToPages {
		pages := make([]PDFPageText, 0, len(docs))
		for idx, doc := range docs {
			pages = append(pages, PDFPageText{
				PageNum: idx + 1,
				Content: doc.Content,
			})
		}
		result.Pages = pages
		textToCache = extractTextForCaching(docs)
	} else {
		var contentBuilder string
		for _, doc := range docs {
			contentBuilder += doc.Content + "\n"
		}
		result.Content = contentBuilder
		textToCache = contentBuilder
	}

	if result.Success && textToCache != "" {
		if err := cache.Set(ctx, req.FilePath, textToCache); err != nil {
			log.Printf("[PDFParser] Failed to cache text: %v", err)
		}
	}

	return &result, nil
}

func reconstructResultFromCache(cachedText string, toPages bool, result *PDFToTextResult) *PDFToTextResult {
	result.Success = true

	if toPages {
		lines := strings.Split(strings.TrimSpace(cachedText), "\n")
		pages := make([]PDFPageText, 0, len(lines))
		for idx, line := range lines {
			if line != "" {
				pages = append(pages, PDFPageText{
					PageNum: idx + 1,
					Content: line,
				})
			}
		}
		result.Pages = pages
		result.TotalPages = len(pages)
	} else {
		result.Content = cachedText
		result.TotalPages = strings.Count(cachedText, "\n") + 1
	}

	return result
}

func extractTextForCaching(docs []*schema.Document) string {
	var textBuilder strings.Builder
	for idx, doc := range docs {
		textBuilder.WriteString(doc.Content)
		if idx < len(docs)-1 {
			textBuilder.WriteString("\n")
		}
	}
	return textBuilder.String()
}

// CreatePDFToTextTool 创建工具实例（供Eino框架注册，大模型识别）
func CreatePDFToTextTool() tool.InvokableTool {
	// 工具元信息：大模型识别的关键（名称、描述、参数定义）
	pdfTool, err := utils.InferTool("pdf_to_text", "将本地PDF文件转换为纯文本，仅支持文本型PDF（可复制文字），不支持扫描件、加密PDF。需传入本地PDF的绝对路径，可选择按页面分割或合并所有页。", ConvertPDFToText)
	if err != nil {
		log.Fatalf("infer tool failed: %v", err)
	}
	fmt.Println("✅ PDF转纯文本工具初始化完成（大模型可调用）")
	return pdfTool
}
