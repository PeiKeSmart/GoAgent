// service_linux.go
//go:build linux

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	serviceName = ServiceName
	serviceFile = "/etc/systemd/system/dhagent.service"
)

func installLinuxService() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("获取可执行文件路径失败: %v", err)
	}

	// 获取绝对路径
	absPath, err := filepath.Abs(exePath)
	if err != nil {
		return fmt.Errorf("获取绝对路径失败: %v", err)
	}

	// 创建 systemd 服务文件内容
	serviceContent := fmt.Sprintf(`[Unit]
Description=%s
After=network.target

[Service]
Type=simple
User=root
ExecStart=%s
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
`, ServiceDisplayName, absPath)

	// 写入服务文件
	if err := os.WriteFile(serviceFile, []byte(serviceContent), 0644); err != nil {
		return fmt.Errorf("创建服务文件失败: %v", err)
	}

	// 重新加载 systemd
	cmd := exec.Command("systemctl", "daemon-reload")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("重新加载 systemd 失败: %v", err)
	}

	// 启用服务
	cmd = exec.Command("systemctl", "enable", serviceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("启用服务失败: %v", err)
	}

	return nil
}

func uninstallLinuxService() error {
	// 停止服务
	stopCmd := exec.Command("systemctl", "stop", serviceName)
	stopCmd.Run()

	// 禁用服务
	disableCmd := exec.Command("systemctl", "disable", serviceName)
	disableCmd.Run()

	// 删除服务文件
	if err := os.Remove(serviceFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("删除服务文件失败: %v", err)
	}

	// 重新加载 systemd
	cmd := exec.Command("systemctl", "daemon-reload")
	cmd.Run()

	return nil
}

func startLinuxService() error {
	cmd := exec.Command("systemctl", "start", serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("启动服务失败: %v, 输出: %s", err, output)
	}
	return nil
}

func stopLinuxService() error {
	cmd := exec.Command("systemctl", "stop", serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("停止服务失败: %v, 输出: %s", err, output)
	}
	return nil
}

func restartLinuxService() error {
	cmd := exec.Command("systemctl", "restart", serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("重启服务失败: %v, 输出: %s", err, output)
	}
	return nil
}

// getServiceStatus 获取Linux服务状态
func getServiceStatus() (string, error) {
	cmd := exec.Command("systemctl", "is-active", serviceName)
	output, err := cmd.CombinedOutput()

	status := string(output)
	status = strings.TrimSpace(status)

	switch status {
	case "active":
		return "运行中", nil
	case "inactive":
		return "已停止", nil
	case "failed":
		return "失败", nil
	case "activating":
		return "启动中", nil
	case "deactivating":
		return "停止中", nil
	default:
		if err != nil {
			// 服务可能不存在
			return "未安装", nil
		}
		return status, nil
	}
}

// getServiceDetails 获取Linux服务详细信息
func getServiceDetails() (map[string]string, error) {
	details := make(map[string]string)

	// 检查服务状态
	statusCmd := exec.Command("systemctl", "is-active", serviceName)
	statusOutput, err := statusCmd.CombinedOutput()
	statusStr := strings.TrimSpace(string(statusOutput))

	if err != nil {
		details["状态"] = "未安装"
		details["启动类型"] = "未配置"
		return details, nil
	}

	switch statusStr {
	case "active":
		details["状态"] = "运行中"
	case "inactive":
		details["状态"] = "已停止"
	case "failed":
		details["状态"] = "失败"
	default:
		details["状态"] = statusStr
	}

	// 检查服务是否启用
	enableCmd := exec.Command("systemctl", "is-enabled", serviceName)
	enableOutput, err := enableCmd.CombinedOutput()
	enableStr := strings.TrimSpace(string(enableOutput))

	if err != nil {
		details["启动类型"] = "未知"
	} else {
		switch enableStr {
		case "enabled":
			details["启动类型"] = "自动"
		case "disabled":
			details["启动类型"] = "手动"
		case "static":
			details["启动类型"] = "静态"
		default:
			details["启动类型"] = enableStr
		}
	}

	// 检查服务文件是否存在
	if _, err := os.Stat(serviceFile); err == nil {
		details["配置"] = "已配置"
	} else {
		details["配置"] = "未配置"
	}

	return details, nil
}

// 通用服务接口实现
func installService() error {
	return installLinuxService()
}

func uninstallService() error {
	return uninstallLinuxService()
}

func startService() error {
	return startLinuxService()
}

func stopService() error {
	return stopLinuxService()
}

func restartService() error {
	return restartLinuxService()
}
