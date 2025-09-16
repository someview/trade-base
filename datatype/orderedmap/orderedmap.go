package orderedmap

import (
	"cmp"
)

type OrderedMap[K cmp.Ordered, V any] struct {
	kv map[K]*Element[K, V]
	ll *list[K, V]
}

func NewOrderedMap[K cmp.Ordered, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		kv: make(map[K]*Element[K, V]),
		ll: newList[K, V](),
	}
}

// Get returns the value for a key. If the key does not exist, the second return
// parameter will be false and the value will be nil.
func (m *OrderedMap[K, V]) Get(key K) (value V, ok bool) {
	v, ok := m.kv[key]
	if ok {
		value = v.Value
	}

	return
}

// Set will set (or replace) a value for a key. If the key was new, then true
// will be returned. The returned value will be false if the value was replaced
// (even if the value was the same).
func (m *OrderedMap[K, V]) Set(key K, value V) bool {
	_, alreadyExist := m.kv[key]
	if alreadyExist {
		m.kv[key].Value = value
		return false
	}

	element := m.ll.Insert(key, value)
	m.kv[key] = element
	return true
}

// GetOrDefault returns the value for a key. If the key does not exist, returns
// the default value instead.
func (m *OrderedMap[K, V]) GetOrDefault(key K, defaultValue V) V {
	if value, ok := m.kv[key]; ok {
		return value.Value
	}

	return defaultValue
}

// GetElement returns the element for a key. If the key does not exist, the
// pointer will be nil.
func (m *OrderedMap[K, V]) GetElement(key K) *Element[K, V] {
	element, ok := m.kv[key]
	if ok {
		return element
	}

	return nil
}

// Len returns the number of elements in the map.
func (m *OrderedMap[K, V]) Len() int {
	return len(m.kv)
}

// Keys returns all of the keys in the order they were inserted. If a key was
// replaced it will retain the same position. To ensure most recently set keys
// are always at the end you must always Delete before Set.
func (m *OrderedMap[K, V]) Keys() (keys []K) {
	keys = make([]K, 0, m.Len())
	for el := m.Front(); el != nil; el = el.Next() {
		keys = append(keys, el.Key)
	}
	return keys
}

// Delete will remove a key from the map. It will return true if the key was
// removed (the key did exist).
func (m *OrderedMap[K, V]) Delete(key K) (didDelete bool) {
	element, ok := m.kv[key]
	if ok {
		m.ll.Remove(element)
		delete(m.kv, key)
	}

	return ok
}

// Front will return the element that is the first (oldest Set element). If
// there are no elements this will return nil.
func (m *OrderedMap[K, V]) Front() *Element[K, V] {
	return m.ll.Front()
}

// Back will return the element that is the last (most recent Set element). If
// there are no elements this will return nil.
func (m *OrderedMap[K, V]) Back() *Element[K, V] {
	return m.ll.Back()
}

// Next will return the value that after the current key
func (m *OrderedMap[K, V]) Next(key K) (v V, exist bool) {
	ele := m.GetElement(key)
	if ele == nil {
		exist = false
		return
	}
	next := ele.Next()
	if next == nil {
		return
	}
	return next.Value, true
}

// Copy returns a new OrderedMap with the same elements.
// Using Copy while there are concurrent writes may mangle the result.
func (m *OrderedMap[K, V]) Copy() *OrderedMap[K, V] {
	m2 := NewOrderedMap[K, V]()
	for el := m.Front(); el != nil; el = el.Next() {
		m2.Set(el.Key, el.Value)
	}
	return m2
}

// DeleteLessThan 删除所有键小于给定阈值的条目，返回删除的数量
func (m *OrderedMap[K, V]) DeleteLessThan(threshold K) {
	// 删除所有小于阈值的元素
	el := m.Front()
	for el != nil && el.Key < threshold {
		// 保存当前元素的值和下一个元素
		key := el.Key
		nextEl := el.Next()

		// 从链表和映射中删除当前元素
		m.ll.Remove(el)
		delete(m.kv, key)

		// 移动到下一个元素
		el = nextEl
	}
}

// SetReverse Set will set (or replace) a value for a key. If the key was new, then true
// will be returned. The returned value will be false if the value was replaced
// (even if the value was the same).
func (m *OrderedMap[K, V]) SetReverse(key K, value V) bool {
	_, alreadyExist := m.kv[key]
	if alreadyExist {
		m.kv[key].Value = value
		return false
	}

	element := m.ll.InsertReverse(key, value)
	m.kv[key] = element
	return true
}

func (m *OrderedMap[K, V]) Range(fn func(K, V) bool) {
	for el := m.Front(); el != nil; el = el.Next() {
		if !fn(el.Key, el.Value) {
			return
		}
	}
}

func (m *OrderedMap[K, V]) RangeReverse(fn func(K, V) bool) {
	for el := m.Back(); el != nil; el = el.Prev() {
		if !fn(el.Key, el.Value) {
			return
		}
	}
}

// DeleteLessThanWithCallback 删除所有键小于给定阈值的条目，并为每个被删除的元素调用回调函数
func (m *OrderedMap[K, V]) DeleteLessThanWithCallback(threshold K, callback func(K, V)) {
	// 从头开始遍历
	el := m.Front()

	// 删除所有小于阈值的元素
	for el != nil && el.Key < threshold {
		// 保存当前元素的值和下一个元素
		key, value := el.Key, el.Value
		nextEl := el.Next()

		// 从链表和映射中删除当前元素
		m.ll.Remove(el)
		delete(m.kv, key)

		// 调用回调函数
		if callback != nil {
			callback(key, value)
		}

		// 移动到下一个元素
		el = nextEl
	}
}
