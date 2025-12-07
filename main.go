package main

import (
	"context"
	"example-eino-pdf-agent/agents/resume"
	"fmt"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
	"log"
	"time"
)

// 获取岗位描述信息,用于分析岗位匹配度
func getJobDescription() string {
	return `
公司名称: 飞牛
岗位名称: C/C++开发工程师
岗位职责:
负责NAS项目的核心功能开发, 包括数据恢复, 网络配置, 虚拟机管理, docker管理等功能模块的开发。

任职要求:
C++/PostgreSQL/STL/Linux开发/部署经验
1）扎实的C++编程技术, 有C++开发经验, 熟悉常用的数据结构、算法；
2) 熟悉Linux操作系统及Linux多线程、进程通信等内容的编程；熟悉计算机网络编程（TCP/IP、UDP等网络通信协议）；
3) 了解Linux文件系统, 有存储类应用开发优先考虑；
4) 有Linux下软RAID使用经验, 了解常用磁盘整列原理(RAID-0, RAID-1, RAID-5等)优先考虑。

公司介绍:
铁刃智造是一家专注于研发安全易用、免费的NAS系统的企业。
我们始终以用户需求为导向，致力于为家庭、工作室等中小企业提供高效、安全、可靠的数据存储和处理解决方案。我们的团队由原UC、PP助手核心成员组成，秉持长期主义的态度，打造出极致体验且安全可靠的NAS系统，为全球市场带来可称之为“国产之光”的产品。
目前主营业务主要有包括：NAS、私有存储周边硬件、私有数据资产+AI解决方案。
`
}

func getUserQuery(userPdf string, jobDesc string) string {
	return fmt.Sprintf(`
【重要】请立即解析以下简历文件并提取关键信息：

简历文件路径：%s
岗位描述(JD): %s

【必须执行的步骤】：
1. 【第一步】立即使用工具解析简历文件，获取完整的简历文本内容
2. 【第二步】从解析的简历文本中提取所有关键信息
3. 【第三步】根据提取的关键信息, 结合 岗位描述(JD) 生成符合要求的 JSON 格式输出

【重要提示】：
- 不要跳过 pdf_to_text 工具调用
- 必须从简历内容中提取真实的信息，不要返回空数据
- 所有JSON字段都必须填充实际内容
- 只返回JSON格式，不要返回其他文本

请返回完整的 JSON 格式结果。
`, userPdf, jobDesc)
}

func main() {
	// PDF文件的绝对路径
	userPdf := "/Users/lucas/work/code/go/example-eino-pdf-agent/examples/test.pdf"

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	agent := resume.NewResumeAgent()
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent: agent,
		// 用于向 Agent 建议其输出模式，但它并非一个强制性约束。
		// 它的核心思想是控制那些同时支持流式和非流式输出的组件的行为，例如 ChatModel
		// 当 EnableStreaming=false 时，对于那些既能流式也能非流式输出的组件，此时会使用一次性返回完整结果的非流式模式。
		EnableStreaming: false,
	})

	// 执行对话
	input := []adk.Message{
		schema.UserMessage(getUserQuery(userPdf, getJobDescription())),
	}
	events := runner.Run(ctx, input)
	output := ""
	for {
		event, ok := events.Next()
		if !ok {
			break
		}

		if event.Err != nil {
			log.Printf("event错误: %v", event.Err)
			break
		}

		if msg, err := event.Output.MessageOutput.GetMessage(); err == nil {
			output = msg.Content
		}
	}

	log.Println(output)

}
