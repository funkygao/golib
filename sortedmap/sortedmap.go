package sortedmap

import (
	"sort"
)

// Not thread safe
type SortedMap struct {
	m map[string]int
	s []string
}

func NewSortedMap() *SortedMap {
	this := new(SortedMap)
	this.m = make(map[string]int)
	return this
}

func (this *SortedMap) Set(key string, val int) {
	this.m[key] = val
}

func (this *SortedMap) Get(key string) int {
	return this.m[key]
}

func (this *SortedMap) Inc(key string, delta int) int {
	v, present := this.m[key]
	if !present {
		v = 0
	}
	v += delta
	this.m[key] = v
	return v
}

func (this *SortedMap) Len() int {
	return len(this.m)
}

func (this *SortedMap) Less(i, j int) bool {
	return this.m[this.s[i]] > this.m[this.s[j]]
}

func (this *SortedMap) Swap(i, j int) {
	this.s[i], this.s[j] = this.s[j], this.s[i]
}

func (this *SortedMap) SortedKeys() []string {
	this.s = make([]string, len(this.m))
	i := 0
	for key, _ := range this.m {
		this.s[i] = key
		i++
	}
	sort.Sort(this)
	return this.s
}
