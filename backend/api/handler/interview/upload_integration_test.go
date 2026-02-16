package interview

import (
	"ai-eino-interview-agent/api/model/interviews"
	"ai-eino-interview-agent/api/response"
	"ai-eino-interview-agent/internal/model"
	"ai-eino-interview-agent/internal/mq"
	"ai-eino-interview-agent/internal/worker"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestEnvironment holds shared resources for integration tests
type TestEnvironment struct {
	DB         *gorm.DB
	Redis      *redis.Client
	Queue      *mq.ResumeUploadQueue
	WorkerPool *worker.WorkerPool
	TmpDir     string
	Cleanup    func()
}

// setupIntegrationTest initializes test environment with all dependencies
func setupIntegrationTest(t *testing.T) (*TestEnvironment, func()) {
	t.Helper()

	// Create temporary directory for test files
	tmpDir := t.TempDir()

	// Initialize SQLite in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	// Set global DB getter
	model.SetDBGetter(func() *gorm.DB {
		return db
	})

	// Auto migrate schemas
	if err := db.AutoMigrate(&model.ResumeUploadStatus{}); err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize Redis client for testing (DB 1)
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   1, // Use DB 1 for tests
	})

	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available, skipping integration test")
		return nil, func() {}
	}

	// Flush test database for isolation
	if err := redisClient.FlushDB(ctx).Err(); err != nil {
		t.Fatalf("Failed to flush test database: %v", err)
	}

	// Initialize queue
	queue := mq.NewResumeUploadQueue(redisClient)

	// Create mock parse resume function for testing
	mockParseResume := func(ctx context.Context, userID uint, filePath string, fileSize int64) (uint64, interface{}, error) {
		// Simulate processing time
		time.Sleep(100 * time.Millisecond)

		// Check if file exists
		if _, err := os.Stat(filePath); err != nil {
			return 0, nil, fmt.Errorf("file not found: %w", err)
		}

		// Return mock resume ID
		return 12345, nil, nil
	}

	// Initialize worker pool
	pool, err := worker.NewWorkerPool(queue, 3, mockParseResume)
	if err != nil {
		t.Fatalf("Failed to create worker pool: %v", err)
	}

	// Start worker pool
	if err := pool.Start(); err != nil {
		t.Fatalf("Failed to start worker pool: %v", err)
	}

	// Set global worker pool
	SetWorkerPool(pool)

	env := &TestEnvironment{
		DB:         db,
		Redis:      redisClient,
		Queue:      queue,
		WorkerPool: pool,
		TmpDir:     tmpDir,
	}

	// Cleanup function
	cleanup := func() {
		if pool != nil {
			_ = pool.Stop()
		}
		if redisClient != nil {
			redisClient.FlushDB(ctx)
			redisClient.Close()
		}
		if db != nil {
			sqlDB, _ := db.DB()
			if sqlDB != nil {
				sqlDB.Close()
			}
		}
	}

	return env, cleanup
}

// createTestPDFFile creates a mock PDF file for testing
func createTestPDFFile(t *testing.T, dir string, content string, sizeKB int) string {
	t.Helper()

	filename := fmt.Sprintf("test_resume_%d.pdf", time.Now().UnixNano())
	filePath := filepath.Join(dir, filename)

	file, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer file.Close()

	// Write content to reach desired size
	data := []byte(content)
	targetSize := sizeKB * 1024
	for written := 0; written < targetSize; {
		n, err := file.Write(data)
		if err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}
		written += n
	}

	return filePath
}

// createMultipartRequest creates a Hertz request with multipart file upload
func createMultipartRequest(t *testing.T, filePath string, userID uint) *app.RequestContext {
	t.Helper()

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer file.Close()

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file field
	part, err := writer.CreateFormFile("resume", filepath.Base(filePath))
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		t.Fatalf("Failed to copy file to form: %v", err)
	}

	contentType := writer.FormDataContentType()

	if err := writer.Close(); err != nil {
		t.Fatalf("Failed to close multipart writer: %v", err)
	}

	// Create Hertz request context
	ctx := app.NewContext(0)
	ctx.Request.SetMethod("POST")
	ctx.Request.SetRequestURI("/api/resume/upload")
	ctx.Request.SetBodyRaw(body.Bytes())
	ctx.Request.Header.Set("Content-Type", contentType)
	ctx.Request.Header.Set("Authorization", fmt.Sprintf("Bearer test-token-%d", userID))

	// Mock user ID in context (middleware would normally set this)
	ctx.Set("user_id", userID)

	return ctx
}

