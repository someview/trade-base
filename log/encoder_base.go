package log

import (
	"reflect"
	"time"
)

// ValueEncoder 是专门用于值编码的接口，只暴露必要的写入方法
// FlatEncoder 和 JSONEncoder 都会提供自己的实现
type ValueEncoder interface {
	WriteString(s string)
	WriteByte(c byte) error
	WriteInt64(n int64)
	WriteUint64(n uint64)
	WriteFloat64(f float64)
	WriteBool(b bool)
	WriteTime(t time.Time)
	WriteQuotedString(s string)

	// 结构化辅助
	WriteBeginObject()   // "{"
	WriteEndObject()     // "}"
	WriteBeginArray()    // "["
	WriteEndArray()      // "]"
	WriteKey(key string) // key= ,or key:
	WriteSeparator()     // ", "
}

// Formatter 允许外部类型自定义编码（依赖具体的 ValueEncoder 实现）
type Formatter interface {
	Format(e ValueEncoder)
}

// Encoder 为整条日志记录的编码器接口
type Encoder interface {
	Encode(r *Record) []byte
}

// isReflectNil checks whether v is invalid or holds a typed nil
func isReflectNil(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}
	switch v.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func:
		return v.IsNil()
	case reflect.Interface:
		if v.IsNil() {
			return true
		}
		return isReflectNil(v.Elem())
	default:
		return false
	}
}
