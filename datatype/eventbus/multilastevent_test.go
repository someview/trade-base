package eventbus

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type MatchEv struct {
	topic string
	count int
}

func (m MatchEv) EventType() string {
	return "test"
}

func (m MatchEv) Topic() string {
	return m.topic
}

func TestMultiLastEvent_Get(t *testing.T) {
	multiLastEvent := NewMultiLastEvent[MatchEv](1)
	ch := multiLastEvent.Get()
	select {
	case <-ch:
	default:
		assert.Empty(t, ch)
	}
	// 连续放入3个事件
	v1 := MatchEv{topic: "topic", count: 1}
	v3 := MatchEv{topic: "topic", count: 3}
	multiLastEvent.Put(v1)
	multiLastEvent.Put(MatchEv{topic: "topic", count: 2})
	multiLastEvent.Put(v3)
	// 首次监听存在事件
	select {
	case v := <-ch:
		assert.Equal(t, v1, v)
	default:
		assert.Failf(t, "channel should return a value", "")
	}
	// 再次监听将会不到事件，除非主动加载
	select {
	case <-ch:
		assert.Failf(t, "channel should return a value", "")
	default:
	}
	multiLastEvent.OnEventProcessed(v1)
	select {
	case v := <-ch:
		assert.Equal(t, v3, v)
	default:
		assert.Failf(t, "channel should not return a value", "")
	}
}

func TestMultiLastEvent_Get_DifferentTopic(t *testing.T) {
	multiLastEvent := NewMultiLastEvent[MatchEv](2)
	ch := multiLastEvent.Get()
	select {
	case <-ch:
	default:
		assert.Empty(t, ch)
	}
	// 连续放入3个事件
	v1 := MatchEv{topic: "topic1", count: 1}
	v3 := MatchEv{topic: "topic3", count: 3}
	multiLastEvent.Put(v1)
	multiLastEvent.Put(v3)
	// 首次监听存在事件
	select {
	case v := <-ch:
		assert.Equal(t, v1, v)
	default:
		assert.Failf(t, "channel should return a value", "")
	}
	// 再次监听将会不到事件，除非主动加载
	select {
	case v := <-ch:
		assert.Equal(t, v, v3)
	default:
		assert.Failf(t, "channel should return a value", "")
	}
}
