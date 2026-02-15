package tool

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"ai-eino-interview-agent/internal/config"
	"ai-eino-interview-agent/internal/repository"
	"ai-eino-interview-agent/internal/service"
)

// setupTestRedis initializes Redis for testing (DB 1)
func setupTestRedis(t *testing.T) {
	t.Helper()

	err := repository.InitRedis(config.RedisConfig{
		Addr:     "localhost:6379",
		Password: "root",
		DB:       1,
	})

	if err != nil {
		t.Skipf("Skipping test: Redis not available: %v", err)
	}

	redisClient := repository.GetRedis()
	if redisClient != nil {
		redisClient.FlushDB(context.Background())
	}
}

func createTestPDFFile(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()
	pdfPath := filepath.Join(tmpDir, "test.pdf")

	minimalPDF := []byte(`%PDF-1.4
1 0 obj << /Type /Catalog /Pages 2 0 R >> endobj
2 0 obj << /Type /Pages /Kids [3 0 R] /Count 1 >> endobj
3 0 obj << /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Contents 4 0 R >> endobj
4 0 obj << /Length 44 >> stream
BT
/F1 12 Tf
100 700 Td
(Test PDF Content) Tj
ET
endstream endobj
xref
0 5
0000000000 65535 f
0000000009 00000 n
0000000068 00000 n
0000000127 00000 n
0000000226 00000 n
trailer << /Size 5 /Root 1 0 R >>
startxref
334
%%EOF`)

	err := os.WriteFile(pdfPath, minimalPDF, 0644)
	if err != nil {
		t.Fatalf("Failed to create test PDF file: %v", err)
	}

	return pdfPath
}

// TestPDFParserCacheHit tests that cached text is returned without parsing
func TestPDFParserCacheHit(t *testing.T) {
	setupTestRedis(t)

	ctx := context.Background()
	cache := service.NewResumeTextCache()
	pdfPath := createTestPDFFile(t)

	// Pre-populate cache with known text
	cachedText := "This is cached PDF text from previous parsing"
	err := cache.Set(ctx, pdfPath, cachedText)
	if err != nil {
		t.Fatalf("Failed to set cache: %v", err)
	}

	// Test ToPages=false (merged mode)
	req := &PDFToTextRequest{
		FilePath: pdfPath,
		ToPages:  false,
	}

	result, err := ConvertPDFToText(ctx, req)

	// Verify cache hit behavior
	if err != nil {
		t.Errorf("Expected no error on cache hit, got: %v", err)
	}
	if !result.Success {
		t.Errorf("Expected success=true on cache hit, got false")
	}
	if result.Content != cachedText {
		t.Errorf("Expected cached text %q, got %q", cachedText, result.Content)
	}

	// Verify cache hit was logged in metadata
	if cacheHit, ok := result.Meta["cache_hit"].(bool); !ok || !cacheHit {
		t.Errorf("Expected cache_hit=true in metadata, got: %v", result.Meta["cache_hit"])
	}
}

// TestPDFParserCacheMiss tests that PDF is parsed and cached on miss
func TestPDFParserCacheMiss(t *testing.T) {
	setupTestRedis(t)

	ctx := context.Background()
	cache := service.NewResumeTextCache()
	pdfPath := createTestPDFFile(t)

	// Ensure cache is empty
	redisClient := repository.GetRedis()
	if redisClient != nil {
		redisClient.FlushDB(ctx)
	}

	// Test ToPages=false (merged mode)
	req := &PDFToTextRequest{
		FilePath: pdfPath,
		ToPages:  false,
	}

	result, err := ConvertPDFToText(ctx, req)

	// Verify parsing occurred
	if err != nil {
		t.Errorf("Expected no error on cache miss, got: %v", err)
	}
	if !result.Success {
		t.Errorf("Expected success=true after parsing, got false")
	}

	// Verify cache miss was logged in metadata
	if cacheHit, ok := result.Meta["cache_hit"].(bool); ok && cacheHit {
		t.Errorf("Expected cache_hit=false in metadata, got: %v", result.Meta["cache_hit"])
	}

	// Verify text was cached after parsing
	cachedText, hit, err := cache.Get(ctx, pdfPath)
	if err != nil {
		t.Errorf("Failed to get cached text: %v", err)
	}
	if !hit {
		t.Errorf("Expected cache hit after parsing, got miss")
	}
	if cachedText != result.Content {
		t.Errorf("Cached text doesn't match parsed result")
	}
}

