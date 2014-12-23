package pool

import (
	log "github.com/funkygao/log4go"
	"time"
)

type resourceTracker struct {
	r Resource
	t time.Time
}

// A diagnostic tracker for the recycle pool
type DiagnosticTracker struct {
	name      string
	capacity  int
	trackings []resourceTracker
	quit      chan bool
}

func NewDiagnosticTracker(name string, capacity int) *DiagnosticTracker {
	return &DiagnosticTracker{
		name:      name,
		capacity:  capacity,
		trackings: make([]resourceTracker, 0, capacity), // FIXME
		quit:      make(chan bool),
	}
}

func (this *DiagnosticTracker) BorrowResource(r Resource) {
	t := resourceTracker{r: r, t: time.Now()}
	this.trackings = append(this.trackings, t)
}

func (this *DiagnosticTracker) ReturnResource(r Resource) {

}

func (this *DiagnosticTracker) Run(interval time.Duration) {
	var (
		ever = true
	)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for ever {
		select {
		case <-ticker.C:
			log.Debug("resource pool[%s]: %+v", this.name, this.trackings)

			for _, t := range this.trackings {
				if time.Now().Sub(t.t).Seconds() > 10 {
					// not return after borrow for >10s
					log.Warn("resource[%+v] not return after 10s", t.r)
				}
			}

			this.trackings = this.trackings[:0] // reset

		case <-this.quit:
			ever = false
		}
	}
}

func (this *DiagnosticTracker) Stop() {
	close(this.quit)
}
