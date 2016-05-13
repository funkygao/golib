package ratelimiter

import (
	"testing"
	"time"
)

func BenchmarkLeakyBucketsPour(b *testing.B) {
	lb := NewLeakyBuckets(1000, time.Minute)
	for i := 1; i < b.N; i++ {
		lb.Pour("bar", 1)
		lb.Pour("foo", 22)
	}
}

func BenchmarkLeakyBucketsDelete(b *testing.B) {
	lb := NewLeakyBuckets(1000, time.Minute)
	for i := 1; i < b.N; i++ {
		lb.Delete("foo")
	}
}

func BenchmarkLeakyBucketPour(b *testing.B) {
	lk := NewLeakyBucket(100, time.Second)
	for i := 0; i < b.N; i++ {
		lk.Pour(1)
	}
}
