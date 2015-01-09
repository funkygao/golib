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

	outstandings map[uint64]resourceWrapper // key is resource id
}

func NewDiagnosticTracker(pool *ResourcePool) *DiagnosticTracker {
	return &DiagnosticTracker{
		pool:         pool,
		quit:         make(chan bool),
		outstandings: make(map[uint64]resourceWrapper),
	}
}

func (this *DiagnosticTracker) BorrowResource(r Resource) {
	this.mutex.Lock()
	this.outstandings[r.Id()] = resourceWrapper{resource: r, timeUsed: time.Now()}
	this.mutex.Unlock()
}

func (this *DiagnosticTracker) ReturnResource(r Resource) {
	this.mutex.Lock()
	delete(this.outstandings, r.Id())
	this.mutex.Unlock()
}

func (this *DiagnosticTracker) Run(interval time.Duration, borrowTimeout int) {
	if interval == 0 {
		log.Warn("ResourcePool[%s] diagnostic disabled", this.pool.name)
		return
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	ever := true
	for ever {
		select {
		case <-ticker.C:
			if int64(len(this.outstandings)) > this.pool.MaxCapacity() {
				log.Warn("ResourcePool[%s] too many outstandings: %d > %d",
					this.pool.name, len(this.outstandings), this.pool.MaxCapacity())
			}

			if borrowTimeout > 0 {
				for _, r := range this.outstandings {
					if int(time.Now().Sub(r.timeUsed).Seconds()) > borrowTimeout {
						log.Warn("ResourcePool[%s] resource:%d killed: borrowed too long",
							this.pool.name, r.resource.Id())

						r.resource.Close() // force resource close
						//this.pool.Put(nil)
						this.ReturnResource(r.resource)
					}
				}
			}

			// FIXME if less returns, this.outstandings will be bigger and bigger

		case <-this.quit:
			ever = false
		}
	}
}

func (this *DiagnosticTracker) Stop() {
	close(this.quit)
}
