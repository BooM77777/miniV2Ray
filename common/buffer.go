package common

import (
	"bytes"
	"fmt"
	"math/bits"
	"sync"
)

var buffer []sync.Pool

func initBuffer() []sync.Pool {
	// 64 * 1024 = 1 << 16
	// (0, 1]
	// (1, 2]
	// (2, 4]
	// (4, 8]
	// (8, 16]
	// ...
	// [32768, 65536)
	buffer = make([]sync.Pool, 17) // 1B ~ 64 KB
	for i := range buffer {
		k := i
		buffer[i].New = func() interface{} {
			return make([]byte, 1<<uint32(k))
		}
	}
	return buffer
}

// 返回最高有效位，即为缓存池中的下标
func msb(size int) int {
	// bits.Len32(x) : 表示x所需要的最小位数
	return bits.Len32(uint32(size)) - 1
}

// GetBuffer 创建一个buffer
func GetBuffer(size int) ([]byte, error) {
	// 避免手动调用init函数，这样可以只对外暴露Get和Put接口
	if buffer == nil || len(buffer) == 0 {
		buffer = initBuffer()
	}
	// (0, 64K]
	if size <= 0 || size > 65536 {
		return nil, fmt.Errorf("Invalid Buffer Size To Malloc - %d", size)
	}
	idx := msb(size)
	// 构建左开右闭的区间，确保最大值
	if size == 1<<idx {
		return buffer[idx].Get().([]byte)[:size], nil
	}
	return buffer[idx+1].Get().([]byte)[:size], nil
}

// PutBuffer 回收一个buffer
func PutBuffer(buf []byte) error {
	idx := msb(cap(buf))
	if cap(buf) == 0 || cap(buf) > 65536 || 1<<idx != cap(buf) {
		return fmt.Errorf("Invalid Buffer Size To Free - %d", cap(buf))
	}
	buffer[idx].Put(buf)
	return nil
}

var writeBuffer *sync.Pool

func initWriteBuffer() *sync.Pool {
	return &sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
}

// GetWriteBuffer 申请WriteBuffer
func GetWriteBuffer() *bytes.Buffer {
	if writeBuffer == nil {
		writeBuffer = initWriteBuffer()
	}
	return writeBuffer.Get().(*bytes.Buffer)
}

// PutWriteBuffer 回收WriteBuffer
func PutWriteBuffer(buf *bytes.Buffer) {
	buf.Reset()
	writeBuffer.Put(buf)
}
