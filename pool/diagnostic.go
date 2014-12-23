package pool

import (
	log "github.com/funkygao/log4go"
	"sync"
	"time"
)

// A diagnostic tracker for the recycle pool
type DiagnosticTracker struct {
	pool  *ResourcePool
	quit  chan bool
	mutex sync.Mutex

	trackings map[uint64]resourceWrapper // key is resource id
}

func NewDiagnosticTracker(pool *ResourcePool) *DiagnosticTracker {
	return &DiagnosticTracker{
		pool:      pool,
		quit:      make(chan bool),
		trackings: make(map[uint64]resourceWrapper),
	}
}

func (this *DiagnosticTracker) BorrowResource(r Resource) {
	this.mutex.Lock()
	this.trackings[r.Id()] = resourceWrapper{resource: r, timeUsed: time.Now()}
	this.mutex.Unlock()
}

func (this *DiagnosticTracker) ReturnResource(r Resource) {
	this.mutex.Lock()
	delete(this.trackings, r.Id())
	this.mutex.Unlock()
}

func (this *DiagnosticTracker) Run(interval time.Duration) {
	var (
		ever          = true
		borrowTimeout = 10 // TODO
	)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for ever {
		select {
		case <-ticker.C:
			if int64(len(this.trackings)) > this.pool.MaxCapacity() {
				log.Warn("ResourcePool[%s] too few returned: %d < %d", this.pool.name,
					len(this.trackings), this.pool.MaxCapacity())
			}

			for _, r := range this.trackings {
				if int(time.Now().Sub(r.timeUsed).Seconds()) > borrowTimeout {
					log.Warn("ResourcePool[%s] not return within %ds, closed",
						this.pool.name, borrowTimeout)

					// force resource close
					r.resource.Close()
				}
			}

			// FIXME if less returns, this.trackings will be bigger and bigger

		case <-this.quit:
			ever = false
		}
	}
}

func (this *DiagnosticTracker) Stop() {
	close(this.quit)
}
