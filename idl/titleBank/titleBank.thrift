namespace go titleBank

// ====================  创建面试题 ====================

// 创建面试题请求
struct CreateInterviewTitleRequest {
    1: required string title (api.form="title")
    2: required string type (api.form="type")
    3: required string domain (api.form="domain")
    4: required string Level (api.form="level")
}

// 创建面试题响应
struct CreateInterviewTitleResponse {
    1: required string state
}

// ====================  创建面试题标签 ====================

// 创建面试题标签请求
struct CreateInterviewLabelRequest {
    1: required string label (api.form="label")
    2: required i64 title_id (api.form="title_id")
}

// 创建面试题标签响应
struct CreateInterviewLabelResponse {
    1: required string state
}

// ====================  创建面试题解析 ====================

// 创建面试题解析请求
struct CreateInterviewParseRequest {
    1: required i64 title_id (api.form="title_id")
    2: required string parse (api.form="parse")
    3: optional string url (api.form="url")
}

// 创建面试题解析响应
struct CreateInterviewParseResponse {
    1: required string state
}

// ====================  创建面试题互动 ====================

// 创建面试题互动请求
struct CreateUserTitleInteractRequest {
    1: required i64 user_id (api.form="user_id")
    2: required i64 title_id (api.form="title_id")
    3: required i32 interact (api.form="interact")

}

// 创建面试题互动响应
struct CreateUserTitleInteractResponse {
    1: required string state
}

// 服务定义
service TitleBankService {
        // 1. 创建面试题请求
       CreateInterviewTitleResponse CreateInterviewTitle(1: CreateInterviewTitleRequest request) (
           api.post="/api/titleBank/create/title",
           api.category="titleBank",
           api.gen_path="titleBank"
       )
        // 2. 创建面试题标签
        CreateInterviewLabelResponse CreateInterviewLabel(1: CreateInterviewLabelRequest request) (
            api.post="/api/titleBank/create/label",
            api.category="titleBank",
            api.gen_path="titleBank"
        )
        // 3. 创建面试题解析
        CreateInterviewParseResponse CreateInterviewParse(1: CreateUserTitleInteractRequest request) (
            api.post="/api/titleBank/create/parse",
            api.category="titleBank",
            api.gen_path="titleBank"
        )
        // 4. 创建面试题互动
        CreateInterviewLabelResponse CreateUserTitleInteract(1: CreateUserTitleInteractResponse request) (
            api.post="/api/titleBank/create/interact",
            api.category="titleBank",
            api.gen_path="titleBank"
        )
}