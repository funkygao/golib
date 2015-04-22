package main

import (
	"log"
	"sync/atomic"
	"time"
)

func recvMsg() {
	atomic.AddInt64(&recvN, 1)
}

func sentMsg() {
	atomic.AddInt64(&sentN, 1)
}

func runStats() {
	ticker := time.NewTicker(time.Second * time.Duration(options.tick))
	defer ticker.Stop()

	var lastRecv, lastSent int64

	for _ = range ticker.C {
		r := atomic.LoadInt64(&recvN)
		s := atomic.LoadInt64(&sentN)
		c := atomic.LoadInt64(&concurrency)
		if lastRecv != 0 || lastSent != 0 {
			log.Printf("c: %d recv qps: %d/%d send qps: %d/%d",
				c,
				(r-lastRecv)/options.tick,
				r,
				(s-lastSent)/options.tick,
				s)
		}

		lastSent = s
		lastRecv = r
	}
}
