package mq

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"ai-eino-interview-agent/internal/repository"

	"github.com/redis/go-redis/v9"
)

// setupTestRedis creates a Redis client for testing using DB 1
func setupTestRedis(t *testing.T) *redis.Client {
	t.Helper()

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1, // Use DB 1 for testing
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available, skipping test")
		return nil
	}

	// Flush DB before each test for isolation
	if err := client.FlushDB(ctx).Err(); err != nil {
		t.Fatalf("Failed to flush test database: %v", err)
	}

	// Cleanup after test
	t.Cleanup(func() {
		client.FlushDB(ctx)
		client.Close()
	})

	return client
}

// TestNewResumeUploadQueue tests queue initialization
func TestNewResumeUploadQueue(t *testing.T) {
	client := setupTestRedis(t)
	if client == nil {
		return
	}

	queue := NewResumeUploadQueue(client)

	if queue == nil {
		t.Fatal("Expected non-nil queue")
	}

	if queue.client == nil {
		t.Error("Expected non-nil Redis client")
	}

	if queue.queueKey != "resume:upload:queue" {
		t.Errorf("Expected queue key 'resume:upload:queue', got %s", queue.queueKey)
	}
}

// TestEnqueueJob tests adding a job to the queue
func TestEnqueueJob(t *testing.T) {
	client := setupTestRedis(t)
	if client == nil {
		return
	}

	queue := NewResumeUploadQueue(client)
	ctx := context.Background()

	job := &ResumeUploadJob{
		UploadID: "test-upload-123",
		UserID:   42,
		FilePath: "/tmp/resume.pdf",
		ResumeID: 100,
	}

	// Enqueue the job
	err := queue.EnqueueJob(ctx, job)
	if err != nil {
		t.Fatalf("EnqueueJob failed: %v", err)
	}

	// Verify the job is in Redis LIST
	length := client.LLen(ctx, "resume:upload:queue").Val()
	if length != 1 {
		t.Errorf("Expected queue length 1, got %d", length)
	}

	// Verify the job data
	data := client.LPop(ctx, "resume:upload:queue").Val()
	var storedJob ResumeUploadJob
	if err := json.Unmarshal([]byte(data), &storedJob); err != nil {
		t.Fatalf("Failed to unmarshal job: %v", err)
	}

	if storedJob.UploadID != job.UploadID {
		t.Errorf("Expected UploadID %s, got %s", job.UploadID, storedJob.UploadID)
	}
	if storedJob.UserID != job.UserID {
		t.Errorf("Expected UserID %d, got %d", job.UserID, storedJob.UserID)
	}
	if storedJob.FilePath != job.FilePath {
		t.Errorf("Expected FilePath %s, got %s", job.FilePath, storedJob.FilePath)
	}
	if storedJob.ResumeID != job.ResumeID {
		t.Errorf("Expected ResumeID %d, got %d", job.ResumeID, storedJob.ResumeID)
	}
}

// TestDequeueJob tests getting a job from the queue
func TestDequeueJob(t *testing.T) {
	client := setupTestRedis(t)
	if client == nil {
		return
	}

	queue := NewResumeUploadQueue(client)
	ctx := context.Background()

	// Enqueue a test job first
	originalJob := &ResumeUploadJob{
		UploadID: "test-upload-456",
		UserID:   99,
		FilePath: "/tmp/test.pdf",
		ResumeID: 200,
	}
	err := queue.EnqueueJob(ctx, originalJob)
	if err != nil {
		t.Fatalf("Failed to enqueue job: %v", err)
	}

	// Dequeue the job
	job, err := queue.DequeueJob(ctx, 2*time.Second)
	if err != nil {
		t.Fatalf("DequeueJob failed: %v", err)
	}

	if job == nil {
		t.Fatal("Expected non-nil job")
	}

	// Verify job data
	if job.UploadID != originalJob.UploadID {
		t.Errorf("Expected UploadID %s, got %s", originalJob.UploadID, job.UploadID)
	}
	if job.UserID != originalJob.UserID {
		t.Errorf("Expected UserID %d, got %d", originalJob.UserID, job.UserID)
	}
	if job.FilePath != originalJob.FilePath {
		t.Errorf("Expected FilePath %s, got %s", originalJob.FilePath, job.FilePath)
	}
	if job.ResumeID != originalJob.ResumeID {
		t.Errorf("Expected ResumeID %d, got %d", originalJob.ResumeID, job.ResumeID)
	}

	// Verify queue is empty
	length := client.LLen(ctx, "resume:upload:queue").Val()
	if length != 0 {
		t.Errorf("Expected empty queue, got length %d", length)
	}
}

// TestDequeueJobTimeout tests blocking dequeue with timeout
func TestDequeueJobTimeout(t *testing.T) {
	client := setupTestRedis(t)
	if client == nil {
		return
	}

	queue := NewResumeUploadQueue(client)
	ctx := context.Background()

	// Try to dequeue from empty queue with short timeout
	start := time.Now()
	job, err := queue.DequeueJob(ctx, 1*time.Second)
	elapsed := time.Since(start)

	// Should return nil job and no error (timeout is expected)
	if job != nil {
		t.Error("Expected nil job for empty queue")
	}
	if err != nil {
		t.Errorf("Expected no error for timeout, got %v", err)
	}

	// Should have waited approximately 1 second
	if elapsed < 900*time.Millisecond || elapsed > 1500*time.Millisecond {
		t.Errorf("Expected ~1s timeout, took %v", elapsed)
	}
}

