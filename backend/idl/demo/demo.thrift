namespace go demo

//创建demo
struct CreateDemoRequest {
    1: required string name (api.form="name")
}

//响应
struct CreateDemoResponse {

}

//定义服务
service DemoService {
//    创建demo
    CreateDemoResponse CreateDemo(1:CreateDemoRequest request)(
       api.post="/api/demo/create",
       api.category="demo",
       api.gen_path="demo"
    )
}