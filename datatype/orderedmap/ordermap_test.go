package orderedmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderedMap_RangeDelete_CallbackBehavior(t *testing.T) {
	m := New[int, string]()
	// 添加一些键值对
	m.Set(1, "one")
	m.Set(2, "two")
	m.Set(3, "three")
	m.Set(4, "four")
	// 使用 RangeDelete 删除 <3 的键
	m.RangeDelete(func(key int, value string) (deleted, next bool) {
		if key < 3 {
			return true, true
		}
		return false, true
	})
	_, ok := m.Get(3)
	assert.True(t, ok)
	_, ok = m.Get(1)
	assert.False(t, ok)
}

func TestOrderedMap_DeleteLessThan(t *testing.T) {
	m := New[int, string]()
	// 添加一些键值对
	m.Set(1, "one")
	m.Set(2, "two")
	m.Set(3, "three")
	m.Set(4, "four")
	m.DeleteLessThan(3)
	_, ok := m.Get(3)
	assert.True(t, ok)
	_, ok = m.Get(2)
	assert.False(t, ok)
}

func TestOrderedMap_SetReverse(t *testing.T) {
	m := New[int, string]()
	// 添加一些键值对
	m.SetReverse(1, "one")
	m.SetReverse(2, "two")
	m.SetReverse(3, "three")
	m.SetReverse(4, "four")
	// 添加测试代码
	_, ok := m.Get(3)
	assert.True(t, ok)
	val := m.Back()
	assert.NotNil(t, val)
	assert.Equal(t, "four", val.Value)
}

func TestOrderedMap_RangeDelete(t *testing.T) {
	m := New[int, string]()
	for i := 1; i <= 5; i++ {
		m.Set(i, string(rune('a'+i-1)))
	}
	// 删除所有小于3的元素
	m.RangeDelete(func(k int, v string) (deleted, next bool) {
		if k < 3 {
			return true, true
		}
		return false, true
	})
	_, ok := m.Get(1)
	assert.False(t, ok)
	_, ok = m.Get(2)
	assert.False(t, ok)
	_, ok = m.Get(3)
	assert.True(t, ok)
	_, ok = m.Get(4)
	assert.True(t, ok)
	_, ok = m.Get(5)
	assert.True(t, ok)
}

func TestOrderedMap_RangeReverseDelete(t *testing.T) {
	m := New[int, string]()
	for i := 1; i <= 5; i++ {
		m.Set(i, string(rune('a'+i-1)))
	}
	// 从尾部删除所有大于3的元素
	m.RangeReverseDelete(func(k int, v string) (deleted, next bool) {
		if k > 3 {
			return true, true
		}
		return false, true
	})
	_, ok := m.Get(4)
	assert.False(t, ok)
	_, ok = m.Get(5)
	assert.False(t, ok)
	_, ok = m.Get(3)
	assert.True(t, ok)
	_, ok = m.Get(2)
	assert.True(t, ok)
	_, ok = m.Get(1)
	assert.True(t, ok)
}
