namespace go demo

struct CreateDemoModelRequest {
    1: required string name (api.form="name")
}

struct CreateDemoModelResponse {

}

// 服务定义
service DemoService {
    // 1. 创建用户模型
       CreateDemoModelResponse CreateDemoModel(1: CreateDemoModelRequest request) (
           api.post="/api/demo/create/model",
           api.category="demo",
           api.gen_path="demo"
       )
}
