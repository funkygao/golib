package cache

import (
	"fmt"
	"math/rand"
	"testing"
)

func BenchmarkCreateRandKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("mc_stress:%d", rand.Int())
	}
}

func BenchmarkLruCacheSetWithRandKey(b *testing.B) {
	b.ReportAllocs()
	lru := NewLruCache(1 << 20)
	var key string
	for i := 0; i < b.N; i++ {
		key = fmt.Sprintf("mc_stress:%d", rand.Int())
		lru.Set(key, 5)
	}
	b.SetBytes(int64(len(key)))
}

func BenchmarkLruCacheGetWithRandKey(b *testing.B) {
	b.ReportAllocs()
	lru := NewLruCache(1 << 10)
	var key string
	for i := 0; i < b.N; i++ {
		key = fmt.Sprintf("mc_stress:%d", rand.Int())
		lru.Set(key, 5)
		lru.Get(key)
	}
	b.SetBytes(int64(len(key)))
}

func BenchmarkSLruCacheGetWithRandKey(b *testing.B) {
	b.ReportAllocs()
	lru := NewSLruCache(1 << 10)
	var key string
	for i := 0; i < b.N; i++ {
		key = fmt.Sprintf("mc_stress:%d", rand.Int())
		lru.Set(key, 5)
		lru.Get(key)
	}
	b.SetBytes(int64(len(key)))
}

func BenchmarkLruCacheSetWithSameKey(b *testing.B) {
	b.ReportAllocs()
	lru := NewLruCache(1 << 10)
	var key string = "stuff"
	for i := 0; i < b.N; i++ {
		lru.Set(key, 1000)
	}
	b.SetBytes(int64(len(key)))
}

func BenchmarkLruCacheGetWithSameKey(b *testing.B) {
	b.ReportAllocs()
	lru := NewLruCache(1 << 10)
	var key string = "stuff"
	lru.Set(key, 1000)
	for i := 0; i < b.N; i++ {
		lru.Get(key)
	}
	b.SetBytes(int64(len(key)))
}