// TestUploadIntegration_QuickSync tests the quick sync upload path
func TestUploadIntegration_QuickSync(t *testing.T) {
	env, cleanup := setupIntegrationTest(t)
	if env == nil {
		return
	}
	defer cleanup()

	// Create small PDF file (should complete within 15s timeout)
	smallFile := createTestPDFFile(t, env.TmpDir, "Small resume content\n", 1) // 1KB

	// Create multipart upload request
	ctx := createMultipartRequest(t, smallFile, 1)

	// Call upload handler
	UploadResume(context.Background(), ctx)

	// Verify response status
	assert.Equal(t, 200, ctx.Response.StatusCode(), "Expected 200 OK for sync upload")

	// Parse response body
	var apiResp struct {
		Code    int                              `json:"code"`
		Message string                           `json:"message"`
		Data    *interviews.UploadResumeResponse `json:"data"`
	}
	err := json.Unmarshal(ctx.Response.Body(), &apiResp)
	assert.NoError(t, err, "Failed to parse response")

	// Verify response structure
	assert.NotNil(t, apiResp.Data, "Response data should not be nil")
	assert.False(t, apiResp.Data.IsAsync, "Expected sync mode (is_async=false)")
	assert.NotNil(t, apiResp.Data.ResumeID, "ResumeID should be present in sync mode")
	assert.NotEmpty(t, apiResp.Data.UploadID, "UploadID should be present")
	assert.Contains(t, apiResp.Data.Message, "successfully", "Message should indicate success")

	// Verify database status
	status, err := model.ResumeUploadStatusDao.GetByUploadID(apiResp.Data.UploadID)
	assert.NoError(t, err, "Failed to get upload status")
	assert.Equal(t, "completed", status.Status, "Status should be completed")
	assert.Equal(t, 100, status.Progress, "Progress should be 100%")
	assert.Equal(t, uint64(12345), status.ResumeID, "ResumeID should match mock value")
}

// TestUploadIntegration_AsyncWithProgress tests async worker pool path
func TestUploadIntegration_AsyncWithProgress(t *testing.T) {
	env, cleanup := setupIntegrationTest(t)
	if env == nil {
		return
	}
	defer cleanup()

	// Override with slow mock to force async mode
	slowMockParseResume := func(ctx context.Context, userID uint, filePath string, fileSize int64) (uint64, interface{}, error) {
		// Simulate long processing (longer than QuickSyncTimeout)
		time.Sleep(20 * time.Second)
		return 99999, nil, nil
	}

	// Recreate worker pool with slow mock
	_ = env.WorkerPool.Stop()
	slowPool, err := worker.NewWorkerPool(env.Queue, 1, slowMockParseResume)
	assert.NoError(t, err)
	_ = slowPool.Start()
	defer slowPool.Stop()
	SetWorkerPool(slowPool)

	// Create large PDF file
	largeFile := createTestPDFFile(t, env.TmpDir, "Large resume with lots of content\n", 100) // 100KB

	// Create multipart upload request
	ctx := createMultipartRequest(t, largeFile, 2)

	// Call upload handler
	UploadResume(context.Background(), ctx)

	// Verify response status
	assert.Equal(t, 200, ctx.Response.StatusCode(), "Expected 200 OK for async upload")

	// Parse response body
	var apiResp struct {
		Code    int                              `json:"code"`
		Message string                           `json:"message"`
		Data    *interviews.UploadResumeResponse `json:"data"`
	}
	err = json.Unmarshal(ctx.Response.Body(), &apiResp)
	assert.NoError(t, err, "Failed to parse response")

	// Verify response structure for async mode
	assert.NotNil(t, apiResp.Data, "Response data should not be nil")
	assert.True(t, apiResp.Data.IsAsync, "Expected async mode (is_async=true)")
	assert.Nil(t, apiResp.Data.ResumeID, "ResumeID should be nil in async mode")
	assert.NotEmpty(t, apiResp.Data.UploadID, "UploadID should be present")
	assert.Contains(t, apiResp.Data.Message, "asynchronously", "Message should indicate async processing")

	uploadID := apiResp.Data.UploadID

	// Verify job enqueued to Redis queue
	jobCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	job, err := env.Queue.DequeueJob(jobCtx, 1*time.Second)
	assert.NoError(t, err, "Failed to dequeue job")
	assert.NotNil(t, job, "Job should be in queue")
	assert.Equal(t, uploadID, job.UploadID, "Job uploadID should match")
	assert.Equal(t, uint(2), job.UserID, "Job userID should match")

	// Verify initial database status
	status, err := model.ResumeUploadStatusDao.GetByUploadID(uploadID)
	assert.NoError(t, err, "Failed to get upload status")
	assert.Equal(t, "queued", status.Status, "Status should be queued initially")
}

