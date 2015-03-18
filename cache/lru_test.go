package cache

import (
	"github.com/funkygao/assert"
	"reflect"
	"testing"
)

type simpleStruct struct {
	int
	string
}

type complexStruct struct {
	int
	simpleStruct
}

var getTests = []struct {
	name       string
	keyToAdd   interface{}
	keyToGet   interface{}
	expectedOk bool
}{
	{"string_hit", "myKey", "myKey", true},
	{"string_miss", "myKey", "nonsense", false},
	{"simple_struct_hit", simpleStruct{1, "two"}, simpleStruct{1, "two"}, true},
	{"simeple_struct_miss", simpleStruct{1, "two"}, simpleStruct{0, "noway"}, false},
	{"complex_struct_hit", complexStruct{1, simpleStruct{2, "three"}},
		complexStruct{1, simpleStruct{2, "three"}}, true},
}

func TestLruCacheGet(t *testing.T) {
	for _, tt := range getTests {
		lru := NewLruCache(0)
		lru.Set(tt.keyToAdd, 1234)
		val, ok := lru.Get(tt.keyToGet)
		if ok != tt.expectedOk {
			t.Fatalf("%s: cache hit = %v; want %v", tt.name, ok, !ok)
		} else if ok && val != 1234 {
			t.Fatalf("%s expected get to return 1234 but got %v", tt.name, val)
		}
	}
}

func TestLruCacheDel(t *testing.T) {
	lru := NewLruCache(0)
	lru.Set("myKey", 1234)
	if val, ok := lru.Get("myKey"); !ok {
		t.Fatal("TestRemove returned no match")
	} else if val != 1234 {
		t.Fatalf("TestRemove failed.  Expected %d, got %v", 1234, val)
	}

	lru.Del("myKey")
	if _, ok := lru.Get("myKey"); ok {
		t.Fatal("TestRemove returned a removed entry")
	}
}

func TestLruCacheInc(t *testing.T) {
	lru := NewLruCache(10)
	counter := lru.Inc("foo", 1)
	assert.Equal(t, 1, counter)
	counter = lru.Inc("foo", 1)
	assert.Equal(t, 2, counter)
	assert.Equal(t, 5, lru.Inc("foo", 3))
	assert.Equal(t, 0, lru.Inc("foo", -5))
	lru.Del("foo")
	counter = lru.Inc("foo", 1)
	assert.Equal(t, 1, counter)
}

func TestLruCacheLenAndPurge(t *testing.T) {
	lru := NewLruCache(0)
	assert.Equal(t, 0, lru.Len())
	lru.Set("myKey", 1234)
	assert.Equal(t, 1, lru.Len())
	lru.Purge()
	assert.Equal(t, 0, lru.Len())
}

func TestLruCacheKeys(t *testing.T) {
	lru := NewLruCache(0)
	lru.Set("myKey", 1234)
	lru.Inc("boo", 1)
	keys := lru.Keys()
	reflect.DeepEqual(keys, []interface{}{"myKey", "boo"})
	assert.Equal(t, true, reflect.DeepEqual(keys, []interface{}{"myKey", "boo"}) ||
		reflect.DeepEqual(keys, []interface{}{"boo", "myKey"}))
}

func TestLruCacheAdd(t *testing.T) {
	lru := NewLruCache(0)
	assert.Equal(t, true, lru.Add("key", 1))
	assert.Equal(t, false, lru.Add("key", 3))
}
