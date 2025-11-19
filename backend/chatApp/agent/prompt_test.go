// 文件名：prompt_test.go（必须以 _test.go 结尾，放在 backend/chatApp/agent/ 下）
package agent

import (
	"context"
	// 这里填上面确认的正确导入路径，比如：
	"go-eino-interview-agent/backend/chatApp/prompt" // 或完整路径 "go-eino-interview-agent/backend/chatApp/prompt"
	"testing"
)

// 函数名必须以 Test 开头，参数是 *testing.T
func TestGetPromptInstruction(t *testing.T) {
	ctx := context.Background()
	// 调用目标函数
	instr := prompt.GetPromptInstruction(ctx, "ResumeReviewAgent")
	// 打印结果
	t.Logf("提示词函数返回结果：%s", instr)

	// 可选：判断结果是否合理（非空）
	if instr == "" {
		t.Warn("返回的提示词为空，可能是Redis没配置或默认模板缺失")
	}
}
