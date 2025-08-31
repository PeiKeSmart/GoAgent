# Linux vs Windows 权限管理对比

## 📋 权限管理能力对比

| 功能 | Windows | Linux | 说明 |
|------|---------|-------|------|
| **检测当前权限** | ✅ 支持 | ✅ 支持 | 都可以准确检测当前进程权限 |
| **主动申请权限** | ✅ 支持 | ❌ 不支持 | Windows有UAC，Linux没有类似机制 |
| **图形化权限提升** | ✅ UAC对话框 | ❌ 不支持 | Linux只能通过命令行 |
| **自动重启提权** | ✅ 支持 | 🔄 部分支持 | Linux可尝试sudo重启 |

## 🔍 具体实现差异

### Windows权限管理
```go
// ✅ 可以主动申请权限
func RequestAdminPrivileges() error {
    // 通过ShellExecute + "runas" 启动UAC
    // 用户点击"是"后，新进程获得管理员权限
    // 原进程自动退出
}
```

**特点：**
- 🎯 **用户友好**：图形化UAC对话框
- 🔒 **安全性高**：用户明确确认权限提升
- ⚡ **自动化强**：程序可以自主完成权限申请流程

### Linux权限管理
```go
// ❌ 无法主动申请权限，但可以智能处理
func RequestAdminPrivileges() error {
    // 检查sudo是否可用
    // 尝试使用sudo重新启动程序
    // 失败时提供友好的命令提示
}
```

**特点：**
- 📝 **命令行为主**：依赖sudo命令
- 🔄 **智能重启**：自动尝试sudo重新启动
- 💡 **友好提示**：失败时提供准确的命令

## 🚀 用户体验差异

### Windows用户体验
```bash
# 用户运行程序
GoAgent.exe install

# 系统自动弹出UAC对话框
# [UAC] 是否允许此应用对设备进行更改？
# 用户点击"是" -> 自动以管理员权限重新启动

# 无需用户手动操作，完全自动化
```

### Linux用户体验
```bash
# 用户运行程序
./goagent install

# 程序检测权限不足，自动尝试sudo
# 需要管理员权限，正在使用sudo重新启动程序...
# 执行命令: sudo ./goagent install
# [sudo] password for user: ___

# 用户输入密码后，程序以root权限运行
```

## 🎯 最佳实践建议

### Windows开发者
```go
// ✅ 推荐做法：让程序自动处理权限
func main() {
    if IsElevationRequired(operation) {
        if err := CheckAdminForServiceOperations(); err != nil {
            // 自动申请权限，用户只需点击UAC确认
            RequestAdminPrivileges()
            return
        }
    }
    // 执行操作...
}
```

### Linux开发者
```go
// ✅ 推荐做法：智能sudo处理 + 友好提示
func main() {
    if IsElevationRequired(operation) {
        if err := CheckAdminForServiceOperations(); err != nil {
            // 尝试自动sudo，失败时提供命令提示
            if err := RequestAdminPrivileges(); err != nil {
                fmt.Printf("请手动运行: %s\n", GetSudoCommand())
            }
            return
        }
    }
    // 执行操作...
}
```

## 📚 技术实现细节

### Windows API调用
```go
// 使用Windows API检测Token权限
procGetTokenInformation.Call(
    uintptr(token),
    TokenElevationType,
    uintptr(unsafe.Pointer(&elevationType)),
    unsafe.Sizeof(elevationType),
    uintptr(unsafe.Pointer(&returnLength)),
)

// 通过ShellExecute启动UAC
procShellExecuteW.Call(
    0,
    uintptr(unsafe.Pointer(verb)), // "runas"
    uintptr(unsafe.Pointer(file)),
    uintptr(unsafe.Pointer(params)),
    0,
    1, // SW_SHOWNORMAL
)
```

### Linux系统调用
```go
// 使用Unix系统调用检测权限
os.Geteuid() == 0  // 检查是否为root

// 尝试sudo重新启动
cmd := exec.Command("sudo", append([]string{exePath}, os.Args[1:]...)...)
cmd.Stdin = os.Stdin
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
cmd.Run()
```

## 🎯 总结

**Windows：完全自动化权限管理**
- ✅ 真正的"主动申请权限"
- ✅ 用户体验最佳
- ✅ 完全自动化流程

**Linux：智能权限辅助**
- ✅ 权限检测准确
- 🔄 智能sudo处理
- 💡 友好用户提示
- ❌ 无法完全自动化（需要用户输入密码）

这正是Linux和Windows在权限管理架构上的根本差异！
