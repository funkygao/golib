package mempool

import (
	"sync"
)

type slab struct {
	size int
	*sync.Pool
}

func newSlab(size int) *slab {
	return &slab{size: size, Pool: &sync.Pool{
		New: func() interface{} {
			return make([]byte, size)
		},
	}}
}

func (this *slab) Get() []byte {
	return this.Pool.Get().([]byte)
}
