package orderedmap

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestOrderedMap_DeleteLessThanWithCallback(t *testing.T) {
    m := New[int, string]()
    // 添加一些键值对
    m.Set(1, "one")
    m.Set(2, "two")
    m.Set(3, "three")
    m.Set(4, "four")
    m.DeleteLessThanWithCallback(3, func(key int, value string) {
        t.Logf("Deleted key: %d, value: %s", key, value)
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
