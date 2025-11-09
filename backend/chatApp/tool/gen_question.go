package tool

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

type questionInterval struct {
	start     int
	end       int
	direction string
}

var questionDirectionPlan = []questionInterval{
	{start: 1, end: 2, direction: "热身开场：引导候选人自我介绍、确认核心技能与当前求职重点"},
	{start: 3, end: 4, direction: "职业动机与规划：了解候选人选择机会的动机、职业路径与未来计划"},
	{start: 5, end: 6, direction: "核心项目亮点：聚焦最能体现技术深度或业务价值的代表性项目"},
	{start: 7, end: 8, direction: "技术架构与设计权衡：探讨系统设计理念、技术选型及关键的trade-off"},
	{start: 9, end: 10, direction: "性能与优化实践：挖掘候选人在性能调优、算法优化等方面的方法论"},
	{start: 11, end: 12, direction: "可靠性与质量保障：覆盖测试策略、监控告警、故障预案等保障体系"},
	{start: 13, end: 14, direction: "安全与合规意识：了解候选人在安全加固、隐私合规等领域的经验"},
	{start: 15, end: 16, direction: "跨团队协作与影响力：评估沟通协作、推动跨职能项目的能力"},
	{start: 17, end: 18, direction: "学习成长与创新：了解自我驱动学习、技术创新或知识沉淀的案例"},
	{start: 19, end: 20, direction: "复盘总结与文化契合：总结面试重点、检验价值观及与团队的契合度"},
}

// InterviewQuestionRequest 定义生成面试问题的入参
type InterviewQuestionRequest struct {
	ResumeAnalysis json.RawMessage   `json:"resume_analysis" jsonschema:"description=上游简历分析的原始JSON内容，可选"`
	HistoryQA      json.RawMessage   `json:"history_qa" jsonschema:"description=历史问答记录的原始JSON内容，可选"`
	QuestionIndex  int               `json:"question_index" jsonschema:"required,description=准备提出的问题序号，从1开始递增"`
	Meta           map[string]string `json:"meta,omitempty" jsonschema:"description=可选的额外上下文字段，例如面试阶段"`
}

// InterviewQARecord 历史问答
type InterviewQARecord struct {
	Question string `json:"question" jsonschema:"description=历史问题"`
	Answer   string `json:"answer" jsonschema:"description=候选人的回答"`
}

type InterviewQuestionResponse struct {
	Success   bool   `json:"success" jsonschema:"description=是否成功"`
	Direction string `json:"direction" jsonschema:"description=问题方向"`
	ErrorMsg  string `json:"error_msg,omitempty" jsonschema:"description=错误信息"`
}

// GenerateInterviewQuestion 根据题号返回给大模型的提问方向
func GenerateInterviewQuestion(_ context.Context, req *InterviewQuestionRequest) (InterviewQuestionResponse, error) {
	if req == nil {
		return InterviewQuestionResponse{
			Success:  false,
			ErrorMsg: "参数错误：请求体不能为空",
		}, errors.New("参数错误：请求体不能为空")
	}
	if req.QuestionIndex <= 0 {
		return InterviewQuestionResponse{
			Success:  false,
			ErrorMsg: "参数错误：question_index 必须从 1 开始",
		}, errors.New("参数错误：question_index 必须从 1 开始")
	}

	direction := directionByIndex(req.QuestionIndex)
	return InterviewQuestionResponse{
		Success:   true,
		Direction: direction,
	}, nil
}

func directionByIndex(index int) string {
	for _, interval := range questionDirectionPlan {
		if index >= interval.start && index <= interval.end {
			return interval.direction
		}
	}
	return "补充提问：根据候选人未覆盖的重点进行追问或给出收尾反馈"
}

// CreateInterviewQuestionTool 创建工具实例
func CreateGenQuestionTool() tool.InvokableTool {
	tool, err := utils.InferTool("gen_question", "根据提问序号返回对应的模拟面试问题方向", GenerateInterviewQuestion)
	if err != nil {
		log.Fatalf("infer tool failed: %v", err)
	}

	fmt.Println("✅  interview_question tool initialized")
	return tool
}
