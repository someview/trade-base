package log

import (
	"bytes"
	"fmt"
	"log/slog"
	"reflect"
	"strconv"
	"time"
)

type Attr = slog.Attr
type Value = slog.Value
type Kind = slog.Kind
type LogValuer = slog.LogValuer

// Err 创建一个错误属性
func Err(err error) Attr {
	return slog.String("error", errString(err))
}

// NamedError 创建一个命名错误属性
func NamedError(key string, err error) Attr {
	return slog.String(key, errString(err))
}

// errString 将错误转换为字符串，处理nil错误
func errString(err error) string {
	if err == nil {
		return "nil"
	}
	return err.Error()
}

// 其他函数直接使用slog的实现
var (
	Bool     = slog.Bool
	Int      = slog.Int
	Int64    = slog.Int64
	Uint64   = slog.Uint64
	Float64  = slog.Float64
	String   = slog.String
	Time     = slog.Time
	Duration = slog.Duration
	Group    = slog.Group
)

// Any 特殊处理错误类型
func Any(key string, val any) Attr {
	if err, ok := val.(error); ok {
		return NamedError(key, err)
	}
	return slog.Any(key, val)
}

func appendValueToBuffer(b *bytes.Buffer, v slog.Value) {
	switch v.Kind() {
	case slog.KindString:
		// 优化字符串写入
		b.WriteString(v.String())
	case slog.KindInt64:
		b.Write(strconv.AppendInt(b.Bytes(), v.Int64(), 10))

	case slog.KindUint64:
		b.Write(strconv.AppendUint(b.Bytes(), v.Uint64(), 10))

	case slog.KindFloat64:
		b.Write(strconv.AppendFloat(b.Bytes(), v.Float64(), 'g', -1, 64))

	case slog.KindBool:
		b.Write(strconv.AppendBool(b.Bytes(), v.Bool()))

	case slog.KindTime:
		// 复用buffer的内存空间
		tmp := b.Bytes()
		tmp = v.Time().AppendFormat(tmp, time.RFC3339Nano)
		b.Reset()
		b.Write(tmp)

	case slog.KindDuration:
		b.Write(strconv.AppendInt(b.Bytes(), v.Duration().Nanoseconds(), 10))

	case slog.KindGroup:
		formatGroupAppend(b, v.Group())

	case slog.KindAny:
		formatAnyAppend(b, v.Any())

	default:
		b.WriteString(fmt.Sprintf("%v", v.Any()))
	}
}

// 优化后的复杂类型处理
func formatAnyAppend(b *bytes.Buffer, val interface{}) {
	if val == nil {
		b.WriteString("nil")
		return
	}

	rv := reflect.ValueOf(val)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			b.WriteString("nil")
			return
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Struct:
		formatStructAppend(b, rv)
	case reflect.Map:
		formatMapAppend(b, rv)
	case reflect.Slice, reflect.Array:
		formatSliceAppend(b, rv)
	default:
		b.Write(strconv.AppendQuote(b.Bytes(), fmt.Sprintf("%v", val)))
	}
}

// 结构体格式化（优化版）
func formatStructAppend(b *bytes.Buffer, rv reflect.Value) {
	b.WriteByte('{')
	for i := 0; i < rv.NumField(); i++ {
		if i > 0 {
			b.WriteString(", ")
		}
		field := rv.Type().Field(i)
		b.WriteString(field.Name)
		b.WriteByte('=')
		appendValueToBuffer(b, slog.AnyValue(rv.Field(i).Interface()))
	}
	b.WriteByte('}')
}

// 映射格式化（优化版）
func formatMapAppend(b *bytes.Buffer, rv reflect.Value) {
	b.WriteByte('{')
	keys := rv.MapKeys()
	for i, k := range keys {
		if i > 0 {
			b.WriteString(", ")
		}
		appendValueToBuffer(b, slog.AnyValue(k.Interface()))
		b.WriteByte('=')
		appendValueToBuffer(b, slog.AnyValue(rv.MapIndex(k).Interface()))
	}
	b.WriteByte('}')
}

// 切片格式化（优化版）
func formatSliceAppend(b *bytes.Buffer, rv reflect.Value) {
	b.WriteByte('[')
	for i := 0; i < rv.Len(); i++ {
		if i > 0 {
			b.WriteString(", ")
		}
		appendValueToBuffer(b, slog.AnyValue(rv.Index(i).Interface()))
	}
	b.WriteByte(']')
}

// 组处理（优化版）
func formatGroupAppend(b *bytes.Buffer, attrs []slog.Attr) {
	b.WriteByte('{')
	for i, attr := range attrs {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(attr.Key)
		b.WriteByte('=')
		appendValueToBuffer(b, attr.Value)
	}
	b.WriteByte('}')
}