// TestPublishProgress tests progress update publishing
func TestPublishProgress(t *testing.T) {
	client := setupTestRedis(t)
	if client == nil {
		return
	}

	queue := NewResumeUploadQueue(client)
	ctx := context.Background()

	uploadID := "test-upload-789"
	progress := &ProgressUpdate{
		Status:   "extracting",
		Progress: 50,
		Stage:    "Extracting text from PDF",
	}

	// Subscribe to the progress channel
	channelName := "resume:progress:" + uploadID
	pubsub := client.Subscribe(ctx, channelName)
	defer pubsub.Close()

	// Wait for subscription to be ready
	_, err := pubsub.Receive(ctx)
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	// Publish progress in a goroutine
	go func() {
		time.Sleep(100 * time.Millisecond)
		if err := queue.PublishProgress(ctx, uploadID, progress); err != nil {
			t.Errorf("PublishProgress failed: %v", err)
		}
	}()

	// Receive the message
	msg, err := pubsub.ReceiveTimeout(ctx, 2*time.Second)
	if err != nil {
		t.Fatalf("Failed to receive message: %v", err)
	}

	// Verify message content
	message, ok := msg.(*redis.Message)
	if !ok {
		t.Fatal("Expected *redis.Message")
	}

	var received ProgressUpdate
	if err := json.Unmarshal([]byte(message.Payload), &received); err != nil {
		t.Fatalf("Failed to unmarshal progress: %v", err)
	}

	if received.Status != progress.Status {
		t.Errorf("Expected status %s, got %s", progress.Status, received.Status)
	}
	if received.Progress != progress.Progress {
		t.Errorf("Expected progress %d, got %d", progress.Progress, received.Progress)
	}
	if received.Stage != progress.Stage {
		t.Errorf("Expected stage %s, got %s", progress.Stage, received.Stage)
	}
}

// TestMultipleJobsOrdering tests FIFO behavior
func TestMultipleJobsOrdering(t *testing.T) {
	client := setupTestRedis(t)
	if client == nil {
		return
	}

	queue := NewResumeUploadQueue(client)
	ctx := context.Background()

	// Enqueue multiple jobs
	jobs := []*ResumeUploadJob{
		{UploadID: "job-1", UserID: 1, FilePath: "/tmp/1.pdf"},
		{UploadID: "job-2", UserID: 2, FilePath: "/tmp/2.pdf"},
		{UploadID: "job-3", UserID: 3, FilePath: "/tmp/3.pdf"},
	}

	for _, job := range jobs {
		if err := queue.EnqueueJob(ctx, job); err != nil {
			t.Fatalf("Failed to enqueue job %s: %v", job.UploadID, err)
		}
	}

	// Dequeue jobs and verify FIFO order
	for i, expected := range jobs {
		job, err := queue.DequeueJob(ctx, 1*time.Second)
		if err != nil {
			t.Fatalf("Failed to dequeue job %d: %v", i, err)
		}
		if job.UploadID != expected.UploadID {
			t.Errorf("Job %d: expected UploadID %s, got %s", i, expected.UploadID, job.UploadID)
		}
	}
}

// TestNilRedisClient tests graceful handling of nil client
func TestNilRedisClient(t *testing.T) {
	queue := NewResumeUploadQueue(nil)
	ctx := context.Background()

	job := &ResumeUploadJob{
		UploadID: "test",
		UserID:   1,
		FilePath: "/tmp/test.pdf",
	}

	// Operations should fail gracefully with nil client
	err := queue.EnqueueJob(ctx, job)
	if err == nil {
		t.Error("Expected error for nil Redis client")
	}

	_, err = queue.DequeueJob(ctx, 1*time.Second)
	if err == nil {
		t.Error("Expected error for nil Redis client")
	}

	progress := &ProgressUpdate{Status: "test", Progress: 0, Stage: "test"}
	err = queue.PublishProgress(ctx, "test", progress)
	if err == nil {
		t.Error("Expected error for nil Redis client")
	}
}

// TestWithRepositoryRedisClient tests integration with repository.GetRedis()
func TestWithRepositoryRedisClient(t *testing.T) {
	// Initialize repository Redis client (if available)
	if repository.GetRedis() == nil {
		t.Skip("Repository Redis client not initialized")
		return
	}

	queue := NewResumeUploadQueue(repository.GetRedis())
	ctx := context.Background()

	// Basic functionality test
	job := &ResumeUploadJob{
		UploadID: "repo-test",
		UserID:   1,
		FilePath: "/tmp/test.pdf",
	}

	err := queue.EnqueueJob(ctx, job)
	if err != nil {
		t.Fatalf("Failed to enqueue with repository client: %v", err)
	}

	// Cleanup
	repository.GetRedis().Del(ctx, "resume:upload:queue")
}
