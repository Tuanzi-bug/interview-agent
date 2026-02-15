package worker

import (
	"ai-eino-interview-agent/internal/mq"
	"ai-eino-interview-agent/internal/repository"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

// Helper function to get Redis client for testing
func getTestRedisClient(t *testing.T) *redis.Client {
	t.Helper()

	client := repository.GetRedis()
	if client == nil {
		t.Skip("Redis not available, skipping test")
		return nil
	}

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skipf("Redis not available: %v", err)
		return nil
	}

	// Use DB 1 for testing
	testClient := redis.NewClient(&redis.Options{
		Addr: client.Options().Addr,
		DB:   1,
	})

	// Flush test database
	if err := testClient.FlushDB(ctx).Err(); err != nil {
		t.Fatalf("Failed to flush test database: %v", err)
	}

	return testClient
}

// Helper function to create test resume file
func createTestResumeFile(t *testing.T) string {
	t.Helper()

	// Create temporary directory
	tmpDir := t.TempDir()

	// Create a test PDF file path
	testFilePath := filepath.Join(tmpDir, "test_resume.pdf")

	// Create the file (empty is fine for testing)
	file, err := os.Create(testFilePath)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer file.Close()

	// Write some content
	_, err = file.WriteString("Test resume content")
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	return testFilePath
}

func mockParseResumeFunc(ctx context.Context, userID uint, filePath string, fileSize int64) (uint64, interface{}, error) {
	return 12345, nil, nil
}

func mockFailingParseResumeFunc(ctx context.Context, userID uint, filePath string, fileSize int64) (uint64, interface{}, error) {
	return 0, nil, fmt.Errorf("mock parse error")
}

func TestNewWorkerPool(t *testing.T) {
	client := getTestRedisClient(t)
	if client == nil {
		return
	}
	defer client.Close()

	tests := []struct {
		name        string
		workerCount int
		wantErr     bool
	}{
		{
			name:        "valid worker count",
			workerCount: 3,
			wantErr:     false,
		},
		{
			name:        "single worker",
			workerCount: 1,
			wantErr:     false,
		},
		{
			name:        "zero workers should fail",
			workerCount: 0,
			wantErr:     true,
		},
		{
			name:        "negative workers should fail",
			workerCount: -1,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue := mq.NewResumeUploadQueue(client)
			pool, err := NewWorkerPool(queue, tt.workerCount, mockParseResumeFunc)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if pool == nil {
				t.Fatal("expected worker pool, got nil")
			}
		})
	}
}

func TestWorkerPoolStartStop(t *testing.T) {
	client := getTestRedisClient(t)
	if client == nil {
		return
	}
	defer client.Close()

	queue := mq.NewResumeUploadQueue(client)
	pool, err := NewWorkerPool(queue, 2, mockParseResumeFunc)
	if err != nil {
		t.Fatalf("failed to create worker pool: %v", err)
	}

	// Test Start
	if err := pool.Start(); err != nil {
		t.Fatalf("failed to start worker pool: %v", err)
	}

	// Give workers time to start
	time.Sleep(100 * time.Millisecond)

	// Test Stop
	if err := pool.Stop(); err != nil {
		t.Fatalf("failed to stop worker pool: %v", err)
	}

	// Test double stop (should be safe)
	if err := pool.Stop(); err != nil {
		t.Fatalf("double stop should not error: %v", err)
	}
}

