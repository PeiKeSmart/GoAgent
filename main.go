package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

// 版本信息 - 支持编译时动态注入
var (
	AppName   = "DHAgent" // 应用名称
	Version   = "dev"     // 版本号，编译时通过 -ldflags 注入
	BuildTime = "unknown" // 构建时间，编译时自动生成
	GitCommit = "unknown" // Git提交哈希，编译时获取
	GitBranch = "unknown" // Git分支，编译时获取
	GoVersion = "unknown" // Go版本，编译时获取
)

// 服务配置
const (
	ServiceName        = "DHAgent"
	ServiceDisplayName = "星尘代理服务"
	ServiceDescription = "星尘，分布式资源调度，部署于每一个节点，连接服务端，支持节点监控、远程发布。"
)

// 全局变量
var (
	ExecutableName string // 可执行文件名（动态获取）
)

// 初始化函数
func init() {
	// 获取可执行文件名
	if exePath, err := os.Executable(); err == nil {
		ExecutableName = filepath.Base(exePath)
	} else {
		// 如果获取失败，使用默认名称
		ExecutableName = "GoAgent.exe"
	}
}

func main() {
	// 检查命令行参数
	if len(os.Args) > 1 {
		operation := os.Args[1]

		// 检查是否需要管理员权限
		if IsElevationRequired(operation) {
			if err := CheckAdminForServiceOperations(); err != nil {
				log.Printf("权限检查失败: %v", err)
				fmt.Println("正在请求管理员权限...")

				if err := RequestAdminPrivileges(); err != nil {
					log.Fatalf("无法获取管理员权限: %v", err)
				}

				fmt.Println("已启动管理员权限进程，当前进程将退出。")
				os.Exit(0)
				return
			}
		}

		switch operation {
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
		case "check-admin":
			if IsRunningAsAdmin() {
				fmt.Println("当前程序正以管理员权限运行")
			} else {
				fmt.Println("当前程序未以管理员权限运行")
			}
			return
		case "status":
			showServiceStatus()
			return
		case "version", "-v", "--version":
			showVersion()
			return
		case "help", "-h", "--help":
			showHelp()
			return
		default:
			fmt.Printf("未知命令: %s\n", operation)
			fmt.Println("使用 'help' 查看可用命令")
			return
		}
	}

	// 在主程序启动时显示服务状态
	fmt.Println("GoAgent 服务管理工具")
	fmt.Println("===================")
	showServiceStatus()
	fmt.Println()
	fmt.Printf("💡 使用 '%s help' 查看所有可用命令\n", ExecutableName)
	fmt.Println("💡 按 Ctrl+C 停止程序")
	fmt.Println()

	// 运行主程序
	runMainProgram()
}

func runMainProgram() {
	// 显示服务启动信息
	fmt.Println("========================================")
	fmt.Printf("服务：星尘代理(%s)\n", ServiceName)
	fmt.Printf("描述：%s\n", ServiceDescription)

	// 获取当前执行路径
	exePath, err := os.Executable()
	if err != nil {
		exePath = ExecutableName
	}

	// 根据不同平台显示不同的状态信息
	if isWindowsService() {
		fmt.Println("状态：Windows 服务运行中")
	} else {
		fmt.Println("状态：程序运行中")
	}
	fmt.Printf("路径：%s\n", exePath)
	fmt.Println("========================================")

	// 显示版本信息
	fmt.Printf("%s       版本：%s   构建时间：%s\n", AppName, Version, BuildTime)
	if GitCommit != "unknown" {
		fmt.Printf("Git提交：%s   分支：%s   Go版本：%s\n", GitCommit, GitBranch, GoVersion)
	}
	fmt.Println()

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

// showServiceStatus 显示服务状态信息
func showServiceStatus() {
	fmt.Println("服务状态信息:")
	fmt.Println("==============")

	// 获取服务状态
	status, err := getServiceStatus()
	if err != nil {
		fmt.Printf("❌ 获取服务状态失败: %v\n", err)
		return
	}

	// 根据状态显示不同的图标和颜色提示
	var statusIcon string
	switch status {
	case "运行中":
		statusIcon = "✅"
	case "已停止":
		statusIcon = "⏹️"
	case "未安装":
		statusIcon = "❓"
	case "失败":
		statusIcon = "❌"
	case "启动中":
		statusIcon = "🔄"
	case "停止中":
		statusIcon = "🔄"
	default:
		statusIcon = "ℹ️"
	}

	fmt.Printf("%s 服务状态: %s\n", statusIcon, status)

	// 获取详细信息
	details, err := getServiceDetails()
	if err != nil {
		fmt.Printf("⚠️  获取详细信息失败: %v\n", err)
		return
	}

	// 显示详细信息
	for key, value := range details {
		fmt.Printf("   %s: %s\n", key, value)
	}

	// 显示可用的操作提示
	if status == "未安装" {
		fmt.Println("\n💡 提示: 使用 'install' 命令安装服务")
	} else if status == "已停止" {
		fmt.Println("\n💡 提示: 使用 'start' 命令启动服务")
	} else if status == "运行中" {
		fmt.Println("\n💡 提示: 服务正在正常运行")
	}
}

// showHelp 显示帮助信息
func showHelp() {
	fmt.Println("GoAgent 服务管理工具")
	fmt.Println("===================")
	fmt.Println()
	fmt.Printf("用法: %s [命令]\n", ExecutableName)
	fmt.Println()
	fmt.Println("可用命令:")
	fmt.Println("  install     安装服务到系统")
	fmt.Println("  uninstall   从系统卸载服务")
	fmt.Println("  start       启动服务")
	fmt.Println("  stop        停止服务")
	fmt.Println("  status      显示服务状态信息")
	fmt.Println("  version     显示版本信息")
	fmt.Println("  check-admin 检查当前权限状态")
	fmt.Println("  help        显示此帮助信息")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Printf("  %s install    # 安装服务\n", ExecutableName)
	fmt.Printf("  %s status     # 查看服务状态\n", ExecutableName)
	fmt.Printf("  %s start      # 启动服务\n", ExecutableName)
	fmt.Println()
	fmt.Println("注意:")
	fmt.Println("  - 服务操作需要管理员权限，程序会自动申请")
	fmt.Println("  - 直接运行程序会显示状态并进入服务模式")
	fmt.Println("  - 按 Ctrl+C 可以优雅地停止服务")
}

// showVersion 显示版本信息
func showVersion() {
	fmt.Printf("%s v%s\n", AppName, Version)
	fmt.Printf("构建时间: %s\n", BuildTime)

	if GitCommit != "unknown" {
		fmt.Printf("Git提交: %s (%s)\n", GitCommit, GitBranch)
	}

	if GoVersion != "unknown" {
		fmt.Printf("Go版本: %s\n", GoVersion)
	}

	fmt.Printf("可执行文件: %s\n", ExecutableName)
}

// isWindowsService 检查当前是否作为Windows服务运行
func isWindowsService() bool {
	// 简单的检查方法：在Windows平台下，检查是否存在Windows特有的环境
	// 这里可以根据实际需要进行更精确的判断
	return os.Getenv("USERPROFILE") != "" && os.Getenv("SYSTEMROOT") != ""
}
