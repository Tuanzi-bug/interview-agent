namespace go questionbank

// ==================== 1. 获取所有大类标签 ====================

// 获取大类标签请求
struct GetMainCategoriesRequest {
}

// 标签项
struct MainCategoryItem {
    1: required i64 id                                    // 对应数据库 ID
    2: required string category_name                       // 对应数据库 CategoryName
    3: optional string image_url                          // 对应数据库 ImageURL
    4: required i64 created_time                          // 对应数据库 CreatedTime（时间戳）
    5: required i64 updated_time                          // 对应数据库 UpdatedTime（时间戳）
}

// 获取大类标签响应
struct GetMainCategoriesResponse {
    1: required list<MainCategoryItem> categories
}

// ==================== 2. 根据大类获取所有小分支 ====================

// 获取小分支请求
struct GetSubCategoriesRequest {
    1: required i64 parent_id (api.query="parent_id")    // 父级ID
}

// 小分支标签项
struct SubCategoryItem {
    1: required i64 id                                    // 对应数据库 ID
    2: required string category_name                       // 对应数据库 CategoryName
    3: optional string image_url                          // 对应数据库 ImageURL
    4: required i64 parent_id                             // 对应数据库 ParentID
    5: required i64 created_time                          // 对应数据库 CreatedTime（时间戳）
    6: required i64 updated_time                          // 对应数据库 UpdatedTime（时间戳）
}

// 获取小分支响应
struct GetSubCategoriesResponse {
    1: required i64 parent_id                              // 父级ID
    2: required string parent_name                         // 父级名称
    3: required list<SubCategoryItem> categories
}

// ==================== 服务定义 ====================

service QuestionBankService {
    // 1. 获取所有大类标签
    GetMainCategoriesResponse GetMainCategories(1: GetMainCategoriesRequest request) (
        api.get="/api/questionbank/categories/main",
        api.category="questionbank",
        api.gen_path="questionbank"
    )

    // 2. 根据大类ID获取所有小分支
    GetSubCategoriesResponse GetSubCategories(1: GetSubCategoriesRequest request) (
        api.get="/api/questionbank/categories/sub",
        api.category="questionbank",
        api.gen_path="questionbank"
    )
}
