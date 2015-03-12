package cache

import (
	"fmt"
	"hash/adler32"
	"hash/crc32"
	"hash/fnv"
	"math/rand"
	"strings"
	"testing"
)

// 4.5 ns/op
func BenchmarkGolangTypeAssert(b *testing.B) {
	var a interface{} = 45
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, ok := a.(int); ok {

		}
	}
}

func BenchmarkCreateRandKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("mc_stress:%d", rand.Int())
	}
}

func BenchmarkSetWithRandKey(b *testing.B) {
	b.ReportAllocs()
	lru := NewLruCache(1 << 20)
	var key string
	for i := 0; i < b.N; i++ {
		key = fmt.Sprintf("mc_stress:%d", rand.Int())
		lru.Set(key, 5)
	}
	b.SetBytes(int64(len(key)))
}

func BenchmarkGetWithRandKey(b *testing.B) {
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

func BenchmarkSetWithSameKey(b *testing.B) {
	b.ReportAllocs()
	lru := NewLruCache(1 << 10)
	var key string = "stuff"
	for i := 0; i < b.N; i++ {
		lru.Set(key, 1000)
	}
	b.SetBytes(int64(len(key)))
}

func BenchmarkGetWithSameKey(b *testing.B) {
	b.ReportAllocs()
	lru := NewLruCache(1 << 10)
	var key string = "stuff"
	lru.Set(key, 1000)
	for i := 0; i < b.N; i++ {
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

func BenchmarkFnv100B(b *testing.B) {
	b.ReportAllocs()
	key := strings.Repeat("s", 100)
	f := fnv.New32()
	for i := 0; i < b.N; i++ {
		f.Write([]byte(key))
		f.Sum32()
	}
	b.SetBytes(100)
}
