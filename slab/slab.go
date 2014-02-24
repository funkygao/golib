package slab

type slabClass struct {
	slabs     []*slab  // A growing array of slabs.
	chunkSize int      // Each slab is sliced into fixed-sized chunks.
	chunkFree chunkLoc // Chunks are tracked in a free-list per slabClass.

	numChunks     int64
	numChunksFree int64
}

type slab struct {
	memory []byte  // len(memory) == slabSize + SLAB_MEMORY_FOOTER_LEN.
	chunks []chunk // Parallel array of chunk metadata.
}

func (this *slabClass) pushFreeChunk(c *chunk) {
	if c.refs != 0 {
		panic(ErrPushChunkRefCount)
	}
	c.next = this.chunkFree
	this.chunkFree = c.self
	this.numChunksFree++
}

func (this *slabClass) popFreeChunk() *chunk {
	if this.chunkFree.isEmpty() {
		panic("Empty chunkFree")
	}
	c := this.chunk(this.chunkFree)
	if c.refs != 0 {
		panic("Unexpected ref")
	}
	c.refs = 1
	this.chunkFree = c.next
	c.next = emptyChunkLoc
	this.numChunksFree--
	if this.numChunksFree < 0 {
		panic("shit")
	}
	debug("%# v", *c)
	return c
}

func (this *slabClass) chunkMem(c *chunk) []byte {
	if c == nil || c.self.isEmpty() {
		return nil
	}
	beg := this.chunkSize * c.self.chunkIndex
	return this.slabs[c.self.slabIndex].memory[beg : beg+this.chunkSize]
}

func (this *slabClass) chunk(l chunkLoc) *chunk {
	if l.isEmpty() {
		return nil
	}
	return &(this.slabs[l.slabIndex].chunks[l.chunkIndex])
}
