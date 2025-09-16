package global

import (
	"fmt"
	"sync"
	"time"

	"github.com/someview/trade-base/log"
)

// ServerShutdown 包含服务器关闭的信息
type ServerShutdown struct {
	Err  error
	Time time.Time
}

var (
	// 关闭通道，用于发送关闭请求
	shutdownCh = make(chan ServerShutdown, 1)
	// 互斥锁，防止并发关闭
	shutdownMutex sync.Mutex
	// 是否正在关闭
	isShuttingDown bool
)

// 已有的常量
const UpdateBtUseGoroutine bool = true
const Version string = "BetaGo v3"
const ExpireTime string = "2025-04-17T08:04:05"
const ValidIpList string = "54.199.241.133"

var startTime = time.Now().UTC()
var ServerIp string // 服务ip

func StartTime() time.Time {
	return startTime
}

func StartTimeUTCString() string {
	return startTime.Format(time.RFC3339)
}

// RequestShutdown 请求服务关闭
func RequestShutdown(err error) bool {
	shutdownMutex.Lock()
	defer shutdownMutex.Unlock()

	// 如果已经在关闭中，则忽略新请求
	if isShuttingDown {
		log.Debug("RequestShutdown 服务已经在关闭中，忽略新请求")
		return false
	}

	isShuttingDown = true
	log.Debug("RequestShutdown: " + err.Error())
	// 发送关闭请求到通道
	shutdownCh <- ServerShutdown{
		Err:  err,
		Time: time.Now().UTC(),
	}
	return true
}

// GetShutdownChannel 获取关闭通道
func GetShutdownChannel() <-chan ServerShutdown {
	return shutdownCh
}

func FormatExitMessage(moduleName string, desc string) string {
	return fmt.Sprintf("\n[EXIT] ServerId=%s | %s | %s\n",
		ServerId,
		moduleName,
		desc)
}

// ResetShutdownStatus 在关闭完成以后重置这个标志
func ResetShutdownStatus() {
	shutdownMutex.Lock()
	defer shutdownMutex.Unlock()
	// 重置标志和通道
	isShuttingDown = false
}
