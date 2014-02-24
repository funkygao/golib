package slab

type arenaStats struct {
	numAllocs         int64
	numAddRefs        int64
	numDecRefs        int64
	numGetNexts       int64
	numSetNexts       int64
	numMallocs        int64
	numMallocErrs     int64
	numTooBigErrs     int64
	numNoChunkMemErrs int64
}

func (this *Arena) Stats(m map[string]int64) map[string]int64 {
	m["numSlabClasses"] = int64(len(this.slabClasses))
	return m
}
