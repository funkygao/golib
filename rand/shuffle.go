package rand

import (
	"math/rand"
)

func ShuffleInts(a []int) []int {
	r := make([]int, len(a))
	perm := rand.Perm(len(a))
	for i, v := range perm {
		r[v] = a[i]
	}

	return r
}
