package tracer

import (
	"sync/atomic"
)

var (
	globalTracerPointer atomic.Pointer[Tracer]
)

func GlobalTracer() *Tracer {
	v := globalTracerPointer.Load()
	if v != nil {
		return v
	}
	v = NewTracer(func(span *Span) {})
	if globalTracerPointer.CompareAndSwap(nil, v) {
		return v
	}
	return globalTracerPointer.Load()
}
