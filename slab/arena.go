package slab

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"sort"
)

type Malloc func(size int) []byte

// Arena tracks one or more slabClass structs.
//
// Each slabClass manages a different chunkSize, where it's computed
// using growth factor.
// Each slabClass also tracks 0 or more slabs, where every slab
// tracked by a slabClass will have the same chunkSize.
//
// A slab manages a large continguous array of memory bytes, and the slab's
// memory is subdivided many fixed-sized chunks of the same chunkSize.
// All the chunks in a new slab are placed on a free-list that's part of
// the slabClass.
type Arena struct {
	growthFactor float64     // Should > 1.0.
	slabClasses  []slabClass // The chunkSizes of slabClasses grows by growthFactor.
	slabMagic    int32       // Magic number at the end of each slab memory []byte.
	slabSize     int         // Each slabClass memory totals this size in bytes.
	stats        arenaStats

	malloc func(size int) []byte
}

func NewArena(startChunkSize int, slabSize int, growthFactor float64,
	malloc Malloc) (this *Arena) {
	if malloc == nil {
		malloc = defaultMalloc
	}

	this = &Arena{
		growthFactor: growthFactor,
		slabMagic:    rand.Int31(),
		slabSize:     slabSize,
		malloc:       malloc,
	}
	this.addSlabClass(startChunkSize)

	return
}

// The first 'large enough' slabClass is found and a chunk from the
// free-list is taken to service the allocation.
func (this *Arena) Alloc(bufSize int) (buf []byte, err error) {
	this.stats.numAllocs++
	if bufSize > this.slabSize {
		this.stats.numTooBigErrs++
		return nil, ErrTooBig
	}
	chunkMem := this.assignChunkMem(this.findSlabClassIndex(bufSize))
	if chunkMem == nil {
		this.stats.numNoChunkMemErrs++
		return nil, ErrNoChunkMem
	}
	debug("%# v", chunkMem)
	return chunkMem[0:bufSize], nil
}

func (this *Arena) AddRef(buf []byte) {
	this.stats.numAddRefs++
	slab, chunk := this.bufContainer(buf)
	if slab == nil || chunk == nil {
		panic(ErrOutsideArena)
	}
	chunk.addRef()
}

func (this *Arena) DecRef(buf []byte) bool {
	this.stats.numDecRefs++
	slab, chunk := this.bufContainer(buf)
	if slab == nil || chunk == nil {
		panic(ErrOutsideArena)
	}
	return this.decRef(slab, chunk)
}

func (this *Arena) Owns(buf []byte) bool {
	slab, chunk := this.bufContainer(buf)
	return slab != nil && chunk != nil
}

func (this *Arena) GetNext(buf []byte) (next []byte) {
	this.stats.numGetNexts++
	slab, chunk := this.bufContainer(buf)
	if slab == nil || chunk == nil {
		panic(ErrOutsideArena)
	}
	if chunk.refs <= 0 {
		panic(ErrInvalidRefCount)
	}
	slabNext, chunkNext := this.chunk(chunk.next)
	if slabNext == nil || chunkNext == nil {
		return nil
	}
	chunkNext.addRef()
	return this.chunkMem(chunkNext)[:chunk.next.chunkSize]
}

func (this *Arena) SetNext(buf, bufNext []byte) {
	this.stats.numSetNexts++
	sc, c := this.bufContainer(buf)
	if sc == nil || c == nil {
		panic("buf not from this arena")
	}
	if c.refs <= 0 {
		panic("unexpected ref-count")
	}
	scOldNext, cOldNext := this.chunk(c.next)
	if scOldNext != nil && cOldNext != nil {
		this.decRef(scOldNext, cOldNext)
	}
	c.next = emptyChunkLoc
	if bufNext != nil {
		scNewNext, cNewNext := this.bufContainer(bufNext)
		if scNewNext == nil || cNewNext == nil {
			panic("bufNext not from this arena")
		}
		cNewNext.addRef()
		c.next = cNewNext.self
		c.next.chunkSize = len(bufNext)
	}
}

func (this *Arena) addSlabClass(chunkSize int) {
	debug("before addSlabClass: %# v", *this)
	this.slabClasses = append(this.slabClasses,
		slabClass{
			chunkSize: chunkSize,
			chunkFree: emptyChunkLoc,
		})
	debug("after addSlabClass: %# v", *this)
}

