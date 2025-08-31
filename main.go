package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
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
	// 检查是否作为 Windows 服务运行
	if isRunningAsService() {
		// 作为 Windows 服务运行
		runAsWindowsService(ServiceName, false)
		return
	}

	// 检查命令行参数
	if len(os.Args) > 1 {
		operation := os.Args[1]

		// 处理带 - 前缀的命令
		switch operation {
		case "-status":
			operation = "status"
		case "-u":
			operation = "uninstall"
		case "-stop":
			operation = "stop"
		case "-restart":
			operation = "restart"
		case "-run":
			operation = "run"
		case "-v", "--version":
			operation = "version"
		case "-h", "--help":
			operation = "help"
		}

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
		case "restart":
			if err := restartService(); err != nil {
				log.Fatalf("重启服务失败: %v", err)
			}
			fmt.Println("服务重启成功！")
			return
		case "run":
			fmt.Println("模拟运行模式启动...")
			startAgentService()
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
		case "version":
			showVersion()
			return
		case "help":
			showHelp()
			return
		default:
			fmt.Printf("未知命令: %s\n", operation)
			fmt.Println("使用 'help' 查看可用命令")
			return
		}
	}

	// 无参数启动时显示交互式菜单
	showInteractiveMenu()
}

// startAgentService 启动星尘代理服务 (业务逻辑层)
func startAgentService() {
	// 显示服务启动信息
	fmt.Println("========================================")
	fmt.Printf("服务：星尘代理(%s)\n", ServiceName)
	fmt.Printf("描述：%s\n", ServiceDescription)

	// 获取当前执行路径
	exePath, err := os.Executable()
	if err != nil {
		exePath = ExecutableName
	}

	// 检查真实的服务状态
	serviceStatus, err := getServiceStatus()
	if err != nil {
		fmt.Println("状态：程序直接运行中（非服务模式）")
	} else if serviceStatus == "运行中" {
		fmt.Println("状态：Windows 服务运行中")
	} else {
		fmt.Println("状态：程序直接运行中（非服务模式）")
	}

	fmt.Printf("路径：%s\n", exePath)
	fmt.Println("========================================")

	// 显示版本信息
	fmt.Printf("%s       版本：%s   构建时间：%s\n", AppName, Version, BuildTime)
	if GitCommit != "unknown" {
		fmt.Printf("Git提交：%s   分支：%s   Go版本：%s\n", GitCommit, GitBranch, GoVersion)
	}
	fmt.Println()

	// 创建并启动代理服务
	agent := NewAgentService()

	// 创建信号通道来处理优雅关闭
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动代理服务 (在goroutine中)
	go func() {
		if err := agent.Start(); err != nil {
			log.Printf("代理服务启动失败: %v", err)
		}
	}()

	// 等待停止信号
	sig := <-sigChan
	log.Printf("收到信号 %v，正在关闭服务...", sig)

	// 优雅停止代理服务
	agent.Stop()
}

