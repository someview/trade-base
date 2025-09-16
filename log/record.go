// record.go
package log

import (
	"fmt"
	"log/slog"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	recordPool = sync.Pool{
		New: func() interface{} {
			return &Record{
				buffer: make([]byte, 0, 256),
				attrs:  make([]Attr, 0, 8),
			}
		},
	}
	defaultMaxDepth = 5
)

func init() {
	for i := 0; i < 100; i++ {
		recordPool.Put(&Record{
			buffer: make([]byte, 0, 256),
			attrs:  make([]Attr, 0, 8),
		})
	}
}

type Record struct {
	level       LogLevel
	msg         string
	buffer      []byte
	attrs       []Attr
	pc          uintptr
	encoded     bool
	time        time.Time
	recycleFunc func()
}

func NewRecord(level LogLevel, msg string) *Record {
	r := recordPool.Get().(*Record)
	r.level = level
	r.msg = msg
	r.attrs = r.attrs[:0]
	r.buffer = r.buffer[:0]
	r.pc = 0
	r.recycleFunc = nil
	return r
}

func NewRecordWithRecycleFunc(level LogLevel, msg string, recycleFunc func()) *Record {
	r := recordPool.Get().(*Record)
	r.level = level
	r.msg = msg
	r.recycleFunc = recycleFunc
	r.attrs = r.attrs[:0]
	r.buffer = r.buffer[:0]
	r.pc = 0
	return r
}

func NewPCRecord(level LogLevel, msg string, pc uintptr) *Record {
	r := recordPool.Get().(*Record)
	r.level = level
	r.msg = msg
	r.pc = pc
	r.buffer = r.buffer[:0]
	r.attrs = r.attrs[:0]
	r.encoded = false
	return r
}

func (r *Record) AddAttr(attrs ...Attr) {
	r.attrs = append(r.attrs, attrs...)
}

func (r *Record) recycle() {
	if r.recycleFunc != nil {
		r.recycleFunc()
	}
	r.recycleFunc = nil
	recordPool.Put(r)
}

func (r *Record) format() {
	r.buffer = r.buffer[:0]

	// 时间戳（精确到微秒）
	r.writeTime(r.time)
	r.writeByte('\t')

	// 日志级别
	r.writeString(string(r.level))
	r.writeByte('\t')

	// 函数调用路径
	if r.pc != 0 {
		file, line := pcLocation(r.pc)
		r.writeString(file)
		r.writeByte(':')
		r.writeInt64(int64(line))
		r.writeByte('\t')
	}

	// 日志消息
	if r.msg != "" {
		r.writeByte('\t')
		r.writeString(r.msg)
	}

	// 结构化属性
	for _, attr := range r.attrs {
		r.writeByte('\t')
		r.writeString(attr.Key)
		r.writeByte('=')
		r.appendValue(attr.Value)
	}
	r.writeByte('\n')
}

func (r *Record) Buffer() []byte {
	return r.buffer
}

// 核心值处理方法
func (r *Record) appendValue(v Value) {
	switch v.Kind() {
	case slog.KindAny:
		r.appendAnyValue(v.Any(), 0, defaultMaxDepth)
	case slog.KindString:
		r.writeString(v.String())
	case slog.KindInt64:
		r.writeInt64(v.Int64())
	case slog.KindUint64:
		r.writeUint64(v.Uint64())
	case slog.KindFloat64:
		r.writeFloat64(v.Float64())
	case slog.KindBool:
		r.writeBool(v.Bool())
	case slog.KindTime:
		r.writeTime(v.Time())
	case slog.KindDuration:
		r.writeDuration(v.Duration())
	case slog.KindGroup:
		r.writeGroup(v.Group())
	default:
		r.writeString(fmt.Sprintf("%+v", v.Any()))
	}
}

