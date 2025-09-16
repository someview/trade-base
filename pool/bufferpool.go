package pool

import (
	"bytes"
	"sync"
)

// bufferPool 是一个bytes.Buffer对象池，用于减少内存分配
type bufferPool struct {
	pool sync.Pool
	size int
}

// newBufferPool 创建一个新的bytes.Buffer池
// size 参数指定每个Buffer的初始容量
func newBufferPool(size int) *bufferPool {
	return &bufferPool{
		pool: sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 0, size))
			},
		},
		size: size,
	}
}

// get 从池中获取一个Buffer
func (p *bufferPool) get() *bytes.Buffer {
	buf := p.pool.Get().(*bytes.Buffer)
	buf.Reset() // 确保Buffer是空的
	return buf
}

// put 将一个Buffer归还到池中
func (p *bufferPool) put(buf *bytes.Buffer) {
	if buf == nil {
		return
	}

	// 如果Buffer容量超过预期太多，就不放回池中
	if buf.Cap() > p.size*2 {
		return // 让GC回收它
	}

	p.pool.Put(buf)
}

// 内部使用的各种大小的Buffer池，从256依次*2到4K，对外不可见
var (
	buffer256Pool = newBufferPool(256)  // 256B
	buffer512Pool = newBufferPool(512)  // 512B
	buffer1KPool  = newBufferPool(1024) // 1KB
	buffer2KPool  = newBufferPool(2048) // 2KB
	buffer4KPool  = newBufferPool(4096) // 4KB
)

// GetBuffer 根据需要的大小选择合适的Buffer
func GetBuffer(minCapacity int) *bytes.Buffer {
	if minCapacity <= 256 {
		return buffer256Pool.get()
	} else if minCapacity <= 512 {
		return buffer512Pool.get()
	} else if minCapacity <= 1024 {
		return buffer1KPool.get()
	} else if minCapacity <= 2048 {
		return buffer2KPool.get()
	} else if minCapacity <= 4096 {
		return buffer4KPool.get()
	} else {
		// 超过4KB的直接创建新的Buffer
		return bytes.NewBuffer(make([]byte, 0, minCapacity))
	}
}

// PutBuffer 归还Buffer到合适的池
func PutBuffer(buf *bytes.Buffer) {
	if buf == nil {
		return
	}

	capacity := buf.Cap()
	if capacity <= 256 {
		buffer256Pool.put(buf)
	} else if capacity <= 512 {
		buffer512Pool.put(buf)
	} else if capacity <= 1024 {
		buffer1KPool.put(buf)
	} else if capacity <= 2048 {
		buffer2KPool.put(buf)
	} else if capacity <= 4096 {
		buffer4KPool.put(buf)
	}
	// 如果容量大于4KB，就不放回池中，让GC处理
}

// CloneBuffer 创建Buffer的副本（用于需要保留数据的场景）
func CloneBuffer(src *bytes.Buffer) *bytes.Buffer {
	if src == nil {
		return nil
	}

	dst := GetBuffer(src.Len())
	dst.Write(src.Bytes())
	return dst
}

// GetBufferWithData 获取一个包含指定数据的Buffer
func GetBufferWithData(data []byte) *bytes.Buffer {
	buf := GetBuffer(len(data))
	buf.Write(data)
	return buf
}

func NewBuffer(minCapacity int) *bytes.Buffer {
	if minCapacity <= 256 {
		return bytes.NewBuffer(make([]byte, 0, 256))
	} else if minCapacity <= 512 {
		return bytes.NewBuffer(make([]byte, 0, 512))
	} else if minCapacity <= 1024 {
		return bytes.NewBuffer(make([]byte, 0, 1024))
	} else if minCapacity <= 2048 {
		return bytes.NewBuffer(make([]byte, 0, 2048))
	} else if minCapacity <= 4096 {
		return bytes.NewBuffer(make([]byte, 0, 4096))
	} else {
		// 超过4KB的直接创建新的Buffer
		return bytes.NewBuffer(make([]byte, 0, minCapacity))
	}
}