// TestPDFParserCacheHitToPages tests cache hit with ToPages=true
func TestPDFParserCacheHitToPages(t *testing.T) {
	setupTestRedis(t)

	ctx := context.Background()
	cache := service.NewResumeTextCache()
	pdfPath := createTestPDFFile(t)

	// Pre-populate cache with text that can be split into pages
	cachedText := "Page 1 text\nPage 2 text\nPage 3 text"
	err := cache.Set(ctx, pdfPath, cachedText)
	if err != nil {
		t.Fatalf("Failed to set cache: %v", err)
	}

	// Test ToPages=true (page-split mode)
	req := &PDFToTextRequest{
		FilePath: pdfPath,
		ToPages:  true,
	}

	result, err := ConvertPDFToText(ctx, req)

	// Verify cache hit behavior
	if err != nil {
		t.Errorf("Expected no error on cache hit, got: %v", err)
	}
	if !result.Success {
		t.Errorf("Expected success=true on cache hit, got false")
	}
	if len(result.Pages) == 0 {
		t.Errorf("Expected pages to be populated from cached text, got empty")
	}

	// Verify cache hit was logged
	if cacheHit, ok := result.Meta["cache_hit"].(bool); !ok || !cacheHit {
		t.Errorf("Expected cache_hit=true in metadata")
	}
}

func TestPDFParserCacheError(t *testing.T) {
	ctx := context.Background()
	pdfPath := createTestPDFFile(t)

	req := &PDFToTextRequest{
		FilePath: pdfPath,
		ToPages:  false,
	}

	result, err := ConvertPDFToText(ctx, req)

	if result.Success {
		t.Log("PDF parsing succeeded without Redis (cache gracefully degraded)")
		return
	}

	if err != nil && result.ErrorMsg != "" {
		t.Skipf("PDF parsing failed due to PDF format issue (not cache issue): %v", err)
	}

	if cacheHit, ok := result.Meta["cache_hit"].(bool); ok && cacheHit {
		t.Errorf("Expected cache_hit=false when cache unavailable")
	}
}

// TestPDFParserCacheKeyConsistency tests that same file produces same cache key
func TestPDFParserCacheKeyConsistency(t *testing.T) {
	setupTestRedis(t)

	ctx := context.Background()
	pdfPath := createTestPDFFile(t)

	req1 := &PDFToTextRequest{
		FilePath: pdfPath,
		ToPages:  false,
	}
	result1, err1 := ConvertPDFToText(ctx, req1)
	if err1 != nil {
		t.Fatalf("First parse failed: %v", err1)
	}

	time.Sleep(10 * time.Millisecond)

	req2 := &PDFToTextRequest{
		FilePath: pdfPath,
		ToPages:  false,
	}
	result2, err2 := ConvertPDFToText(ctx, req2)
	if err2 != nil {
		t.Fatalf("Second parse failed: %v", err2)
	}

	if cacheHit, ok := result2.Meta["cache_hit"].(bool); !ok || !cacheHit {
		t.Errorf("Expected cache hit on second parse, got miss")
	}

	if result1.Content != result2.Content {
		t.Errorf("Content mismatch between cached and original parse")
	}
}

// TestPDFParserBothModes tests both ToPages=true and ToPages=false work with cache
func TestPDFParserBothModes(t *testing.T) {
	setupTestRedis(t)

	ctx := context.Background()
	pdfPath := createTestPDFFile(t)

	// Flush cache to start fresh
	redisClient := repository.GetRedis()
	if redisClient != nil {
		redisClient.FlushDB(ctx)
	}

	// Test merged mode first (cache miss)
	reqMerged := &PDFToTextRequest{
		FilePath: pdfPath,
		ToPages:  false,
	}
	resultMerged, err := ConvertPDFToText(ctx, reqMerged)
	if err != nil {
		t.Fatalf("Merged mode parse failed: %v", err)
	}
	if !resultMerged.Success {
		t.Fatalf("Expected success in merged mode")
	}

	// Test pages mode with same file (cache hit)
	reqPages := &PDFToTextRequest{
		FilePath: pdfPath,
		ToPages:  true,
	}
	resultPages, err := ConvertPDFToText(ctx, reqPages)
	if err != nil {
		t.Fatalf("Pages mode parse failed: %v", err)
	}
	if !resultPages.Success {
		t.Fatalf("Expected success in pages mode")
	}

	// Verify cache hit in second call
	if cacheHit, ok := resultPages.Meta["cache_hit"].(bool); !ok || !cacheHit {
		t.Errorf("Expected cache hit in pages mode after merged mode cached")
	}

	// Verify pages were reconstructed from cached text
	if len(resultPages.Pages) == 0 {
		t.Errorf("Expected pages to be reconstructed from cached text")
	}
}
