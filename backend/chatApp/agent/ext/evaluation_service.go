package ext

import (
	interviewsapi "ai-eino-interview-agent/api/model/interviews"
	"ai-eino-interview-agent/chatApp/agent/evaluation"
	"ai-eino-interview-agent/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

// GenerateInterviewEvaluation 调用评估智能体生成面试评估
// 返回评估响应数据，并将评估结果保存到数据库
func GenerateInterviewEvaluation(ctx context.Context, userId uint, reportId uint64) (*interviewsapi.GetInterviewEvaluationResponse, error) {
	// 添加 120 秒超时
	timeoutCtx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	// 创建评估智能体
	agent := evaluation.NewEvaluationAgent(userId)

	// 创建 runner
	runner := adk.NewRunner(timeoutCtx, adk.RunnerConfig{
		Agent: agent,
	})

	// 构建查询消息
	query := fmt.Sprintf(`请对用户ID为 %d、报告ID为 %d 的面试进行全面评估。

请按照以下步骤进行：
1. 首先调用 get_interviews_data 工具获取面试的完整问题和对话记录
2. 仔细分析每个问题的回答质量
3. 对每个评估维度进行详细分析和评分
4. 使用 score_extraction 工具生成结构化的评分数据
5. 生成一份专业的评估报告

评估报告应包含：
- 总体评分
- 各维度评分和分析
- 候选人的优势
- 需要改进的方面
- 总体建议`, userId, reportId)

	// 创建用户消息
	userMsg := &schema.Message{
		Role:    schema.User,
		Content: query,
	}

	messages := []adk.Message{
		userMsg,
	}

	// 运行智能体
	iter := runner.Run(timeoutCtx, messages)

	var lastMessage string
	for {
		select {
		case <-timeoutCtx.Done():
			log.Printf("[GenerateInterviewEvaluation] 超时：等待智能体响应超过 120 秒")
			return nil, fmt.Errorf("timeout waiting for evaluation (120s)")
		default:
		}

		event, ok := iter.Next()
		if !ok {
			break
		}

		if event.Err != nil {
			log.Printf("[GenerateInterviewEvaluation] 错误: %v", event.Err)
			return nil, fmt.Errorf("error during evaluation: %w", event.Err)
		}

		// 收集最后一条消息
		if event.Output != nil && event.Output.MessageOutput != nil {
			lastMessage = event.Output.MessageOutput.Message.Content
		}
	}

	// 构建评估响应
	response := buildEvaluationResponse(lastMessage)

	// 保存评估数据到数据库
	if err := saveEvaluationToDatabase(ctx, userId, reportId, response); err != nil {
		log.Printf("Warning: Failed to save evaluation: %v", err)
	}

	return response, nil
}

// buildEvaluationResponse 从智能体响应构建评估响应
// 直接反序列化智能体返回的 JSON
func buildEvaluationResponse(agentResponse string) *interviewsapi.GetInterviewEvaluationResponse {
	response := &interviewsapi.GetInterviewEvaluationResponse{
		Comment:    "",
		Dimensions: make([]*interviewsapi.EvaluationDimension, 0),
	}

	// 尝试直接解析 JSON
	if err := json.Unmarshal([]byte(agentResponse), response); err != nil {
		// 尝试从文本中提取 JSON
		jsonStr := extractJSONFromResponse(agentResponse)
		if jsonStr == "" {
			log.Printf("[buildEvaluationResponse] 无法提取 JSON，使用默认响应")
			return buildDefaultResponse()
		}

		// 尝试解析提取的 JSON
		if err := json.Unmarshal([]byte(jsonStr), response); err != nil {
			log.Printf("[buildEvaluationResponse] 解析提取的 JSON 失败: %v", err)
			return buildDefaultResponse()
		}
	}

	return response
}

// extractJSONFromResponse 从文本中提取 JSON 字符串
func extractJSONFromResponse(text string) string {
	// 查找对象格式 {...}
	start := -1
	braceCount := 0

	for i := 0; i < len(text); i++ {
		if text[i] == '{' {
			if start == -1 {
				start = i
			}
			braceCount++
		} else if text[i] == '}' {
			braceCount--
			if start != -1 && braceCount == 0 {
				jsonStr := text[start : i+1]
				// 尝试验证 JSON 是否有效
				var temp interface{}
				if err := json.Unmarshal([]byte(jsonStr), &temp); err == nil {
					return jsonStr
				}
			}
		}
	}

	return ""
}

// buildDefaultResponse 构建默认响应
func buildDefaultResponse() *interviewsapi.GetInterviewEvaluationResponse {
	return &interviewsapi.GetInterviewEvaluationResponse{
		Comment: "评估处理失败，无法解析智能体响应",
		Dimensions: []*interviewsapi.EvaluationDimension{
			{DimensionName: "专业领域", Evaluation: "无法评估", Score: 0},
			{DimensionName: "项目经历", Evaluation: "无法评估", Score: 0},
			{DimensionName: "技术深度", Evaluation: "无法评估", Score: 0},
			{DimensionName: "技术基础", Evaluation: "无法评估", Score: 0},
			{DimensionName: "团队协作", Evaluation: "无法评估", Score: 0},
			{DimensionName: "系统架构设计", Evaluation: "无法评估", Score: 0},
		},
	}
}

// saveEvaluationToDatabase 将评估数据保存到数据库
func saveEvaluationToDatabase(ctx context.Context, userId uint, reportId uint64, response *interviewsapi.GetInterviewEvaluationResponse) error {
	// 将维度数据转换为 []*model.EvaluationDimension
	var dimensionList []*model.EvaluationDimension
	for _, dim := range response.Dimensions {
		dimensionList = append(dimensionList, &model.EvaluationDimension{
			DimensionName: dim.DimensionName,
			Evaluation:    dim.Evaluation,
			Score:         float64(dim.Score),
		})
	}

	// 计算总体评分（各维度评分的平均值）
	var totalScore float64
	if len(response.Dimensions) > 0 {
		for _, dim := range response.Dimensions {
			totalScore += float64(dim.Score)
		}
		totalScore = totalScore / float64(len(response.Dimensions))
	}

	// 创建评估记录
	evaluation := &model.InterviewEvaluation{
		UserID:     userId,
		ReportID:   reportId,
		Comment:    response.Comment,
		Score:      totalScore,
		Dimensions: dimensionList,
		Deleted:    0,
	}

	// 直接调用 DAO 方法保存到数据库
	err := model.InterviewEvaluationDao.CreateEvaluation(evaluation)
	if err != nil {
		log.Printf("[saveEvaluationToDatabase] 保存评估失败: %v", err)
		return fmt.Errorf("failed to save evaluation: %w", err)
	}

	return nil
}
