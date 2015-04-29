package timewheel

import (
	"sync"
	"time"
)

type TimeWheel struct {
	sync.Mutex

	interval time.Duration

	ticker *time.Ticker
	quit   chan struct{}

	maxTimeout time.Duration

	cs []chan struct{}

	pos int
}

func NewTimeWheel(interval time.Duration, buckets int) *TimeWheel {
	w := new(TimeWheel)

	w.interval = interval

	w.quit = make(chan struct{})
	w.pos = 0

	w.maxTimeout = time.Duration(interval * (time.Duration(buckets)))

	w.cs = make([]chan struct{}, buckets)

	for i := range w.cs {
		w.cs[i] = make(chan struct{})
	}

	w.ticker = time.NewTicker(interval)
	go w.run()

	return w
}

func (w *TimeWheel) Stop() {
	close(w.quit)
}

func (w *TimeWheel) After(timeout time.Duration) <-chan struct{} {
	if timeout >= w.maxTimeout {
		panic("timeout too much, over maxtimeout")
	}

	w.Lock()

	index := (w.pos + int(timeout/w.interval)) % len(w.cs)

	b := w.cs[index]

	w.Unlock()

	return b
}

func (w *TimeWheel) run() {
	for {
		select {
		case <-w.ticker.C:
			w.onTicker()
		case <-w.quit:
			w.ticker.Stop()
			return
		}
	}
}

func (w *TimeWheel) onTicker() {
	w.Lock()

	lastC := w.cs[w.pos]
	w.cs[w.pos] = make(chan struct{})

	w.pos = (w.pos + 1) % len(w.cs)

	w.Unlock()

	close(lastC)
}
