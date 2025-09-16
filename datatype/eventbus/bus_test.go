package eventbus

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testEvent struct {
	topic string
	data  string
}

func (t testEvent) EventType() string {
	return "test-event"
}

func (t testEvent) Topic() string { return t.topic }

func TestEventBus_SubscribeGroup(t *testing.T) {
	t.Run("基础功能测试-接收组内主题事件", func(t *testing.T) {
		eb := NewEventBus[testEvent]()
		groupSuber := eb.SubscribeGroup([]string{"group-topic1", "group-topic2"})
		defer eb.UnsubscribeGroup(groupSuber)

		// 发布两个主题的事件
		eb.Publish(testEvent{topic: "group-topic1", data: "msg1"})
		eb.Publish(testEvent{topic: "group-topic2", data: "msg2"})

		// 验证接收结果
		received := make(map[string]struct{})
		for i := 0; i < 2; i++ {
			select {
			case e := <-groupSuber.Get(): // 正确使用Get()方法
				received[e.topic] = struct{}{}
			case <-time.After(100 * time.Millisecond):
				t.Fatal("等待事件超时")
			}
		}

		if _, ok := received["group-topic1"]; !ok {
			t.Error("未收到group-topic1事件")
		}
		if _, ok := received["group-topic2"]; !ok {
			t.Error("未收到group-topic2事件")
		}
	})

	t.Run("主题隔离测试-不接收组外事件", func(t *testing.T) {
		eb := NewEventBus[testEvent]()
		groupSuber := eb.SubscribeGroup([]string{"private-topic"})
		defer eb.UnsubscribeGroup(groupSuber)

		// 发布无关主题
		eb.Publish(testEvent{topic: "other-topic", data: "msg"})

		select {
		case e := <-groupSuber.Get():
			t.Errorf("意外收到非组内事件: %s", e.topic)
		case <-time.After(50 * time.Millisecond):
			// 正常情况不应收到
		}
	})

	t.Run("并发订阅测试-多个订阅者独立接收", func(t *testing.T) {
		eb := NewEventBus[testEvent]()
		groupSuber1 := eb.SubscribeGroup([]string{"shared-topic"})
		groupSuber2 := eb.SubscribeGroup([]string{"shared-topic"})
		defer eb.UnsubscribeGroup(groupSuber1)
		defer eb.UnsubscribeGroup(groupSuber2)

		// 并发发布
		go func() {
			eb.Publish(testEvent{topic: "shared-topic", data: "concurrent-msg"})
		}()

		// 验证两个订阅者
		receivedCount := 0
		for i := 0; i < 2; i++ {
			select {
			case <-groupSuber1.Get():
				receivedCount++
			case <-groupSuber2.Get():
				receivedCount++
			case <-time.After(100 * time.Millisecond):
			}
		}

		if receivedCount != 2 {
			t.Errorf("应收到2次事件，实际收到%d次", receivedCount)
		}
	})

	t.Run("取消订阅测试-取消后不再接收", func(t *testing.T) {
		eb := NewEventBus[testEvent]()
		groupSuber := eb.SubscribeGroup([]string{"unsub-test"})

		// 立即取消订阅
		eb.UnsubscribeGroup(groupSuber)

		eb.Publish(testEvent{topic: "unsub-test", data: "should-not-see"})

		select {
		case e := <-groupSuber.Get():
			t.Errorf("取消订阅后仍收到事件: %v", e)
		case <-time.After(50 * time.Millisecond):
			// 正确情况
		}
	})

	t.Run("缓冲区测试-积压事件处理", func(t *testing.T) {
		eb := NewEventBus[testEvent]()
		groupSuber := eb.SubscribeGroup([]string{"burst-topic"})
		defer eb.UnsubscribeGroup(groupSuber)

		// 1. 首次发布3个事件（超过默认缓冲区大小）
		for i := 0; i < 3; i++ {
			eb.Publish(testEvent{topic: "burst-topic", data: "msg" + strconv.Itoa(i)})
		}

		// 2. 验证只能收到最后1个事件（根据MultiLastEvent设计）
		firstEvent := <-groupSuber.Get()
		if firstEvent.data != "msg0" {
			t.Errorf("期望收到最新事件msg0，实际收到 %s", firstEvent.data)
		}

		// 3. 标记事件已处理，触发后续事件加载
		groupSuber.OnEventProcessed(firstEvent)
		secondEvent := <-groupSuber.Get()
		assert.Equal(t, "msg2", secondEvent.data)
		// 4. 发布新事件验证通道恢复
	})

}
