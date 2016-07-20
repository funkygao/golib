package sync2

import (
	"sync"
	"time"
)

// WaitGroupTimeout is a wrapper of sync.WaitGroup that
// supports wait with timeout.
type WaitGroupTimeout struct {
	sync.WaitGroup
}

func (this *WaitGroupTimeout) WaitTimeout(timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		this.WaitGroup.Wait()
		close(c)
	}()

	select {
	case <-c:
		return false // completed normally

	case <-time.After(timeout):
		return true
	}
}
