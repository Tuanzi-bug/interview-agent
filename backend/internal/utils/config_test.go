package utils

import (
	"testing"
	"time"
)

func TestParseDurationWithDefault(t *testing.T) {
	tests := []struct {
		name         string
		value        string
		defaultValue time.Duration
		fieldName    string
		expected     time.Duration
	}{
		{
			name:         "valid duration",
			value:        "5m",
			defaultValue: 3 * time.Minute,
			fieldName:    "test_timeout",
			expected:     5 * time.Minute,
		},
		{
			name:         "invalid duration returns default",
			value:        "invalid",
			defaultValue: 3 * time.Minute,
			fieldName:    "test_timeout",
			expected:     3 * time.Minute,
		},
		{
			name:         "empty string returns default",
			value:        "",
			defaultValue: 10 * time.Second,
			fieldName:    "test_timeout",
			expected:     10 * time.Second,
		},
		{
			name:         "complex duration",
			value:        "1h30m45s",
			defaultValue: 1 * time.Hour,
			fieldName:    "test_timeout",
			expected:     time.Hour + 30*time.Minute + 45*time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseDurationWithDefault(tt.value, tt.defaultValue, tt.fieldName)
			if result != tt.expected {
				t.Errorf("ParseDurationWithDefault() = %v, want %v", result, tt.expected)
			}
		})
	}
}
