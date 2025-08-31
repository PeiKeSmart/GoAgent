# GoAgent

一个使用 Go 语言开发的跨平台系统服务代理程序，支持 Windows 和 Linux 系统的自启动服务管理。

## 📋 项目介绍

GoAgent 是一个轻量级的系统服务管理工具，它可以作为系统服务运行在后台，并提供完整的服务生命周期管理功能。该项目采用跨平台设计，能够在 Windows 和 Linux 系统上无缝运行。

## ✨ 功能特点

- **跨平台支持**：同时支持 Windows 和 Linux 操作系统
- **系统服务集成**：完整的系统服务安装、卸载、启动、停止功能
- **自启动能力**：支持系统启动时自动运行
- **优雅关闭**：响应系统信号，实现优雅的服务关闭
- **日志记录**：内置日志系统，便于监控和调试
- **轻量级设计**：单文件部署，无外部依赖
- **可扩展架构**：易于添加自定义业务逻辑

## 🚀 安装和使用

### 系统要求

- **Windows**: Windows 7/8/10/11 或 Windows Server 2012 及以上版本
- **Linux**: 支持 systemd 的发行版（如 Ubuntu 16.04+, CentOS 7+, Debian 8+ 等）
- **权限要求**: 
  - Windows: 需要管理员权限
  - Linux: 需要 root 权限

### 编译安装

#### 1. 克隆项目
```bash
git clone https://github.com/PeiKeSmart/GoAgent.git
cd GoAgent
```

#### 2. 编译项目

**Windows 编译：**
```bash
# 编译 Windows 版本
go build -o GoAgent.exe .

# 交叉编译（在其他平台编译 Windows 版本）
GOOS=windows GOARCH=amd64 go build -o GoAgent.exe .
```

**Linux 编译：**
```bash
# 编译 Linux 版本
go build -o goagent .

# 交叉编译（在其他平台编译 Linux 版本）
GOOS=linux GOARCH=amd64 go build -o goagent .
```

### 使用教程

#### Windows 系统

1. **以管理员身份运行命令提示符或 PowerShell**

2. **安装服务**
   ```cmd
   .\GoAgent.exe install
   ```
   成功输出：`服务安装成功！`

3. **启动服务**
   ```cmd
   .\GoAgent.exe start
   ```
   成功输出：`服务启动成功！`

4. **停止服务**
   ```cmd
   .\GoAgent.exe stop
   ```
   成功输出：`服务停止成功！`

5. **卸载服务**
   ```cmd
   .\GoAgent.exe uninstall
   ```
   成功输出：`服务卸载成功！`

6. **直接运行（前台模式）**
   ```cmd
   .\GoAgent.exe
   ```

#### Linux 系统

1. **使用 root 权限运行**

2. **安装服务**
   ```bash
   sudo ./goagent install
   ```
   成功输出：`服务安装成功！`

3. **启动服务**
   ```bash
   sudo ./goagent start
   # 或使用 systemctl
   sudo systemctl start goagent
   ```
   成功输出：`服务启动成功！`

4. **停止服务**
   ```bash
   sudo ./goagent stop
   # 或使用 systemctl
   sudo systemctl stop goagent
   ```
   成功输出：`服务停止成功！`

5. **卸载服务**
   ```bash
   sudo ./goagent uninstall
   ```
   成功输出：`服务卸载成功！`

6. **直接运行（前台模式）**
   ```bash
   ./goagent
   ```

### 验证服务状态

#### Windows
```cmd
# 查看服务状态
sc query GoAgent

# 查看服务配置
sc qc GoAgent
```

#### Linux
```bash
# 查看服务状态
sudo systemctl status goagent

# 查看服务日志
sudo journalctl -u goagent -f

# 查看服务是否开机自启
sudo systemctl is-enabled goagent
```

## 🔧 技术实现

### 架构设计

项目采用模块化设计，主要包含以下组件：

- **main.go**: 主程序入口，处理命令行参数和业务逻辑
- **service_windows.go**: Windows 平台特定的服务管理实现
- **service_linux.go**: Linux 平台特定的服务管理实现

### 核心技术栈

- **编程语言**: Go 1.25+
- **构建标签**: 使用 Go 的构建标签实现平台特定编译
- **Windows 服务**: 基于 Windows SC (Service Control) 命令
- **Linux 服务**: 基于 systemd 服务管理器
- **信号处理**: 实现优雅关闭机制
- **并发控制**: 使用 Go 的 channel 和 goroutine

### 跨平台实现原理

项目使用 Go 的构建标签（Build Tags）实现跨平台编译：

