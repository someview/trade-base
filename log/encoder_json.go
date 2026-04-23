package log

import (
	stdjson "encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"
)

// JSONEncoder 以 JSON 格式编码日志
type JSONEncoder struct {
	TimeKey   string
	LevelKey  string
	MsgKey    string
	CallerKey string
}

func NewJSONEncoder() *JSONEncoder {
	return &JSONEncoder{
		TimeKey:   "ts",
		LevelKey:  "level",
		MsgKey:    "msg",
		CallerKey: "caller",
	}
}

func (e *JSONEncoder) Encode(r *Record) []byte {
	slog.Info(
		"JSONEncoder Encode begin",
		slog.String("level", string(r.level)),
		slog.String("msg", r.msg),
		slog.Time("time", r.time),
		slog.Uint64("pc", uint64(r.pc)),
		slog.Int("attrs_len", len(r.attrs)),
	)

	// 直接写入 Record 的缓冲，避免额外分配
	r.buffer = r.buffer[:0]
	r.writeByte('{')
	writeJSONFieldString(r, e.TimeKey, r.time.UTC().Format(time.RFC3339Nano))
	// level
	r.writeByte(',')
	writeJSONFieldString(r, e.LevelKey, string(r.level))
	// msg
	if r.msg != "" {
		r.writeByte(',')
		writeJSONFieldString(r, e.MsgKey, r.msg)
	}
	// caller
	if r.pc != 0 {
		file, line := pcLocation(r.pc)
		r.writeByte(',')
		writeJSONFieldKeyPrefix(r, e.CallerKey)
		r.buffer = strconv.AppendQuote(r.buffer, fmt.Sprintf("%s_%d", file, line))
	}
	// attrs
	for _, attr := range r.attrs {
		r.writeByte(',')
		writeJSONFieldKeyPrefix(r, attr.Key)
		appendSlogValueAsJSONToRecord(r, attr.Value)
	}
	r.writeByte('}')
	r.writeByte('\n')
	out := r.Buffer()
	slog.Info("JSONEncoder Encode end", slog.Int("bytes_len", len(out)))
	return out
}

// JSON 值编码器，供外部 Formatter 使用以生成 JSON 结构
type jsonValueEncoder struct {
	r *Record
}

func newJSONValueEncoder(r *Record) *jsonValueEncoder { return &jsonValueEncoder{r: r} }

func (e *jsonValueEncoder) WriteString(s string)   { e.r.buffer = strconv.AppendQuote(e.r.buffer, s) }
func (e *jsonValueEncoder) WriteByte(c byte) error { e.r.writeByte(c); return nil }
func (e *jsonValueEncoder) WriteInt64(n int64)     { e.r.writeInt64(n) }
func (e *jsonValueEncoder) WriteUint64(n uint64)   { e.r.writeUint64(n) }
func (e *jsonValueEncoder) WriteFloat64(f float64) { e.r.writeFloat64(f) }
func (e *jsonValueEncoder) WriteBool(b bool)       { e.r.writeBool(b) }
func (e *jsonValueEncoder) WriteTime(t time.Time) {
	e.r.buffer = strconv.AppendQuote(e.r.buffer, t.UTC().Format(time.RFC3339Nano))
}
func (e *jsonValueEncoder) WriteQuotedString(s string) {
	e.r.buffer = strconv.AppendQuote(e.r.buffer, s)
}
func (e *jsonValueEncoder) WriteBeginObject() { e.r.writeByte('{') }
func (e *jsonValueEncoder) WriteEndObject()   { e.r.writeByte('}') }
func (e *jsonValueEncoder) WriteBeginArray()  { e.r.writeByte('[') }
func (e *jsonValueEncoder) WriteEndArray()    { e.r.writeByte(']') }
func (e *jsonValueEncoder) WriteKey(key string) {
	e.r.writeByte('"')
	e.r.writeString(key)
	e.r.writeByte('"')
	e.r.writeByte(':')
}
func (e *jsonValueEncoder) WriteSeparator() { e.r.writeByte(',') }

// 控制字段分隔
func writeJSONFieldString(r *Record, key, val string) {
	writeJSONFieldKeyPrefix(r, key)
	r.buffer = strconv.AppendQuote(r.buffer, val)
}

func writeJSONFieldKeyPrefix(r *Record, key string) {
	r.writeByte('"')
	r.writeString(key)
	r.writeByte('"')
	r.writeByte(':')
}

func appendGroupToRecord(r *Record, attrs []slog.Attr) {
	r.writeByte('{')
	for i, attr := range attrs {
		if i > 0 {
			r.writeByte(',')
		}
		r.writeByte('"')
		r.writeString(attr.Key)
		r.writeByte('"')
		r.writeByte(':')
		appendSlogValueAsJSONToRecord(r, attr.Value)
	}
	r.writeByte('}')
}

func appendAnyToRecord(r *Record, v any) {
	if b, err := stdjson.Marshal(v); err == nil {
		r.buffer = append(r.buffer, b...)
		return
	}
	r.buffer = strconv.AppendQuote(r.buffer, fmt.Sprintf("%v", v))
}

// 以 JSON 格式写 slog.Value
func appendSlogValueAsJSONToRecord(r *Record, v slog.Value) {
	switch v.Kind() {
	case slog.KindString:
		r.buffer = strconv.AppendQuote(r.buffer, v.String())
	case slog.KindInt64:
		r.writeInt64(v.Int64())
	case slog.KindUint64:
		r.writeUint64(v.Uint64())
	case slog.KindFloat64:
		r.writeFloat64(v.Float64())
	case slog.KindBool:
		r.writeBool(v.Bool())
	case slog.KindTime:
		r.buffer = strconv.AppendQuote(r.buffer, v.Time().UTC().Format(time.RFC3339Nano))
	case slog.KindDuration:
		r.writeInt64(v.Duration().Nanoseconds())
	case slog.KindGroup:
		appendGroupToRecord(r, v.Group())
	case slog.KindAny:
		any := v.Any()
		if f, ok := any.(Formatter); ok {
			jenc := newJSONValueEncoder(r)
			f.Format(jenc)
			return
		}
		appendAnyToRecord(r, any)
	default:
		appendAnyToRecord(r, v.Any())
	}
}
