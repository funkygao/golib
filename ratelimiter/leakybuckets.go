package ratelimiter

import (
	"sync"
	"time"
)

type LeakyBuckets struct {
	buckets map[string]*LeakyBucket
	mu      sync.RWMutex

	size     int
	interval time.Duration
}

func NewLeakyBuckets(size int, leakInterval time.Duration) *LeakyBuckets {
	return &LeakyBuckets{
		buckets:  make(map[string]*LeakyBucket),
		size:     size,
		interval: leakInterval,
	}
}

func (this *LeakyBuckets) Pour(key string, amount int) bool {
	this.mu.Lock()
	if b, present := this.buckets[key]; present {
		this.mu.Unlock()
		return b.Pour(amount)
	}

	// key not present
	this.buckets[key] = NewLeakyBucket(this.size, this.interval)
	this.mu.Unlock()

	return this.buckets[key].Pour(amount)
}

func (this *LeakyBuckets) Delete(key string) {
	this.mu.Lock()
	this.buckets[key] = nil
	delete(this.buckets, key)
	this.mu.Unlock()
}
