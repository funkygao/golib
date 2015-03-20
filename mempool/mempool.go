package mempool

import (
	"sync"
)

type MemPool struct {
	maxSize int
	minSize int
	chunks  []*chunkPool
}

type chunkPool struct {
	size int
	*sync.Pool
}

func newChunkPool(size int) *chunkPool {
	p := &chunkPool{size: size, Pool: &sync.Pool{
		New: func() interface{} {
			return make([]byte, size)
		},
	}}
	return p
}

func (this *chunkPool) Get() []byte {
	return this.Pool.Get().([]byte)
}
