package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// LogLevel 日志级别
type LogLevel int

const (
	LogLevelDebug LogLevel = iota // 调试信息
	LogLevelInfo                  // 一般信息
	LogLevelWarn                  // 警告信息
	LogLevelError                 // 错误信息
)

// AgentService 星尘代理服务核心业务逻辑
type AgentService struct {
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	config   *AgentConfig
	logger   *log.Logger
	logFile  *os.File
	logLevel LogLevel
}

// AgentConfig 代理配置
type AgentConfig struct {
	// 服务端连接配置
	ServerURL  string `json:"server_url"`
	ServerPort int    `json:"server_port"`

	// 节点配置
	NodeID     string `json:"node_id"`
	NodeName   string `json:"node_name"`
	NodeRegion string `json:"node_region"`

	// 监控配置
	MonitorInterval time.Duration `json:"monitor_interval"`
	ReportInterval  time.Duration `json:"report_interval"`

	// 资源限制
	MaxCPU    float64 `json:"max_cpu"`
	MaxMemory int64   `json:"max_memory"`
}

// NewAgentService 创建新的代理服务实例
func NewAgentService() *AgentService {
	ctx, cancel := context.WithCancel(context.Background())

	service := &AgentService{
		ctx:      ctx,
		cancel:   cancel,
		config:   getDefaultConfig(),
		logLevel: LogLevelInfo, // 默认信息级别
	}

	// 初始化日志
	service.initLogger()

	return service
}

// initLogger 初始化日志系统
func (a *AgentService) initLogger() {
	// 获取可执行文件目录
	exePath, err := os.Executable()
	if err != nil {
		log.Printf("❌ 无法获取可执行文件路径: %v", err)
		a.logger = log.New(os.Stdout, "[GoAgent] ", log.LstdFlags|log.Lshortfile)
		return
	}

	// 在可执行文件同目录创建日志文件
	exeDir := filepath.Dir(exePath)
	logPath := filepath.Join(exeDir, "goagent.log")

	// 检查日志文件大小，如果超过10MB则轮转
	if info, err := os.Stat(logPath); err == nil {
		if info.Size() > 10*1024*1024 { // 10MB
			backupPath := filepath.Join(exeDir, fmt.Sprintf("goagent.log.%s", time.Now().Format("20060102-150405")))
			os.Rename(logPath, backupPath)
		}
	}

	// 打开或创建日志文件
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("❌ 无法创建日志文件 %s: %v", logPath, err)
		a.logger = log.New(os.Stdout, "[GoAgent] ", log.LstdFlags|log.Lshortfile)
		return
	}

	a.logFile = logFile

	// 检查是否作为Windows服务运行
	var writer io.Writer
	if isRunningAsService() {
		// 服务模式：只输出到文件
		writer = logFile
		a.logger = log.New(writer, "[GoAgent Service] ", log.LstdFlags|log.Lshortfile)
	} else {
		// 普通模式：同时输出到控制台和文件
		multiWriter := io.MultiWriter(os.Stdout, logFile)
		writer = multiWriter
		a.logger = log.New(writer, "[GoAgent] ", log.LstdFlags|log.Lshortfile)
	}

	a.logger.Printf("📝 日志系统已初始化，日志文件: %s", logPath)
	if isRunningAsService() {
		a.logger.Printf("🔧 运行模式: Windows服务 (仅文件日志)")
	} else {
		a.logger.Printf("🔧 运行模式: 普通程序 (控制台+文件日志)")
	}
} // getDefaultConfig 获取默认配置
func getDefaultConfig() *AgentConfig {
	return &AgentConfig{
		ServerURL:       "https://star.newlifex.com",
		ServerPort:      443,
		NodeID:          generateNodeID(),
		NodeName:        getHostname(),
		NodeRegion:      "default",
		MonitorInterval: 30 * time.Second,
		ReportInterval:  60 * time.Second,
		MaxCPU:          80.0,
		MaxMemory:       2 * 1024 * 1024 * 1024, // 2GB
	}
}

// Start 启动代理服务
func (a *AgentService) Start() error {
	a.logInfo("🚀 星尘代理服务启动中...")
	a.logInfo("节点ID: %s", a.config.NodeID)
	a.logInfo("节点名称: %s", a.config.NodeName)
	a.logInfo("服务端: %s:%d", a.config.ServerURL, a.config.ServerPort)
	a.logInfo("监控间隔: %v, 报告间隔: %v", a.config.MonitorInterval, a.config.ReportInterval)

	// 启动各个服务组件
	a.wg.Add(4)

	// 1. 节点监控服务
	go a.runNodeMonitor()

	// 2. 服务端连接管理
	go a.runServerConnection()

	// 3. 资源调度器
	go a.runResourceScheduler()

	// 4. 状态报告器
	go a.runStatusReporter()

	a.logInfo("✅ 星尘代理服务已启动")

	// 等待所有服务停止
	a.wg.Wait()

	a.logInfo("🛑 星尘代理服务已停止")
	return nil
}

