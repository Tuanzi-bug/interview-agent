package tool

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/google/uuid"
)

// ContextKey 用于在 context 中传递值的类型
type ContextKey string

const (
	// MaxQuestionsKey 最大提问数量的 context key
	MaxQuestionsKey ContextKey = "max_questions"
	// QuestionCountKey 当前提问计数的 context key
	QuestionCountKey ContextKey = "question_count"
	// SessionIDKey session ID 的 context key
	SessionIDKey ContextKey = "session_id"
)

// 全局问题计数器（线程安全）
var (
	questionCounterMutex sync.Mutex
	questionCounters     = make(map[string]int) // 使用 session ID 作为 key
)

// GenerateSessionID 生成唯一的 session ID
func GenerateSessionID() string {
	return uuid.New().String()
}

type AskForClarificationInput struct {
	Question string `json:"question" jsonschema_description:"The specific question you want to ask the user to get the missing information"`
}

func askForInput(ctx context.Context, input *AskForClarificationInput, opts ...tool.Option) (string, error) {
	fmt.Printf("\nQuestion: %s\n", input.Question)

	// 检查是否达到最大提问数量限制
	maxQuestions, ok := ctx.Value(MaxQuestionsKey).(int)
	sessionID := "default"
	var questionCount int

	if ok && maxQuestions > 0 {
		// 从 context 中获取 session ID
		sid, ok := ctx.Value(SessionIDKey).(string)
		if ok {
			sessionID = sid
		}

		// 线程安全地更新计数器
		questionCounterMutex.Lock()
		questionCounters[sessionID]++
		questionCount = questionCounters[sessionID]
		questionCounterMutex.Unlock()

		// 显示当前进度
		fmt.Printf("[问题 %d/%d]\n", questionCount, maxQuestions)

		// 如果达到最大数量，提示这是最后一个问题
		if questionCount >= maxQuestions {
			fmt.Println("很好，下面是最后的问题")
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("\nyour input here: ")
	scanner.Scan()
	fmt.Println()
	nInput := scanner.Text()

	// 如果达到最大提问数量，在返回值中添加特殊标记
	if ok && maxQuestions > 0 && questionCount >= maxQuestions {
		// 返回特殊标记，告诉 agent 已达到限制
		// 清理计数器
		questionCounterMutex.Lock()
		delete(questionCounters, sessionID)
		questionCounterMutex.Unlock()
		return nInput + "\n[INTERVIEW_LIMIT_REACHED]", nil
	}

	return nInput, nil
}

func NewAskForInputTool() tool.InvokableTool {
	t, err := utils.InferOptionableTool(
		"ask_for_input",
		"等待用户输入,并返回用户的输入",
		askForInput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✅ Ask for input tool initialized")
	return t
}

type ExitInput struct {
	FinalResult string `json:"final_result" jsonschema_description:"The final result of the tool"`
}

func exit(ctx context.Context, input *ExitInput, opts ...tool.Option) (string, error) {
	fmt.Printf("\nFinalResult: %s\n", input.FinalResult)
	os.Exit(0)
	return input.FinalResult, nil
}

func NewExitTool() tool.InvokableTool {
	t, err := utils.InferOptionableTool(
		"exit",
		"退出当前的智能体",
		exit)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✅ Exit tool initialized")
	return t
}
