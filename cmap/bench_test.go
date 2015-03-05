package cmap

import (
	"github.com/funkygao/golib/fixture"
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
	cm := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cm.GetShard("user.121212")
	}
}

func BenchmarkSetAndGetWithShard32(b *testing.B) {
	cm := New()
	for i := 0; i < ('~'-'!')*('~'-'!'); i++ {
		cm.Set(fixture.RandomString(2), 1)
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		key := fixture.RandomString(2)
		cm.Set(key, 1)
		_, _ = cm.Get(key)
	}
}

func BenchmarkSetAndGetWithShard1(b *testing.B) {
	SHARD_COUNT = 1
	cm := New()
	for i := 0; i < ('~'-'!')*('~'-'!'); i++ {
		cm.Set(fixture.RandomString(2), 1)
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		key := fixture.RandomString(2)
		cm.Set(key, 1)
		_, _ = cm.Get(key)
	}
}

func BenchmarkNotSafeMap(b *testing.B) {
	m := make(map[string]interface{})
	for i := 0; i < ('~'-'!')*('~'-'!'); i++ {
		m[fixture.RandomString(2)] = 1
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		key := fixture.RandomString(2)
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
