package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 检查命令行参数
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "install":
			if err := installService(); err != nil {
				log.Fatalf("安装服务失败: %v", err)
			}
			fmt.Println("服务安装成功！")
			return
		case "uninstall":
			if err := uninstallService(); err != nil {
				log.Fatalf("卸载服务失败: %v", err)
			}
			fmt.Println("服务卸载成功！")
			return
		case "start":
			if err := startService(); err != nil {
				log.Fatalf("启动服务失败: %v", err)
			}
			fmt.Println("服务启动成功！")
			return
		case "stop":
			if err := stopService(); err != nil {
				log.Fatalf("停止服务失败: %v", err)
			}
			fmt.Println("服务停止成功！")
			return
		}
	}

	// 运行主程序
	runMainProgram()
}

func runMainProgram() {
	log.Println("GoAgent 服务已启动")

	// 创建信号通道来处理优雅关闭
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 主循环
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 这里执行您的主要业务逻辑
			log.Println("GoAgent 正在运行...")
			// 可以在这里添加您的具体功能

		case sig := <-sigChan:
			log.Printf("收到信号 %v，正在关闭服务...", sig)
			return
		}
	}
}