// Stop 停止代理服务
func (a *AgentService) Stop() {
	a.logger.Printf("🛑 正在停止星尘代理服务...")
	a.cancel()
	a.wg.Wait()

	// 关闭日志文件
	if a.logFile != nil {
		a.logFile.Close()
	}
}

// runNodeMonitor 运行节点监控
func (a *AgentService) runNodeMonitor() {
	defer a.wg.Done()

	ticker := time.NewTicker(a.config.MonitorInterval)
	defer ticker.Stop()

	a.logger.Printf("📊 节点监控服务已启动 (间隔: %v)", a.config.MonitorInterval)

	for {
		select {
		case <-a.ctx.Done():
			a.logger.Printf("📊 节点监控服务已停止")
			return
		case <-ticker.C:
			a.collectNodeMetrics()
		}
	}
}

// runServerConnection 运行服务端连接管理
func (a *AgentService) runServerConnection() {
	defer a.wg.Done()

	a.logger.Printf("🔗 服务端连接管理器已启动")

	// 连接重试逻辑
	for {
		select {
		case <-a.ctx.Done():
			a.logger.Printf("🔗 服务端连接管理器已停止")
			return
		default:
			if err := a.connectToServer(); err != nil {
				a.logger.Printf("❌ 连接服务端失败: %v, 5秒后重试", err)
				time.Sleep(5 * time.Second)
			} else {
				a.logger.Printf("✅ 已连接到服务端")
				// 保持连接，监听服务端指令
				a.handleServerCommands()
			}
		}
	}
}

// runResourceScheduler 运行资源调度器
func (a *AgentService) runResourceScheduler() {
	defer a.wg.Done()

	a.logger.Printf("⚡ 资源调度器已启动")

	for {
		select {
		case <-a.ctx.Done():
			a.logger.Printf("⚡ 资源调度器已停止")
			return
		default:
			// 处理调度任务
			a.processScheduledTasks()
			time.Sleep(10 * time.Second)
		}
	}
}

// runStatusReporter 运行状态报告器
func (a *AgentService) runStatusReporter() {
	defer a.wg.Done()

	ticker := time.NewTicker(a.config.ReportInterval)
	defer ticker.Stop()

	a.logger.Printf("📡 状态报告器已启动 (间隔: %v)", a.config.ReportInterval)

	for {
		select {
		case <-a.ctx.Done():
			a.logger.Printf("📡 状态报告器已停止")
			return
		case <-ticker.C:
			a.reportNodeStatus()
		}
	}
}

// collectNodeMetrics 收集节点指标
func (a *AgentService) collectNodeMetrics() {
	// 模拟收集CPU、内存、磁盘等指标
	a.logger.Printf("📊 收集节点指标: CPU=25%%, 内存=60%%, 磁盘=45%%")
}

// connectToServer 连接到服务端
func (a *AgentService) connectToServer() error {
	// 模拟连接逻辑
	a.logger.Printf("🔗 正在连接服务端 %s:%d...", a.config.ServerURL, a.config.ServerPort)

	// 模拟连接延迟
	time.Sleep(2 * time.Second)

	// 模拟连接成功/失败
	// 这里可以实现真实的HTTP/gRPC/WebSocket连接
	return nil
}

// handleServerCommands 处理服务端指令
func (a *AgentService) handleServerCommands() {
	a.logger.Printf("👂 开始监听服务端指令...")

	// 模拟保持连接并处理指令
	for {
		select {
		case <-a.ctx.Done():
			return
		default:
			// 模拟接收和处理服务端指令
			time.Sleep(30 * time.Second)
			a.logger.Printf("📨 处理服务端指令: 心跳检测")
		}
	}
}

// processScheduledTasks 处理调度任务
func (a *AgentService) processScheduledTasks() {
	// 模拟任务调度逻辑
	a.logger.Printf("⚡ 检查待调度任务...")
}

// reportNodeStatus 报告节点状态
func (a *AgentService) reportNodeStatus() {
	a.logger.Printf("📡 向服务端报告节点状态")
}

// 辅助函数
func generateNodeID() string {
	return fmt.Sprintf("node-%d", time.Now().Unix())
}

func getHostname() string {
	// 简化实现，实际应该获取真实主机名
	return "localhost"
}

// 日志方法
func (a *AgentService) logDebug(format string, v ...interface{}) {
	if a.logLevel <= LogLevelDebug && a.logger != nil {
		a.logger.Printf("[DEBUG] "+format, v...)
	}
}

func (a *AgentService) logInfo(format string, v ...interface{}) {
	if a.logLevel <= LogLevelInfo && a.logger != nil {
		a.logger.Printf("[INFO] "+format, v...)
	}
}

func (a *AgentService) logWarn(format string, v ...interface{}) {
	if a.logLevel <= LogLevelWarn && a.logger != nil {
		a.logger.Printf("[WARN] "+format, v...)
	}
}

func (a *AgentService) logError(format string, v ...interface{}) {
	if a.logLevel <= LogLevelError && a.logger != nil {
		a.logger.Printf("[ERROR] "+format, v...)
	}
}
