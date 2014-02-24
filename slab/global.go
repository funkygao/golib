package slab

const (
	// slabClassIndex + slabIndex + slabMagic
	SLAB_MEMORY_FOOTER_LEN int = 4 + 4 + 4
)

var (
	emptyChunkLoc = chunkLoc{-1, -1, -1, -1} // A sentinel.
)

func defaultMalloc(size int) []byte {
	return make([]byte, size)
}
