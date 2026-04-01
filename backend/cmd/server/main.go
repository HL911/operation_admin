package main

import (
	"flag"
	"log"

	"operation_admin/backend/internal/app"
)

// main 负责解析启动参数并拉起整个后端应用。
func main() {
	// configPath 用于接收外部显式传入的配置文件路径。
	configPath := flag.String("config", "", "配置文件路径，默认自动探测")
	flag.Parse()

	// application 负责统一管理配置、日志、数据库与 HTTP 服务生命周期。
	application, err := app.New(*configPath)
	if err != nil {
		log.Fatalf("初始化应用失败: %v", err)
	}

	defer func() {
		if closeErr := application.Close(); closeErr != nil {
			log.Printf("释放应用资源失败: %v", closeErr)
		}
	}()

	if err := application.Run(); err != nil {
		log.Fatalf("运行应用失败: %v", err)
	}
}