func TestWorkerPoolProcessJob(t *testing.T) {
	client := getTestRedisClient(t)
	if client == nil {
		return
	}
	defer client.Close()

	// Create test file
	testFilePath := createTestResumeFile(t)

	tests := []struct {
		name    string
		job     *mq.ResumeUploadJob
		wantErr bool
	}{
		{
			name: "valid job",
			job: &mq.ResumeUploadJob{
				UploadID: "test-upload-123",
				UserID:   1,
				FilePath: testFilePath,
			},
			wantErr: false,
		},
		{
			name: "missing file",
			job: &mq.ResumeUploadJob{
				UploadID: "test-upload-456",
				UserID:   1,
				FilePath: "/nonexistent/file.pdf",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue := mq.NewResumeUploadQueue(client)
			pool, err := NewWorkerPool(queue, 1, mockParseResumeFunc)
			if err != nil {
				t.Fatalf("failed to create worker pool: %v", err)
			}

			ctx := context.Background()
			err = pool.processJob(ctx, tt.job)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestWorkerPoolJobProcessing(t *testing.T) {
	client := getTestRedisClient(t)
	if client == nil {
		return
	}
	defer client.Close()

	// Create test file
	testFilePath := createTestResumeFile(t)

	queue := mq.NewResumeUploadQueue(client)
	pool, err := NewWorkerPool(queue, 2, mockParseResumeFunc)
	if err != nil {
		t.Fatalf("failed to create worker pool: %v", err)
	}

	// Start worker pool
	if err := pool.Start(); err != nil {
		t.Fatalf("failed to start worker pool: %v", err)
	}
	defer pool.Stop()

	// Enqueue test job
	ctx := context.Background()
	job := &mq.ResumeUploadJob{
		UploadID: "test-job-123",
		UserID:   1,
		FilePath: testFilePath,
	}

	if err := queue.EnqueueJob(ctx, job); err != nil {
		t.Fatalf("failed to enqueue job: %v", err)
	}

	// Give worker time to process job
	// Note: This will fail in GREEN phase until we implement actual processing
	time.Sleep(2 * time.Second)

	// In a real test, we would verify the job was processed
	// For now, we just verify no panic occurred
}

func TestWorkerPoolProgressPublishing(t *testing.T) {
	client := getTestRedisClient(t)
	if client == nil {
		return
	}
	defer client.Close()

	queue := mq.NewResumeUploadQueue(client)
	_, err := NewWorkerPool(queue, 1, mockParseResumeFunc)
	if err != nil {
		t.Fatalf("failed to create worker pool: %v", err)
	}

	ctx := context.Background()
	uploadID := "test-progress-123"

	// Subscribe to progress updates
	pubsub := client.Subscribe(ctx, fmt.Sprintf("resume:progress:%s", uploadID))
	defer pubsub.Close()

	// Wait for subscription to be ready
	time.Sleep(50 * time.Millisecond)

	// Test progress publishing (direct call to internal method)
	progress := &mq.ProgressUpdate{
		Status:   "extracting",
		Progress: 25,
		Stage:    "Extracting text from PDF",
	}

	if err := queue.PublishProgress(ctx, uploadID, progress); err != nil {
		t.Fatalf("failed to publish progress: %v", err)
	}

	// Verify progress update received
	select {
	case msg := <-pubsub.Channel():
		if msg.Payload == "" {
			t.Error("received empty progress update")
		}
	case <-time.After(2 * time.Second):
		t.Error("timeout waiting for progress update")
	}
}

func TestWorkerPoolGracefulShutdown(t *testing.T) {
	client := getTestRedisClient(t)
	if client == nil {
		return
	}
	defer client.Close()

	queue := mq.NewResumeUploadQueue(client)
	pool, err := NewWorkerPool(queue, 3, mockParseResumeFunc)
	if err != nil {
		t.Fatalf("failed to create worker pool: %v", err)
	}

	// Start workers
	if err := pool.Start(); err != nil {
		t.Fatalf("failed to start worker pool: %v", err)
	}

	// Give workers time to start
	time.Sleep(100 * time.Millisecond)

	// Test graceful shutdown
	stopStart := time.Now()
	if err := pool.Stop(); err != nil {
		t.Fatalf("failed to stop worker pool: %v", err)
	}
	stopDuration := time.Since(stopStart)

	// Shutdown should be quick since no jobs are processing
	if stopDuration > 5*time.Second {
		t.Errorf("shutdown took too long: %v", stopDuration)
	}

	// Verify workers actually stopped by trying to stop again
	if err := pool.Stop(); err != nil {
		t.Errorf("stop after stop should be safe: %v", err)
	}
}

func TestWorkerPoolConcurrentProcessing(t *testing.T) {
	client := getTestRedisClient(t)
	if client == nil {
		return
	}
	defer client.Close()

	// Create test files
	testFile1 := createTestResumeFile(t)
	testFile2 := createTestResumeFile(t)
	testFile3 := createTestResumeFile(t)

	queue := mq.NewResumeUploadQueue(client)
	pool, err := NewWorkerPool(queue, 3, mockParseResumeFunc)
	if err != nil {
		t.Fatalf("failed to create worker pool: %v", err)
	}

	// Start worker pool
	if err := pool.Start(); err != nil {
		t.Fatalf("failed to start worker pool: %v", err)
	}
	defer pool.Stop()

	// Enqueue multiple jobs
	ctx := context.Background()
	jobs := []*mq.ResumeUploadJob{
		{UploadID: "concurrent-1", UserID: 1, FilePath: testFile1},
		{UploadID: "concurrent-2", UserID: 2, FilePath: testFile2},
		{UploadID: "concurrent-3", UserID: 3, FilePath: testFile3},
	}

	for _, job := range jobs {
		if err := queue.EnqueueJob(ctx, job); err != nil {
			t.Fatalf("failed to enqueue job: %v", err)
		}
	}

	// Give workers time to process jobs
	time.Sleep(3 * time.Second)

	// In a real test, we would verify all jobs were processed
	// For now, we just verify no panic occurred
}

func TestWorkerPoolErrorHandling(t *testing.T) {
	client := getTestRedisClient(t)
	if client == nil {
		return
	}
	defer client.Close()

	queue := mq.NewResumeUploadQueue(client)
	pool, err := NewWorkerPool(queue, 1, mockParseResumeFunc)
	if err != nil {
		t.Fatalf("failed to create worker pool: %v", err)
	}

	// Start worker pool
	if err := pool.Start(); err != nil {
		t.Fatalf("failed to start worker pool: %v", err)
	}
	defer pool.Stop()

	// Enqueue job with invalid file path
	ctx := context.Background()
	job := &mq.ResumeUploadJob{
		UploadID: "error-test",
		UserID:   1,
		FilePath: "/invalid/path/to/file.pdf",
	}

	if err := queue.EnqueueJob(ctx, job); err != nil {
		t.Fatalf("failed to enqueue job: %v", err)
	}

	// Give worker time to process (and fail)
	time.Sleep(2 * time.Second)

	// Worker should continue running despite error
	// Verify by enqueuing valid job after error
	testFile := createTestResumeFile(t)
	validJob := &mq.ResumeUploadJob{
		UploadID: "valid-after-error",
		UserID:   1,
		FilePath: testFile,
	}

	if err := queue.EnqueueJob(ctx, validJob); err != nil {
		t.Fatalf("failed to enqueue valid job: %v", err)
	}

	time.Sleep(2 * time.Second)

	// If we reach here without panic, error handling works
}
