package ratelimiter

import (
	"sync"
	"time"
)

type LeakyBuckets struct {
	buckets map[string]*LeakyBucket
	mu      sync.RWMutex

	capacity int64
	rate     time.Duration
}

func NewLeakyBuckets(capacity int64, rate time.Duration) *LeakyBuckets {
	return &LeakyBuckets{
		buckets:  make(map[string]*LeakyBucket, 20),
		capacity: capacity,
		rate:     rate,
	}
}

// When amount is 0, just check if the bucket is full.
func (this *LeakyBuckets) Pour(key string, amount int) bool {
	this.mu.Lock()
	if b, present := this.buckets[key]; present {
		r := b.Pour(amount)
		this.mu.Unlock()
		return r
	}

	// key not present
	this.buckets[key] = NewLeakyBucket(this.capacity, this.rate)
	r := this.buckets[key].Pour(amount)
	this.mu.Unlock()
	return r
}

func (this *LeakyBuckets) Delete(key string) {
	this.mu.Lock()
	this.buckets[key] = nil
	delete(this.buckets, key)
	this.mu.Unlock()
}
