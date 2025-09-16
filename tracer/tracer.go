package tracer

import (
	"sync"
	"sync/atomic"
	"time"
)

type Tracer struct {
	idCounter   uint64
	spanHandler func(span *Span)
	spanPool    sync.Pool
}

func NewTracer(traceHandler func(span *Span)) *Tracer {
	if traceHandler == nil {
		traceHandler = func(span *Span) {}
	}
	res := &Tracer{
		spanPool: sync.Pool{
			New: func() interface{} {
				return &Span{
					Attributes: [4][2]string{},
				}
			},
		},
		spanHandler: traceHandler,
	}
	for i := 0; i < 1000; i++ {
		res.spanPool.Put(&Span{
			Attributes: [4][2]string{},
		})
	}
	return res
}

func (t *Tracer) SetSpanHandler(h func(span *Span)) {
	t.spanHandler = h
}

type Stage int

type Span struct {
	SpanID     uint64
	ParentID   uint64
	TraceID    string
	Stage      Stage        // 每个Span对应唯一Stage
	StartTime  time.Time    // 阶段开始时间
	EndTime    time.Time    // 阶段结束时间
	Attributes [4][2]string // 允许设置的最大属性数量为4,不使用symbol pair
	tracer     *Tracer
}

func (t *Tracer) StartSpan(parent *Span, stage Stage) *Span {
	span := t.spanPool.Get().(*Span)
	id := atomic.AddUint64(&t.idCounter, 1)

	span.tracer = t
	span.SpanID = id
	span.Stage = stage
	span.StartTime = time.Now().UTC()

	if parent != nil {
		span.ParentID = parent.SpanID
		span.TraceID = parent.TraceID
	} else {
		span.TraceID = NewTraceId()
	}
	return span
}

func (s *Span) End() {
	s.EndTime = time.Now().UTC()
	s.tracer.spanHandler(s)
	s.reset()
	s.tracer.spanPool.Put(s)
}

func (s *Span) Child(state Stage) *Span {
	return s.tracer.StartSpan(s, state)
}

func (s *Span) SetAttribute(key, value string) {
	for i := range s.Attributes {
		if s.Attributes[i][0] == "" || s.Attributes[i][0] == key {
			s.Attributes[i] = [2]string{key, value}
			return
		}
	}
}

func (s *Span) reset() {
	s.ParentID = 0
	s.Attributes = [4][2]string{}
}
