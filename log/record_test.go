package log

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"reflect"
	"runtime"
	"strconv"
	"testing"
)

// goos: windows
// goarch: amd64
// pkg: ai-flow/pkg/log
// cpu: 13th Gen Intel(R) Core(TM) i5-13400F
// BenchmarkRecordBuildWithRecycle
// BenchmarkRecordBuildWithRecycle-16    	 8967745	       130.3 ns/op	       0 B/op	       0 allocs/op
// BenchmarkBuildString
// BenchmarkBuildString-16               	  829525	      1459 ns/op	    1360 B/op	      18 allocs/op
// BenchmarkRecordBuild
// BenchmarkRecordBuild-16               	 3217929	       356.6 ns/op	     664 B/op	       7 allocs/op
// BenchmarkCaller2
// BenchmarkCaller2-16                   	 4662724	       259.1 ns/op
// BenchmarkReflectFuncPointer
// BenchmarkReflectFuncPointer-16        	53396696	        24.73 ns/op
// PASS
func BenchmarkBuildString(b *testing.B) {
	b.ReportAllocs()

	testMsg := "SubFastWsStreams bestPrice"
	usedIPs := "192.168.1.1,192.168.1.2"
	symbolstrs := "BTCUSDT,ETHUSDT"
	symbolen := 2
	localIPCount := 2
	args := []string{
		"msg=" + testMsg + " usedIPlen=" + strconv.Itoa(localIPCount),
		"usedIPs=" + usedIPs,
		"symbols=" + symbolstrs,
		"len=" + strconv.Itoa(symbolen),
	}
	SetLogPathWithCallerPath()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 包含完整字符串拼接过程
		_ = buildString(InfoLevel, args...)
	}
}

func BenchmarkRecordBuildWithRecycle(b *testing.B) {
	b.ReportAllocs()

	testMsg := "SubFastWsStreams bestPrice"
	usedIPs := "192.168.1.1,192.168.1.2"
	symbolstrs := "BTCUSDT,ETHUSDT"
	symbolen := 2
	localIPCount := 2

	attrs := []Attr{
		Int("usedIPlen", localIPCount),
		String("usedIPs", usedIPs),
		String("symbols", symbolstrs),
		Int("len", symbolen),
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// 完整Record构建过程
		record := NewRecord(InfoLevel, testMsg)
		record.attrs = append(record.attrs, attrs...)
		record.pc = CallerPC(ttt)
		record.format()
		record.recycle()
	}
}

func BenchmarkRecordBuild(b *testing.B) {
	b.ReportAllocs()

	testMsg := "SubFastWsStreams bestPrice"
	usedIPs := "192.168.1.1,192.168.1.2"
	symbolstrs := "BTCUSDT,ETHUSDT"
	symbolen := 2
	localIPCount := 2

	attrs := []Attr{
		Int("usedIPlen", localIPCount),
		String("usedIPs", usedIPs),
		String("symbols", symbolstrs),
		Int("len", symbolen),
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// 完整Record构建过程
		record := &Record{level: DebugLevel, msg: testMsg}
		record.attrs = append(record.attrs, attrs...)
		record.pc = CallerPC(ttt)
		record.format()
	}
}

func ttt() {
	runtime.Caller(0)
}

func BenchmarkCaller2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.Caller(2)
	}
}

func ttt1() {
	ttt2()
}

func ttt2() {
	pc := reflect.ValueOf(ttt).Pointer() //2ns
	runtime.FuncForPC(pc).FileLine(pc)   // 15ns
	slog.Info("123456")
}

func BenchmarkReflectFuncPointer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pc := reflect.ValueOf(ttt2).Pointer() //2ns
		runtime.FuncForPC(pc).FileLine(pc)    // 15ns
	}
	_ = fmt.Sprintf("123")
}

func TestReflectFuncPointer(t *testing.T) {
	pc := reflect.ValueOf(ttt1).Pointer()              //2ns
	fileName, ok := runtime.FuncForPC(pc).FileLine(pc) // 15ns
	t.Log(fileName, ok)
}

type L3 struct{ Data int }
type L2 struct{ Child L3 }
type L1 struct{ Child L2 }

// 修改后的测试用例

// 修改后的测试用例
func TestDataTypeFormatting(t *testing.T) {
	// 创建混合数据结构
	type NamedStruct struct {
		Field string
	}

	mixedData := map[string]interface{}{
		"mapValue":        map[string]int{"key1": 1, "key2": 2},
		"sliceValue":      []int{1, 2, 3},
		"namedStruct":     NamedStruct{Field: "value"},
		"anonymousStruct": struct{ X int }{X: 42},
	}

	r := NewRecord(InfoLevel, "test")
	r.appendAnyValue(mixedData, 0, 3)

	output := string(r.Buffer())
	t.Logf("格式化输出: %s", output)

	// 验证各种类型的格式
	assert.Contains(t, output, "iter{")        // Map 格式
	assert.Contains(t, output, "NamedStruct{") // 命名结构体
	assert.Contains(t, output, "Field=value")  // 结构体字段
	assert.Contains(t, output, "[1, 2, 3]")    // 切片格式
}
