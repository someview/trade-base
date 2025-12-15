package log

import "time"

type valueEncoderImpl struct {
	r *Record
}

func (e *valueEncoderImpl) WriteString(s string)       { e.r.writeString(s) }
func (e *valueEncoderImpl) WriteByte(c byte) error     { e.r.writeByte(c); return nil }
func (e *valueEncoderImpl) WriteInt64(n int64)         { e.r.writeInt64(n) }
func (e *valueEncoderImpl) WriteUint64(n uint64)       { e.r.writeUint64(n) }
func (e *valueEncoderImpl) WriteFloat64(f float64)     { e.r.writeFloat64(f) }
func (e *valueEncoderImpl) WriteBool(b bool)           { e.r.writeBool(b) }
func (e *valueEncoderImpl) WriteTime(t time.Time)      { e.r.writeTime(t) }
func (e *valueEncoderImpl) WriteQuotedString(s string) { e.r.writeQuotedString(s) }
func (e *valueEncoderImpl) WriteBeginObject()          { e.r.writeByte('{') }
func (e *valueEncoderImpl) WriteEndObject()            { e.r.writeByte('}') }
func (e *valueEncoderImpl) WriteBeginArray()           { e.r.writeByte('[') }
func (e *valueEncoderImpl) WriteEndArray()             { e.r.writeByte(']') }
func (e *valueEncoderImpl) WriteKey(key string)        { e.r.writeString(key); e.r.writeByte('=') }
func (e *valueEncoderImpl) WriteSeparator()            { e.r.writeString(", ") }

func (r *Record) Encoder() ValueEncoder { return &valueEncoderImpl{r: r} }

// FlatEncoder 整条记录的平铺文本编码器（调用已有格式化逻辑）
type FlatEncoder struct{}

func NewFlatEncoder() *FlatEncoder { return &FlatEncoder{} }

func (e *FlatEncoder) Encode(r *Record) []byte {
	r.format()
	return r.Buffer()
}
