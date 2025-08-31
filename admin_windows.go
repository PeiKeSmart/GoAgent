//go:build windows

package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"unsafe"
)

var (
	advapi32                = syscall.NewLazyDLL("advapi32.dll")
	kernel32                = syscall.NewLazyDLL("kernel32.dll")
	shell32                 = syscall.NewLazyDLL("shell32.dll")
	procGetCurrentProcess   = kernel32.NewProc("GetCurrentProcess")
	procOpenProcessToken    = advapi32.NewProc("OpenProcessToken")
	procGetTokenInformation = advapi32.NewProc("GetTokenInformation")
	procShellExecuteW       = shell32.NewProc("ShellExecuteW")
)

const (
	TOKEN_QUERY            = 0x0008
	TokenElevationType     = 18
	TokenElevationTypeFull = 2
)

// IsRunningAsAdmin 检查当前进程是否以管理员权限运行
func IsRunningAsAdmin() bool {
	var token syscall.Handle

	// 获取当前进程句柄
	currentProcess, _, _ := procGetCurrentProcess.Call()

	// 打开进程令牌
	ret, _, _ := procOpenProcessToken.Call(
		currentProcess,
		TOKEN_QUERY,
		uintptr(unsafe.Pointer(&token)),
	)

	if ret == 0 {
		return false
	}
	defer syscall.CloseHandle(token)

	// 获取令牌提升类型
	var elevationType uint32
	var returnLength uint32

	ret, _, _ = procGetTokenInformation.Call(
		uintptr(token),
		TokenElevationType,
		uintptr(unsafe.Pointer(&elevationType)),
		unsafe.Sizeof(elevationType),
		uintptr(unsafe.Pointer(&returnLength)),
	)

	if ret == 0 {
		return false
	}

	// 检查是否为完全提升（管理员权限）
	return elevationType == TokenElevationTypeFull
}

// RequestAdminPrivileges 请求管理员权限
// 通过UAC提示重新启动程序
func RequestAdminPrivileges() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("获取可执行文件路径失败: %v", err)
	}

	// 获取命令行参数
	args := ""
	if len(os.Args) > 1 {
		for i, arg := range os.Args[1:] {
			if i > 0 {
				args += " "
			}
			args += arg
		}
	}

	// 使用ShellExecute以管理员权限启动
	verb := syscall.StringToUTF16Ptr("runas")
	file := syscall.StringToUTF16Ptr(exePath)
	params := syscall.StringToUTF16Ptr(args)

	ret, _, _ := procShellExecuteW.Call(
		0,
		uintptr(unsafe.Pointer(verb)),
		uintptr(unsafe.Pointer(file)),
		uintptr(unsafe.Pointer(params)),
		0,
		1, // SW_SHOWNORMAL
	)

	// ShellExecute 返回值大于32表示成功
	if ret <= 32 {
		return fmt.Errorf("启动管理员进程失败，错误代码: %d", ret)
	}

	return nil
}

// EnsureAdminPrivileges 确保程序以管理员权限运行
// 如果没有管理员权限，会提示用户并尝试重新启动
func EnsureAdminPrivileges() error {
	if IsRunningAsAdmin() {
		log.Println("程序已以管理员权限运行")
		return nil
	}

	log.Println("检测到程序未以管理员权限运行")
	fmt.Println("此程序需要管理员权限才能正常运行。")
	fmt.Println("正在请求管理员权限...")

	err := RequestAdminPrivileges()
	if err != nil {
		return fmt.Errorf("请求管理员权限失败: %v", err)
	}

	fmt.Println("已启动管理员权限进程，当前进程将退出。")
	os.Exit(0)
	return nil
}

// CheckAdminForServiceOperations 检查服务操作是否需要管理员权限
func CheckAdminForServiceOperations() error {
	if !IsRunningAsAdmin() {
		return fmt.Errorf("服务操作需要管理员权限，请以管理员身份运行此程序")
	}
	return nil
}

// RunAsAdminIfNeeded 如果需要，以管理员权限运行程序
// 这是一个便捷函数，结合了检查和请求权限的功能
func RunAsAdminIfNeeded() {
	if !IsRunningAsAdmin() {
		log.Println("程序需要管理员权限，正在请求...")
		err := RequestAdminPrivileges()
		if err != nil {
			log.Fatalf("无法获取管理员权限: %v", err)
		}
		// 如果成功启动了管理员进程，当前进程应该退出
		os.Exit(0)
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
