package tool

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

type AnswerEvalRequest struct {
	ResumeAnalysis json.RawMessage `json:"resume_analysis" jsonschema:"description=上游简历分析的原始JSON内容，可选"`
	HistoryQA      json.RawMessage `json:"history_qa" jsonschema:"description=历史问答记录的原始JSON内容，可选"`
}

type AnswerEvalQARecord struct {
	Question string `json:"question" jsonschema:"description=问题"`
	Answer   string `json:"answer" jsonschema:"description=答案"`
}

type AnswerEvalResponse struct {
	Success  bool   `json:"success" jsonschema:"description=是否成功"`
	Score    int    `json:"score" jsonschema:"description=评估结果"`
	Comment  string `json:"comment" jsonschema:"description=评估结果"`
	ErrorMsg string `json:"error_msg,omitempty" jsonschema:"description=错误信息"`
}

func AnswerEval(_ context.Context, req *AnswerEvalRequest) (*AnswerEvalResponse, error) {
	if req == nil {
		return &AnswerEvalResponse{
			Success:  false,
			ErrorMsg: "参数错误：请求体不能为空",
		}, errors.New("参数错误：请求体不能为空") // 返回错误信息
	}

	return &AnswerEvalResponse{
		Success: true,
		Score:   85,
		Comment: "候选人表现优秀，回答问题清晰，思路清晰，回答问题准确",
	}, nil
}

func CreateAnswerEvalTool() tool.InvokableTool {
	tool, err := utils.InferTool("answer_eval", "根据候选人的简历、问题和答案去评估候选人的面试表现 并给出评估结果", AnswerEval)
	if err != nil {
		log.Fatalf("infer tool failed: %v", err)
	}
	return tool
}
