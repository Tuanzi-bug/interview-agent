package worker

import (
	"ai-eino-interview-agent/internal/model"
	"ai-eino-interview-agent/internal/mq"
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type ParseResumeFunc func(ctx context.Context, userID uint, filePath string, fileSize int64) (uint64, interface{}, error)

type WorkerPool struct {
	queue           *mq.ResumeUploadQueue
	workerCount     int
	parseResumeFunc ParseResumeFunc
	ctx             context.Context
	cancel          context.CancelFunc
	wg              sync.WaitGroup
	stopOnce        sync.Once
}

func NewWorkerPool(queue *mq.ResumeUploadQueue, workerCount int, parseResumeFunc ParseResumeFunc) (*WorkerPool, error) {
	if workerCount <= 0 {
		return nil, fmt.Errorf("worker count must be positive, got: %d", workerCount)
	}

	if queue == nil {
		return nil, fmt.Errorf("queue cannot be nil")
	}

	if parseResumeFunc == nil {
		return nil, fmt.Errorf("parseResumeFunc cannot be nil")
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &WorkerPool{
		queue:           queue,
		workerCount:     workerCount,
		parseResumeFunc: parseResumeFunc,
		ctx:             ctx,
		cancel:          cancel,
	}, nil
}

func (wp *WorkerPool) Start() error {
	log.Printf("[WorkerPool] Starting %d workers", wp.workerCount)

	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}

	log.Printf("[WorkerPool] All %d workers started", wp.workerCount)
	return nil
}

func (wp *WorkerPool) Stop() error {
	var stopErr error

	wp.stopOnce.Do(func() {
		log.Printf("[WorkerPool] Initiating graceful shutdown")

		wp.cancel()

		done := make(chan struct{})
		go func() {
			wp.wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			log.Printf("[WorkerPool] All workers stopped gracefully")
		case <-time.After(30 * time.Second):
			log.Printf("[WorkerPool] Shutdown timeout - some workers may not have exited")
			stopErr = fmt.Errorf("shutdown timeout after 30 seconds")
		}
	})

	return stopErr
}

func (wp *WorkerPool) worker(workerID int) {
	defer wp.wg.Done()
	log.Printf("[WorkerPool] Worker %d started", workerID)

	for {
		select {
		case <-wp.ctx.Done():
			log.Printf("[WorkerPool] Worker %d shutting down", workerID)
			return
		default:
			job, err := wp.queue.DequeueJob(wp.ctx, 5*time.Second)
			if err != nil {
				log.Printf("[WorkerPool] Worker %d dequeue error: %v", workerID, err)
				continue
			}

			if job == nil {
				continue
			}

			if err := wp.processJob(wp.ctx, job); err != nil {
				log.Printf("[WorkerPool] Worker %d job processing error: %v", workerID, err)
			}
		}
	}
}

func (wp *WorkerPool) processJob(ctx context.Context, job *mq.ResumeUploadJob) error {
	log.Printf("[WorkerPool] Processing job: uploadID=%s, userID=%d, filePath=%s",
		job.UploadID, job.UserID, job.FilePath)

	wp.publishProgress(ctx, job.UploadID, &mq.ProgressUpdate{
		Status:   "pending",
		Progress: 0,
		Stage:    "Job received, waiting to process",
	})

	wp.updateDatabase(job.UploadID, map[string]interface{}{
		"status":   "pending",
		"progress": 0,
		"stage":    "Job received, waiting to process",
	})

	fileInfo, err := os.Stat(job.FilePath)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get file info: %v", err)
		wp.publishProgress(ctx, job.UploadID, &mq.ProgressUpdate{
			Status:   "failed",
			Progress: 0,
			Stage:    "Error: " + errMsg,
			ErrorMsg: errMsg,
		})

		wp.updateDatabase(job.UploadID, map[string]interface{}{
			"status":    "failed",
			"error_msg": errMsg,
		})

		return fmt.Errorf("failed to get file info: %w", err)
	}

	fileSize := fileInfo.Size()

	wp.publishProgress(ctx, job.UploadID, &mq.ProgressUpdate{
		Status:   "extracting",
		Progress: 25,
		Stage:    "Extracting text from PDF",
	})

	wp.updateDatabase(job.UploadID, map[string]interface{}{
		"status":   "extracting",
		"progress": 25,
		"stage":    "Extracting text from PDF",
	})

	extractStart := time.Now()

	wp.publishProgress(ctx, job.UploadID, &mq.ProgressUpdate{
		Status:   "extracted",
		Progress: 50,
		Stage:    "Text extraction complete, starting AI analysis",
	})

	extractDuration := time.Since(extractStart).Milliseconds()

	wp.updateDatabase(job.UploadID, map[string]interface{}{
		"status":           "extracted",
		"progress":         50,
		"stage":            "Text extraction complete, starting AI analysis",
		"extract_duration": extractDuration,
	})

	wp.publishProgress(ctx, job.UploadID, &mq.ProgressUpdate{
		Status:   "analyzing",
		Progress: 75,
		Stage:    "Analyzing resume with AI",
	})

	wp.updateDatabase(job.UploadID, map[string]interface{}{
		"status":   "analyzing",
		"progress": 75,
		"stage":    "Analyzing resume with AI",
	})

	analyzeStart := time.Now()

	dbResumeID, _, err := wp.parseResumeFunc(ctx, job.UserID, job.FilePath, fileSize)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to parse resume: %v", err)
		wp.publishProgress(ctx, job.UploadID, &mq.ProgressUpdate{
			Status:   "failed",
			Progress: 75,
			Stage:    "Error: " + errMsg,
			ErrorMsg: errMsg,
		})

		wp.updateDatabase(job.UploadID, map[string]interface{}{
			"status":    "failed",
			"error_msg": errMsg,
		})

		return fmt.Errorf("failed to process resume upload job: %w", err)
	}

	analyzeDuration := time.Since(analyzeStart).Milliseconds()

	wp.publishProgress(ctx, job.UploadID, &mq.ProgressUpdate{
		Status:   "completed",
		Progress: 100,
		Stage:    "Resume analysis complete",
	})

	wp.updateDatabase(job.UploadID, map[string]interface{}{
		"status":           "completed",
		"progress":         100,
		"stage":            "Resume analysis complete",
		"resume_id":        dbResumeID,
		"analyze_duration": analyzeDuration,
	})

	log.Printf("[WorkerPool] Job completed successfully: uploadID=%s, resumeID=%d",
		job.UploadID, dbResumeID)

	return nil
}

func (wp *WorkerPool) publishProgress(ctx context.Context, uploadID string, progress *mq.ProgressUpdate) {
	if err := wp.queue.PublishProgress(ctx, uploadID, progress); err != nil {
		log.Printf("[WorkerPool] Failed to publish progress for uploadID=%s: %v", uploadID, err)
	}
}

func (wp *WorkerPool) updateDatabase(uploadID string, updates map[string]interface{}) {
	if err := model.ResumeUploadStatusDao.UpdateStatus(uploadID, updates); err != nil {
		log.Printf("[WorkerPool] Failed to update database for uploadID=%s: %v", uploadID, err)
	}
}
