package cmap

import (
	"github.com/funkygao/golib/hack"
	"hash/fnv"
)

var (
	SHARD_COUNT = 32
)

func New() ConcurrentMap {
	this := make(ConcurrentMap, SHARD_COUNT)
	for i := 0; i < SHARD_COUNT; i++ {
		this[i] = &concurrentMapSharded{items: make(map[string]interface{})}
	}
	return this
}

// Returns shard under given key.
func (m ConcurrentMap) GetShard(key string) *concurrentMapSharded {
	hasher := fnv.New32()
	hasher.Write(hack.Byte(key))
	return m[uint(hasher.Sum32())%uint(SHARD_COUNT)]
}

func (m *ConcurrentMap) Set(key string, value interface{}) {
	shard := m.GetShard(key)
	shard.Lock()
	shard.items[key] = value
	shard.Unlock()
}

func (m ConcurrentMap) Get(key string) (interface{}, bool) {
	shard := m.GetShard(key)
	shard.RLock()
	val, ok := shard.items[key]
	shard.RUnlock()
	return val, ok
}

func (m ConcurrentMap) Count() int {
	count := 0
	for i := 0; i < SHARD_COUNT; i++ {
		shard := m[i]
		shard.RLock()
		count += len(shard.items)
		shard.RUnlock()
	}
	return count
}

func (m *ConcurrentMap) Has(key string) bool {
	shard := m.GetShard(key)
	shard.RLock()
	_, ok := shard.items[key]
	shard.RUnlock()
	return ok
}

func (m *ConcurrentMap) Remove(key string) {
	shard := m.GetShard(key)
	shard.Lock()
	delete(shard.items, key)
	shard.Unlock()
}

// Returns an iterator which could be used in a for range loop.
func (m ConcurrentMap) Iter() <-chan Tuple {
	ch := make(chan Tuple)
	go func() {
		for _, shard := range m {
			shard.RLock()
			for key, val := range shard.items {
				ch <- Tuple{key, val}
			}
			shard.RUnlock()
		}
		close(ch)
	}()
	return ch
}

// Returns a buffered iterator which could be used in a for range loop.
func (m ConcurrentMap) IterBuffered() <-chan Tuple {
	ch := make(chan Tuple, m.Count())
	go func() {
		// Foreach shard.
		for _, shard := range m {
			// Foreach key, value pair.
			shard.RLock()
			for key, val := range shard.items {
				ch <- Tuple{key, val}
			}
			shard.RUnlock()
		}
		close(ch)
	}()
	return ch
}
