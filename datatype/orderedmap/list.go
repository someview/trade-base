package orderedmap

import "cmp"

// Element 是双向链表中的一个元素
type Element[K comparable, V any] struct {
	// 双向链表中的前驱和后继指针
	next, prev *Element[K, V]

	// 该元素在有序映射中对应的键
	Key K

	// 该元素存储的值
	Value V
}

// Next 返回链表中的下一个元素，如果没有下一个元素则返回nil
func (e *Element[K, V]) Next() *Element[K, V] {
	return e.next
}

// Prev 返回链表中的前一个元素，如果没有前一个元素则返回nil
func (e *Element[K, V]) Prev() *Element[K, V] {
	return e.prev
}

// list 表示双向链表
// 链表在实例化后立即可用，无需额外初始化
type list[K cmp.Ordered, V any] struct {
	// root 指向链表的头和尾
	root *Element[K, V]
}

// NewList 创建并返回一个新的双向链表
func newList[K cmp.Ordered, V any]() *list[K, V] {
	return &list[K, V]{root: &Element[K, V]{}}
}

// IsEmpty 判断链表是否为空
func (l *list[K, V]) IsEmpty() bool {
	return l.root.next == nil
}

// Front 返回链表的第一个元素，如果链表为空则返回nil
func (l *list[K, V]) Front() *Element[K, V] {
	return l.root.next
}

// Back 返回链表的最后一个元素，如果链表为空则返回nil
func (l *list[K, V]) Back() *Element[K, V] {
	return l.root.prev
}

// Remove 从链表中移除指定元素
func (l *list[K, V]) Remove(e *Element[K, V]) {
	if e.prev == nil {
		// 移除头元素
		l.root.next = e.next
	} else {
		e.prev.next = e.next
	}

	if e.next == nil {
		// 移除尾元素
		l.root.prev = e.prev
	} else {
		e.next.prev = e.prev
	}

	// 防止内存泄漏
	e.next = nil
	e.prev = nil
}

// Find 在链表中查找指定键的元素，如果找不到则返回nil
func (l *list[K, V]) Find(key K) *Element[K, V] {
	for e := l.root.next; e != nil; e = e.next {
		if e.Key == key {
			return e
		}
		// 如果键是有序的，可以提前终止搜索
		if e.Key > key {
			break
		}
	}

	return nil
}

// PushFront 在链表头部插入一个新元素，并返回该元素
func (l *list[K, V]) PushFront(key K, value V) *Element[K, V] {
	e := &Element[K, V]{Key: key, Value: value}
	if l.root.next == nil {
		l.root.next = e
		l.root.prev = e
		return e
	}

	oldHead := l.root.next
	oldHead.prev = e
	e.next = oldHead

	l.root.next = e
	return e
}

// PushBack 在链表尾部插入一个新元素，并返回该元素
func (l *list[K, V]) PushBack(key K, value V) *Element[K, V] {
	e := &Element[K, V]{Key: key, Value: value}
	if l.root.prev == nil {
		l.root.next = e
		l.root.prev = e
		return e
	}

	oldTail := l.root.prev
	oldTail.next = e
	e.prev = oldTail

	l.root.prev = e
	return e
}

// Insert 在有序链表中插入一个新元素，保持链表有序，并返回该元素
func (l *list[K, V]) Insert(key K, value V) *Element[K, V] {
	e := &Element[K, V]{Key: key, Value: value}

	if l.root.next == nil {
		l.root.next = e
		l.root.prev = e
		return e
	}

	// 查找插入位置
	current := l.root.next
	for current != nil && current.Key < key {
		current = current.next
	}

	if current == nil {
		// 插入到链表末尾
		oldTail := l.root.prev
		oldTail.next = e
		e.prev = oldTail
		l.root.prev = e
		return e
	}

	// 在current前插入新元素
	oldPrev := current.prev
	oldNext := current

	e.next = oldNext
	oldNext.prev = e

	if oldPrev != nil {
		oldPrev.next = e
		e.prev = oldPrev
	} else {
		// 新元素成为头部
		l.root.next = e
	}

	return e
}

// InsertReverse 在有序链表中插入一个新元素，从链表尾部开始查找插入位置，并返回该元素
// 当链表较长且新元素键值较大时，这种方法可能更高效
func (l *list[K, V]) InsertReverse(key K, value V) *Element[K, V] {
	e := &Element[K, V]{Key: key, Value: value}

	if l.root.next == nil {
		l.root.next = e
		l.root.prev = e
		return e
	}

	// 从尾部开始查找插入位置
	current := l.root.prev
	for current != nil && current.Key > key {
		current = current.prev
	}

	if current == nil {
		// 插入到链表头部
		oldHead := l.root.next
		oldHead.prev = e
		e.next = oldHead
		l.root.next = e
		return e
	}

	// 在current后插入新元素
	oldPrev := current
	oldNext := current.next

	e.prev = oldPrev
	oldPrev.next = e

	if oldNext != nil {
		oldNext.prev = e
		e.next = oldNext
	} else {
		// 新元素成为尾部
		l.root.prev = e
	}
	return e
}
