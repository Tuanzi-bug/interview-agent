package impl

import "time"

const (
	// Redis Key前缀
	PromptKeyPrefix = "prompt:agent:"
)

// PromptData 提示词数据结构
type PromptData struct {
	AgentName   string    `json:"agent_name"`
	Instruction string    `json:"instruction"`
	Version     string    `json:"version"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	Status      string    `json:"status"` // active, inactive, deprecated
}
