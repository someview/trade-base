package log

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJSONEncoder_BasicFields(t *testing.T) {
	r := NewRecord(InfoLevel, "hello")
	r.time = time.Date(2025, 1, 2, 3, 4, 5, 6, time.UTC)
	r.pc = CallerPC(TestJSONEncoder_BasicFields)
	r.AddAttr(String("k", "v"))
	enc := NewJSONEncoder()
	out := enc.Encode(r)
	fmt.Println(string(out))
	var m map[string]any
	err := json.Unmarshal(out, &m)
	assert.NoError(t, err)
	assert.Equal(t, "INFO", m["level"])
	assert.Equal(t, "hello", m["msg"])
	assert.Equal(t, "2025-01-02T03:04:05.000000006Z", m["ts"])
	assert.Equal(t, "v", m["k"])
	assert.Contains(t, m, "caller")
}

func TestJSONEncoder_Types(t *testing.T) {
	r := NewRecord(InfoLevel, "types")
	now := time.Date(2025, 2, 3, 4, 5, 6, 0, time.UTC)
	r.time = now
	r.AddAttr(String("s", "x"))
	r.AddAttr(Int("i", 123))
	r.AddAttr(Uint64("u", 456))
	r.AddAttr(Float64("f", 7.5))
	r.AddAttr(Bool("b", true))
	r.AddAttr(Time("t", now))
	r.AddAttr(Duration("d", time.Second))
	enc := NewJSONEncoder()
	out := enc.Encode(r)
	fmt.Println(string(out))
	var m map[string]any
	err := json.Unmarshal(out, &m)
	assert.NoError(t, err)
	assert.Equal(t, "x", m["s"])
	assert.Equal(t, float64(123), m["i"])
	assert.Equal(t, float64(456), m["u"])
	assert.Equal(t, 7.5, m["f"])
	assert.Equal(t, true, m["b"])
	assert.Equal(t, "2025-02-03T04:05:06Z", m["t"])
	assert.Equal(t, float64(time.Second.Nanoseconds()), m["d"])
}

func TestJSONEncoder_GroupAndAny(t *testing.T) {
	type S struct {
		A int
	}
	r := NewRecord(InfoLevel, "group")
	r.AddAttr(Group("g", Int("a", 1), String("b", "y")))
	r.AddAttr(Any("any", S{A: 3}))
	enc := NewJSONEncoder()
	out := enc.Encode(r)
	fmt.Println(string(out))
	var m map[string]any
	err := json.Unmarshal(out, &m)
	assert.NoError(t, err)
	g := m["g"].(map[string]any)
	assert.Equal(t, float64(1), g["a"])
	assert.Equal(t, "y", g["b"])
	am := m["any"].(map[string]any)
	assert.Equal(t, float64(3), am["A"])
}

func TestJSONEncoder_NestedStructures_Print(t *testing.T) {
	type Inner struct {
		X int
		Y string
	}
	type Outer struct {
		Name string
		In   Inner
		Arr  []Inner
		M    map[string]Inner
	}
	r := NewRecord(InfoLevel, "nested")
	r.time = time.Now().UTC()
	r.AddAttr(Any("slice", []any{1, "s", map[string]int{"a": 1}, []int{2, 3}, Inner{X: 42, Y: "hi"}}))
	m := map[string]any{
		"list": []int{1, 2},
		"obj":  Inner{X: 7, Y: "ok"},
		"deep": map[string]any{"x": []Inner{{X: 1, Y: "a"}, {X: 2, Y: "b"}}},
	}
	r.AddAttr(Any("map", m))
	r.AddAttr(Any("struct", Outer{
		Name: "outer",
		In:   Inner{X: 3, Y: "z"},
		Arr:  []Inner{{X: 4, Y: "q"}},
		M:    map[string]Inner{"k": {X: 5, Y: "w"}},
	}))
	enc := NewJSONEncoder()
	out := enc.Encode(r)
	fmt.Println(string(out))
	var mm map[string]any
	err := json.Unmarshal(out, &mm)
	assert.NoError(t, err)
	s := mm["slice"].([]interface{})
	assert.Equal(t, float64(1), s[0])
	assert.Equal(t, "s", s[1])
	mp := mm["map"].(map[string]any)
	deep := mp["deep"].(map[string]any)
	xarr := deep["x"].([]interface{})
	x0 := xarr[0].(map[string]any)
	assert.Equal(t, float64(1), x0["X"])
	st := mm["struct"].(map[string]any)
	in := st["In"].(map[string]any)
	assert.Equal(t, "outer", st["Name"])
	assert.Equal(t, float64(3), in["X"])
}
