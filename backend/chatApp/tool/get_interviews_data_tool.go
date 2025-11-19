package tool

import (
	"ai-eino-interview-agent/internal/model"
	"context"
	"log"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

// GetInterviewsDataRequest 获取面试数据的请求结构体
type GetInterviewsDataRequest struct {
	UserID   uint   `json:"user_id" jsonschema:"description=用户ID"`
	ReportID uint64 `json:"report_id" jsonschema:"description=报告ID"`
}

// GetInterviewsDataResponse 获取面试数据的响应结构体
type GetInterviewsDataResponse struct {
	Data []*model.QuestionTopicWithDialogues `json:"data" jsonschema:"description=面试问题和对话数据"`
}

// GetInterviewsData 获取面试数据
func GetInterviewsData(_ context.Context, req *GetInterviewsDataRequest) (*GetInterviewsDataResponse, error) {
	if req == nil {
		return nil, nil
	}
	data, err := model.InterviewQuestionTopicDao.GetWithDialoguesByUserIDAndReportID(req.UserID, req.ReportID)
	if err != nil {
		log.Printf("get db data failed: %v", err)
		return nil, err
	}
	return &GetInterviewsDataResponse{
		Data: data,
	}, nil
}

// GetInterviewsDataTool 创建获取面试数据的工具
func GetInterviewsDataTool() tool.InvokableTool {
	t, err := utils.InferTool(
		"get_interviews_data",
		"获取用户的面试问题和对话数据，包括问题主题和对应的对话记录",
		GetInterviewsData,
	)
	if err != nil {
		log.Fatalf("infer tool failed: %v", err)
	}
	return t
}
