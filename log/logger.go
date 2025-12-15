package log

import (
	"archive/zip"
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

// LogLevel 定义日志级别
type LogLevel string

const (
	DebugLevel    LogLevel = "DEBUG"
	InfoLevel     LogLevel = "INFO"
	WarnLevel     LogLevel = "WARN"
	ErrorLevel    LogLevel = "ERROR"
	CriticalLevel LogLevel = "CRITICAL"
)

// UpdateNewDayRecord 表示一条日志记录

// Logger 主日志结构体
type Logger struct {
	config     Config
	file       *os.File
	filePath   string
	writer     *bufio.Writer
	encoder    Encoder
	recordChan chan *Record
	wg         sync.WaitGroup
	shutdown   atomic.Bool
	ctx        context.Context
	cancel     context.CancelFunc
}

type LoggerOption func(*Logger)

func WithEncoder(enc Encoder) LoggerOption {
	return func(l *Logger) {
		if enc != nil {
			l.encoder = enc
		}
	}
}

// NewLogger 创建一个新的日志器
func NewLogger(config Config, opts ...LoggerOption) (*Logger, error) {
	// 创建日志目录
	if err := os.MkdirAll(config.LogPath, 0755); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %w", err)
	}

	// 打开日志文件
	filePath := newLogFilePath(config.LogPath)
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
		encoder:    NewFlatEncoder(),
		recordChan: make(chan *Record, 2000),
		ctx:        ctx,
		cancel:     cancel,
	}

	for _, opt := range opts {
		opt(logger)
	}

	// 启动日志处理工作器
	logger.wg.Add(1)
	go logger.run()

	return logger, nil
}

// NewLoggerWithEncoder 创建一个使用自定义编码器的日志器
func NewLoggerWithEncoder(config Config, enc Encoder) (*Logger, error) {
	return NewLogger(config, WithEncoder(enc))
}

// NewRecord 创建一个新的日志记录，可供外部使用

// formatRecord 格式化日志记录到buffer
// run 处理日志的主循环
func (l *Logger) run() {
	defer l.wg.Done()
	ticker := time.NewTicker(l.config.FlushInterval)
	defer ticker.Stop()

	// 初始化每日定时压缩定时器（主协程管理）
	var dailyCompressTimer *time.Timer
	if l.config.CompressTime != "" { // 注意：配置字段已重命名
		dailyCompressTimer = l.calculateNextDailyCompressTimer()
	}

	for {
		select {
		case <-l.ctx.Done():
			l.processRemainLogs()
			return

		case record := <-l.recordChan:
			l.handleRecord(record)

		case <-ticker.C:
			_ = l.writer.Flush()
			l.checkFileRotation()

		// 每日定时压缩触发（主协程处理）
		case <-dailyCompressTimer.C:
			l.rotateFile()
			dailyCompressTimer = l.calculateNextDailyCompressTimer() // 重置定时器
		}
	}
}

// 计算下一次每日压缩的定时器（重命名后）
func (l *Logger) calculateNextDailyCompressTimer() *time.Timer {
	// 解析配置时间（如"08:00"）
	t, err := time.Parse("15:04", l.config.CompressTime) // 注意：配置字段已重命名
	if err != nil {
		return nil // 无效时间，不触发
	}

	now := time.Now().UTC()
	// 计算下一个触发时间（今日8点，若已过则为次日8点）
	nextTime := time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, now.Location())
	if now.After(nextTime) {
		nextTime = nextTime.Add(24 * time.Hour)
	}

	return time.NewTimer(nextTime.Sub(now))
}

// handleRecord 处理单条日志记录
func (l *Logger) handleRecord(record *Record) {
	if !record.encoded {
		record.buffer = l.encoder.Encode(record)
		record.encoded = true
	}
	// 写入缓冲区
	_, err := l.writer.Write(record.Buffer())
	record.recycle()
	if err != nil {
		slog.Error("Logger.handleRecord writer write", String("err", err.Error()))
		return
	}
	// 更新文件大小
	// 刷新检查
	if l.writer.Available() < 1024 {
		if err := l.writer.Flush(); err != nil {
			slog.Error("Logger.handleRecord writer flush", String("err", err.Error()))
		}
	}
}

// processRemainLogs 处理剩余日志
// processRemainLogs 处理剩余日志
func (l *Logger) processRemainLogs() {
	// 处理其他级别日志
	process := func(ch <-chan *Record) {
		for {
			select {
			case <-time.After(2 * time.Second): // todo 添加到配置文件, 设置一个优雅退出的时间, 避免程序退出时, channel里面的日志没有写完
				return
			case record := <-ch:
				l.handleRecord(record)
			}
		}
	}

	process(l.recordChan)
	if err := l.writer.Flush(); err != nil {
		slog.Error("Logger.processRemainLogs flush", String("err", err.Error()))
		return
	}
	if err := l.file.Close(); err != nil {
		slog.Error("Logger.rotateFile close old file", String("err", err.Error()))
		return
	}
	// 在测试或者开发环境中不需要进行压缩,在生产环境中通常需要进行压缩
	if l.config.CloseWithoutCompressed {
		return
	}
	compressFile(l.filePath) // 在程序关闭时,压缩
}

// checkFileRotation 检查是否需要轮转文件

// checkFileRotation 检查是否需要轮转文件
// checkFileRotation 检查是否需要轮转文件
func (l *Logger) checkFileRotation() {
	stat, err := l.file.Stat()
	if err != nil {
		slog.Error("Logger.checkFileRotation Stat", String("err", err.Error()))
		return
	}

	if stat.Size() > l.config.MaxSize {
		l.rotateFile()
	}
}

