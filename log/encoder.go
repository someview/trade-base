package log

import "time"

// ValueEncoder 是专门用于值编码的接口，只暴露必要的写入方法
type ValueEncoder interface {
	WriteString(s string)
	WriteByte(c byte) error
	WriteInt64(n int64)
	WriteUint64(n uint64)
	WriteFloat64(f float64)
	WriteBool(b bool)
	WriteTime(t time.Time)
	WriteQuotedString(s string)

	// BeginObject 可以添加一些辅助方法
	WriteBeginObject()   // 写入 "{"
	WriteEndObject()     // 写入 "}"
	WriteBeginArray()    // 写入 "["
	WriteEndArray()      // 写入 "]"
	WriteKey(key string) // 写入 "key="
	WriteSeparator()     // 写入 ", "
}

// valueEncoderImpl 是 ValueEncoder 的内部实现
type valueEncoderImpl struct {
	r *Record
}

// 实现 ValueEncoder 接口
func (e *valueEncoderImpl) WriteString(s string) {
	e.r.writeString(s)
}

func (e *valueEncoderImpl) WriteByte(c byte) error {
	e.r.writeByte(c)
	return nil
}

func (e *valueEncoderImpl) WriteInt64(n int64) {
	e.r.writeInt64(n)
}

func (e *valueEncoderImpl) WriteUint64(n uint64) {
	e.r.writeUint64(n)
}

func (e *valueEncoderImpl) WriteFloat64(f float64) {
	e.r.writeFloat64(f)
}

func (e *valueEncoderImpl) WriteBool(b bool) {
	e.r.writeBool(b)
}

func (e *valueEncoderImpl) WriteTime(t time.Time) {
	e.r.writeTime(t)
}

func (e *valueEncoderImpl) WriteQuotedString(s string) {
	e.r.writeQuotedString(s)
}

// 辅助方法
func (e *valueEncoderImpl) WriteBeginObject() { e.r.writeByte('{') }
func (e *valueEncoderImpl) WriteEndObject()   { e.r.writeByte('}') }
func (e *valueEncoderImpl) WriteBeginArray()  { e.r.writeByte('[') }
func (e *valueEncoderImpl) WriteEndArray()    { e.r.writeByte(']') }
func (e *valueEncoderImpl) WriteKey(key string) {
	e.r.writeString(key)
	e.r.writeByte('=')
}
func (e *valueEncoderImpl) WriteSeparator() { e.r.writeString(", ") }

type Formatter interface {
	Format(e ValueEncoder)
}

func (r *Record) Encoder() ValueEncoder {
	return &valueEncoderImpl{r: r}
}
