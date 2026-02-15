package interview

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"regexp"
	"time"

	"ai-eino-interview-agent/api/handler/interview/mianshi"
	"ai-eino-interview-agent/api/response"
	"ai-eino-interview-agent/internal/mq"
	"ai-eino-interview-agent/internal/repository"

	"github.com/cloudwego/hertz/pkg/app"
)

// UUID validation regex (RFC 4122)
var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// isValidUUID validates if the string is a valid UUID format
func isValidUUID(uuid string) bool {
	return uuidRegex.MatchString(uuid)
}

// GetUploadProgress streams real-time progress updates for a resume upload via SSE
// @router /api/v1/resume/upload/progress/:uploadID [GET]
func GetUploadProgress(ctx context.Context, c *app.RequestContext) {
	uploadID := c.Param("uploadID")

	// Validate uploadID format (UUID)
	if !isValidUUID(uploadID) {
		response.BadRequest(ctx, c, "Invalid upload ID format. Expected UUID format.")
		return
	}

	log.Printf("[GetUploadProgress] Starting SSE stream for uploadID: %s", uploadID)

	// Setup SSE headers (reuse existing function)
	mianshi.SetupSSEResponse(c)

	// Create pipe for streaming
	pipeReader, pipeWriter := io.Pipe()
	c.SetBodyStream(pipeReader, -1)

	// Start streaming goroutine
	go func() {
		defer func() {
			if err := pipeWriter.Close(); err != nil {
				log.Printf("[GetUploadProgress] Error closing pipe writer: %v", err)
			}
			log.Printf("[GetUploadProgress] SSE stream closed for uploadID: %s", uploadID)
		}()

		// Create timeout context (5 minutes)
		timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()

		// Subscribe to Redis progress channel
		redisClient := repository.GetRedis()
		if redisClient == nil {
			log.Printf("[GetUploadProgress] Redis client not available")
			_ = mianshi.SendSSEEvent(pipeWriter, map[string]interface{}{
				"type":    "error",
				"message": "Progress tracking unavailable",
			})
			return
		}

		channel := fmt.Sprintf("resume:progress:%s", uploadID)
		pubsub := redisClient.Subscribe(timeoutCtx, channel)
		defer func() {
			if err := pubsub.Close(); err != nil {
				log.Printf("[GetUploadProgress] Error closing Redis subscription: %v", err)
			}
		}()

		log.Printf("[GetUploadProgress] Subscribed to Redis channel: %s", channel)

		// Send initial connection event
		if err := mianshi.SendSSEEvent(pipeWriter, map[string]interface{}{
			"type":      "connected",
			"message":   "Progress stream connected",
			"upload_id": uploadID,
		}); err != nil {
			log.Printf("[GetUploadProgress] Error sending connected event: %v", err)
			return
		}

		// Stream progress updates
		progressComplete := false
		for {
			select {
			case msg, ok := <-pubsub.Channel():
				if !ok {
					log.Printf("[GetUploadProgress] Channel closed for uploadID: %s", uploadID)
					return
				}

				// Parse progress update
				var progress mq.ProgressUpdate
				if err := json.Unmarshal([]byte(msg.Payload), &progress); err != nil {
					log.Printf("[GetUploadProgress] Error parsing progress: %v", err)
					continue
				}

				log.Printf("[GetUploadProgress] Received progress: status=%s, progress=%d%%, stage=%s",
					progress.Status, progress.Progress, progress.Stage)

				// Send progress event
				event := map[string]interface{}{
					"type":      "progress",
					"status":    progress.Status,
					"progress":  progress.Progress,
					"stage":     progress.Stage,
					"upload_id": uploadID,
				}

				if progress.ErrorMsg != "" {
					event["error"] = progress.ErrorMsg
				}

				if err := mianshi.SendSSEEvent(pipeWriter, event); err != nil {
					log.Printf("[GetUploadProgress] Error sending progress event: %v", err)
					return
				}

				// Check if upload is complete (100% or error status)
				if progress.Progress >= 100 || progress.Status == "completed" || progress.Status == "failed" {
					progressComplete = true
					log.Printf("[GetUploadProgress] Upload finished with status: %s", progress.Status)

					// Send completion event
					_ = mianshi.SendSSEEvent(pipeWriter, map[string]interface{}{
						"type":      "complete",
						"message":   "Progress stream completed",
						"status":    progress.Status,
						"upload_id": uploadID,
					})
					return
				}

			case <-timeoutCtx.Done():
				// Timeout - close connection
				log.Printf("[GetUploadProgress] Timeout after 5 minutes for uploadID: %s, progressComplete=%v",
					uploadID, progressComplete)

				// Send timeout event
				_ = mianshi.SendSSEEvent(pipeWriter, map[string]interface{}{
					"type":      "timeout",
					"message":   "Progress stream timeout",
					"upload_id": uploadID,
				})
				return
			}
		}
	}()
}
