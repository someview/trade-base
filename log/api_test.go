package log

import (
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
)

func TestLogger_Debug(t *testing.T) {
	logger, err := NewLogger(NewDefaultConfig())
	assert.Nil(t, err)
	logger.Debug("test debug", Int("int", 1), String("123", "456"))
	logger.Close()
}

func ddd() {

}

func TestLogger_LogWithPC(t *testing.T) {
	logger, err := NewLogger(NewDefaultConfig())
	assert.Nil(t, err)
	logger.LogWithPC(CallerPC(ddd), DebugLevel, "test LogWithPC")
	logger.Close()
}

func TestLogger_Log_AnyStruct(t *testing.T) {
	SetLogPathWithCallerPath()
	logger, err := NewLogger(NewDefaultConfig())
	assert.Nil(t, err)
	logger.LogWithPC(CallerPC(ddd), DebugLevel, "test LogWithPC",
		slog.Any("struct", struct {
			name int
			Age  string
		}{
			name: 123,
			Age:  "456",
		}))
	logger.Close()
}

type sss1 struct {
	name int
	Age  string
	TTT2 *sss2
}
type sss2 struct {
	name int
	Age  string
}

func TestLogger_Log_AnyMap(t *testing.T) {
	m := make(map[string]any)
	m["111"] = &sss1{Age: "234", TTT2: &sss2{Age: "567"}}
	m["222"] = &sss1{Age: "345", TTT2: &sss2{Age: "678"}}
	SetLogPathWithCallerPath()
	logger, err := NewLogger(NewDefaultConfig())
	assert.Nil(t, err)
	logger.LogWithPC(CallerPC(ddd), DebugLevel, "test LogWithPC",
		slog.Any("MAP", m))
	logger.Close()
}

func TestLogger_LogAnySlice(t *testing.T) {
	m := make([]any, 10)
	m[0] = 1
	m[1] = &sss1{Age: "345", TTT2: &sss2{Age: "678"}}
	SetLogPathWithCallerPath()
	logger, err := NewLogger(NewDefaultConfig())
	assert.Nil(t, err)
	logger.LogWithPC(CallerPC(ddd), DebugLevel, "test LogWithPC",
		slog.Any("slice", m))
	logger.Close()
}
