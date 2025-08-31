// service_windows.go
//go:build windows

package main

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	serviceName        = ServiceName
	serviceDisplayName = ServiceDisplayName
	serviceDescription = ServiceDescription
)

func installWindowsService() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("获取可执行文件路径失败: %v", err)
	}

	// 使用 sc 命令创建 Windows 服务
	cmd := exec.Command("sc", "create", serviceName,
		"binPath=", fmt.Sprintf("\"%s\"", exePath),
		"DisplayName=", serviceDisplayName,
		"start=", "auto")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("创建服务失败: %v, 输出: %s", err, output)
	}

	// 设置服务描述
	descCmd := exec.Command("sc", "description", serviceName, serviceDescription)
	descCmd.Run()

	return nil
}

func uninstallWindowsService() error {
	// 先停止服务
	stopCmd := exec.Command("sc", "stop", serviceName)
	stopCmd.Run()

	// 删除服务
	cmd := exec.Command("sc", "delete", serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("删除服务失败: %v, 输出: %s", err, output)
	}

	return nil
}

func startWindowsService() error {
	cmd := exec.Command("sc", "start", serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("启动服务失败: %v, 输出: %s", err, output)
	}
	return nil
}

func stopWindowsService() error {
	cmd := exec.Command("sc", "stop", serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("停止服务失败: %v, 输出: %s", err, output)
	}
	return nil
}

// getServiceStatus 获取Windows服务状态
func getServiceStatus() (string, error) {
	cmd := exec.Command("sc", "query", serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// 如果服务不存在，sc query会返回错误
		return "未安装", nil
	}

	outputStr := string(output)

	// 解析服务状态
	if contains(outputStr, "RUNNING") {
		return "运行中", nil
	} else if contains(outputStr, "STOPPED") {
		return "已停止", nil
	} else if contains(outputStr, "START_PENDING") {
		return "启动中", nil
	} else if contains(outputStr, "STOP_PENDING") {
		return "停止中", nil
	} else if contains(outputStr, "PAUSED") {
		return "已暂停", nil
	} else {
		return "未知状态", nil
	}
}

// getServiceDetails 获取服务详细信息
func getServiceDetails() (map[string]string, error) {
	details := make(map[string]string)

	// 获取服务配置信息
	cmd := exec.Command("sc", "qc", serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return details, fmt.Errorf("获取服务配置失败: %v", err)
	}

	outputStr := string(output)
	details["配置"] = "已配置"

	// 获取服务状态信息
	statusCmd := exec.Command("sc", "query", serviceName)
	statusOutput, err := statusCmd.CombinedOutput()
	if err != nil {
		details["状态"] = "未安装"
		return details, nil
	}

	statusStr := string(statusOutput)

	if contains(statusStr, "RUNNING") {
		details["状态"] = "运行中"
	} else if contains(statusStr, "STOPPED") {
		details["状态"] = "已停止"
	} else {
		details["状态"] = "其他状态"
	}

	// 获取启动类型
	if contains(outputStr, "AUTO_START") {
		details["启动类型"] = "自动"
	} else if contains(outputStr, "DEMAND_START") {
		details["启动类型"] = "手动"
	} else if contains(outputStr, "DISABLED") {
		details["启动类型"] = "禁用"
	} else {
		details["启动类型"] = "未知"
	}

	return details, nil
}

// contains 检查字符串是否包含子串（简单实现）
func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// installService 安装当前程序为 Windows 服务。
// 返回：
//
//	error - 如果安装过程中出现错误，则返回错误信息；否则返回 nil。
func installService() error {
	return installWindowsService()
}

// uninstallService 卸载当前服务。
// 该函数封装了卸载Windows服务的逻辑，返回操作过程中可能遇到的错误。
// 适用于需要卸载服务的场景，通常用于服务管理或清理操作。
func uninstallService() error {
	return uninstallWindowsService()
}

// startService 启动服务函数
//
// 返回值：
// - error：如果启动服务失败，则返回错误信息；否则返回nil
func startService() error {
	return startWindowsService()
}

// stopService 停止当前服务。
// 返回：
//
//	error - 如果停止服务过程中发生错误，则返回错误信息；否则返回 nil。
func stopService() error {
	return stopWindowsService()
}

// restartService 重启当前服务
// 返回：
//
//	error - 如果重启服务过程中发生错误，则返回错误信息；否则返回 nil。
func restartService() error {
	// 先停止服务
	if err := stopWindowsService(); err != nil {
		// 如果停止失败但不是因为服务未运行，则返回错误
		if !contains(err.Error(), "1062") { // 服务未启动的错误代码
			return fmt.Errorf("停止服务失败: %v", err)
		}
	}

	// 等待一段时间确保服务完全停止
	// time.Sleep(2 * time.Second)

	// 启动服务
	return startWindowsService()
}
