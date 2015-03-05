package cmap

import (
	"hash/fnv"
	"testing"
)

func BenchmarkHashFnvNew32(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = fnv.New32()
	}
}

func BenchmarkGetShard(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {

	}
}

func BenchmarkSetAndGetWithShard32(b *testing.B) {
	const key = "key"
	cm := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cm.Set(key, 1)
		_, _ = cm.Get(key)
	}
}

func BenchmarkSetAndGetWithShard1(b *testing.B) {
	SHARD_COUNT = 1
	const key = "key"
	cm := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cm.Set(key, 1)
		_, _ = cm.Get(key)
	}
}

func BenchmarkNotSafeMap(b *testing.B) {
	m := make(map[string]interface{})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		m["key"] = 1
		_, _ = m["key"]
	}
}

func BenchmarkHas(b *testing.B) {
	cm := New()
	cm.Set("key", 1)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = cm.Has("key")
	}

}

func BenchmarkHasNot(b *testing.B) {
	cm := New()
	cm.Set("key", 1)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = cm.Has("key_not_exist")
	}

}
