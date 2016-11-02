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
	lastQps := make(Counter)
	minCounter := make(Counter)
	maxCounter := make(Counter)
	for range ticker.C {
		counterMutex.RLock()
		c := atomic.LoadInt32(&concurrency)
		gn := runtime.NumGoroutine()
		slaves := atomic.LoadInt32(&activeSlaves)
		counters := make(Counter, len(defaultCounter))
		for k, v := range defaultCounter {
			counters[k] = v
		}
		for _, stat := range globalStats {
			c += stat.C
			gn += stat.G
			for k, v := range stat.Counter {
				counters[k] += v
			}
		}

		// reset active slaves counter each tick
		atomic.StoreInt32(&activeSlaves, 0)

		sortedKeys := make([]string, 0, len(counters))
		for k, _ := range counters {
			sortedKeys = append(sortedKeys, k)
		}
		sort.Strings(sortedKeys)
		s := ""
		for _, k := range sortedKeys {
			v := counters[k]

			// calc min/max
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

			qps := (v - lastCounter[k]) / Flags.Tick
			qpsDeltaDirection := "▲"
			qpsDelta := qps - lastQps[k]
			if qpsDelta < 0 {
				qpsDelta = -qpsDelta
				qpsDeltaDirection = "▾"
			}

			s += fmt.Sprintf("%s: %s%-5d %5d/%-9d(%4d-%-5d) ", k, qpsDeltaDirection, qpsDelta, qps, v, min, max)

			lastQps[k] = qps
			lastCounter[k] = v
		}
		log.Printf("slave:%-2d C:%-5d G:%-5d qps: {%s}", slaves, c, gn, s)
		counterMutex.RUnlock()
	}

}