// rotateFile 轮转日志文件
// rotateFile 轮转日志文件（直接创建新文件版）
func (l *Logger) rotateFile() {
	// 刷新当前缓冲区（确保原文件数据写入磁盘）
	if err := l.writer.Flush(); err != nil {
		slog.Error("Logger.rotateFile flush", String("err", err.Error()))
		return
	}
	// 新增：强制将文件内容刷入磁盘（关键修复）
	if err := l.file.Sync(); err != nil {
		slog.Error("Logger.rotateFile sync old file", String("err", err.Error()))
		return
	}
	oldPath := l.filePath
	newPath := newLogFilePath(l.config.LogPath)
	if newPath == oldPath { // 时间相差太小,直接返回
		return
	}
	// 创建新日志文件（直接使用 newPath 作为新写入路径）
	newFile, err := os.OpenFile(newPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		slog.Error("Logger.rotateFile create new file", String("err", err.Error()))
		return
	}
	if err := l.file.Close(); err != nil {
		slog.Error("Logger.rotateFile close old file", String("err", err.Error()))
		return
	}
	// 触发压缩原文件（oldPath）
	l.compressFileAsync(oldPath)
	// 关闭原文件（释放句柄，允许后续压缩）

	// 更新 Logger 状态（后续日志写入 newPath）
	l.file = newFile
	l.filePath = newPath
	l.writer = bufio.NewWriterSize(newFile, l.config.BufferSize)
}

func newLogFilePath(dirPath string) string {
	return filepath.Join(dirPath, fmt.Sprintf("%s.debug.log", time.Now().Format("20060102-150405")))
}

func (l *Logger) compressFileAsync(path string) {
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		compressFile(path)
	}()
}

// compressFile 压缩指定的文件
// 通用压缩函数（合并实现）
// logger.go 修改compressFile函数
func compressFile(sourcePath string) {
	zipPath := sourcePath + ".zip"

	// 创建ZIP文件
	zipFile, err := os.Create(zipPath)
	if err != nil {
		slog.Error("Logger.compressFile create entry", String("err", err.Error()))
		return // 失败直接退出
	}

	// 创建ZIP写入器
	zipWriter := zip.NewWriter(zipFile)

	// 添加原始文件
	writer, err := zipWriter.Create(filepath.Base(sourcePath)) // 原代码忽略了此错误，需检查
	if err != nil {
		slog.Error("Logger.compressFile create writer", String("err", err.Error()))
		_ = zipWriter.Close() // 提前关闭避免资源泄漏
		_ = zipFile.Close()
		_ = os.Remove(zipPath) // 清理无效ZIP文件
		return
	}

	// 打开源文件
	srcFile, err := os.Open(sourcePath)
	if err != nil {
		slog.Error("Logger.compressFile open source", String("err", err.Error()))
		_ = zipWriter.Close()
		_ = zipFile.Close()
		_ = os.Remove(zipPath)
		return
	}
	defer srcFile.Close()

	// 复制内容并检查错误
	_, copyErr := io.Copy(writer, srcFile)

	// 仅当所有步骤成功时，才删除原文件（关键修复）
	if copyErr != nil {
		_ = zipWriter.Close()
		_ = zipFile.Close()
		_ = os.Remove(zipPath)
		slog.Error("Logger.compressFile copy content", String("err", copyErr.Error()))
		return
	}

	// 手动关闭ZIP写入器并检查错误（关键修复）
	if zipErr := zipWriter.Close(); zipErr != nil {
		slog.Error("Logger.compressFile zipWriter close", String("err", zipErr.Error()))
		_ = zipFile.Close()
		_ = os.Remove(zipPath) // 清理不完整的ZIP文件
		return
	}

	// 关闭ZIP文件（确保数据落盘）
	if err := zipFile.Close(); err != nil {
		slog.Error("Logger.compressFile zipFile close", String("err", err.Error()))
		_ = os.Remove(zipPath)
		return
	}

	// 删除源文件
	if err = os.Remove(sourcePath); err != nil {
		slog.Info("compressFile remove origin file error", String("err", err.Error()), String("path", sourcePath))
	}
}

// LevelEnabled 检查给定级别是否被启用
// 全局定义日志级别顺序映射，创建一次，多次使用
var levelOrder = map[LogLevel]int{
	DebugLevel:    0,
	InfoLevel:     1,
	WarnLevel:     2,
	ErrorLevel:    3,
	CriticalLevel: 4,
}

func (l *Logger) LevelEnabled(level LogLevel) bool {

	configLevel, exists := levelOrder[l.config.LogLevel]
	if !exists {
		configLevel = 0 // 默认为Debug级别
	}

	recordLevel, exists := levelOrder[level]
	if !exists {
		return false // 未知级别不记录
	}

	return recordLevel >= configLevel
}

// log 通用日志方法
func (l *Logger) writeAsync(r *Record) {
	// 发送到处理通道
	r.time = time.Now()
	select {
	case l.recordChan <- r:
		// 成功发送
	default:
		r.buffer = l.encoder.Encode(r)
		r.encoded = true
		l.recordChan <- r
		slog.Info("logger write loop is too busy")
	}
}

// Debug 输出Debug级别日志

// Close 关闭日志器
func (l *Logger) Close() {
	if !l.shutdown.CompareAndSwap(false, true) {
		return // 已经关闭
	}

	l.cancel()
	l.wg.Wait()
}

func (l *Logger) SetLogLevel(level LogLevel) {
	l.config.LogLevel = level
}

func (l *Logger) SetCloseWithoutCompressed(closeWithoutCompressed bool) {
	l.config.CloseWithoutCompressed = closeWithoutCompressed
}
