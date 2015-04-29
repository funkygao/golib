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
	for _ = range ticker.C {
		counterMutex.RLock()
		s := ""
		c := atomic.LoadInt32(&concurrency)
		gn := runtime.NumGoroutine()
		for k, v := range defaultCounter {
			s += fmt.Sprintf("%s:%d/%d ", k, (v-lastCounter[k])/flags.tick, v)
			lastCounter[k] = v
		}
		log.Printf("c:%d go:%d qps: {%s}", c, gn, s)
		counterMutex.RUnlock()
	}

}
