// Package sequence provides a simple sequence of consecutive numbers on
// which you can add new number and find where the consecutiveness is broken.
package sequence

import (
	"fmt"
	"math"
	"sort"
	"sync"

	"github.com/funkygao/golib/color"
)

type Sequence struct {
	datas []int
	mu    sync.Mutex
	min   int
}

func New() *Sequence {
	return &Sequence{
		datas: []int{},
		min:   int(math.MaxInt64),
	}
}

func (s *Sequence) Add(v int) {
	s.mu.Lock()
	s.datas = append(s.datas, v)
	s.mu.Unlock()

	if v < s.min {
		s.min = v
	}
}

func (s *Sequence) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	r := "["
	var last = s.datas[0]
	for _, d := range s.datas[1:] {
		if d > last {
			r += fmt.Sprintf("%d ", d)
		} else {
			r += color.Red("%d ", d)
		}

		last = d
	}

	return r[:len(r)-1] + "]"
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
	if min != s.min {
		loss = append(loss, s.min)
	}

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
