package service

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"ai-eino-interview-agent/internal/repository"

	"github.com/redis/go-redis/v9"
)

// setupTestRedis initializes Redis client for testing
func setupTestRedis(t *testing.T) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1, // Use separate DB for testing
	})

	// Verify connection
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		t.Skipf("Redis not available for testing: %v", err)
	}

	// Clean up test data
	client.FlushDB(ctx)

	// Set global Redis client for service
	repository.RedisClient = client

	return client
}

// createTestPDFFile creates a temporary PDF file for testing
func createTestPDFFile(t *testing.T, content string) string {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.pdf")

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	return filePath
}

func TestGetCacheKey(t *testing.T) {
	cache := NewResumeTextCache()

	// Create two test files with same content
	content := "test PDF content"
	file1 := createTestPDFFile(t, content)
	file2 := createTestPDFFile(t, content)

	// Create a file with different content
	file3 := createTestPDFFile(t, "different content")

	key1, err := cache.GetCacheKey(file1)
	if err != nil {
		t.Fatalf("GetCacheKey failed: %v", err)
	}

	key2, err := cache.GetCacheKey(file2)
	if err != nil {
		t.Fatalf("GetCacheKey failed: %v", err)
	}

	key3, err := cache.GetCacheKey(file3)
	if err != nil {
		t.Fatalf("GetCacheKey failed: %v", err)
	}

	// Same content should generate same hash
	if key1 != key2 {
		t.Errorf("Same content should generate same key, got %s and %s", key1, key2)
	}

	// Different content should generate different hash
	if key1 == key3 {
		t.Errorf("Different content should generate different keys")
	}

	// Key should have correct prefix
	if len(key1) <= len(PDFTextCacheKeyPrefix) {
		t.Errorf("Key should contain prefix and hash")
	}
}

func TestGetCacheKey_NonExistentFile(t *testing.T) {
	cache := NewResumeTextCache()

	_, err := cache.GetCacheKey("/nonexistent/file.pdf")
	if err == nil {
		t.Error("GetCacheKey should return error for non-existent file")
	}
}

func TestGet_CacheMiss(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	cache := NewResumeTextCache()
	ctx := context.Background()

	testFile := createTestPDFFile(t, "test content")

	text, hit, err := cache.Get(ctx, testFile)

	if err != nil {
		t.Errorf("Get should not return error on cache miss: %v", err)
	}

	if hit {
		t.Error("Cache hit should be false when key doesn't exist")
	}

	if text != "" {
		t.Errorf("Expected empty text on cache miss, got: %s", text)
	}
}

func TestGet_CacheHit(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	cache := NewResumeTextCache()
	ctx := context.Background()

	testFile := createTestPDFFile(t, "test content")
	expectedText := "cached PDF text content"

	// Set cache first
	err := cache.Set(ctx, testFile, expectedText)
	if err != nil {
		t.Fatalf("Failed to set cache: %v", err)
	}

	// Now get it
	text, hit, err := cache.Get(ctx, testFile)

	if err != nil {
		t.Errorf("Get failed: %v", err)
	}

	if !hit {
		t.Error("Cache hit should be true")
	}

	if text != expectedText {
		t.Errorf("Expected text %s, got %s", expectedText, text)
	}
}

func TestSet(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	cache := NewResumeTextCache()
	ctx := context.Background()

	testFile := createTestPDFFile(t, "test content")
	testText := "PDF text to cache"

	err := cache.Set(ctx, testFile, testText)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	// Verify it was actually cached
	cacheKey, _ := cache.GetCacheKey(testFile)
	cachedValue, err := client.Get(ctx, cacheKey).Result()

	if err != nil {
		t.Errorf("Failed to retrieve cached value: %v", err)
	}

	if cachedValue != testText {
		t.Errorf("Expected cached value %s, got %s", testText, cachedValue)
	}

	// Verify TTL is set correctly
	ttl, err := client.TTL(ctx, cacheKey).Result()
	if err != nil {
		t.Errorf("Failed to get TTL: %v", err)
	}

	// TTL should be approximately 30 days (allow 1 minute variance)
	expectedTTL := PDFTextCacheTTL
	if ttl < expectedTTL-time.Minute || ttl > expectedTTL {
		t.Errorf("Expected TTL around %v, got %v", expectedTTL, ttl)
	}
}

func TestDelete(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	cache := NewResumeTextCache()
	ctx := context.Background()

	testFile := createTestPDFFile(t, "test content")
	testText := "PDF text to cache"

	// Set cache first
	err := cache.Set(ctx, testFile, testText)
	if err != nil {
		t.Fatalf("Failed to set cache: %v", err)
	}

	// Verify it exists
	_, hit, _ := cache.Get(ctx, testFile)
	if !hit {
		t.Fatal("Cache should be set before deletion")
	}

	// Delete it
	err = cache.Delete(ctx, testFile)
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}

	// Verify it's gone
	_, hit, _ = cache.Get(ctx, testFile)
	if hit {
		t.Error("Cache should not exist after deletion")
	}
}

func TestGet_RedisNotInitialized(t *testing.T) {
	// Save current Redis client
	oldClient := repository.RedisClient
	defer func() {
		repository.RedisClient = oldClient
	}()

	// Set Redis client to nil
	repository.RedisClient = nil

	cache := NewResumeTextCache()
	ctx := context.Background()

	testFile := createTestPDFFile(t, "test content")

	text, hit, err := cache.Get(ctx, testFile)

	if err == nil {
		t.Error("Get should return error when Redis not initialized")
	}

	if hit {
		t.Error("Cache hit should be false when Redis not initialized")
	}

	if text != "" {
		t.Error("Text should be empty when Redis not initialized")
	}
}

func TestSet_RedisNotInitialized(t *testing.T) {
	// Save current Redis client
	oldClient := repository.RedisClient
	defer func() {
		repository.RedisClient = oldClient
	}()

	// Set Redis client to nil
	repository.RedisClient = nil

	cache := NewResumeTextCache()
	ctx := context.Background()

	testFile := createTestPDFFile(t, "test content")

	err := cache.Set(ctx, testFile, "test text")

	// Should return error
	if err == nil {
		t.Error("Set should return error when Redis not initialized")
	}
}