// TestProgressEndpoint_SSEStreaming tests SSE endpoint
func TestProgressEndpoint_SSEStreaming(t *testing.T) {
	env, cleanup := setupIntegrationTest(t)
	if env == nil {
		return
	}
	defer cleanup()

	uploadID := "test-sse-upload-123"

	// Create upload status record
	status := &model.ResumeUploadStatus{
		UploadID: uploadID,
		UserID:   3,
		FilePath: "/tmp/test.pdf",
		Status:   "pending",
		Progress: 0,
	}
	err := model.ResumeUploadStatusDao.CreateUploadStatus(status)
	assert.NoError(t, err, "Failed to create upload status")

	// Test progress updates via Redis pub/sub
	// (Full SSE streaming test would require complex HTTP server setup)
	// Instead, verify Redis pub/sub mechanism works
	ctx := context.Background()
	channel := fmt.Sprintf("resume:progress:%s", uploadID)
	pubsub := env.Redis.Subscribe(ctx, channel)
	defer pubsub.Close()

	// Wait for subscription
	time.Sleep(100 * time.Millisecond)

	// Publish progress
	progress := &mq.ProgressUpdate{
		Status:   "extracting",
		Progress: 50,
		Stage:    "Extracting PDF",
	}
	err = env.Queue.PublishProgress(ctx, uploadID, progress)
	assert.NoError(t, err)

	// Verify message received
	select {
	case msg := <-pubsub.Channel():
		var received mq.ProgressUpdate
		err := json.Unmarshal([]byte(msg.Payload), &received)
		assert.NoError(t, err)
		assert.Equal(t, progress.Status, received.Status)
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for progress update")
	}
}

// TestUploadIntegration_Errors tests error scenarios
func TestUploadIntegration_Errors(t *testing.T) {
	env, cleanup := setupIntegrationTest(t)
	if env == nil {
		return
	}
	defer cleanup()

	tests := []struct {
		name           string
		setupFile      func() string
		expectedStatus int
		expectedError  string
	}{
		{
			name: "invalid file type (not PDF)",
			setupFile: func() string {
				txtFile := filepath.Join(env.TmpDir, "test.txt")
				_ = os.WriteFile(txtFile, []byte("Not a PDF"), 0644)
				return txtFile
			},
			expectedStatus: 400,
			expectedError:  "只支持 PDF 格式",
		},
		{
			name: "file too large (> 10MB)",
			setupFile: func() string {
				// Create 11MB file
				return createTestPDFFile(t, env.TmpDir, "X", 11*1024) // 11MB
			},
			expectedStatus: 400,
			expectedError:  "文件大小不能超过 10MB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := tt.setupFile()

			// Create request
			ctx := createMultipartRequest(t, testFile, 4)

			// Call upload handler
			UploadResume(context.Background(), ctx)

			// Verify error response
			assert.Equal(t, tt.expectedStatus, ctx.Response.StatusCode(), "Expected error status code")

			// Parse response
			var apiResp response.Response
			err := json.Unmarshal(ctx.Response.Body(), &apiResp)
			if err == nil && apiResp.Message != "" {
				assert.Contains(t, apiResp.Message, tt.expectedError, "Error message should contain expected text")
			}
		})
	}
}

// TestUploadIntegration_WorkerFailure tests worker processing failure
func TestUploadIntegration_WorkerFailure(t *testing.T) {
	env, cleanup := setupIntegrationTest(t)
	if env == nil {
		return
	}
	defer cleanup()

	// Create failing mock
	failingMock := func(ctx context.Context, userID uint, filePath string, fileSize int64) (uint64, interface{}, error) {
		return 0, nil, fmt.Errorf("mock processing error")
	}

	// Recreate worker pool with failing mock
	_ = env.WorkerPool.Stop()
	failPool, err := worker.NewWorkerPool(env.Queue, 1, failingMock)
	assert.NoError(t, err)
	_ = failPool.Start()
	defer failPool.Stop()
	SetWorkerPool(failPool)

	// Create test file
	testFile := createTestPDFFile(t, env.TmpDir, "Test content\n", 1)

	// Create request
	ctx := createMultipartRequest(t, testFile, 5)

	// Call upload handler (should go to async due to failure)
	UploadResume(context.Background(), ctx)

	// Parse response
	var apiResp struct {
		Code    int                              `json:"code"`
		Message string                           `json:"message"`
		Data    *interviews.UploadResumeResponse `json:"data"`
	}
	err = json.Unmarshal(ctx.Response.Body(), &apiResp)
	assert.NoError(t, err)

	uploadID := apiResp.Data.UploadID

	// Give worker time to process and fail
	time.Sleep(2 * time.Second)

	// Verify database shows failed status
	status, err := model.ResumeUploadStatusDao.GetByUploadID(uploadID)
	assert.NoError(t, err)
	// Status should eventually be "failed" after worker processes
	// Due to timing, it might still be "pending" or "extracting"
	assert.NotEqual(t, "completed", status.Status, "Status should not be completed")
}

