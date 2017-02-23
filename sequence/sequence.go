// Package sequence provides a simple sequence of consecutive numbers on
// which you can add new number and find where the consecutiveness is broken.
package sequence

import (
	"sort"
	"sync"
)

type Sequence struct {
	datas []int
	mu    sync.Mutex
}

func New() *Sequence {
	return &Sequence{
		datas: []int{},
	}
}

func (s *Sequence) Add(v int) {
	s.mu.Lock()
	s.datas = append(s.datas, v)
	s.mu.Unlock()
}

func (s *Sequence) Length() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.datas)
}

func (s *Sequence) Summary() (min, max int, loss []int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.datas) == 0 {
		panic("empty sequence")
	}

	// sort
	sort.Ints(s.datas)
	min = s.datas[0]
	max = s.datas[len(s.datas)-1]

	missingSuccessors := make(map[int]struct{})
	for i, d := range s.datas {
		if i < len(s.datas)-1 {
			successorFound := false
			for _, n := range s.datas[i+1:] {
				if n == d+1 || n == d {
					successorFound = true
					break
				}
			}

			if !successorFound {
				missingSuccessors[d+1] = struct{}{}
			}
		}
	}

	for k := range missingSuccessors {
		loss = append(loss, k)
	}
	sort.Ints(loss)

	return
}
