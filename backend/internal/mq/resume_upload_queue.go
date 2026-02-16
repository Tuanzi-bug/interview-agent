package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	ResumeUploadQueueKey = "resume:upload:queue"
)

type ResumeUploadJob struct {
	UploadID    string `json:"upload_id"`
	UserID      uint   `json:"user_id"`
	FilePath    string `json:"file_path"`
	ResumeID    uint64 `json:"resume_id,omitempty"`
	TimeoutSecs int    `json:"timeout_secs,omitempty"` // Per-job timeout in seconds; 0 means use worker context
}

type ProgressUpdate struct {
	Status   string `json:"status"`
	Progress int    `json:"progress"`
	Stage    string `json:"stage"`
	ErrorMsg string `json:"error_msg,omitempty"`
}

type ResumeUploadQueue struct {
	client   *redis.Client
	queueKey string
}

func NewResumeUploadQueue(client *redis.Client) *ResumeUploadQueue {
	return &ResumeUploadQueue{
		client:   client,
		queueKey: ResumeUploadQueueKey,
	}
}

func (q *ResumeUploadQueue) EnqueueJob(ctx context.Context, job *ResumeUploadJob) error {
	if q.client == nil {
		return fmt.Errorf("redis client not initialized")
	}

	jobJSON, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	result := q.client.LPush(ctx, q.queueKey, string(jobJSON))
	if err := result.Err(); err != nil {
		return fmt.Errorf("failed to enqueue job: %w", err)
	}

	log.Printf("[ResumeUploadQueue] Job enqueued: uploadID=%s, userID=%d, filePath=%s",
		job.UploadID, job.UserID, job.FilePath)

	return nil
}

func (q *ResumeUploadQueue) DequeueJob(ctx context.Context, timeout time.Duration) (*ResumeUploadJob, error) {
	if q.client == nil {
		return nil, fmt.Errorf("redis client not initialized")
	}

	result := q.client.BRPop(ctx, timeout, q.queueKey)
	if err := result.Err(); err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to dequeue job: %w", err)
	}

	values := result.Val()
	if len(values) != 2 {
		return nil, fmt.Errorf("unexpected BRPOP result length: %d", len(values))
	}

	jobJSON := values[1]

	var job ResumeUploadJob
	if err := json.Unmarshal([]byte(jobJSON), &job); err != nil {
		return nil, fmt.Errorf("failed to unmarshal job: %w", err)
	}

	log.Printf("[ResumeUploadQueue] Job dequeued: uploadID=%s, userID=%d, filePath=%s",
		job.UploadID, job.UserID, job.FilePath)

	return &job, nil
}

func (q *ResumeUploadQueue) PublishProgress(ctx context.Context, uploadID string, progress *ProgressUpdate) error {
	if q.client == nil {
		return fmt.Errorf("redis client not initialized")
	}

	progressJSON, err := json.Marshal(progress)
	if err != nil {
		return fmt.Errorf("failed to marshal progress: %w", err)
	}

	channel := fmt.Sprintf("resume:progress:%s", uploadID)

	result := q.client.Publish(ctx, channel, string(progressJSON))
	if err := result.Err(); err != nil {
		return fmt.Errorf("failed to publish progress: %w", err)
	}

	numSubscribers := result.Val()
	log.Printf("[ResumeUploadQueue] Progress published to channel %s: status=%s, progress=%d%%, subscribers=%d",
		channel, progress.Status, progress.Progress, numSubscribers)

	return nil
}