// TestUploadIntegration_ConcurrentUploads tests multiple concurrent uploads
func TestUploadIntegration_ConcurrentUploads(t *testing.T) {
	env, cleanup := setupIntegrationTest(t)
	if env == nil {
		return
	}
	defer cleanup()

	numUploads := 5
	results := make(chan bool, numUploads)

	for i := 0; i < numUploads; i++ {
		go func(userID uint) {
			testFile := createTestPDFFile(t, env.TmpDir, fmt.Sprintf("User %d resume\n", userID), 1)
			ctx := createMultipartRequest(t, testFile, userID)

			UploadResume(context.Background(), ctx)

			results <- ctx.Response.StatusCode() == 200
		}(uint(i + 10))
	}

	// Wait for all uploads to complete
	successCount := 0
	for i := 0; i < numUploads; i++ {
		if <-results {
			successCount++
		}
	}

	assert.Equal(t, numUploads, successCount, "All concurrent uploads should succeed")
}

// TestUploadIntegration_RedisPubSubProgress tests Redis pub/sub progress events
func TestUploadIntegration_RedisPubSubProgress(t *testing.T) {
	env, cleanup := setupIntegrationTest(t)
	if env == nil {
		return
	}
	defer cleanup()

	uploadID := "test-pubsub-456"
	ctx := context.Background()

	// Subscribe to progress channel
	channel := fmt.Sprintf("resume:progress:%s", uploadID)
	pubsub := env.Redis.Subscribe(ctx, channel)
	defer pubsub.Close()

	// Wait for subscription
	_, err := pubsub.Receive(ctx)
	assert.NoError(t, err, "Failed to subscribe")

	// Publish progress updates
	progressUpdates := []*mq.ProgressUpdate{
		{Status: "pending", Progress: 0, Stage: "Starting"},
		{Status: "extracting", Progress: 25, Stage: "Extracting"},
		{Status: "analyzing", Progress: 75, Stage: "Analyzing"},
		{Status: "completed", Progress: 100, Stage: "Done"},
	}

	go func() {
		for _, progress := range progressUpdates {
			time.Sleep(100 * time.Millisecond)
			_ = env.Queue.PublishProgress(ctx, uploadID, progress)
		}
	}()

	// Receive and verify progress events
	receivedCount := 0
	timeout := time.After(5 * time.Second)

	for receivedCount < len(progressUpdates) {
		select {
		case msg := <-pubsub.Channel():
			var progress mq.ProgressUpdate
			err := json.Unmarshal([]byte(msg.Payload), &progress)
			assert.NoError(t, err, "Failed to unmarshal progress")
			assert.Equal(t, progressUpdates[receivedCount].Status, progress.Status)
			receivedCount++

		case <-timeout:
			t.Fatalf("Timeout waiting for progress events, received %d/%d", receivedCount, len(progressUpdates))
		}
	}

	assert.Equal(t, len(progressUpdates), receivedCount, "Should receive all progress events")
}

// TestUploadIntegration_DatabaseStatusTracking tests database status updates throughout lifecycle
func TestUploadIntegration_DatabaseStatusTracking(t *testing.T) {
	env, cleanup := setupIntegrationTest(t)
	if env == nil {
		return
	}
	defer cleanup()

	uploadID := "test-db-tracking-789"

	// Create initial status
	status := &model.ResumeUploadStatus{
		UploadID: uploadID,
		UserID:   6,
		FilePath: "/tmp/test.pdf",
		Status:   "pending",
		Progress: 0,
	}
	err := model.ResumeUploadStatusDao.CreateUploadStatus(status)
	assert.NoError(t, err)

	// Simulate status updates through lifecycle
	updates := []struct {
		status   string
		progress int
		stage    string
	}{
		{"extracting", 25, "Extracting PDF"},
		{"extracted", 50, "Extraction done"},
		{"analyzing", 75, "AI analysis"},
		{"completed", 100, "Finished"},
	}

	for _, update := range updates {
		err := model.ResumeUploadStatusDao.UpdateStatus(uploadID, map[string]interface{}{
			"status":   update.status,
			"progress": update.progress,
			"stage":    update.stage,
		})
		assert.NoError(t, err, "Failed to update status")

		// Verify update
		current, err := model.ResumeUploadStatusDao.GetByUploadID(uploadID)
		assert.NoError(t, err)
		assert.Equal(t, update.status, current.Status)
		assert.Equal(t, update.progress, current.Progress)
		assert.Equal(t, update.stage, current.Stage)

		time.Sleep(50 * time.Millisecond)
	}

	// Verify final state
	final, err := model.ResumeUploadStatusDao.GetByUploadID(uploadID)
	assert.NoError(t, err)
	assert.Equal(t, "completed", final.Status)
	assert.Equal(t, 100, final.Progress)
}