func runMainProgram() {
	// 兼容性函数，现在直接调用代理服务
	startAgentService()
} // showServiceStatus 显示服务状态信息
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

	// 只有在服务存在时才获取详细信息
	if status != "未安装" {
		// 获取详细信息
		details, err := getServiceDetails()
		if err != nil {
			fmt.Printf("⚠️  获取详细信息失败: %v\n", err)
		} else {
			// 显示详细信息
			for key, value := range details {
				fmt.Printf("   %s: %s\n", key, value)
			}
		}
	}

	// 显示可用的操作提示
	fmt.Println()
	if status == "未安装" {
		fmt.Println("💡 建议操作:")
		fmt.Println("   - 使用菜单选项 6 或命令 'install' 安装服务")
		fmt.Println("   - 使用菜单选项 5 或命令 '-run' 进行模拟运行")
	} else if status == "已停止" {
		fmt.Println("💡 建议操作:")
		fmt.Println("   - 使用菜单选项 7 或命令 'start' 启动服务")
		fmt.Println("   - 使用菜单选项 2 或命令 'uninstall' 卸载服务")
	} else if status == "运行中" {
		fmt.Println("💡 可用操作:")
		fmt.Println("   - 使用菜单选项 3 或命令 'stop' 停止服务")
		fmt.Println("   - 使用菜单选项 4 或命令 'restart' 重启服务")
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
	fmt.Println("  install           安装服务到系统")
	fmt.Println("  uninstall (-u)    从系统卸载服务")
	fmt.Println("  start             启动服务")
	fmt.Println("  stop (-stop)      停止服务")
	fmt.Println("  restart (-restart) 重启服务")
	fmt.Println("  status (-status)  显示服务状态信息")
	fmt.Println("  run (-run)        模拟运行模式（非服务）")
	fmt.Println("  version (-v)      显示版本信息")
	fmt.Println("  check-admin       检查当前权限状态")
	fmt.Println("  help (-h)         显示此帮助信息")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Printf("  %s install       # 安装服务\n", ExecutableName)
	fmt.Printf("  %s -status       # 查看服务状态\n", ExecutableName)
	fmt.Printf("  %s start         # 启动服务\n", ExecutableName)
	fmt.Printf("  %s -run          # 模拟运行模式\n", ExecutableName)
	fmt.Println()
	fmt.Println("交互模式:")
	fmt.Printf("  %s               # 启动交互式菜单\n", ExecutableName)
	fmt.Println()
	fmt.Println("注意:")
	fmt.Println("  - 服务操作需要管理员权限，程序会自动申请")
	fmt.Println("  - 直接运行程序会进入交互式菜单")
	fmt.Println("  - 使用 -run 参数可以在非服务模式下运行")
	fmt.Println("  - 按 Ctrl+C 可以优雅地停止程序")
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

// showInteractiveMenu 显示交互式菜单
func showInteractiveMenu() {
	// 显示服务基本信息
	fmt.Println("========================================")
	fmt.Printf("服务：星尘代理(%s)\n", ServiceName)
	fmt.Printf("描述：%s\n", ServiceDescription)

	// 获取当前执行路径
	exePath, err := os.Executable()
	if err != nil {
		exePath = ExecutableName
	}

	// 检查真实的服务状态
	serviceStatus, err := getServiceStatus()
	if err != nil {
		fmt.Println("状态：程序直接运行中（非服务模式）")
	} else if serviceStatus == "运行中" {
		fmt.Println("状态：Windows 服务运行中")
	} else {
		fmt.Println("状态：程序直接运行中（非服务模式）")
	}

	fmt.Printf("路径：%s\n", exePath)
	fmt.Println("========================================")

	// 显示版本信息
	fmt.Printf("%s       版本：%s   构建时间：%s\n", AppName, Version, BuildTime)
	if GitCommit != "unknown" {
		fmt.Printf("Git提交：%s   分支：%s   Go版本：%s\n", GitCommit, GitBranch, GoVersion)
	}
	fmt.Println()

	fmt.Println("GoAgent 服务管理工具")
	fmt.Println("===================")
	showServiceStatus()
	fmt.Println()

	for {
		showMenu()
		choice := getUserInput()

		if !handleMenuChoice(choice) {
			break
		}

		fmt.Println()
	}
}

// showMenu 显示菜单选项
func showMenu() {
	fmt.Println("序号 功能名称   命令行参数")
	fmt.Println(" 1、 显示状态   -status")
	fmt.Println(" 2、 卸载服务   -u")
	fmt.Println(" 3、 停止服务   -stop")
	fmt.Println(" 4、 重启服务   -restart")
	fmt.Println(" 5、 模拟运行   -run")
	fmt.Println(" 6、 安装服务   install")
	fmt.Println(" 7、 启动服务   start")
	fmt.Println(" v、 版本信息   version")
	fmt.Println(" h、 帮助信息   help")
	fmt.Println(" 0、 退出")
	fmt.Print("请选择操作 (输入序号或字母): ")
}

// getUserInput 获取用户输入
func getUserInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text())
	}
	return ""
}

// handleMenuChoice 处理菜单选择
func handleMenuChoice(choice string) bool {
	switch choice {
	case "1":
		fmt.Println("\n正在显示服务状态...")
		showServiceStatus()
	case "2":
		fmt.Println("\n正在卸载服务...")
		handlePrivilegedOperation("uninstall", func() error {
			return uninstallService()
		})
	case "3":
		fmt.Println("\n正在停止服务...")
		handlePrivilegedOperation("stop", func() error {
			return stopService()
		})
	case "4":
		fmt.Println("\n正在重启服务...")
		handlePrivilegedOperation("restart", func() error {
			return restartService()
		})
	case "5":
		fmt.Println("\n启动模拟运行模式...")
		fmt.Println("按 Ctrl+C 停止运行")
		startAgentService()
	case "6":
		fmt.Println("\n正在安装服务...")
		handlePrivilegedOperation("install", func() error {
			return installService()
		})
	case "7":
		fmt.Println("\n正在启动服务...")
		handlePrivilegedOperation("start", func() error {
			return startService()
		})
	case "v", "V":
		fmt.Println()
		showVersion()
	case "h", "H":
		fmt.Println()
		showHelp()
	case "0":
		fmt.Println("退出程序...")
		return false
	default:
		fmt.Printf("❌ 无效选择: %s\n", choice)
	}
	return true
}

// handlePrivilegedOperation 处理需要权限的操作
func handlePrivilegedOperation(operation string, fn func() error) {
	// 检查是否需要管理员权限
	if IsElevationRequired(operation) {
		if err := CheckAdminForServiceOperations(); err != nil {
			fmt.Printf("❌ 操作失败: %v\n", err)
			fmt.Println("💡 提示: 请以管理员身份重新启动程序")
			fmt.Printf("💡 或者在管理员命令提示符中运行: %s %s\n", ExecutableName, operation)
			return
		}
	}

	// 执行操作
	if err := fn(); err != nil {
		fmt.Printf("❌ 操作失败: %v\n", err)
	} else {
		var successMsg string
		switch operation {
		case "install":
			successMsg = "✅ 服务安装成功！"
		case "uninstall":
			successMsg = "✅ 服务卸载成功！"
		case "start":
			successMsg = "✅ 服务启动成功！"
		case "stop":
			successMsg = "✅ 服务停止成功！"
		case "restart":
			successMsg = "✅ 服务重启成功！"
		default:
			successMsg = "✅ 操作成功！"
		}
		fmt.Println(successMsg)
	}
}

// isWindowsService 检查当前是否作为Windows服务运行
func isWindowsService() bool {
	// 更准确的检查方法：检查是否有控制台窗口
	// 如果没有控制台窗口且没有用户交互环境，通常表示作为服务运行
	return os.Getenv("USERNAME") == "" || os.Getenv("SESSIONNAME") == ""
}
