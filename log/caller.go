package log

import (
	"reflect"
	"runtime"
	"strings"
	"sync"
)

// location 存储位置信息
type location struct {
	file string
	line int
}

// locationCache 缓存PC到位置信息的映射

// 全局缓存实例
var _callerCache = &sync.Map{}

// pcLocation 根据PC获取位置信息，不使用runtime.Caller
func pcLocation(pc uintptr) (file string, line int) {
	// 查找缓存
	if v, ok := _callerCache.Load(pc); ok {
		loc := v.(location)
		return loc.file, loc.line
	}
	// 缓存未命中，通过函数信息获取位置
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown", 0
	}

	// 获取文件和行号
	file, line = fn.FileLine(pc)

	// 简化文件路径
	file = trimPath(file)

	// 缓存结果
	_callerCache.Store(pc, location{
		file: file,
		line: line,
	})
	return file, line
}

// 回退方案，仅当PC无效时使用
func fallbackCaller(skip int) (file string, line int) {
	// 这是最后的回退方案，只在PC无效时使用
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown", 0
	}

	// 简化文件路径
	file = trimPath(file)
	return file, line
}

// trimPath 简化文件路径
func trimPath(path string) string {
	// 简化Go标准库路径
	if idx := strings.LastIndex(path, "/src/"); idx != -1 {
		path = path[idx+5:]
	}

	// 保留最后3个路径段
	segments := strings.Split(path, "/")
	if len(segments) <= 3 {
		return path
	}

	return strings.Join(segments[len(segments)-3:], "/")
}

func CallerPC(f any) uintptr {
	return reflect.ValueOf(f).Pointer()
}
