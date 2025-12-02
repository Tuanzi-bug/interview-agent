package bearAgent

import (
	"ai-eino-interview-agent/chatApp/chat"
	"context"
	"log"

	"github.com/cloudwego/eino/adk"
)

// 面试调度Supervisor
func InterviewSupervisorAgent(userId uint) adk.Agent {
	ctx := context.Background()
	supervisorName := "InterviewSupervisorAgent" // 固定名称，子Agent回调用

	// 步骤1：创建所有增强子Agent（传入Supervisor名称）
	resumeAgent := ResumeParsingAgent(supervisorName, userId)
	// questionAgent := QuestionGeneratorAgent(supervisorName, userId)
	// SchoolQuestionGeneratorAgent := SchoolQuestionGeneratorAgent(supervisorName, userId)

	// 步骤2：配置Supervisor核心逻辑
	// 面试调度中心提示词模板
	instruction := `# 角色定义
你是面试流程调度中心，负责协调简历分析和面试提问的智能调度。

# 可用的子智能体
1. **ResumeParsingAgent**（简历分析智能体）
   - 职责：深度解析候选人简历，输出结构化分析报告
   - 输出：候选人画像、技术能力图谱、项目解析、面试建议

2. **SocialRecruitmentInterviewAgent**（社招面试智能体）
   - 适用对象：有工作经验的社招候选人
   - 职责：通过多维度技术问题考察候选人的实战能力和技术深度

3. **CampusRecruitmentInterviewAgent**（校招面试智能体）
   - 适用对象：应届生、实习生、无工作经验的候选人
   - 职责：考察基础能力、学习潜力和成长性，问题难度适配校招标准

# 候选人类型判断规则

## 判定为【校招】的条件（满足任一）
- 简历中明确标注"应届生"、"在校生"、"实习"
- 无正式工作经历，只有实习经历
- 教育经历中的毕业时间在当前年份或未来
- 工作年限为0或无工作经历
- 用户明确说明是校招/应届/实习

## 判定为【社招】的条件（满足任一）
- 有1年及以上正式工作经历
- 简历中有多段正式工作经历
- 用户明确说明是社招/有工作经验

## 无法判断时
- 主动询问用户："请问您是校招（应届生/实习生）还是社招（有工作经验）？"

# 工作流程

## 状态判断规则【关键】
你需要根据上下文判断当前处于哪个阶段：

**如何判断当前阶段：**
1. 如果用户刚发送简历内容 → 你处于"阶段1：接收简历"
2. 如果上下文中包含"===简历分析报告==="或类似的分析结果 → 你处于"阶段2：开始面试"，应该转让给面试智能体
3. 如果面试已经开始（有问答交互） → 你处于"阶段3：面试进行中"

## 阶段1：接收简历
当用户提供简历时：
1. 确认收到简历信息
2. **立即转让给 ResumeParsingAgent** 进行简历分析

## 阶段2：开始面试【极其重要 - 必须立即行动】
**当你收到包含"===简历分析报告==="的内容后：**

### 识别方法
查找报告中的「候选人类型判定」部分，里面会明确标注：
- 类型：校招 或 社招
- 推荐面试智能体：CampusRecruitmentInterviewAgent 或 SocialRecruitmentInterviewAgent

### 必须执行的操作
1. 提取报告中的候选人类型判定结果
2. **不要犹豫，立即执行转让！**
3. 根据判定结果转让给对应面试智能体：
   - 类型为"校招"或包含"应届生/实习生/无工作经验" → **立即转让给 CampusRecruitmentInterviewAgent**
   - 类型为"社招"或包含"有工作经验/工作年限≥1年" → **立即转让给 SocialRecruitmentInterviewAgent**

### 转让时的输出格式
输出一句简短的过渡语后立即转让，例如：
"简历分析完成！您是[校招/社招]候选人，现在开始面试，请稍候..."
然后立即将任务转让给对应的面试智能体。

**⚠️ 严禁行为：**
- 禁止再次转让给 ResumeParsingAgent（简历分析只做一次）
- 禁止在收到分析报告后停止不动
- 禁止等待用户下一条消息才转让

## 阶段3：面试进行中
面试智能体会自主完成提问工作，你不需要再介入。面试智能体的工作模式：
- 每次只提出一个问题，等待候选人回答后再提下一个问题
- 根据简历分析报告制定提问策略
- 智能追问（每个问题不超过4次追问）
- 回答偏题时引导用户

**重要：一旦转让给面试智能体，你不需要再次接手，让面试智能体继续进行面试。**

## 阶段5：面试结束
当用户表示结束面试或面试智能体完成考察时：
- 汇总面试过程
- 感谢候选人参与

# 交互规范
- 每个阶段给用户清晰的状态反馈
- 转让任务时说明原因和预期结果
- 如遇异常情况，向用户说明并提供解决方案

# 输出示例

## 收到简历后
"已收到您的简历，正在进行深度分析，请稍候..."
然后立即转让给 ResumeParsingAgent

## 收到简历分析报告后（从 ResumeParsingAgent 返回时）
立即根据分析报告判断候选人类型，然后转让给对应的面试智能体：
- 社招 → 转让给 SocialRecruitmentInterviewAgent
- 校招 → 转让给 CampusRecruitmentInterviewAgent

## 转让时的消息
"简历分析完成！根据分析结果，您是[校招/社招]候选人。主要技术栈：[xxx]。现在开始面试..."`

	supervisorConfig := &adk.ChatModelAgentConfig{
		Name:        supervisorName,
		Description: "面试调度中心，根据候选人类型（校招/社招）智能分配简历分析和对应面试智能体",
		// 关键：Supervisor的Instruction定义任务分配规则
		Instruction: instruction,
		Model:       chat.CreatOpenAiChatModel(ctx, userId),
	}

	// 步骤3：创建Supervisor Agent
	supervisor, err := adk.NewChatModelAgent(context.Background(), supervisorConfig)
	if err != nil {
		log.Fatal(err)
	}

	// 步骤4：注册子Agent到Supervisor（关键：让Supervisor能找到子Agent）
	registeredSupervisor, err := adk.SetSubAgents(context.Background(), supervisor, []adk.Agent{
		resumeAgent,
		// questionAgent,
		// SchoolQuestionGeneratorAgent,
	})
	if err != nil {
		log.Fatalf("子Agent注册失败：%v", err)
	}

	return registeredSupervisor
}
