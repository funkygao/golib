// Package ratelimiter implements the Leaky Bucket ratelimiting algorithm.
package ratelimiter

import (
	"time"
)

type LeakyBucket struct {
	capacity  int64
	remaining int64
	reset     time.Time
	rate      time.Duration
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

// When amount is 0, just check if the bucket is full.
func (b *LeakyBucket) Pour(amount int) bool {
	if time.Now().After(b.reset) {
		b.reset = time.Now().Add(b.rate)
		b.remaining = b.capacity
	}

	if amount == 0 {
		// check whether more pour permitted
		return b.remaining > 0
	}

	amount64 := int64(amount)
	if amount64 > b.remaining {
		return false
	} else {
		b.remaining -= amount64
		return true
	}
}
