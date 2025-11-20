package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/cloudwego/eino/schema"

	"ai-eino-interview-agent/internal/config"
	"ai-eino-interview-agent/internal/eino/milvus"
)

func main() {
	cmd := flag.String("cmd", "help", "Command to run: retrieve|import-<source>|batch-<source>|import-<source>-batch")
	text := flag.String("text", "", "Text content for embed/index (optional)")
	filePath := flag.String("file", "", "File path for split/index/import-file")
	dirPath := flag.String("dir", "", "Directory path for import-dir")
	paths := flag.String("paths", "", "Comma-separated file/dir paths for batch-import")
	query := flag.String("query", "", "Query text for retrieve")
	langFlag := flag.String("language", "", "Language for imports: golang|java|middleware (empty=auto)")
	catFlag := flag.String("category", "", "Category for imports: basic|specialized|comprehensive (empty=auto)")
	topK := flag.Int("topk", 5, "TopK for retrieval")
	source := flag.String("source", "cli", "Source metadata for imports")
	dsName := flag.String("ds", "", "Data source name under backend/internal/eino/milvus/data for ds-* commands")
	recursive := flag.Bool("recursive", true, "Recursive when importing directory")
	flag.Parse()

	ctx := context.Background()
	cfg := getTestConfig()
	// Resolve data directory relative to this source file so it works from any CWD:
	// <repo>/backend/internal/eino/milvus/cmd/milvusctl/main.go
	// data directory is at: <repo>/backend/internal/eino/milvus/data
	_, thisFile, _, _ := runtime.Caller(0)
	milvusRootDir := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", ".."))
	dataBaseDir := filepath.Join(milvusRootDir, "data")

	switch *cmd {
	case "init":
		manager, err := milvus.InitMilvusManager(ctx, cfg)
		exitIfErr(err)
		defer manager.Close()
		fmt.Println("Milvus manager.go initialized")
	case "health":
		manager, err := milvus.InitMilvusManager(ctx, cfg)
		exitIfErr(err)
		defer manager.Close()
		err = manager.HealthCheck(ctx)
		exitIfErr(err)
		fmt.Println("Milvus health check passed")
	case "embed":
		if *text == "" {
			exitIfErr(fmt.Errorf("-text is required for embed"))
		}
		manager, err := milvus.InitMilvusManager(ctx, cfg)
		exitIfErr(err)
		defer manager.Close()
		embs, err := manager.EmbeddingService.GetEmbedder().EmbedStrings(ctx, []string{*text})
		exitIfErr(err)
		fmt.Printf("Embedded 1 text, dim=%d, first5=%v\n", len(embs[0]), embs[0][:min(5, len(embs[0]))])
	case "split":
		if *filePath == "" {
			exitIfErr(fmt.Errorf("-file is required for split"))
		}
		data, err := os.ReadFile(*filePath)
		exitIfErr(err)
		manager, err := milvus.InitMilvusManager(ctx, cfg)
		exitIfErr(err)
		defer manager.Close()
		chunks, err := manager.SplitterService.SplitMarkdown(ctx, string(data))
		exitIfErr(err)
		fmt.Printf("Split into %d chunks\n", len(chunks))
		for i, c := range chunks {
			if i >= 3 {
				break
			}
			fmt.Printf("Chunk %d: %s\n", i+1, truncate(c.Content, 120))
		}
	case "index":
		if *text == "" && *filePath == "" {
			exitIfErr(fmt.Errorf("either -text or -file is required for index"))
		}
		manager, err := milvus.InitMilvusManager(ctx, cfg)
		exitIfErr(err)
		defer manager.Close()
		var docs []*schema.Document
		if *text != "" {
			docs = []*schema.Document{
				{Content: *text, MetaData: map[string]any{"source": "cli"}},
			}
		} else {
			data, err := os.ReadFile(*filePath)
			exitIfErr(err)
			chunks, err := manager.SplitterService.SplitMarkdown(ctx, string(data))
			exitIfErr(err)
			docs = chunks
		}
		ids, err := manager.IndexerService.Store(ctx, docs)
		exitIfErr(err)
		fmt.Printf("Indexed %d documents\n", len(ids))
	case "retrieve":
		if *query == "" {
			exitIfErr(fmt.Errorf("-query is required for retrieve"))
		}
		cfg.Milvus.TopK = *topK
		manager, err := milvus.InitMilvusManager(ctx, cfg)
		exitIfErr(err)
		defer manager.Close()
		res, err := manager.RetrieverService.Retrieve(ctx, *query)
		exitIfErr(err)
		fmt.Printf("Query: %s\n", *query)
		for i, d := range res {
			score := d.MetaData["score"]
			fmt.Printf("%d) score=%v content=%s\n", i+1, score, truncate(d.Content, 160))
		}
	case "import-file":
		if *filePath == "" {
			exitIfErr(fmt.Errorf("-file is required for import-file"))
		}
		manager, err := milvus.InitMilvusManager(ctx, cfg)
		exitIfErr(err)
		defer manager.Close()
		importer, err := milvus.NewMarkdownImporter(manager)
		exitIfErr(err)
		opts := milvus.DefaultImportOptions()
		// language/category from flags (optional; empty means auto-infer)
		if l, ok := parseLanguage(*langFlag); ok {
			opts.Language = l
		}
		if c, ok := parseCategory(*catFlag); ok {
			opts.Category = c
		}
		opts.Source = *source
		res, err := importer.ImportFile(ctx, *filePath, opts)
		exitIfErr(err)
		fmt.Printf("Imported file. chunks=%d ids=%d\n", res.TotalChunks, len(res.DocumentIDs))
	case "import-dir":
		if *dirPath == "" {
			exitIfErr(fmt.Errorf("-dir is required for import-dir"))
		}
		manager, err := milvus.InitMilvusManager(ctx, cfg)
		exitIfErr(err)
		defer manager.Close()
		importer, err := milvus.NewMarkdownImporter(manager)
		exitIfErr(err)
		opts := milvus.DefaultImportOptions()
		if l, ok := parseLanguage(*langFlag); ok {
			opts.Language = l
		}
		if c, ok := parseCategory(*catFlag); ok {
			opts.Category = c
		}
		opts.Source = *source
		opts.Recursive = *recursive
		res, err := importer.ImportDirectory(ctx, *dirPath, opts)
		exitIfErr(err)
		fmt.Printf("Imported dir. files=%d chunks=%d ids=%d\n", res.SuccessFiles, res.TotalChunks, len(res.DocumentIDs))
	case "batch-import":
		if *paths == "" {
			exitIfErr(fmt.Errorf("-paths (comma-separated) is required for batch-import"))
		}
		manager, err := milvus.InitMilvusManager(ctx, cfg)
		exitIfErr(err)
		defer manager.Close()
		importer, err := milvus.NewMarkdownImporter(manager)
		exitIfErr(err)
		opts := milvus.DefaultImportOptions()
		if l, ok := parseLanguage(*langFlag); ok {
			opts.Language = l
		}
		if c, ok := parseCategory(*catFlag); ok {
			opts.Category = c
		}
		opts.Source = *source
		var list []string
		for _, p := range strings.Split(*paths, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				list = append(list, p)
			}
		}
		res, err := importer.BatchImport(ctx, list, opts)
		exitIfErr(err)
		fmt.Printf("Batch imported. files=%d chunks=%d ids=%d\n", res.SuccessFiles, res.TotalChunks, len(res.DocumentIDs))
	case "ds-import":
		// Single file under specific data source directory: -ds <name> and -file <filename>
		if *dsName == "" || *filePath == "" {
			exitIfErr(fmt.Errorf("-ds and -file (filename only) are required for ds-import"))
		}
		full := filepath.Join(dataBaseDir, *dsName, *filePath)
		manager, err := milvus.InitMilvusManager(ctx, cfg)
		exitIfErr(err)
		defer manager.Close()
		importer, err := milvus.NewMarkdownImporter(manager)
		exitIfErr(err)
		opts := milvus.DefaultImportOptions()
		if l, ok := parseLanguage(*langFlag); ok {
			opts.Language = l
		}
		if c, ok := parseCategory(*catFlag); ok {
			opts.Category = c
		}
		opts.Source = *dsName
		res, err := importer.ImportFile(ctx, full, opts)
		exitIfErr(err)
		fmt.Printf("Imported data source file. ds=%s file=%s chunks=%d ids=%d\n", *dsName, full, res.TotalChunks, len(res.DocumentIDs))
	case "ds-batch":
		// Batch import entire data source directory: -ds <name>
		if *dsName == "" {
			exitIfErr(fmt.Errorf("-ds is required for ds-batch"))
		}
		fullDir := filepath.Join(dataBaseDir, *dsName)
		manager, err := milvus.InitMilvusManager(ctx, cfg)
		exitIfErr(err)
		defer manager.Close()
		importer, err := milvus.NewMarkdownImporter(manager)
		exitIfErr(err)
		opts := milvus.DefaultImportOptions()
		if l, ok := parseLanguage(*langFlag); ok {
			opts.Language = l
		}
		if c, ok := parseCategory(*catFlag); ok {
			opts.Category = c
		}
		opts.Source = *dsName
		opts.Recursive = true
		res, err := importer.ImportDirectory(ctx, fullDir, opts)
		exitIfErr(err)
		fmt.Printf("Imported data source dir. ds=%s dir=%s files=%d chunks=%d ids=%d\n", *dsName, fullDir, res.SuccessFiles, res.TotalChunks, len(res.DocumentIDs))
	case "full-workflow":
		manager, err := milvus.InitMilvusManager(ctx, cfg)
		exitIfErr(err)
		defer manager.Close()
		longDoc := `Go 语言简介

Go（又称 Golang）是 Google 开发的一种静态强类型、编译型、并发型编程语言。
主要特点：
1. 简洁性
2. 并发性
3. 性能
4. 垃圾回收
`
		doc := &schema.Document{
			Content: longDoc,
			MetaData: map[string]any{
				"title":  "Go 语言介绍",
				"source": "cli",
			},
		}
		chunks, err := manager.SplitterService.Split(ctx, []*schema.Document{doc})
		exitIfErr(err)
		fmt.Printf("Split into %d chunks\n", len(chunks))
		ids, err := manager.IndexerService.Store(ctx, chunks)
		exitIfErr(err)
		fmt.Printf("Stored %d chunks\n", len(ids))
		time.Sleep(2 * time.Second)
		res, err := manager.RetrieverService.Retrieve(ctx, "Go 语言的并发特性是什么？")
		exitIfErr(err)
		for i, d := range res {
			fmt.Printf("%d) %s\n", i+1, truncate(d.Content, 160))
		}
	default:
		// Dynamic data source commands:
		// - import-<source>: single-file import, requires -file <filename only>
		// - batch-<source> or import-<source>-batch: batch import whole directory
		if strings.HasPrefix(*cmd, "import-") && *cmd != "import-file" && *cmd != "import-dir" {
			ds := strings.TrimPrefix(*cmd, "import-")
			// allow both "import-<ds>-batch" as batch
			if strings.HasSuffix(ds, "-batch") {
				ds = strings.TrimSuffix(ds, "-batch")
				// honor -language/-category flags for import-<ds>-batch alias
				var lang milvus.DocumentLanguage
				if l, ok := parseLanguage(*langFlag); ok {
					lang = l
				}
				var cat milvus.DocumentCategory
				if c, ok := parseCategory(*catFlag); ok {
					cat = c
				}
				runBatchForSource(ctx, cfg, dataBaseDir, ds, lang, cat)
				return
			}
			// single file import must provide filename only in -file
			if *filePath == "" {
				exitIfErr(fmt.Errorf("for %q you must provide -file <filename in data/%s>", *cmd, ds))
			}
			full := filepath.Join(dataBaseDir, ds, *filePath)
			// honor -language/-category flags for dynamic import-<source>
			var lang milvus.DocumentLanguage
			if l, ok := parseLanguage(*langFlag); ok {
				lang = l
			}
			var cat milvus.DocumentCategory
			if c, ok := parseCategory(*catFlag); ok {
				cat = c
			}
			runSingleForSource(ctx, cfg, full, ds, lang, cat)
			return
		}
		if strings.HasPrefix(*cmd, "batch-") {
			ds := strings.TrimPrefix(*cmd, "batch-")
			// honor -language/-category flags for dynamic batch-<source>
			var lang milvus.DocumentLanguage
			if l, ok := parseLanguage(*langFlag); ok {
				lang = l
			}
			var cat milvus.DocumentCategory
			if c, ok := parseCategory(*catFlag); ok {
				cat = c
			}
			runBatchForSource(ctx, cfg, dataBaseDir, ds, lang, cat)
			return
		}
		usage()
	}
}

