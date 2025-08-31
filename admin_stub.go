//go:build !windows && !linux

package main

// 这个文件是为了解决编辑器中的权限管理函数"undefined"错误而创建的
// 当编辑器无法识别构建标签时，会使用这些函数定义
// 实际构建时不会包含这个文件，因为有构建标签限制

import "fmt"

func IsRunningAsAdmin() bool {
	return false
}

func RequestAdminPrivileges() error {
	return fmt.Errorf("not implemented for this platform")
}

func EnsureAdminPrivileges() error {
	return fmt.Errorf("not implemented for this platform")
}

func CheckAdminForServiceOperations() error {
	return fmt.Errorf("not implemented for this platform")
}

func RunAsAdminIfNeeded() {
	// stub implementation
}

func IsElevationRequired(operation string) bool {
	return false
}