func (this *Arena) findSlabClassIndex(bufSize int) int {
	// binary search
	idx := sort.Search(len(this.slabClasses), func(i int) bool {
		return bufSize <= this.slabClasses[i].chunkSize
	})
	if idx >= len(this.slabClasses) {
		// Will be at tail of my slabClasses
		tailSlabClass := &(this.slabClasses[len(this.slabClasses)-1])
		nextChunkSize := float64(tailSlabClass.chunkSize) * this.growthFactor
		this.addSlabClass(int(math.Ceil(nextChunkSize)))
		return this.findSlabClassIndex(bufSize)
	}
	debug("findSlabClassIndex for %d: %d", bufSize, idx)
	return idx
}

func (this *Arena) assignChunkMem(slabClassIndex int) (chunkMem []byte) {
	slabClass := &this.slabClasses[slabClassIndex]
	if slabClass.chunkFree.isEmpty() {
		debug("chunkFree empty :%# v", *slabClass)
		if !this.addSlab(slabClassIndex, this.slabSize, this.slabMagic) {
			return nil
		}
	}
	return this.chunkMem(slabClass.popFreeChunk())
}

func (this *Arena) addSlab(slabClassIndex, slabSize int, slabMagic int32) bool {
	slabClass := &(this.slabClasses[slabClassIndex])
	chunksPerSlab := slabSize / slabClass.chunkSize
	debug("chunkPerSlab for class:%d: %d", slabClassIndex, chunksPerSlab)
	if chunksPerSlab <= 0 {
		chunksPerSlab = 1
	}
	slabIndex := len(slabClass.slabs)
	memorySize := (slabClass.chunkSize * chunksPerSlab) + SLAB_MEMORY_FOOTER_LEN
	this.stats.numMallocs++
	memory := this.malloc(memorySize)
	if memory == nil {
		this.stats.numMallocErrs++
		return false
	}

	slab := &slab{
		memory: memory,
		chunks: make([]chunk, chunksPerSlab),
	}
	footer := slab.memory[len(slab.memory)-SLAB_MEMORY_FOOTER_LEN:]
	binary.BigEndian.PutUint32(footer[0:4], uint32(slabClassIndex))
	binary.BigEndian.PutUint32(footer[4:8], uint32(slabIndex))
	binary.BigEndian.PutUint32(footer[8:12], uint32(slabMagic))
	slabClass.slabs = append(slabClass.slabs, slab)
	for i := 0; i < len(slab.chunks); i++ {
		c := &(slab.chunks[i])
		c.self.slabClassIndex = slabClassIndex
		c.self.slabIndex = slabIndex
		c.self.chunkIndex = i
		c.self.chunkSize = slabClass.chunkSize
		slabClass.pushFreeChunk(c)
	}
	slabClass.numChunks += int64(len(slab.chunks))
	debug("%# v", *this)
	return true
}

func (this *Arena) chunkMem(c *chunk) []byte {
	if c == nil || c.self.isEmpty() {
		return nil
	}
	return this.slabClasses[c.self.slabClassIndex].chunkMem(c)
}

func (this *Arena) chunk(l chunkLoc) (*slabClass, *chunk) {
	if l.isEmpty() {
		return nil, nil
	}
	sc := &(this.slabClasses[l.slabClassIndex])
	return sc, sc.chunk(l)
}

func (this *Arena) bufContainer(buf []byte) (*slabClass, *chunk) {
	if buf == nil || cap(buf) <= SLAB_MEMORY_FOOTER_LEN {
		return nil, nil
	}
	debug("buf: %# v", buf)
	rest := buf[:cap(buf)]
	footerDistance := len(rest) - SLAB_MEMORY_FOOTER_LEN
	footer := rest[footerDistance:]
	debug("rest: %# v, footer: %# v", rest, footer)
	slabClassIndex := binary.BigEndian.Uint32(footer[0:4])
	slabIndex := binary.BigEndian.Uint32(footer[4:8])
	slabMagic := binary.BigEndian.Uint32(footer[8:12])
	if slabMagic != uint32(this.slabMagic) {
		return nil, nil
	}
	sc := &(this.slabClasses[slabClassIndex])
	slab := sc.slabs[slabIndex]
	chunkIndex := len(slab.chunks) - (footerDistance / sc.chunkSize)
	return sc, &(slab.chunks[chunkIndex])
}

func (this *Arena) decRef(sc *slabClass, c *chunk) bool {
	c.refs--
	if c.refs < 0 {
		panic(fmt.Sprintf("unexpected ref-count during decRef: %#v", c))
	}
	if c.refs == 0 {
		scNext, cNext := this.chunk(c.next)
		if scNext != nil && cNext != nil {
			this.decRef(scNext, cNext)
		}
		c.next = emptyChunkLoc
		sc.pushFreeChunk(c)
		return true
	}
	return false
}
