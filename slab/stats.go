package slab

func (this *Arena) Stats(m map[string]int64) map[string]int64 {
	m["numSlabClasses"] = int64(len(this.slabClasses))
	return m
}
