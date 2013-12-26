package cache

import (
	"container/list"
	"sync"
)

// LRU cache implementation.
type LruCache struct {
	Cacheable
	HasLength

	*sync.RWMutex

	// MaxEntries is the maximum number of cache entries before
	// an item is evicted. Zero means no limit.
	MaxEntries int

	// OnEvicted optionally specificies a callback function to be
	// executed when an entry is purged from the cache.
	OnEvicted func(key Key, value interface{})

	ll    *list.List // double linked list
	cache map[interface{}]*list.Element
}

// New creates a new LruCache.
// If maxEntries is zero, the cache has no limit and it's assumed
// that eviction is done by the caller.
func NewLruCache(maxEntries int) *LruCache {
	return &LruCache{
		MaxEntries: maxEntries,
		ll:         list.New(),
		cache:      make(map[interface{}]*list.Element),
		RWMutex:    new(sync.RWMutex),
	}
}

// Add adds a value to the cache.
func (c *LruCache) Set(key Key, value interface{}) {
	c.Lock()
	defer c.Unlock()

	if c.cache == nil {
		c.cache = make(map[interface{}]*list.Element)
		c.ll = list.New()
	}
	if ee, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ee)
		ee.Value.(*entry).value = value
		return
	}
	ele := c.ll.PushFront(&entry{key, value})
	c.cache[key] = ele
	if c.MaxEntries != 0 && c.ll.Len() > c.MaxEntries {
		// evict olded element
		c.removeOldest()
	}
}

// Get looks up a key's value from the cache.
func (c *LruCache) Get(key Key) (value interface{}, ok bool) {
	c.RLock()
	defer c.RUnlock()

	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		c.ll.MoveToFront(ele)
		return ele.Value.(*entry).value, true
	}
	return
}

func (c *LruCache) Del(key Key) {
	c.Lock()
	defer c.Unlock()

	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		c.removeElement(ele)
	}
}

// RemoveOldest removes the oldest item from the cache.
func (c *LruCache) removeOldest() {
	if c.cache == nil {
		return
	}
	ele := c.ll.Back()
	if ele != nil {
		c.removeElement(ele)
	}
}

func (c *LruCache) removeElement(e *list.Element) {
	c.ll.Remove(e)
	kv := e.Value.(*entry)
	delete(c.cache, kv.key)
	if c.OnEvicted != nil {
		c.OnEvicted(kv.key, kv.value)
	}
}

// Len returns the number of items in the cache.
func (c *LruCache) Len() int {
	c.RLock()
	defer c.RUnlock()

	if c.cache == nil {
		return 0
	}
	return c.ll.Len()
}
