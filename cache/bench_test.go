package cache

import (
	"fmt"
	"hash/adler32"
	"hash/crc32"
	"math/rand"
	"strings"
	"testing"
)

func BenchmarkCreateKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("mc_stress:%d", rand.Int())
	}
}

func BenchmarkSet(b *testing.B) {
	b.ReportAllocs()
	lru := NewLruCache(1 << 20)
	var key string
	for i := 0; i < b.N; i++ {
		key = fmt.Sprintf("mc_stress:%d", rand.Int())
		lru.Set(key, 5)
	}
	b.SetBytes(int64(len(key)))
}

func BenchmarkGet(b *testing.B) {
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

func BenchmarkCrc32Of100B(b *testing.B) {
	b.ReportAllocs()
	key := strings.Repeat("s", 100)
	for i := 0; i < b.N; i++ {
		crc32.ChecksumIEEE([]byte(key))
	}
	b.SetBytes(100)
}

func BenchmarkAdler32Of100B(b *testing.B) {
	b.ReportAllocs()
	key := strings.Repeat("s", 100)
	for i := 0; i < b.N; i++ {
		adler32.Checksum([]byte(key))
	}
	b.SetBytes(100)
}
