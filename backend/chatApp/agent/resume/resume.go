package resume

import (
	"ai-eino-interview-agent/chatApp/chat"
	tool2 "ai-eino-interview-agent/chatApp/tool"
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/adk"
	componenttool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

// NewResumeParserAgent 创建简历解析智能体
// 用于解析简历内容，提取关键信息用于面试准备
func NewResumeParserAgent(userId uint) adk.Agent {
	ctx := context.Background()
	baseAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "ResumeParserAgent",
		Description: "一个专业的简历解析智能体，用于提取简历中的关键信息",
		Instruction: `你是一个专业的简历分析专家。你的任务是解析候选人的简历，提取关键信息用于面试准备。

任务说明：
1. 如果提供了PDF路径，使用pdf_to_text工具解析简历内容
2. 从简历中提取以下关键信息：
   - 基本信息（姓名、联系方式、工作年限等）
   - 教育背景（学校、专业、学位等）
   - 工作经历（公司、职位、工作时间、主要职责等）
   - 技术栈（编程语言、框架、工具等）
   - 项目经验（项目名称、技术栈、个人贡献等）
   - 技能特长（核心竞争力、专业技能等）
   - 证书资格（获得的证书、资格认证等）

3. 分析候选人的背景特点：
   - 主要技术方向
   - 行业经验
   - 职业发展轨迹
   - 核心竞争力

4. 生成面试建议：
   - 重点关注的技术领域
   - 可能的深入提问方向
   - 候选人的优势和潜在弱点
   - 推荐的面试难度级别

5. 只返回JSON，不要返回其他文本

必须返回的JSON格式：
{
  "basic_info": {
    "name": "姓名",
    "work_years": "工作年限（数字）",
    "contact": "联系方式"
  },
  "education": [
    {
      "school": "学校名称",
      "major": "专业",
      "degree": "学位",
      "graduation_year": "毕业年份"
    }
  ],
  "work_experience": [
    {
      "company": "公司名称",
      "position": "职位",
      "duration": "工作时间段",
      "responsibilities": "主要职责"
    }
  ],
  "tech_stack": ["技术1", "技术2", "技术3"],
  "projects": [
    {
      "name": "项目名称",
      "description": "项目描述",
      "tech_stack": ["技术1", "技术2"],
      "contribution": "个人贡献"
    }
  ],
  "skills": ["技能1", "技能2", "技能3"],
  "certifications": ["证书1", "证书2"],
  "strengths": "候选人的核心优势",
  "potential_weaknesses": "可能的弱点或不足",
  "recommended_difficulty": "推荐面试难度（初级/中级/高级）",
  "interview_focus_areas": ["重点关注领域1", "重点关注领域2"],
  "suggested_questions_directions": ["提问方向1", "提问方向2"]
}`,

		Model: chat.CreatOpenAiChatModel(ctx, userId),
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []componenttool.BaseTool{
					tool2.CreatePDFToTextTool(),
				},
			},
		},
		MaxIterations: 20,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create resume parser agent: %w", err))
	}
	return baseAgent
}
