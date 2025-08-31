# GoAgent 管理员权限管理

GoAgent 现在支持自动检测和申请管理员权限。这对于Windows服务的安装、启动、停止等操作是必需的。

## 功能特性

### 自动权限检测
- 程序会自动检测当前是否以管理员权限运行
- 支持Windows和Linux/Unix系统

### 自动权限申请
- 在Windows系统中，程序会通过UAC对话框请求管理员权限
- 在Linux/Unix系统中，程序会提示用户使用sudo命令

### 智能权限管理
- 只有需要管理员权限的操作才会触发权限检查
- 普通操作不会要求管理员权限

## 使用方法

### 检查当前权限状态
```bash
GoAgent.exe check-admin
```

### 服务操作（自动申请权限）
以下命令会自动检查权限，如果权限不足会申请管理员权限：

```bash
# 安装服务
GoAgent.exe install

# 卸载服务
GoAgent.exe uninstall

# 启动服务
GoAgent.exe start

# 停止服务
GoAgent.exe stop
```

### 普通运行（不需要管理员权限）
```bash
# 直接运行程序
GoAgent.exe
```

## 权限检查流程

1. **操作识别**: 程序首先识别用户要执行的操作
2. **权限检查**: 如果操作需要管理员权限，检查当前进程权限
3. **权限申请**: 如果权限不足，自动弹出UAC对话框（Windows）或提示使用sudo（Linux）
4. **重新启动**: 用户确认后，程序会以管理员权限重新启动
5. **执行操作**: 新进程以管理员权限执行请求的操作

## API 说明

### 主要函数

#### `IsRunningAsAdmin() bool`
检查当前进程是否以管理员权限运行。

**返回值:**
- `true`: 当前进程具有管理员权限
- `false`: 当前进程没有管理员权限

#### `RequestAdminPrivileges() error`
请求管理员权限，会启动一个新的具有管理员权限的进程。

**返回值:**
- `nil`: 成功启动管理员权限进程
- `error`: 启动失败的错误信息

#### `EnsureAdminPrivileges() error`
确保程序以管理员权限运行，如果没有权限会自动申请。

**返回值:**
- `nil`: 已经具有管理员权限
- `error`: 无法获取管理员权限的错误信息

#### `CheckAdminForServiceOperations() error`
专门用于服务操作的权限检查。

**返回值:**
- `nil`: 具有执行服务操作的权限
- `error`: 权限不足的错误信息

#### `IsElevationRequired(operation string) bool`
检查指定操作是否需要管理员权限。

**参数:**
- `operation`: 操作名称（如 "install", "uninstall", "start", "stop"）

**返回值:**
- `true`: 操作需要管理员权限
- `false`: 操作不需要管理员权限

## 跨平台支持

### Windows系统
- 使用Windows API检测Token权限
- 通过ShellExecute和"runas"动词启动UAC
- **可以主动申请权限**：弹出UAC对话框让用户确认
- 支持Windows Vista及以上版本

### Linux/Unix系统  
- 使用`os.Geteuid()`检查是否为root用户
- **可以检测权限但无法主动申请**：Linux没有图形化权限提升机制
- **智能sudo处理**：程序会尝试使用sudo自动重新启动
- 提供友好的sudo命令提示
- 兼容所有Linux发行版和Unix系统

### Linux权限管理特点

#### 🔍 **权限检测**
```go
// 可以准确检测是否为root用户
if IsRunningAsAdmin() {
    fmt.Println("当前以root权限运行")
} else {
    fmt.Println("当前没有root权限")
}
```

#### 🔄 **智能sudo重启**
```go
// 程序会自动尝试使用sudo重新启动
err := RequestAdminPrivileges()
// 相当于执行: sudo ./goagent install
```

#### 💡 **用户友好提示**
如果sudo不可用或失败，程序会提供清晰的命令提示：
```bash
请手动使用以下命令运行:
sudo ./goagent install
```

## 错误处理

程序包含完善的错误处理机制：

1. **权限检查失败**: 提供详细的错误信息和解决建议
2. **UAC取消**: 用户取消UAC对话框时的优雅处理
3. **系统API调用失败**: 详细记录API调用错误
4. **跨平台兼容性**: 针对不同操作系统的适配处理

## 安全考虑

1. **最小权限原则**: 只在必要时才请求管理员权限
2. **透明度**: 明确告知用户为什么需要管理员权限
3. **用户控制**: 用户可以拒绝权限申请，程序会优雅退出
4. **安全API**: 使用操作系统提供的标准安全API

## 示例代码

```go
package main

import (
    "fmt"
    "log"
)

func main() {
    // 检查是否需要管理员权限
    if IsElevationRequired("install") {
        // 确保有管理员权限
        if err := EnsureAdminPrivileges(); err != nil {
            log.Fatalf("无法获取管理员权限: %v", err)
        }
    }
    
    // 执行需要权限的操作
    fmt.Println("正在执行服务安装...")
}
```

## 注意事项

1. 权限申请会导致当前进程退出，新进程以管理员权限启动
2. 在自动化脚本中使用时，需要考虑UAC对话框的交互
3. 某些企业环境可能限制UAC权限申请
4. 建议在文档中明确说明程序需要管理员权限的原因
