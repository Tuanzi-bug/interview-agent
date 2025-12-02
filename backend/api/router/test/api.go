package test

import (
	testHandler "ai-eino-interview-agent/api/handler/test"

	"github.com/cloudwego/hertz/pkg/app/server"
)

// Register 注册测试相关的路由
func Register(r *server.Hertz) {
	root := r.Group("/")
	{
		api := root.Group("/api")
		{
			test := api.Group("/test")
			{
				// 健康检查
				test.GET("/health", testHandler.HealthCheck)

				// Agent测试
				agent := test.Group("/agent")
				{
					// 同步聊天接口 - 适合Postman测试
					agent.POST("/chat", testHandler.TestAgentChat)

					// 流式聊天接口 - SSE响应
					agent.POST("/chat/stream", testHandler.TestAgentChatStream)
				}
			}
		}
	}
}
