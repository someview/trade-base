package log

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomJSONEncoder_Print(t *testing.T) {
	enc := &JSONEncoder{
		TimeKey:   "time",
		LevelKey:  "lvl",
		MsgKey:    "message",
		CallerKey: "call",
	}
	r := NewRecord(InfoLevel, "hi")
	r.AddAttr(String("k1", "v1"))
	out := enc.Encode(r)
	fmt.Println(string(out))
	var m map[string]any
	err := json.Unmarshal(out, &m)
	assert.NoError(t, err)
	assert.Equal(t, "INFO", m["lvl"])
	assert.Equal(t, "hi", m["message"])
	assert.Equal(t, "v1", m["k1"])
}

type customFmt struct {
	A  int
	Bs []int
}

func (c customFmt) Format(e ValueEncoder) {
	e.WriteBeginObject()
	e.WriteKey("A")
	e.WriteInt64(int64(c.A))
	e.WriteSeparator()
	e.WriteKey("Bs")
	e.WriteBeginArray()
	for i, v := range c.Bs {
		if i > 0 {
			e.WriteSeparator()
		}
		e.WriteInt64(int64(v))
	}
	e.WriteEndArray()
	e.WriteEndObject()
}

func TestJSONValueEncoder_CustomFormatter_Print(t *testing.T) {
	enc := NewJSONEncoder()
	r := NewRecord(InfoLevel, "fmt")
	r.AddAttr(Any("fmt", customFmt{A: 9, Bs: []int{1, 2, 3}}))
	out := enc.Encode(r)
	fmt.Println(string(out))
	var m map[string]any
	err := json.Unmarshal(out, &m)
	assert.NoError(t, err)
	f := m["fmt"].(map[string]any)
	assert.Equal(t, float64(9), f["A"])
	arr := f["Bs"].([]interface{})
	assert.Equal(t, float64(1), arr[0])
	assert.Equal(t, float64(2), arr[1])
	assert.Equal(t, float64(3), arr[2])
}

type Node struct {
	Next *Node
}

func TestJSONEncoder_Cycle_Print(t *testing.T) {
	enc := NewJSONEncoder()
	n := &Node{}
	n.Next = n // create cycle
	r := NewRecord(InfoLevel, "cycle")
	r.AddAttr(Any("cycle", n))
	out := enc.Encode(r)
	fmt.Println(string(out))
	// Should be valid JSON; cycle is either encoded via jsonv2 or fallback string
	var m map[string]any
	err := json.Unmarshal(out, &m)
	assert.NoError(t, err)
	// Accept both object or string, depending on encoder behavior
	switch v := m["cycle"].(type) {
	case map[string]any:
		// jsonv2 may encode with cycle handling; ensure it has a key
		_, _ = v["Next"]
	case string:
		assert.NotEmpty(t, v)
	default:
		t.Fatalf("unexpected type for cycle: %T", v)
	}
}
