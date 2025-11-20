package milvus

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// MarkdownImporter Markdown 文档导入器
// 用于将 Markdown 文档解析、分割并存储到 Milvus
type MarkdownImporter struct {
	manager *MilvusManager
}

// NewMarkdownImporter 创建新的 Markdown 导入器
func NewMarkdownImporter(manager *MilvusManager) (*MarkdownImporter, error) {
	if manager == nil {
		return nil, fmt.Errorf("milvus manager.go is nil")
	}
	return &MarkdownImporter{
		manager: manager,
	}, nil
}

// ImportOptions 导入选项
type ImportOptions struct {
	// 文件路径
	Path string
	// 语言类型
	Language DocumentLanguage
	// 文档分类
	Category DocumentCategory
	// 来源
	Source string
	// 是否递归处理子目录
	Recursive bool
	// 文件扩展名过滤（默认：.md, .markdown）
	Extensions []string
	// 是否跳过隐藏文件
	SkipHidden bool
	// 最大文件大小（字节，0 表示无限制）
	MaxFileSize int64
}

// DefaultImportOptions 返回默认的导入选项
func DefaultImportOptions() *ImportOptions {
	return &ImportOptions{
		Recursive:   true,
		Extensions:  []string{".md", ".markdown"},
		SkipHidden:  true,
		MaxFileSize: 10 * 1024 * 1024, // 10MB
	}
}

// ImportResult 导入结果
type ImportResult struct {
	// 处理的文件总数
	TotalFiles int

	// 成功处理的文件数
	SuccessFiles int

	// 失败的文件数
	FailedFiles int

	// 总文档块数
	TotalChunks int

	// 存储的文档 ID 列表
	DocumentIDs []string

	// 错误信息
	Errors []error
}

// ImportFile 导入单个 Markdown 文件
func (mi *MarkdownImporter) ImportFile(ctx context.Context, filePath string, opts *ImportOptions) (*ImportResult, error) {
	if opts == nil {
		opts = DefaultImportOptions()
	}

	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// 检查文件大小
	if opts.MaxFileSize > 0 && int64(len(content)) > opts.MaxFileSize {
		return nil, fmt.Errorf("file %s exceeds max size %d bytes", filePath, opts.MaxFileSize)
	}

	// 不再进行自动推断：如果未指定，保持为空
	language := opts.Language
	category := opts.Category

	// 创建元数据
	metadata := NewDocumentMetadata(filePath, language, category)
	if opts.Source != "" {
		metadata.Source = opts.Source
	}

	// 尝试从 Markdown 内容提取标题
	title := extractTitleFromMarkdown(string(content))
	if title != "" {
		metadata.Title = title
	}

	// 使用 SplitMarkdown 分割文档
	chunks, err := mi.manager.SplitterService.SplitMarkdown(ctx, string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to split markdown file %s: %w", filePath, err)
	}

	// 为每个块添加元数据和唯一 ID
	enrichedChunks := EnrichDocumentsWithMetadata(chunks, metadata)
	for i, chunk := range enrichedChunks {
		if chunk.ID == "" {
			// 生成唯一 ID：语言类型 + 文档分类 + 时间戳 + 块索引
			chunk.ID = generateChunkID(metadata.Language, metadata.Category, i)
		}
	}

	// 存储到 Milvus
	docIDs, err := mi.manager.IndexerService.Store(ctx, enrichedChunks)
	if err != nil {
		return nil, fmt.Errorf("failed to store documents to Milvus: %w", err)
	}

	return &ImportResult{
		TotalFiles:   1,
		SuccessFiles: 1,
		FailedFiles:  0,
		TotalChunks:  len(enrichedChunks),
		DocumentIDs:  docIDs,
		Errors:       nil,
	}, nil
}

