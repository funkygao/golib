package sortmap

import (
	"sort"
)

type sortedMap struct {
	m map[string]int
	s []string
}

func NewSortedMap() *sortedMap {
	this := new(sortedMap)
	this.m = make(map[string]int)
	return this
}

func (this *sortedMap) Set(key string, val int) {
	this.m[key] = val
}

func (this *sortedMap) Len() int {
	return len(this.m)
}

func (this *sortedMap) Less(i, j int) bool {
	return this.m[this.s[i]] > this.m[this.s[j]]
}

func (this *sortedMap) Swap(i, j int) {
	this.s[i], this.s[j] = this.s[j], this.s[i]
}

func (this *sortedMap) SortedKeys() []string {
	this.s = make([]string, len(this.m))
	i := 0
	for key, _ := range this.m {
		this.s[i] = key
		i++
	}
	sort.Sort(this)
	return this.s
}
