package tracer

import (
	"github.com/stretchr/testify/assert"
	"log/slog"
	"sync"
	"sync/atomic"
	"testing"
)

func TestRootSpanCreation(t *testing.T) {
	tracer := NewTracer(func(span *Span) {
		slog.Info("span", slog.Any("span", span))
	})
	span := tracer.StartSpan(nil, 1)
	assert.NotZero(t, span.SpanID, "SpanID应该自动生成")
	assert.Zero(t, span.ParentID, "根Span的ParentID应为0")
	assert.NotZero(t, span.StartTime, "StartTime应自动设置")
	assert.Equal(t, Stage(1), span.Stage, "Stage应正确设置")
	span.End()
}

func TestChildSpan(t *testing.T) {
	tracer := NewTracer(func(span *Span) {
		slog.Info("span", slog.Any("span", span))
	})
	parent := tracer.StartSpan(nil, 1)
	child := tracer.StartSpan(parent, 2)
	assert.Equal(t, parent.SpanID, child.ParentID, "子Span的ParentID应指向父Span")
	assert.Equal(t, parent.TraceID, child.TraceID, "子Span应继承父Span的TraceID")
	parent.End()
	child.End()
}

func TestEndCallsHandler(t *testing.T) {
	var called bool
	var lastSpan *Span

	tracer := NewTracer(func(s *Span) {
		called = true
		lastSpan = s
	})
	span := tracer.StartSpan(nil, 1)

	span.End()

	assert.True(t, called, "spanHandler应该被调用")
	assert.Same(t, span, lastSpan, "应该传递正确的Span实例")
	assert.NotZero(t, span.EndTime, "EndTime应该被设置")
}

func TestSetAttribute(t *testing.T) {
	tracer := NewTracer(nil)
	span := tracer.StartSpan(nil, 1)

	span.SetAttribute("user", "alice")
	span.SetAttribute("db", "mysql")

	found := 0
	for _, attr := range span.Attributes {
		if attr[0] == "user" {
			assert.Equal(t, "alice", attr[1])
			found++
		}
		if attr[0] == "db" {
			assert.Equal(t, "mysql", attr[1])
			found++
		}
	}
	assert.Equal(t, 2, found, "应该找到所有设置的属性")
}

func TestSpanReuse(t *testing.T) {
	tracer := NewTracer(nil)
	first := tracer.StartSpan(nil, 1)
	firstID := first.SpanID
	first.End()

	// 新Span应该复用对象
	second := tracer.StartSpan(nil, 2)

	assert.NotEqual(t, firstID, second.SpanID, "应该生成新的SpanID")
	assert.Zero(t, second.ParentID, "重置后的ParentID应为0")
}

func TestConcurrentSpans(t *testing.T) {
	tracer := NewTracer(nil)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			span := tracer.StartSpan(nil, Stage(i%5))
			span.SetAttribute("key", "value")
			span.End()
		}()
	}
	wg.Wait()

	// 验证ID生成器原子性
	lastID := atomic.LoadUint64(&tracer.idCounter)
	assert.Equal(t, uint64(100), lastID, "应该正确生成100个SpanID")
}

// 新增基准测试，测试10个Span的性能表现
// goos: windows
// goarch: amd64
// pkg: github.com/jcdtechnology/betago/base/tracer
// cpu: 13th Gen Intel(R) Core(TM) i5-13400F
// Benchmark10Spans
// Benchmark10Spans-16    	 2279658	       522.3 ns/op	      56 B/op	       4 allocs/op
// PASS
func Benchmark10Spans(b *testing.B) {
	tracer := NewTracer(nil) // 使用空handler避免日志干扰
	b.ReportAllocs()         // 报告内存分配
	b.ResetTimer()           // 重置计时器，排除初始化影响

	for i := 0; i < b.N; i++ {
		// 创建10级Span链
		var spans [10]*Span
		spans[0] = tracer.StartSpan(nil, 1)
		for j := 1; j < 10; j++ {
			spans[j] = tracer.StartSpan(spans[j-1], Stage(j))
		}

		// 逆序结束Span（模拟真实调用栈展开）
		for j := 9; j >= 0; j-- {
			spans[j].End()
		}
	}
}

// goos: windows
// goarch: amd64
// pkg: github.com/jcdtechnology/betago/base/tracer
// cpu: 13th Gen Intel(R) Core(TM) i5-13400F
// Benchmark1Span
// Benchmark1Span-16    	60220104	        20.52 ns/op	       0 B/op	       0 allocs/op
// PASS
func BenchmarkNon1Span(b *testing.B) {
	tracer := NewTracer(nil) // 使用空handler避免日志干扰
	b.ReportAllocs()         // 报告内存分配
	b.ResetTimer()           // 重置计时器，排除初始化影响
	rootSpan := tracer.StartSpan(nil, 1)
	for i := 0; i < b.N; i++ {
		span := tracer.StartSpan(rootSpan, Stage(i))
		span.End()
	}
}

// goos: windows
// goarch: amd64
// pkg: github.com/jcdtechnology/betago/base/tracer
// cpu: 13th Gen Intel(R) Core(TM) i5-13400F
// BenchmarkNewTraceID
// BenchmarkNewTraceID-16    	36279423	        32.63 ns/op
// PASS
func BenchmarkNewTraceID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewTraceId()
	}
}

func TestGenUniqueID(t *testing.T) {
	traceId := NewTraceId()
	slog.Info("traceId", slog.String("traceId", traceId))
}
