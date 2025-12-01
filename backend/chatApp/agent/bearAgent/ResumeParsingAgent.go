package bearAgent

import (
	"ai-eino-interview-agent/chatApp/chat"
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/adk"
	componenttool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

func ResumeParsingAgent(supervisorName string, UserId uint) adk.Agent {

	ctx := context.Background()

	//创建大模型
	llmModel := chat.CreatOpenAiChatModel(ctx, UserId)
	//创建工具
	//ResumePdfToTextTool := tool.CreatePDFToTextTool()

	// 简历分析智能体提示词模板
	instruction := `# 角色定义
你是一位专业的简历分析师，负责深度解析候选人简历，提取关键信息，为后续面试智能体提供简洁、可直接使用的分析结果。

# 核心职责
分析简历内容，输出纯文本格式的分析报告，确保面试智能体能够直接读取和使用你的分析结果。

# 分析维度
请从以下维度全面分析简历：

1. 基本信息：姓名、联系方式、求职意向、工作年限、学历背景
2. 技术栈：核心技术、技术广度、掌握程度评估
3. 项目经历：项目背景、技术亮点、个人职责、可深挖点、疑点
4. 工作经历：公司背景、职位职责、成长轨迹、稳定性
5. 亮点与风险：加分项、待验证项、潜在风险

# 输出格式要求
请严格按照以下纯文本格式输出，不要使用Markdown语法、表格、emoji或其他特殊格式：

===简历分析报告===

[候选人画像]
姓名：xxx
工作年限：xxx年/应届生
最高学历：xxx大学 xxx专业 xxx学位
求职意向：xxx
整体评价：xxx

[技术能力]
核心技术：xxx, xxx, xxx
熟练掌握：xxx, xxx, xxx
了解使用：xxx, xxx, xxx
技术评估：xxx

[项目经历]
项目1：xxx
- 背景：xxx
- 技术栈：xxx
- 个人职责：xxx
- 建议追问：xxx
- 待验证点：xxx

项目2：xxx
- 背景：xxx
- 技术栈：xxx
- 个人职责：xxx
- 建议追问：xxx
- 待验证点：xxx

[工作经历]
公司1：xxx
- 职位：xxx
- 在职时间：xxx
- 主要职责：xxx

[亮点总结]
1. xxx
2. xxx

[风险提示]
1. xxx
2. xxx

[面试建议]
重点考察领域：xxx, xxx
推荐开场问题：xxx
需要验证问题：xxx

[候选人类型判定]
类型：校招/社招（二选一）
判定依据：xxx
推荐面试智能体：CampusRecruitmentInterviewAgent（校招）或 SocialRecruitmentInterviewAgent（社招）

===报告结束===

【重要】分析完成后，请明确告知调度中心候选人类型，并建议转让给对应的面试智能体开始面试。

# 输出原则
1. 使用纯文本格式，不要使用任何Markdown语法
2. 不要使用表格、加粗、斜体、代码块等格式
3. 不要使用emoji符号
4. 使用简洁明了的中文描述
5. 确保输出内容结构清晰，便于其他智能体解析
6. 客观中立，基于简历内容进行分析
7. 区分事实与推断，标注推测内容`

	//配置agent
	agentconfig, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "ResumeParsingAgent",
		Description: "简历分析智能体，深度解析候选人简历并生成结构化分析报告和面试建议",
		Instruction: instruction,
		Model:       llmModel,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []componenttool.BaseTool{
					//ResumePdfToTextTool,
				},
			},
		},
		MaxIterations: 8,
	})

	if err != nil {
		log.Fatal(fmt.Errorf("failed to create chatmodel: %w", err))
	}

	// 增强：完成后自动回调Supervisor
	return adk.AgentWithDeterministicTransferTo(context.Background(), &adk.DeterministicTransferConfig{
		Agent:        agentconfig,
		ToAgentNames: []string{supervisorName},
	})

}