// 任意类型处理 - 修改后的版本
func (r *Record) appendAnyValue(val interface{}, currentDepth, maxDepth int) {
	// 立即检查深度
	if currentDepth >= maxDepth {
		slog.Info("Max depth reached, skipping remaining fields.", slog.Any("123", "456"))
		r.writeString("...")
		return
	}

	if val == nil {
		r.writeString("nil")
		return
	}

	// 支持自定义编码
	if encodeVal, ok := val.(Formatter); ok {
		encoder := r.Encoder()
		encodeVal.Format(encoder)
		return
	}

	rv := reflect.ValueOf(val)
	for rv.Kind() == reflect.Ptr && !rv.IsNil() {
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Struct:
		// 传递增加后的深度
		r.formatStruct(rv, currentDepth+1, maxDepth)
	case reflect.Map:
		// 传递增加后的深度
		r.formatMap(rv, currentDepth+1, maxDepth)
	case reflect.Slice, reflect.Array:
		// 传递增加后的深度
		r.formatSlice(rv, currentDepth+1, maxDepth)
	case reflect.String:
		r.writeQuotedString(rv.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r.writeInt64(rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		r.writeUint64(rv.Uint())
	case reflect.Float32, reflect.Float64:
		r.writeFloat64(rv.Float())
	case reflect.Bool:
		r.writeBool(rv.Bool())
	case reflect.Interface:
		if rv.IsNil() {
			r.writeString("nil")
		} else {
			// 接口类型，不增加深度
			r.appendAnyValue(rv.Interface(), currentDepth, maxDepth)
		}
	default:
		r.writeString(fmt.Sprintf("%T(%v)", val, val))
	}
}

// 结构体格式化 - 修改后的版本
func (r *Record) formatStruct(rv reflect.Value, currentDepth, maxDepth int) {
	// 添加结构体类型名称（如果有）
	typeName := rv.Type().Name()
	if typeName != "" {
		r.writeString(typeName)
	}

	r.writeByte('{')

	// 直接判断当前深度
	if currentDepth >= maxDepth {
		r.writeString("...")
	} else {
		var hasField bool
		for i := 0; i < rv.NumField(); i++ {
			field := rv.Type().Field(i)
			if !field.IsExported() {
				continue
			}

			if hasField {
				r.writeString(", ")
			}
			hasField = true

			r.writeString(field.Name)
			r.writeByte('=')
			// 递归调用，已经在appendAnyValue中增加深度
			r.appendAnyValue(rv.Field(i).Interface(), currentDepth, maxDepth)
		}
	}

	r.writeByte('}')
}

// Map格式化 - 修改后的版本
func (r *Record) formatMap(rv reflect.Value, currentDepth, maxDepth int) {
	r.writeString("iter{") // 使用 "iter{" 作为Map的开始标记

	// 直接判断当前深度
	if currentDepth >= maxDepth {
		r.writeString("...")
	} else {
		written := 0
		keys := rv.MapKeys()
		for _, k := range keys {
			if !k.CanInterface() || !rv.MapIndex(k).CanInterface() {
				continue
			}

			if written > 0 {
				r.writeString(", ")
			}

			// 递归调用，已经在appendAnyValue中增加深度
			r.appendAnyValue(k.Interface(), currentDepth, maxDepth)
			r.writeByte('=')
			r.appendAnyValue(rv.MapIndex(k).Interface(), currentDepth, maxDepth)

			written++
		}
	}

	r.writeByte('}')
}

// 切片格式化 - 修改后的版本
func (r *Record) formatSlice(rv reflect.Value, currentDepth, maxDepth int) {
	r.writeByte('[')

	// 直接判断当前深度
	if currentDepth >= maxDepth {
		r.writeString("...")
	} else {
		written := 0
		for i := 0; i < rv.Len(); i++ {
			elem := rv.Index(i)
			if !elem.CanInterface() {
				continue
			}

			if written > 0 {
				r.writeString(", ")
			}

			// 递归调用，已经在appendAnyValue中增加深度
			r.appendAnyValue(elem.Interface(), currentDepth, maxDepth)
			written++
		}
	}

	r.writeByte(']')
}

func (r *Record) writeTime(t time.Time) {
	r.buffer = t.UTC().AppendFormat(r.buffer, time.RFC3339Nano)
}

func (r *Record) writeDuration(d time.Duration) {
	r.buffer = strconv.AppendInt(r.buffer, int64(d), 10)
}

func (r *Record) writeGroup(attrs []slog.Attr) {
	r.writeByte('{')
	for i, attr := range attrs {
		if i > 0 {
			r.writeByte(' ')
		}
		r.writeString(attr.Key)
		r.writeByte('=')
		r.appendValue(attr.Value)
	}
	r.writeByte('}')
}

func (r *Record) writeQuotedString(s string) {
	if needsQuote(s) {
		r.buffer = strconv.AppendQuote(r.buffer, s)
	} else {
		r.writeString(s)
	}
}

func (r *Record) writeString(s string) {
	r.buffer = append(r.buffer, s...)
}

func (r *Record) writeByte(c byte) {
	r.buffer = append(r.buffer, c)
}

func (r *Record) writeInt64(n int64) {
	r.buffer = strconv.AppendInt(r.buffer, n, 10)
}

func (r *Record) writeUint64(n uint64) {
	r.buffer = strconv.AppendUint(r.buffer, n, 10)
}

func (r *Record) writeFloat64(f float64) {
	r.buffer = strconv.AppendFloat(r.buffer, f, 'g', -1, 64)
}

func (r *Record) writeBool(b bool) {
	r.buffer = strconv.AppendBool(r.buffer, b)
}

func (r *Record) WriteValue(value interface{}) {
	// 不涉及深度控制，只用于外部调用
	r.appendAnyValue(value, 0, defaultMaxDepth)
}

type LogFormatter interface {
	// 不需要深度参数
	FormatLog(r *Record)
}

func needsQuote(s string) bool {
	return strings.ContainsAny(s, " \t\n\r\"=,{}[]")
}
