//go:build !windows

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// IsRunningAsAdmin 在非Windows系统中检查是否以root权限运行
func IsRunningAsAdmin() bool {
	return os.Geteuid() == 0
}

// RequestAdminPrivileges 在非Windows系统中尝试重新启动程序（使用sudo）
func RequestAdminPrivileges() error {
	// 检查是否已经是root用户
	if IsRunningAsAdmin() {
		return nil
	}

	// 检查sudo是否可用
	if !isSudoAvailable() {
		return fmt.Errorf("需要root权限，但sudo不可用。请直接以root用户运行此程序")
	}

	// 获取当前执行文件路径
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("无法获取可执行文件路径: %v", err)
	}

	// 构建sudo命令
	args := []string{"sudo", exePath}
	args = append(args, os.Args[1:]...)

	fmt.Printf("需要管理员权限，正在使用sudo重新启动程序...\n")
	fmt.Printf("执行命令: %s\n", strings.Join(args, " "))

	// 使用sudo重新执行程序
	cmd := exec.Command("sudo", append([]string{exePath}, os.Args[1:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// 获取退出码
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			}
		}
		return fmt.Errorf("sudo执行失败: %v", err)
	}

	// sudo命令成功执行后，当前进程应该退出
	os.Exit(0)
	return nil
}

// EnsureAdminPrivileges 确保程序以管理员权限运行
func EnsureAdminPrivileges() error {
	if IsRunningAsAdmin() {
		return nil
	}

	// 尝试使用sudo重新启动
	return RequestAdminPrivileges()
}

// CheckAdminForServiceOperations 检查服务操作是否需要管理员权限
func CheckAdminForServiceOperations() error {
	if !IsRunningAsAdmin() {
		return fmt.Errorf("服务操作需要root权限，请使用 sudo 运行此程序")
	}
	return nil
}

// RunAsAdminIfNeeded 在非Windows系统中尝试使用sudo
func RunAsAdminIfNeeded() {
	if !IsRunningAsAdmin() {
		err := RequestAdminPrivileges()
		if err != nil {
			fmt.Printf("无法获取管理员权限: %v\n", err)
			fmt.Println("请手动使用以下命令运行:")
			fmt.Printf("sudo %s %s\n", os.Args[0], strings.Join(os.Args[1:], " "))
			os.Exit(1)
		}
	}
}

// IsElevationRequired 检查特定操作是否需要提升权限
func IsElevationRequired(operation string) bool {
	elevationRequiredOps := map[string]bool{
		"install":   true,
		"uninstall": true,
		"start":     true,
		"stop":      true,
		"restart":   true,
	}

	return elevationRequiredOps[operation]
}

// isSudoAvailable 检查sudo命令是否可用
func isSudoAvailable() bool {
	_, err := exec.LookPath("sudo")
	return err == nil
}

// GetSudoCommand 返回完整的sudo命令字符串（用于提示用户）
func GetSudoCommand() string {
	exePath, err := os.Executable()
	if err != nil {
		exePath = os.Args[0]
	}

	// 使用绝对路径
	if !filepath.IsAbs(exePath) {
		if absPath, err := filepath.Abs(exePath); err == nil {
			exePath = absPath
		}
	}

	args := []string{"sudo", exePath}
	args = append(args, os.Args[1:]...)

	return strings.Join(args, " ")
}

// CheckSudoPermission 检查用户是否有sudo权限
func CheckSudoPermission() bool {
	if !isSudoAvailable() {
		return false
	}

	// 尝试执行sudo -n true来检查是否可以无密码执行sudo
	cmd := exec.Command("sudo", "-n", "true")
	err := cmd.Run()
	return err == nil
}
