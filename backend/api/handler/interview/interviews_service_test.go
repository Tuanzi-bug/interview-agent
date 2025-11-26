package interview_test

import (
	"ai-eino-interview-agent/api/handler/interview"
	"context"
	"errors"
	"testing"
)

// TestIsRetryableError 验证错误分类逻辑是否符合预期
func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "network timeout",
			err:  errors.New("read tcp 1.2.3.4:1234: i/o timeout"),
			want: true,
		},
		{
			name: "connection refused",
			err:  errors.New("dial tcp 127.0.0.1:3306: connection refused"),
			want: true,
		},
		{
			name: "deadlock",
			err:  errors.New("deadlock found when trying to get lock"),
			want: true,
		},
		{
			name: "context canceled",
			err:  context.Canceled,
			want: false,
		},
		{
			name: "context deadline exceeded",
			err:  context.DeadlineExceeded,
			want: false,
		},
		{
			name: "business error",
			err:  errors.New("invalid input data"),
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := interview.IsRetryableError(tc.err)
			if got != tc.want {
				t.Fatalf("isRetryableError(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}
