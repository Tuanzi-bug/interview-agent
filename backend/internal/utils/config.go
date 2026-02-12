package utils

import (
	"log"
	"time"
)

func ParseDurationWithDefault(value string, defaultValue time.Duration, fieldName string) time.Duration {
	duration, err := time.ParseDuration(value)
	if err != nil {
		log.Printf("Warning: Invalid %s in config, using default %v: %v", fieldName, defaultValue, err)
		return defaultValue
	}
	return duration
}
