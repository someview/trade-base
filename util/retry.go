package util

import (
	"fmt"
	"time"
)

// ConstantRetry 等间隔重试函数
// attempts: 最大尝试次数(包含首次执行)
// interval: 每次重试间隔
// fn: 需要重试的函数
func ConstantRetry(attempts int, interval time.Duration, fn func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		if err = fn(); err == nil {
			return nil
		}
		if i < attempts-1 { // 最后一次不需要sleep
			time.Sleep(interval)
		}
	}
	return fmt.Errorf("after %d attempts: %v", attempts, err)
}

type RetryController struct {
	MaxAttempts    int
	CurrentAttempt int
}

func NewRetryController(max int) *RetryController {
	return &RetryController{
		MaxAttempts:    max,
		CurrentAttempt: 1,
	}
}

func (rc *RetryController) Next() bool {
	rc.CurrentAttempt++
	return rc.CurrentAttempt <= rc.MaxAttempts
}

func (rc *RetryController) Wait() {
	if rc.CurrentAttempt > 1 {
		sleepDuration := time.Duration(rc.CurrentAttempt*10) * time.Millisecond
		time.Sleep(sleepDuration)
	}
}
