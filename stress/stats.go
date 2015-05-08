package stress

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var (
	counterMutex   sync.RWMutex
	defaultCounter Counter = make(Counter)
)

func IncCounter(key string, delta int64) {
	defaultCounter.add(key, delta)
}

type Counter map[string]int64

func (this Counter) add(key string, delta int64) {
	counterMutex.Lock()
	if _, present := this[key]; present {
		this[key] += delta
	} else {
		this[key] = delta
	}
	counterMutex.Unlock()
}

func runConsoleStats() {
	ticker := time.NewTicker(time.Second * time.Duration(flags.tick))
	defer ticker.Stop()

	lastCounter := make(Counter)
	minCounter := make(Counter)
	maxCounter := make(Counter)
	for _ = range ticker.C {
		counterMutex.RLock()
		s := ""
		c := atomic.LoadInt32(&concurrency)
		gn := runtime.NumGoroutine()
		for k, v := range defaultCounter {
			min := (v - lastCounter[k]) / flags.tick
			max := (v - lastCounter[k]) / flags.tick
			x, present := minCounter[k]
			if !present {
				minCounter[k] = min
			} else if x <= min {
				min = x
			} else {
				minCounter[k] = min
			}
			x, present = maxCounter[k]
			if !present {
				maxCounter[k] = max
			} else if x >= max {
				max = x
			} else {
				maxCounter[k] = max
			}

			s += fmt.Sprintf("%s:%d/%d,%d-%d ", k, (v-lastCounter[k])/flags.tick, v,
				min, max)
			lastCounter[k] = v
		}
		log.Printf("c:%d go:%d qps: {%s}", c, gn, s)
		counterMutex.RUnlock()
	}

}
