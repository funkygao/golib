package wg

import (
	"sync"
)

type WaitGroupWrapper struct {
	sync.WaitGroup
}

func (this *WaitGroupWrapper) Wrap(cb func()) {
	this.Add(1)
	go func() {
		cb()
		this.Done()
	}()
}