```go
//go:build windows  // Windows 特定代码
//go:build linux    // Linux 特定代码
```

在编译时，Go 编译器会根据目标平台自动选择相应的文件进行编译，确保生成的二进制文件只包含目标平台需要的代码。

### 服务管理机制

#### Windows 实现
- 使用 `sc create` 命令创建 Windows 服务
- 配置服务自启动（`start=auto`）
- 支持服务描述和显示名称设置

#### Linux 实现
- 生成 systemd 服务单元文件
- 配置服务自启动和重启策略
- 集成系统日志记录

## 🛠️ 自定义功能

### 添加业务逻辑

在 `main.go` 的 `runMainProgram()` 函数中添加你的业务代码：

```go
func runMainProgram() {
    log.Println("GoAgent 服务已启动")

    // 创建信号通道来处理优雅关闭
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    // 主循环
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            // 在这里添加你的业务逻辑
            log.Println("GoAgent 正在运行...")
            // 示例：添加自定义功能
            // performCustomTask()
            
        case sig := <-sigChan:
            log.Printf("收到信号 %v，正在关闭服务...", sig)
            // 在这里添加清理逻辑
            // cleanup()
            return
        }
    }
}
```

### 配置服务参数

可以修改服务相关常量来自定义服务信息：

**Windows (service_windows.go):**
```go
const (
    serviceName        = "YourServiceName"      // 服务名称
    serviceDisplayName = "Your Service Display" // 服务显示名称
    serviceDescription = "Your Service Description" // 服务描述
)
```

**Linux (service_linux.go):**
```go
const (
    serviceName = "yourservice"                    // 服务名称
    serviceFile = "/etc/systemd/system/yourservice.service" // 服务文件路径
)
```

### 添加命令行参数

在 `main()` 函数中添加新的命令行参数处理：

```go
func main() {
    if len(os.Args) > 1 {
        switch os.Args[1] {
        case "install":
            // 安装服务逻辑
        case "status":
            // 添加查看服务状态功能
            if err := checkServiceStatus(); err != nil {
                log.Fatalf("查看服务状态失败: %v", err)
            }
            return
        // 添加更多命令...
        }
    }
    
    runMainProgram()
}
```

## ⚠️ 注意事项

### 权限要求
- **Windows**: 必须以管理员身份运行服务管理命令
- **Linux**: 必须使用 `sudo` 或 root 权限执行服务管理操作

### 防火墙和安全
- 如果服务需要网络访问，请确保防火墙配置正确
- 建议运行在非特权用户下（可修改 Linux 服务配置中的 User 字段）

### 文件路径
- 确保可执行文件路径不包含空格或特殊字符
- Linux 系统建议将可执行文件放在 `/usr/local/bin/` 目录下

### 服务冲突
- 安装前请确保没有同名服务存在
- 如果需要重新安装，请先卸载旧服务

### 日志管理
- Windows 服务日志会记录在 Windows 事件查看器中
- Linux 服务日志通过 `journalctl` 查看
- 建议实现日志轮转以避免日志文件过大

### 性能注意事项
- 主循环默认每 30 秒执行一次，可根据需要调整时间间隔
- 避免在主循环中执行耗时操作，以免阻塞服务
- 对于 CPU 密集型任务，建议使用单独的 goroutine

## 🐛 故障排除

### 常见问题

1. **权限不足错误**
   ```
   解决方案：确保以管理员/root 权限运行
   ```

2. **服务已存在错误**
   ```
   解决方案：先卸载现有服务，再重新安装
   ```

3. **可执行文件路径错误**
   ```
   解决方案：使用绝对路径，避免路径中包含空格
   ```

4. **Linux systemd 服务无法启动**
   ```bash
   # 查看详细错误信息
   sudo journalctl -u goagent -n 50
   
   # 检查服务文件语法
   sudo systemd-analyze verify /etc/systemd/system/goagent.service
   ```

### 调试模式

在开发阶段，建议使用前台模式运行程序以便调试：

```bash
# 直接运行，不安装为服务
./GoAgent    # Windows
./goagent    # Linux
```

## 📝 版本历史

- **v1.0.0**: 初始版本，支持基本的跨平台服务管理功能

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request 来帮助改进这个项目。

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📄 许可证

本项目使用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 📞 联系方式

- 项目地址: [https://github.com/PeiKeSmart/GoAgent](https://github.com/PeiKeSmart/GoAgent)
- 问题反馈: [Issues](https://github.com/PeiKeSmart/GoAgent/issues)

---

⭐ 如果这个项目对你有帮助，请给一个星标支持！
