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
	refs int32
	self chunkLoc
	next chunkLoc
}

func (this *chunkLoc) isEmpty() bool {
	return this.slabClassIndex == -1 && this.slabIndex == -1 &&
		this.chunkSize == -1 && this.chunkIndex == -1
}

func (this *chunk) addRef() *chunk {
	this.refs++
	if this.refs <= 1 {
		panic(fmt.Sprintf("unexpected ref-count during addRef: %#v", c))
	}
}
