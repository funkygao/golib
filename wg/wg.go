package wg

import (
	"sync"
)

type WaitGroupWrapper struct {
	wg sync.WaitGroup
}

func (this *WaitGroupWrapper) Wrap(cb func()) {
	this.wg.Add(1)
	go func() {
		cb()
		this.wg.Done()
	}()
}

func (this *WaitGroupWrapper) Wait() {
	this.wg.Wait()
}
