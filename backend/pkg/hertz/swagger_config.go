package hertz

import (
	"context"
	"log"

	"github.com/cloudwego/hertz/pkg/app"
)

// SetupSwagger 配置并注册Swagger路由
func SetupSwagger() {
	// 获取Hertz服务器实例
	server := GetHertzServer()

	// 提供Swagger UI页面
	server.GET("/swagger", func(c context.Context, ctx *app.RequestContext) {
		// 简单的Swagger UI HTML页面
		swaggerHTML := `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@3/swagger-ui.css">
    <style>
        body { margin: 0; padding: 0; }
        .swagger-ui .topbar { display: none; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@3/swagger-ui-bundle.js"></script>
    <script>
        window.onload = function() {
            SwaggerUIBundle({
                url: "/swagger/doc.json",
                dom_id: '#swagger-ui',
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIBundle.SwaggerUIStandalonePreset
                ],
                layout: "BaseLayout",
                validatorUrl: null
            });
        };
    </script>
</body>
</html>
		`

		ctx.Header("Content-Type", "text/html")
		ctx.String(200, swaggerHTML)
	})

	// 提供API文档JSON
	server.GET("/swagger/doc.json", func(c context.Context, ctx *app.RequestContext) {
		// 这里应该提供完整的Swagger规范JSON
		// 为了简化，我们只提供一个基本的结构
		swaggerJSON := `{
  "swagger": "2.0",
  "info": {
    "title": "AI-Eino 智能面试系统 API文档",
    "description": "智能面试助手系统，提供用户管理、简历管理、面试管理等功能",
    "version": "1.0.0"
  },
  "host": "localhost:8000",
  "basePath": "/api/v1",
  "paths": {
    "/api/v1/interview": {
      "get": {
        "summary": "获取用户面试列表",
        "description": "获取当前用户的所有面试记录",
        "responses": {
          "200": { "description": "获取成功" },
          "401": { "description": "未授权访问" }
        }
      }
    }
  }
}`

		ctx.Header("Content-Type", "application/json")
		ctx.String(200, swaggerJSON)
	})

	// 输出Swagger访问地
	log.Printf("Swagger UI is available at: http://localhost:8888/swagger")
}
