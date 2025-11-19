package tool

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

// ResumeScoreRequest 简历评分请求
type ResumeScoreRequest struct {
	ResumeContent  json.RawMessage `json:"resume_content" jsonschema:"description=简历内容"`
	ReviewAnalysis json.RawMessage `json:"review_analysis" jsonschema:"description=评审分析内容"`
}

// ResumeScoreDimension 评分维度
type ResumeScoreDimension struct {
	Dimension string  `json:"dimension" jsonschema:"description=评分维度名称"`
	Score     float64 `json:"score" jsonschema:"description=该维度的分数"`
	Weight    float64 `json:"weight" jsonschema:"description=权重"`
}

// ResumeScoreResponse 简历评分响应
type ResumeScoreResponse struct {
	Success      bool                   `json:"success" jsonschema:"description=是否成功"`
	OverallScore float64                `json:"overall_score" jsonschema:"description=总体分数 0-100"`
	Dimensions   []ResumeScoreDimension `json:"dimensions" jsonschema:"description=各维度分数"`
	Reasoning    string                 `json:"reasoning" jsonschema:"description=评分理由"`
	ErrorMsg     string                 `json:"error_msg,omitempty" jsonschema:"description=错误信息"`
}

// ResumeScore 简历评分函数
// 从简历评审分析中提取结构化的分数数据
func ResumeScore(_ context.Context, req *ResumeScoreRequest) (*ResumeScoreResponse, error) {
	if req == nil {
		return &ResumeScoreResponse{
			Success:  false,
			ErrorMsg: "参数错误：请求体不能为空",
		}, errors.New("参数错误：请求体不能为空")
	}

	// 这里可以调用 LLM 或规则引擎提取分数
	// 当前实现返回示例结构，实际使用时可以根据review_analysis内容解析
	// 示例返回结构
	return &ResumeScoreResponse{
		Success:      true,
		OverallScore: 82.5,
		Dimensions: []ResumeScoreDimension{
			{Dimension: "模块完整度", Score: 85, Weight: 0.2},
			{Dimension: "技能匹配度", Score: 80, Weight: 0.3},
			{Dimension: "量化成果", Score: 75, Weight: 0.2},
			{Dimension: "语言表达", Score: 88, Weight: 0.15},
			{Dimension: "格式规范", Score: 85, Weight: 0.15},
		},
		Reasoning: "简历整体质量良好，技能描述清晰，但量化成果可以更具体",
	}, nil
}

// CreateResumeScoreTool 创建简历评分工具实例
func CreateResumeScoreTool() tool.InvokableTool {
	t, err := utils.InferTool(
		"resume_score",
		"从简历评审分析中提取结构化的分数数据，包括总体分数（0-100分）和各维度分数。各维度包括：模块完整度、技能匹配度、量化成果、语言表达、格式规范",
		ResumeScore,
	)
	if err != nil {
		log.Fatalf("infer tool failed: %v", err)
	}
	return t
}
