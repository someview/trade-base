package eventbus

import "sync"

// MultiLastEvent 管理多个主题的最新事件（简化版）
type MultiLastEvent struct {
	mu     sync.RWMutex
	c      chan Event         // 所有主题共享的事件通道
	values map[EventKey]Event // key eventType value event
}

// NewMultiLastEvent 创建一个新的多主题事件管理器（简化版）
func NewMultiLastEvent(eventChNum int) *MultiLastEvent {
	return &MultiLastEvent{
		c:      make(chan Event, eventChNum), // 通道容量设为主题数（避免过大）
		values: make(map[EventKey]Event),
	}
}

// Put 存入一个主题的最新事件（简化逻辑）
func (m *MultiLastEvent) Put(event Event) {
	m.mu.Lock()
	defer m.mu.Unlock()

	topic := event.Key()
	m.values[topic] = event
	// 尝试将当前主题的最新事件发送到通道（通道有空间时发送）
	select {
	case m.c <- event:
		// 发送成功，从values中删除（避免重复发送）
		delete(m.values, topic)
	default:
		// 通道已满，保留values中的最新事件（等待后续Load或OnEventProcessed处理）
	}
}

// Get 获取事件通道
func (m *MultiLastEvent) Get() chan Event {
	return m.c
}

// OnEventProcessed 当事件被处理后调用（简化逻辑）
func (m *MultiLastEvent) OnEventProcessed(event Event) {
	m.mu.Lock()
	defer m.mu.Unlock()
	eventKey := event.Key()
	// 事件已被处理，检查该主题是否有新事件需要发送
	if newEvent, exists := m.values[eventKey]; exists {
		// 尝试发送新事件（通道有空间时发送）
		select {
		case m.c <- newEvent:
			delete(m.values, eventKey) // 发送成功，清理缓存
		default:
			// 通道仍满，保留缓存（等待下次机会）
		}
	}
}

// Clear 清空所有主题的最新事件（简化逻辑）
func (m *MultiLastEvent) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	clear(m.values)
}
