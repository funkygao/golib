package cmap

import (
	"github.com/funkygao/golib/rand"
	"hash/fnv"
	"strings"
	"testing"
)

func BenchmarkHashFnvNew32(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = fnv.New32()
	}
}

func BenchmarkGetShardWithKeyLen10(b *testing.B) {
	cm := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cm.GetShard("user.12121")
	}
}

func BenchmarkGetShardWithKeyLen100(b *testing.B) {
	cm := New()
	b.ReportAllocs()
	key := strings.Repeat("a", 100)
	for i := 0; i < b.N; i++ {
		cm.GetShard(key)
	}
}

func BenchmarkSetAndGetWithShard32(b *testing.B) {
	cm := New()
	for i := 0; i < ('~'-'!')*('~'-'!'); i++ {
		cm.Set(rand.RandomString(2), 1)
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		key := rand.RandomString(2)
		cm.Set(key, 1)
		_, _ = cm.Get(key)
	}
}

func BenchmarkSetAndGetWithShard1(b *testing.B) {
	SHARD_COUNT = 1
	cm := New()
	for i := 0; i < ('~'-'!')*('~'-'!'); i++ {
		cm.Set(rand.RandomString(2), 1)
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		key := rand.RandomString(2)
		cm.Set(key, 1)
		_, _ = cm.Get(key)
	}
}

func BenchmarkBuiltinUnsafeMap(b *testing.B) {
	m := make(map[string]interface{})
	for i := 0; i < ('~'-'!')*('~'-'!'); i++ {
		m[rand.RandomString(2)] = 1
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		key := rand.RandomString(2)
		m[key] = 1
		_, _ = m[key]
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
