package eventbus

import (
	"sync"
)

type Event interface {
	Key() EventKey
}

type EventKey struct {
	Topic     string
	EventType string
}

type EventBus struct {
	mu               sync.RWMutex
	groupSubscribers map[int64]GroupSuber            // 组Id -> 订阅者事件发布器
	topicToGroups    map[EventKey]map[int64]struct{} // 主题 -> 所属的组Id集合
	nextGroupId      int64                           // 下一个可用的组Id
}

type GroupSuber struct {
	*MultiLastEvent
	groupId int64
	keys    []EventKey
}

func (g GroupSuber) GroupId() int64 {
	return g.groupId
}

func (g GroupSuber) EventKeys() []EventKey {
	return g.EventKeys()
}

func NewEventBus() *EventBus {
	return &EventBus{
		groupSubscribers: make(map[int64]GroupSuber), // groupId是唯一的
		topicToGroups:    make(map[EventKey]map[int64]struct{}),
		nextGroupId:      1, // 从1开始分配组Id
	}
}

// Publish 发布事件
func (eb *EventBus) Publish(t Event) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	key := t.Key()
	// 发送给包含此主题的所有组的订阅者
	if groups, exists := eb.topicToGroups[key]; exists {
		for groupId := range groups {
			if subscriber, ok := eb.groupSubscribers[groupId]; ok {
				subscriber.Put(t)
			}
		}
	}
}

// SubscribeGroup 创建并订阅主题组，返回组Id和订阅者
func (eb *EventBus) Subscribe(keys []EventKey) GroupSuber {
	subscriber := NewMultiLastEvent(len(keys))
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.nextGroupId++
	groupId := eb.nextGroupId
	// 添加主题到组
	for _, topic := range keys {
		topicGroupSet := eb.topicToGroups[topic]
		if topicGroupSet == nil {
			topicGroupSet = make(map[int64]struct{})
			eb.topicToGroups[topic] = topicGroupSet
		}
		topicGroupSet[groupId] = struct{}{}
	}
	eb.groupSubscribers[groupId] = GroupSuber{
		groupId:        groupId,
		MultiLastEvent: subscriber,
		keys:           keys,
	}
	return eb.groupSubscribers[groupId]
}

// UnsubscribeGroup 取消订阅主题组
func (eb *EventBus) Unsubscribe(suber GroupSuber) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	delete(eb.groupSubscribers, suber.groupId)
	// 从所有的topic中删除该组的订阅者
	for _, topic := range suber.keys {
		if topicGroups, exists := eb.topicToGroups[topic]; exists {
			delete(topicGroups, suber.groupId)
		}
	}
}
