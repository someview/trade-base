package log

func (l *Logger) Debug(msg string, attrs ...Attr) {
	l.Log(DebugLevel, msg, attrs...)
}

func (l *Logger) Info(msg string, attrs ...Attr) {
	l.Log(InfoLevel, msg, attrs...)
}

func (l *Logger) Warn(msg string, attrs ...Attr) {
	l.Log(WarnLevel, msg, attrs...)
}

func (l *Logger) Error(msg string, attrs ...Attr) {
	l.Log(ErrorLevel, msg, attrs...)
}

func (l *Logger) Critical(msg string, attrs ...Attr) {
	l.Log(CriticalLevel, msg, attrs...)
}

func (l *Logger) Log(level LogLevel, msg string, attrs ...Attr) {
	if !l.LevelEnabled(level) {
		return
	}
	record := NewRecord(level, msg)
	record.attrs = append(record.attrs, attrs...)
	l.writeAsync(record)
}

func (l *Logger) LogWithPC(pc uintptr, level LogLevel, msg string, attrs ...Attr) {
	if !l.LevelEnabled(level) {
		return
	}
	record := NewRecord(ErrorLevel, msg)
	record.pc = pc
	record.attrs = append(record.attrs, attrs...)
	l.writeAsync(record)
}

func (l *Logger) LogRecord(record *Record) {
	if !l.LevelEnabled(record.level) {
		record.recycle()
		return
	}
	// 格式化日志记录
	// 发送到处理通道
	l.writeAsync(record)
}

// LogString 性能很差，不建议使用
// 建议采用Record.WriteString来代替
func (l *Logger) LogString(logL LogLevel, msg ...string) {
	str := buildString(logL, msg...)
	l.Log(logL, str)
}
