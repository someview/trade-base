package eventbus

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testEvent2 struct {
	num int
}

func TestLastEvent_Get(t *testing.T) {
	t.Run("多个旧值", func(t *testing.T) {
		e := NewLastEvent[testEvent2]()
		e.Put(testEvent2{1})
		e.Put(testEvent2{2})
		v1 := <-e.Get()
		assert.Equal(t, 1, v1.num)
		e.Load()
		v2 := <-e.Get()
		assert.Equal(t, 2, v2.num)
	})
}

func BenchmarkMultiLastEvent_Get(b *testing.B) {
	e := NewLastEvent[testEvent2]()
	go func() {
		for {
			<-e.Get()
			e.Load()
		}
	}()
	for i := 0; i < b.N; i++ {
		e.Put(testEvent2{1})
	}
}

// goos: windows
// goarch: amd64
// pkg: github.com/jcdtechnology/betago/base/eventbus
// cpu: 13th Gen Intel(R) Core(TM) i5-13400F
// BenchmarkMultiLastEvent_LoadToGetLatency
// BenchmarkMultiLastEvent_LoadToGetLatency-16    	13961523	        81.65 ns/op	        80.64 ns/load_to_get
// PASS
func BenchmarkMultiLastEvent_LoadToGetLatency(b *testing.B) {
	e := NewLastEvent[int]()
	var totalDuration time.Duration
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 先放入一个事件
		// 获取事件
		start := time.Now().UTC()
		// 执行Load
		e.Put(i)
		// 记录时间点
		<-e.Get()
		e.Load()
		// 计算延迟
		totalDuration += time.Since(start)
	}

	avgNs := float64(totalDuration.Nanoseconds()) / float64(b.N)
	b.ReportMetric(avgNs, "ns/load_to_get")
}

// goos: windows
// goarch: amd64
// pkg: github.com/jcdtechnology/betago/base/eventbus
// cpu: 13th Gen Intel(R) Core(TM) i5-13400F
// BenchmarkTimeSince
// BenchmarkTimeSince-16    	247418497	         4.794 ns/op
// PASS
func BenchmarkTimeSince(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Since(time.Now())
	}
}

// goos: windows
// goarch: amd64
// pkg: github.com/jcdtechnology/betago/base/eventbus
// cpu: 13th Gen Intel(R) Core(TM) i5-13400F
// BenchmarkCh
// BenchmarkCh-16    	31629387	        37.16 ns/op	        35.71 ns/ch_send_to_receive
// PASS
func BenchmarkCh(b *testing.B) {
	var totalDuration time.Duration
	ch := make(chan int, 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		start := time.Now().UTC()
		ch <- i
		<-ch
		totalDuration += time.Since(start)
	}
	avgNs := float64(totalDuration.Nanoseconds()) / float64(b.N)
	b.ReportMetric(avgNs, "ns/ch_send_to_receive")

}
