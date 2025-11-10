package tool

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

type AskForClarificationInput struct {
	Question string `json:"question" jsonschema_description:"The specific question you want to ask the user to get the missing information"`
}

func askForInput(ctx context.Context, input *AskForClarificationInput, opts ...tool.Option) (string, error) {
	fmt.Printf("\nQuestion: %s\n", input.Question)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("\nyour input here: ")
	scanner.Scan()
	fmt.Println()
	nInput := scanner.Text()
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