func runSingleForSource(ctx context.Context, cfg *config.Config, fullPath string, ds string, lang milvus.DocumentLanguage, cat milvus.DocumentCategory) {
	manager, err := milvus.InitMilvusManager(ctx, cfg)
	exitIfErr(err)
	defer manager.Close()
	importer, err := milvus.NewMarkdownImporter(manager)
	exitIfErr(err)
	opts := milvus.DefaultImportOptions()
	opts.Source = ds
	// apply overrides if provided (non-empty)
	if lang != "" {
		opts.Language = lang
	}
	if cat != "" {
		opts.Category = cat
	}
	res, err := importer.ImportFile(ctx, fullPath, opts)
	exitIfErr(err)
	fmt.Printf("Imported data source file. ds=%s file=%s chunks=%d ids=%d\n", ds, fullPath, res.TotalChunks, len(res.DocumentIDs))
}

func runBatchForSource(ctx context.Context, cfg *config.Config, dataBaseDir, ds string, lang milvus.DocumentLanguage, cat milvus.DocumentCategory) {
	if ds == "" {
		exitIfErr(fmt.Errorf("data source name is required"))
	}
	fullDir := filepath.Join(dataBaseDir, ds)
	manager, err := milvus.InitMilvusManager(ctx, cfg)
	exitIfErr(err)
	defer manager.Close()
	importer, err := milvus.NewMarkdownImporter(manager)
	exitIfErr(err)
	opts := milvus.DefaultImportOptions()
	opts.Source = ds
	opts.Recursive = true
	// apply overrides if provided (non-empty)
	if lang != "" {
		opts.Language = lang
	}
	if cat != "" {
		opts.Category = cat
	}
	res, err := importer.ImportDirectory(ctx, fullDir, opts)
	exitIfErr(err)
	fmt.Printf("Imported data source dir. ds=%s dir=%s files=%d chunks=%d ids=%d\n", ds, fullDir, res.SuccessFiles, res.TotalChunks, len(res.DocumentIDs))
}

func exitIfErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// usage prints simple help text.
func usage() {
	fmt.Println("Usage:")
	fmt.Println("  From repo root:  go run ./backend/internal/eino/milvus/cmd/milvusctl -cmd <command> [flags]")
	fmt.Println("  From milvusctl/: go run . -cmd <command> [flags]")
	fmt.Println("Commands:")
	fmt.Println("  retrieve          -query <text>                 Retrieve similar docs [-topk N]")
	fmt.Println("  import-<source>   -file <filename>              Import one file under data/<source>/")
	fmt.Println("  batch-<source>                                  Import entire data/<source>/ directory")
	fmt.Println("  import-<source>-batch                           Same as batch-<source>")
	fmt.Println("Common import flags:")
	fmt.Println("  -language   golang|java|middleware   (optional, default: auto infer)")
	fmt.Println("  -category   basic|specialized|comprehensive (optional, default: auto infer)")
}

// parseLanguage maps user flag to milvus.DocumentLanguage.
func parseLanguage(s string) (milvus.DocumentLanguage, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "", "auto":
		return "", false
	case "go", "golang":
		return milvus.LanguageGolang, true
	case "java":
		return milvus.LanguageJava, true
	case "middleware", "中间件":
		return milvus.LanguageMiddleware, true
	default:
		return "", false
	}
}

// parseCategory maps user flag to milvus.DocumentCategory.
func parseCategory(s string) (milvus.DocumentCategory, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "", "auto":
		return "", false
	case "basic", "基础":
		return milvus.CategoryBasic, true
	case "specialized", "专项":
		return milvus.CategorySpecialized, true
	case "comprehensive", "综合":
		return milvus.CategoryComprehensive, true
	default:
		return "", false
	}
}
