package stress

import (
	"fmt"
	"log"
	"runtime"
	"sort"
	"sync/atomic"
	"time"
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

func runSlaveReporter() {
	ticker := time.NewTicker(time.Second * time.Duration(Flags.Tick))
	defer ticker.Stop()

	for range ticker.C {
		reportToMaster()
	}
}

func runMasterReporter() {
	ticker := time.NewTicker(time.Second * time.Duration(Flags.Tick))
	defer ticker.Stop()

	lastCounter := make(Counter)
	minCounter := make(Counter)
	maxCounter := make(Counter)
	for range ticker.C {
		counterMutex.RLock()
		s := ""
		c := atomic.LoadInt32(&concurrency)
		gn := runtime.NumGoroutine()
		sortedKeys := make([]string, 0, len(defaultCounter))
		for k, _ := range defaultCounter {
			sortedKeys = append(sortedKeys, k)
		}
		sort.Strings(sortedKeys)
		for _, k := range sortedKeys {
			v := defaultCounter[k]

			min := (v - lastCounter[k]) / Flags.Tick
			max := (v - lastCounter[k]) / Flags.Tick
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

			s += fmt.Sprintf("%s:%6d/%-9d(%4d-%-6d) ", k, (v-lastCounter[k])/Flags.Tick, v,
				min, max)
			lastCounter[k] = v
		}
		log.Printf("c:%-5d go:%-5d qps: {%s}", c, gn, s)
		counterMutex.RUnlock()
	}

}
