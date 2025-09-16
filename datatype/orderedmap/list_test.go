package orderedmap

import "testing"

func TestList_Insert(t *testing.T) {
	l := newList[int, string]()
	l.Insert(2, "2")
	l.Insert(1, "1")
	// 遍历,应该得到的是1,2,而不是2,1
	for e := l.Front(); e != nil; e = e.Next() {
		t.Log(e.Value)
	}
}

func TestList_InsertReverse(t *testing.T) {
	l := newList[int, string]()
	l.InsertReverse(2, "2")
	l.InsertReverse(1, "1")
	l.InsertReverse(5, "5")
	l.InsertReverse(4, "4")
	// 遍历,应该得到的是1,2,而不是2,1
	for e := l.Front(); e != nil; e = e.Next() {
		t.Log(e.Value)
	}
}
