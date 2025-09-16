package eventbus

// 描述：这是一个最新事件的包装器
// 消费者只消费最新的事件，只有初始化或者消费者明确load的时候，才会尝试从新加载事件
// 生成者产生的事件最多只有一个最新的事件，尚未被消费的历史事件会被最新事件覆盖掉
import "sync"

type LastEvent struct {
	mu     sync.RWMutex
	closed bool
	c      chan Event
	v      Event // 直接存储值
	ok     bool  // 标记是否有未消费的事件
}

func (l *LastEvent) Put(t Event) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.closed {
		return
	}
	select {
	case l.c <- t:
		l.ok = false // 已消费
	default:
		l.v = t
		l.ok = true // 设置标志，表示有未消费的值
	}
}

func (l *LastEvent) Get() chan Event {
	return l.c
}

func (l *LastEvent) Load() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.closed {
		return
	}
	if l.ok { // 如果有未消费的值
		select {
		case l.c <- l.v:
			l.ok = false // 已消费
		default:
			// 通道已满，保持当前状态
		}
	}
}

func (l *LastEvent) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.closed = true
}

func NewLastEvent() *LastEvent {
	return &LastEvent{
		c:  make(chan Event, 1), // 用缓冲区长度为1的管道
		ok: false,               // 初始状态没有未消费的值
	}
}
