package main

// func main() {
// 	ctx := context.Background()

// 	// 1. 创建Supervisor Agent（包含所有子Agent）
// 	interviewSupervisor := example.NewInterviewSupervisorAgent()

// 	// 2. 创建Runner
// 	runner := adk.NewRunner(ctx, adk.RunnerConfig{
// 		Agent: interviewSupervisor,
// 	})

// 	// 3. 启动面试流程（用户输入简历）
// 	log.Println("====== 面试流程启动 ======")
// 	iter := runner.Query(ctx, "请面试我，我的简历：熟悉Go语言并发编程、Redis缓存、微服务架构，参与过电商订单系统开发。")

// 	// 4. 处理事件流
// 	for {
// 		event, ok := iter.Next()
// 		if !ok {
// 			break
// 		}
// 		if event.Err != nil {
// 			log.Fatalf("面试流程出错：%v", event.Err)
// 		}

// 		// 打印转让事件
// 		if event.Action != nil && event.Action.TransferToAgent != nil {
// 			log.Printf("\n🔄 调度中心转让任务给：%s", event.Action.TransferToAgent.DestAgentName)
// 			continue
// 		}

// 		// 打印Agent输出
// 		if event.Output != nil && event.Output.MessageOutput != nil {
// 			log.Printf("\n📢 %s 输出：\n%s", event.AgentName, event.Output.MessageOutput.Message.Content)
// 			if event.AgentName == "InterviewReportAgent" {
// 				log.Println("\n====== 面试流程结束 ======")
// 			}
// 		}
// 	}
// }
