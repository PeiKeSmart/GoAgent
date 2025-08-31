// service_windows.go
//go:build windows

package main

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	serviceName        = "GoAgent"
	serviceDisplayName = "Go Agent Service"
	serviceDescription = "Go Agent 自启动服务"
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
