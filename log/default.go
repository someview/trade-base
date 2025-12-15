package log

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync/atomic"
)

// 以下保持原有全局函数兼容性
var defaultLogger atomic.Pointer[Logger]

func GlobalLogger() *Logger {
	return defaultLogger.Load()
	// 快速路径
	//if logger := defaultLogger.Load(); logger != nil {
	//	return logger
	//}
	//// 慢路径
	//newLogger, err := NewLogger(NewDefaultConfig())
	//if err != nil {
	//	slog.Info("init global logger error", "err", err)
	//	return nil
	//}
	//if defaultLogger.CompareAndSwap(nil, newLogger) {
	//	return newLogger
	//}
	//return defaultLogger.Load()
}

// SetGlobalLogger 允许外部设置全局的logger
func SetGlobalLogger(logger *Logger) {
	defaultLogger.Store(logger)
}

// 使用方式与方案一相同

func Debug(msg string, attr ...Attr) {
	GlobalLogger().Log(DebugLevel, msg, attr...)
}

func Info(msg string, attr ...Attr) {
	GlobalLogger().Log(InfoLevel, msg, attr...)
}

func Warn(msg string, attr ...Attr) {
	GlobalLogger().Log(WarnLevel, msg, attr...)
}

func Error(msg string, attr ...Attr) {
	GlobalLogger().Log(ErrorLevel, msg, attr...)
}

func Critical(msg string, attr ...Attr) {
	GlobalLogger().Log(CriticalLevel, msg, attr...)
}

func DebugWithPC(pc uintptr, msg string, attr ...Attr) {
	GlobalLogger().LogWithPC(pc, DebugLevel, msg, attr...)
}

func InfoWithPC(pc uintptr, msg string, attr ...Attr) {
	GlobalLogger().LogWithPC(pc, InfoLevel, msg, attr...)
}

func WarnWithPC(pc uintptr, msg string, attr ...Attr) {
	GlobalLogger().LogWithPC(pc, WarnLevel, msg, attr...)
}

func ErrorWithPC(pc uintptr, msg string, attr ...Attr) {
	GlobalLogger().LogWithPC(pc, ErrorLevel, msg, attr...)
}

func CriticalWithPC(pc uintptr, msg string, attr ...Attr) {
	GlobalLogger().LogWithPC(pc, CriticalLevel, msg, attr...)
}

func LogWithPC(pc uintptr, level LogLevel, msg string, attr ...Attr) {
	GlobalLogger().LogWithPC(pc, level, msg, attr...)
}

func LogRecord(r *Record) {
	GlobalLogger().LogRecord(r)
}

func LogString(level LogLevel, msg ...string) {
	str := buildString(level, msg...)
	GlobalLogger().LogString(level, str)
}

func NewGlobalTestLogger() *Logger {
	_, file, _, _ := runtime.Caller(1)
	os.Setenv("LOG_PATH", filepath.Dir(file))
	conf := NewDefaultConfig()
	conf.CloseWithoutCompressed = true
	logger, err := newTestLogger(conf)
	if err != nil {
		panic(err)
	}
	SetGlobalLogger(logger)
	return logger
}

func newTestLogger(config Config) (*Logger, error) {
	// 创建日志目录
	if err := os.MkdirAll(config.LogPath, 0755); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %w", err)
	}

	// 打开日志文件
	filePath := filepath.Join(config.LogPath, "test.log")
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("打开日志文件失败: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// 创建日志器
	logger := &Logger{
		config:     config,
		file:       file,
		filePath:   filePath,
		writer:     bufio.NewWriterSize(file, config.BufferSize),
		encoder:    NewJSONEncoder(),
		recordChan: make(chan *Record, 2000),
		ctx:        ctx,
		cancel:     cancel,
	}

	// 启动日志处理工作器
	logger.wg.Add(1)
	go logger.run()
	return logger, nil
}
