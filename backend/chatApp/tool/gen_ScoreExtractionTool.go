package tool

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

// ScoreExtractionRequest 分数提取请求
type ScoreExtractionRequest struct {
	EvaluationContent json.RawMessage `json:"evaluation_content" jsonschema:"description=评估内容的原始JSON"`
	HistoryQA         json.RawMessage `json:"history_qa" jsonschema:"description=历史问答记录的原始JSON内容，可选"`
	ResumeAnalysis    json.RawMessage `json:"resume_analysis" jsonschema:"description=简历分析结果，可选"`
}

// ScoreDimension 评分维度
type ScoreDimension struct {
	Dimension string  `json:"dimension" jsonschema:"description=评分维度名称"`
	Score     float64 `json:"score" jsonschema:"description=该维度的分数"`
	Weight    float64 `json:"weight" jsonschema:"description=权重"`
}

// ScoreExtractionResponse 分数提取响应
type ScoreExtractionResponse struct {
	Success      bool             `json:"success" jsonschema:"description=是否成功"`
	OverallScore float64          `json:"overall_score" jsonschema:"description=总体分数 0-100"`
	Dimensions   []ScoreDimension `json:"dimensions" jsonschema:"description=各维度分数"`
	Reasoning    string           `json:"reasoning" jsonschema:"description=评分理由"`
	ErrorMsg     string           `json:"error_msg,omitempty" jsonschema:"description=错误信息"`
}

// ScoreExtraction 分数提取函数
func ScoreExtraction(_ context.Context, req *ScoreExtractionRequest) (*ScoreExtractionResponse, error) {
	if req == nil {
		return &ScoreExtractionResponse{
			Success:  false,
			ErrorMsg: "参数错误：请求体不能为空",
		}, errors.New("参数错误：请求体不能为空")
	}

	// 这里可以调用 LLM 或规则引擎提取分数
	// 示例返回结构
	return &ScoreExtractionResponse{
		Success:      true,
		OverallScore: 82.5,
		Dimensions: []ScoreDimension{
			{Dimension: "技术深度", Score: 85, Weight: 0.3},
			{Dimension: "沟通表达", Score: 80, Weight: 0.2},
			{Dimension: "问题解决", Score: 82, Weight: 0.3},
			{Dimension: "团队协作", Score: 81, Weight: 0.2},
		},
		Reasoning: "候选人在技术深度和问题解决能力上表现突出，沟通表达清晰，整体表现优秀",
	}, nil
}

// CreateScoreExtractionTool 创建工具实例
func CreateScoreExtractionTool() tool.InvokableTool {
	t, err := utils.InferTool(
		"score_extraction",
		"从面试评估内容中提取结构化的分数数据，包括总体分数和各维度分数",
		ScoreExtraction,
	)
	if err != nil {
		log.Fatalf("infer tool failed: %v", err)
	}
	return t
}
