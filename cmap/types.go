package cmap

import (
	"sync"
)

// A thread safe map.
// To avoid lock bottlenecks this map is dived to several (SHARD_COUNT) map shards.
type ConcurrentMap []*concurrentMapSharded

type concurrentMapSharded struct {
	items map[string]interface{}
	sync.RWMutex
}

// Used by the Iter & IterBuffered functions to wrap two variables together over a channel,
type Tuple struct {
	Key string
	Val interface{}
}
