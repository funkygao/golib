package cache

import (
	"github.com/funkygao/golib/hack"
	"hash/fnv"
)

const (
	shardN = 32
)

// Sharded LruCache.
type SLruCache []*LruCache

func NewSLruCache(maxItems int) SLruCache {
	c := make(SLruCache, shardN)
	for i := 0; i < shardN; i++ {
		c[i] = NewLruCache(maxItems)
	}
	return c
}

// Returns shard LruCache under given key.
// Only string key type supported.
// TODO
func (c SLruCache) GetShard(key Key) *LruCache {
	hasher := fnv.New32()
	if k, ok := key.(string); ok {
		hasher.Write(hack.Byte(k))
		return c[uint(hasher.Sum32())%uint(shardN)]
	}

	return c[0] // unable to shard, always use the 1st slot
}

func (c SLruCache) Purge() {
	for _, lru := range c {
		lru.Purge()
	}
}

// Set adds a value to the cache.
// If key already exists, its value gets overwritten.
func (c SLruCache) Set(key Key, value interface{}) {
	shard := c.GetShard(key)
	shard.Set(key, value)
}

// Add will return true and set the key to cache if key not existent, else return false.
func (c SLruCache) Add(key Key, value interface{}) bool {
	shard := c.GetShard(key)
	return shard.Add(key, value)
}

// Get looks up a key's value from the cache.
func (c SLruCache) Get(key Key) (value interface{}, ok bool) {
	shard := c.GetShard(key)
	return shard.Get(key)
}

// Keys return active keys in the cache.
// Order is not garranteed.
func (c SLruCache) Keys() []interface{} {
	keys := make([]interface{}, 0, c.Len())
	for _, lru := range c {
		for _, k := range lru.Keys() {
			keys = append(keys, k)
		}
	}
	return keys
}

// Len returns the number of items in the cache.
func (c SLruCache) Len() int {
	l := 0
	for _, lru := range c {
		l += lru.Len()
	}
	return l
}

func (c SLruCache) Inc(key Key, delta int) (newVal int) {
	shard := c.GetShard(key)
	return shard.Inc(key, delta)
}

func (c SLruCache) Del(key Key) {
	shard := c.GetShard(key)
	shard.Del(key)
}
