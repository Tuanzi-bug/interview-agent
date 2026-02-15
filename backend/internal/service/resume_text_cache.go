package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"ai-eino-interview-agent/internal/repository"
)

const (
	PDFTextCacheTTL       = 30 * 24 * time.Hour
	PDFTextCacheKeyPrefix = "resume:pdf_text:"
)

type ResumeTextCache struct{}

func NewResumeTextCache() *ResumeTextCache {
	return &ResumeTextCache{}
}

func (c *ResumeTextCache) GetCacheKey(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file for hashing: %w", err)
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("failed to hash file: %w", err)
	}

	hash := hex.EncodeToString(hasher.Sum(nil))
	return PDFTextCacheKeyPrefix + hash, nil
}

func (c *ResumeTextCache) Get(ctx context.Context, filePath string) (string, bool, error) {
	cacheKey, err := c.GetCacheKey(filePath)
	if err != nil {
		log.Printf("[ResumeTextCache] Failed to generate cache key: %v", err)
		return "", false, err
	}

	redisClient := repository.GetRedis()
	if redisClient == nil {
		log.Printf("[ResumeTextCache] Redis client not initialized")
		return "", false, nil
	}

	text, err := redisClient.Get(ctx, cacheKey).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			log.Printf("[ResumeTextCache] Cache miss for key: %s", cacheKey)
			return "", false, nil
		}
		log.Printf("[ResumeTextCache] Redis error: %v", err)
		return "", false, err
	}

	log.Printf("[ResumeTextCache] Cache HIT for key: %s (text length: %d bytes)", cacheKey, len(text))
	return text, true, nil
}

func (c *ResumeTextCache) Set(ctx context.Context, filePath, text string) error {
	cacheKey, err := c.GetCacheKey(filePath)
	if err != nil {
		return err
	}

	redisClient := repository.GetRedis()
	if redisClient == nil {
		return fmt.Errorf("redis client not initialized")
	}

	err = redisClient.Set(ctx, cacheKey, text, PDFTextCacheTTL).Err()
	if err != nil {
		log.Printf("[ResumeTextCache] Failed to set cache: %v", err)
		return err
	}

	log.Printf("[ResumeTextCache] Cached text for key: %s (length: %d bytes, TTL: %v)",
		cacheKey, len(text), PDFTextCacheTTL)
	return nil
}

func (c *ResumeTextCache) Delete(ctx context.Context, filePath string) error {
	cacheKey, err := c.GetCacheKey(filePath)
	if err != nil {
		return err
	}

	redisClient := repository.GetRedis()
	if redisClient == nil {
		return fmt.Errorf("redis client not initialized")
	}

	return redisClient.Del(ctx, cacheKey).Err()
}
