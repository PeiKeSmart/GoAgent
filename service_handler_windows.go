//go:build windows

package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
)

var elog debug.Log

// isRunningAsService 检查是否作为 Windows 服务运行
func isRunningAsService() bool {
	isService, err := svc.IsWindowsService()
	if err != nil {
		return false
	}
	return isService
}

// windowsService 实现 Windows 服务接口
type windowsService struct {
	agent  *AgentService
	ctx    context.Context
	cancel context.CancelFunc
}

// Execute 实现 svc.Handler 接口
func (m *windowsService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}

	elog.Info(1, "服务正在启动")

	// 创建代理服务实例
	m.agent = NewAgentService()
	m.ctx, m.cancel = context.WithCancel(context.Background())

	// 在 goroutine 中启动代理服务
	go func() {
		if err := m.agent.Start(); err != nil {
			elog.Error(1, fmt.Sprintf("代理服务启动失败: %v", err))
		}
	}()

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
	elog.Info(1, "服务已启动")

loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
				// 测试 deadlock from https://code.google.com/p/winsvc/issues/detail?id=4
				time.Sleep(100 * time.Millisecond)
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				elog.Info(1, "服务正在停止")
				changes <- svc.Status{State: svc.StopPending}

				// 停止代理服务
				if m.agent != nil {
					m.agent.Stop()
				}
				if m.cancel != nil {
					m.cancel()
				}

				break loop
			case svc.Pause:
				changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
				elog.Info(1, "服务已暂停")
			case svc.Continue:
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
				elog.Info(1, "服务已恢复")
			default:
				elog.Error(1, fmt.Sprintf("意外的控制请求 #%d", c))
			}
		}
	}

	changes <- svc.Status{State: svc.Stopped}
	elog.Info(1, "服务已停止")
	return
}

// runAsWindowsService 作为 Windows 服务运行
func runAsWindowsService(name string, isDebug bool) {
	var err error
	if isDebug {
		elog = debug.New(name)
	} else {
		elog, err = eventlog.Open(name)
		if err != nil {
			return
		}
	}
	defer elog.Close()

	elog.Info(1, fmt.Sprintf("正在启动 %s 服务", name))
	run := svc.Run
	if isDebug {
		run = debug.Run
	}

	err = run(name, &windowsService{})
	if err != nil {
		elog.Error(1, fmt.Sprintf("服务运行失败: %v", err))
		return
	}

	elog.Info(1, fmt.Sprintf("%s 服务已停止", name))
}
