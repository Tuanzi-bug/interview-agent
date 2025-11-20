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
    7: optional string report          // 最终报告（当 Type 为 done 时）
    8: optional double score           // 面试评分（当 Type 为 done 时）
    9: optional i64 duration           // 面试时长（秒）（当 Type 为 done 时）
    10: optional string feedback       // 反馈信息（当 Type 为 done 时）
    11: optional string messages       // 对话历史（JSON 格式）（当 Type 为 done 时）
}

// 面试记录 DTO（对应 interview_record 表）
struct InterviewRecordDTO {
    1: required i64   id                    // 记录ID
    2: required i32   user_id               // 用户ID
    3: required string title                // 面试标题
    4: required string type                 // 面试类型(综合面试、专项面试)
    5: required string difficulty           // 面试难度(简单、中等、困难)
    6: required string domain               // 面试领域(校招、社招；java、golang)
    7: optional string company_name         // 公司名称
    8: optional string position_name        // 岗位名称
    9: optional string interview_duration   // 面试时长
    10: required string status              // 面试状态(pending/completed)
    11: optional i64   duration             // 面试耗时（秒）
    12: optional i64   created_at           // 创建时间（毫秒时间戳）
    13: optional i64   updated_at           // 更新时间（毫秒时间戳）
}

// ==================== 请求和响应结构 ====================

// 启动面试请求
struct StartInterviewRequest {
    1: required string title (api.body="title")  // 面试标题
    2: required string query (api.body="query")  // 用户输入的查询
    3: required string type (api.body="type")  // 面试类型(综合面试、专项面试)
    4: required string domain (api.body="domain")  // 面试领域（综合面试对应：校招、社招；专项面试对应java、golang等)
    5: required string difficulty (api.body="difficulty")  // 难度级别（简单、中等、困难）
    6: optional string company_name (api.body="company_name")  // 公司名称（专项面试不用填写）
    7: optional string position_name (api.body="position_name")  // 岗位名称（专项面试不用填写）
}

// 启动面试响应
struct StartInterviewResponse {
    1: required InterviewResult result
}

// 继续面试请求
struct ContinueInterviewRequest {
    1: required string query (api.body="query")  // 用户输入的查询
}

// 获取面试记录列表请求
struct ListInterviewRecordsRequest {
    1: optional i32 page      (api.query="page")       // 页码，默认 1
    2: optional i32 page_size (api.query="page_size")  // 每页数量，默认 10
}

// 获取面试记录列表响应
struct ListInterviewRecordsResponse {
    1: required list<InterviewRecordDTO> records   // 面试记录列表
    2: required i64 total                          // 总条数
    3: required i32 page                           // 当前页码
    4: required i32 page_size                      // 每页数量
}

// 获取单个面试记录详情请求
struct GetInterviewRecordRequest {
    1: required i64 id (api.path="id")  // 面试记录ID
}

// 获取单个面试记录详情响应
struct GetInterviewRecordResponse {
    1: required InterviewRecordDTO record
}

// 提交面试答案请求
struct SubmitInterviewAnswerRequest {
    1: required string session_id (api.body="session_id")  // 会话ID
    2: required string answer     (api.body="answer")      // 用户的答案内容
    3: optional string action     (api.body="action")      // 操作类型：answer, continue, quit（默认为 answer）
}

// 提交面试答案响应
struct SubmitInterviewAnswerResponse {
    1: required string status      // 状态：received, error
    2: optional string message     // 消息说明
    3: optional string session_id  // 会话ID
}

// 评估维度
struct EvaluationDimension {
    1: required string dimension_name  // 维度名称（如：技术能力、沟通能力等）
    2: required string evaluation      // 该维度的评估内容
    3: required i32 score              // 该维度的评分（0-100）
}

// 获取面试评估请求
struct GetInterviewEvaluationRequest {
    1: required i64 report_id (api.query="report_id")  // 面试报告ID
}

// 获取面试评估响应
struct GetInterviewEvaluationResponse {
    1: required string comment                          // 整体评价
    2: required list<EvaluationDimension> dimensions   // 各维度评估列表
}

// 答题记录中的单条对话
struct AnswerRecordMessage {
    1: required i32 order       // 对话顺序
    2: required string question // 提问内容
    3: required string answer   // 回答内容
}

// 答题记录中的评论信息
struct AnswerRecordComment {
    1: required i32 score           // 评分
    2: required string key_points   // 关键点
    3: required string difficulty   // 难度等级
    4: required string strengths    // 优势
    5: required string weaknesses   // 不足
    6: required string suggestion   // 建议
    7: required string know_points  // 知识点
    8: required string thinking     // 思考过程
    9: required string reference    // 参考答案
}

// 单个答题记录
struct AnswerRecord {
    1: required i32 order                           // 问题顺序
    2: required string content                      // 问题内容
    3: required AnswerRecordComment comment         // 评论信息
    4: required list<AnswerRecordMessage> message   // 对话列表
}

// 获取答题记录请求
struct GetAnswerRecordRequest {
    1: required i64 report_id (api.query="report_id")  // 面试报告ID
}

// 获取答题记录响应
struct GetAnswerRecordResponse {
    1: required list<AnswerRecord> records  // 答题记录列表
}



// ==================== 服务定义 ====================

service InterviewsService {
    // 启动面试流程（流式）
    StartInterviewResponse StartInterviewStream(1: StartInterviewRequest request) (
        api.post="/api/interview/start/stream",
        api.category="interviews",
        api.gen_path="interviews"
    )

    // 提交面试回答
    SubmitInterviewAnswerResponse SubmitInterviewAnswer(1: SubmitInterviewAnswerRequest request) (
        api.post="/api/interview/submit/answer",
        api.category="interviews",
        api.gen_path="interviews"
    )

    // 获取面试评估
    GetInterviewEvaluationResponse GetInterviewEvaluation(1: GetInterviewEvaluationRequest request) (
        api.get="/api/interview/evaluation",
        api.category="interviews",
        api.gen_path="interviews"
    )

    // 获取答题记录
    GetAnswerRecordResponse GetAnswerRecord(1: GetAnswerRecordRequest request) (
        api.get="/api/interview/answer-record",
        api.category="interviews",
        api.gen_path="interviews"
    )

    // 获取面试记录列表
    ListInterviewRecordsResponse GetInterviewRecords(1: ListInterviewRecordsRequest request) (
        api.get="/api/interview/records",
        api.category="interviews",
        api.gen_path="interviews"
    )
}