// ImportDirectory 导入目录中的所有 Markdown 文件
func (mi *MarkdownImporter) ImportDirectory(ctx context.Context, dirPath string, opts *ImportOptions) (*ImportResult, error) {
	if opts == nil {
		opts = DefaultImportOptions()
	}

	result := &ImportResult{
		DocumentIDs: make([]string, 0),
		Errors:      make([]error, 0),
	}

	// 遍历目录
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("error accessing path %s: %w", path, err))
			return nil // 继续处理其他文件
		}

		// 跳过目录
		if d.IsDir() {
			// 如果不递归，跳过子目录
			if !opts.Recursive && path != dirPath {
				return fs.SkipDir
			}
			// 跳过隐藏目录
			if opts.SkipHidden && strings.HasPrefix(d.Name(), ".") {
				return fs.SkipDir
			}
			return nil
		}

		// 检查文件扩展名
		if !hasValidExtension(path, opts.Extensions) {
			return nil
		}

		// 跳过隐藏文件
		if opts.SkipHidden && strings.HasPrefix(d.Name(), ".") {
			return nil
		}

		result.TotalFiles++

		// 导入文件
		fileResult, err := mi.ImportFile(ctx, path, opts)
		if err != nil {
			result.FailedFiles++
			result.Errors = append(result.Errors, fmt.Errorf("failed to import file %s: %w", path, err))
			return nil // 继续处理其他文件
		}

		result.SuccessFiles++
		result.TotalChunks += fileResult.TotalChunks
		result.DocumentIDs = append(result.DocumentIDs, fileResult.DocumentIDs...)

		return nil
	})

	if err != nil {
		result.Errors = append(result.Errors, fmt.Errorf("error walking directory %s: %w", dirPath, err))
	}

	return result, nil
}

// Import 导入文件或目录（自动判断）
func (mi *MarkdownImporter) Import(ctx context.Context, path string, opts *ImportOptions) (*ImportResult, error) {
	if opts == nil {
		opts = DefaultImportOptions()
	}

	// 检查路径是否存在
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("path %s does not exist: %w", path, err)
	}

	if info.IsDir() {
		return mi.ImportDirectory(ctx, path, opts)
	}

	return mi.ImportFile(ctx, path, opts)
}

// BatchImport 批量导入多个文件或目录
func (mi *MarkdownImporter) BatchImport(ctx context.Context, paths []string, opts *ImportOptions) (*ImportResult, error) {
	if opts == nil {
		opts = DefaultImportOptions()
	}

	result := &ImportResult{
		DocumentIDs: make([]string, 0),
		Errors:      make([]error, 0),
	}

	for _, path := range paths {
		pathResult, err := mi.Import(ctx, path, opts)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("failed to import %s: %w", path, err))
			result.FailedFiles += pathResult.FailedFiles
			continue
		}

		result.TotalFiles += pathResult.TotalFiles
		result.SuccessFiles += pathResult.SuccessFiles
		result.FailedFiles += pathResult.FailedFiles
		result.TotalChunks += pathResult.TotalChunks
		result.DocumentIDs = append(result.DocumentIDs, pathResult.DocumentIDs...)
		result.Errors = append(result.Errors, pathResult.Errors...)
	}

	return result, nil
}

// GetManager 获取 Milvus 管理器
func (mi *MarkdownImporter) GetManager() *MilvusManager {
	return mi.manager
}

// extractTitleFromMarkdown 从 Markdown 内容中提取标题
// 优先提取第一个一级标题，如果没有则提取第一个二级标题
func extractTitleFromMarkdown(content string) string {
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// 检查一级标题
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
	}

	// 如果没有一级标题，查找二级标题
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "## ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "## "))
		}
	}

	return ""
}

// hasValidExtension 检查文件是否有有效的扩展名
func hasValidExtension(filePath string, extensions []string) bool {
	if len(extensions) == 0 {
		return true // 如果没有指定扩展名，接受所有文件
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	for _, validExt := range extensions {
		if ext == strings.ToLower(validExt) {
			return true
		}
	}

	return false
}

// generateChunkID 生成文档块的唯一 ID：语言类型_文档分类_时间戳_块索引
func generateChunkID(language DocumentLanguage, category DocumentCategory, chunkIndex int) string {
	timestamp := time.Now().Unix()
	// 如果未提供语言和分类，则直接返回时间戳
	if string(language) == "" && string(category) == "" {
		return fmt.Sprintf("%d", timestamp)
	}
	return fmt.Sprintf("%s_%s_%d_%d", string(language), string(category), timestamp, chunkIndex)
}
