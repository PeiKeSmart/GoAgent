// service_linux.go
//go:build linux

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	serviceName = "goagent"
	serviceFile = "/etc/systemd/system/goagent.service"
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
Description=Go Agent Service
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
`, absPath)

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
