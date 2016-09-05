// Package timewheel implements a timewheel algorithm that
// is suitable for large numbers of timers.
//
// It is based on golang channel broadcast mechanism.
package timewheel

import (
	"sync"
	"time"
)

type TimeWheel struct {
	mu sync.Mutex

	interval   time.Duration
	maxTimeout time.Duration

	ticker *time.Ticker
	quit   chan struct{}

	cs []chan struct{}

	pos int // current time tick pointer
}

func NewTimeWheel(interval time.Duration, buckets int) *TimeWheel {
	this := new(TimeWheel)

	this.interval = interval
	this.maxTimeout = time.Duration(interval * (time.Duration(buckets)))

	this.quit = make(chan struct{})
	this.pos = 0
	this.cs = make([]chan struct{}, buckets)
	for i := range this.cs {
		this.cs[i] = make(chan struct{})
	}

	this.ticker = time.NewTicker(interval)
	go this.run()

	return this
}

func (this *TimeWheel) Stop() {
	close(this.quit)
}

func (this *TimeWheel) After(timeout time.Duration) <-chan struct{} {
	if timeout >= this.maxTimeout {
		panic("timeout too much, over maxtimeout")
	} else if timeout < time.Second {
		timeout = time.Second
	}

	this.mu.Lock()

	index := (this.pos + int(timeout/this.interval)) % len(this.cs)
	broadcastCh := this.cs[index]

	this.mu.Unlock()

	return broadcastCh
}

func (this *TimeWheel) run() {
	for {
		select {
		case <-this.ticker.C:
			this.onTicker()

		case <-this.quit:
			this.ticker.Stop()
			return
		}
	}
}

func (this *TimeWheel) onTicker() {
	this.mu.Lock()

	this.pos = (this.pos + 1) % len(this.cs) // move the time pointer ahead
	broadcastCh := this.cs[this.pos]
	this.cs[this.pos] = make(chan struct{})

	this.mu.Unlock()

	// broadcast the timers: time is up!
	close(broadcastCh)
}
