package cache

import (
	"testing"
)

func TestShardedGet(t *testing.T) {
	for _, tt := range getTests {
		slru := NewSLruCache(0)
		slru.Set(tt.keyToAdd, 1234)
		val, ok := slru.Get(tt.keyToGet)
		if ok != tt.expectedOk {
			t.Fatalf("%s: cache hit = %v; want %v", tt.name, ok, !ok)
		} else if ok && val != 1234 {
			t.Fatalf("%s expected get to return 1234 but got %v", tt.name, val)
		}
	}
}
