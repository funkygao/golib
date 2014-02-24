package slab

import (
	"fmt"
)

type chunkLoc struct {
	slabClassIndex int
	slabIndex      int
	chunkIndex     int
	chunkSize      int
}

type chunk struct {
	refs int32    // Reference count.
	self chunkLoc // The chunkLoc for this chunk.
	next chunkLoc // Used when the chunk is in the free-list or when chained.
}

func (this *chunkLoc) isEmpty() bool {
	return this.slabClassIndex == emptyChunkLoc.slabClassIndex &&
		this.slabIndex == emptyChunkLoc.slabIndex &&
		this.chunkSize == emptyChunkLoc.chunkSize &&
		this.chunkIndex == emptyChunkLoc.chunkIndex
}

func (this *chunk) addRef() *chunk {
	this.refs++
	if this.refs <= 1 {
		panic(fmt.Sprintf("unexpected ref-count during addRef: %#v", *this))
	}
	return this
}
