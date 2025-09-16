package log

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func buildString(logL LogLevel, msg ...string) string {
	builder := strings.Builder{}
	builder.WriteString("\n")
	builder.WriteString(time.Now().UTC().Format("2006-01-02T15:04:05.000000Z"))
	builder.WriteString("\t")
	builder.WriteString(string(logL))
	builder.WriteString("\t")

	_, file, line, _ := runtime.Caller(2)
	relativePath := trimPath(file)
	builder.WriteString(relativePath + ":" + strconv.Itoa(line))
	for _, s := range msg {
		builder.WriteString("\t")
		builder.WriteString(s)
	}

	return builder.String()
}

func LogPath() string {
	if path := os.Getenv("LOG_PATH"); path != "" {
		return path
	}
	// 获取当前可执行文件路径
	ext, err := os.Executable()
	if err != nil {
		panic(fmt.Sprintf("获取可执行文件路径失败: %v", err))
	}

	// 尝试解析符号链接获取真实路径
	// 若解析失败，使用原始路径的目录作为回退方案
	realPath, err := filepath.EvalSymlinks(ext)
	if err != nil {
		return filepath.Dir(ext)
	}
	return filepath.Dir(realPath)
}

func SetLogPath(path string) {
	os.Setenv("LOG_PATH", path)
}

func SetLogPathWithCallerPath() {
	_, file, _, _ := runtime.Caller(1)
	os.Setenv("LOG_PATH", filepath.Dir(file))
}
