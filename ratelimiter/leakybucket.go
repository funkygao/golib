// Package ratelimiter implements the Leaky Bucket ratelimiting algorithm.
package ratelimiter

import (
	"sync"
	"time"
)

type LeakyBucket struct {
	capacity  int64
	remaining int64
	reset     time.Time
	rate      time.Duration
	mutex     sync.Mutex
}

func NewLeakyBucket(capacity int64, rate time.Duration) *LeakyBucket {
	bucket := LeakyBucket{
		capacity:  capacity,
		remaining: capacity,
		reset:     time.Now().Add(rate),
		rate:      rate,
	}

	return &bucket
}

func (b *LeakyBucket) Pour(amount int) bool {
	//b.mutex.Lock()
	//defer b.mutex.Unlock()

	if time.Now().After(b.reset) {
		b.reset = time.Now().Add(b.rate)
		b.remaining = b.capacity
	}

	amount64 := int64(amount)
	if amount64 > b.remaining {
		return false
	} else {
		b.remaining -= amount64
		return true
	}
}
