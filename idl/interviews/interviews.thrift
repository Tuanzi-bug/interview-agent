namespace go interviews

// ==================== 数据结构定义 ====================

// 消息结构
struct Message {
    1: required string role        // user, assistant, system
    2: required string content     // 消息内容
    3: optional string agent       // 发送消息的 Agent 名称
}

// 面试结果
struct InterviewResult {
    1: required list<Message> messages    // 完整的对话历史
    2: optional string report              // 最终报告
    3: required string status              // 状态：resume_analysis, question_generation, answer_evaluation, report_generation, completed
    4: optional string current_agent       // 当前活跃的 Agent 名称
}

// 面试事件（用于流式响应）
struct InterviewEvent {
    1: required string type            // 事件类型：message, transfer, error, done
    2: optional string agent_name      // Agent 名称
    3: optional string message         // 消息内容（当 Type 为 message 时）
    4: optional string transfer_to     // 转让目标（当 Type 为 transfer 时）
    5: optional string error           // 错误信息（当 Type 为 error 时）
    6: optional string status          // 状态更新
}

// ==================== 请求和响应结构 ====================

// 启动面试请求
struct StartInterviewRequest {
    1: required string query (api.body="query")  // 用户输入的查询
}

// 启动面试响应
struct StartInterviewResponse {
    1: required InterviewResult result
}

// 继续面试请求
struct ContinueInterviewRequest {
    1: required string query (api.body="query")  // 用户输入的查询
}

// ==================== 服务定义 ====================

service InterviewsService {
    // 启动面试流程（流式）
    StartInterviewResponse StartInterviewStream(1: StartInterviewRequest request) (
        api.post="/api/interview/start/stream",
        api.category="interviews",
        api.gen_path="interviews"
    )

    // 继续面试流程（用于多轮对话）
    StartInterviewResponse ContinueInterview(1: ContinueInterviewRequest request) (
        api.post="/api/interview/continue",
        api.category="interviews",
        api.gen_path="interviews"
    )
}

