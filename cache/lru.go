package cache

import (
	"container/list"
	"sync"
)

// goroutine safe LRU cache implementation.
type LruCache struct {
	Cacheable
	HasLength

	*sync.Mutex

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
		cache:      make(map[interface{}]*list.Element, maxEntries),
		Mutex:      new(sync.Mutex),
	}
}

// Add adds a value to the cache.
func (c *LruCache) Set(key Key, value interface{}) {
	c.Lock()

	if ee, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ee)
		ee.Value.(*entry).value = value
		c.Unlock()
		return
	}

	ele := c.ll.PushFront(&entry{key, value})
	c.cache[key] = ele
	if c.MaxEntries != 0 && c.ll.Len() > c.MaxEntries {
		// evict olded element
		c.removeOldest()
	}

	c.Unlock()
}

// Get looks up a key's value from the cache.
func (c *LruCache) Get(key Key) (value interface{}, ok bool) {
	c.Lock()

	if ele, hit := c.cache[key]; hit {
		c.ll.MoveToFront(ele)
		c.Unlock()
		return ele.Value.(*entry).value, true
	}

	c.Unlock()
	return
}

func (c *LruCache) Decr(key Key) (value int) {
	c.Lock()

	if ee, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ee)
		counter := ee.Value.(*entry).value.(int)
		ee.Value.(*entry).value = counter - 1
		c.Unlock()
		return counter - 1
	}

	// 1st element
	ele := c.ll.PushFront(&entry{key, 0})
	c.cache[key] = ele
	if c.MaxEntries != 0 && c.ll.Len() > c.MaxEntries {
		// evict olded element
		c.removeOldest()
	}

	c.Unlock()
	return 0
}

func (c *LruCache) Inc(key Key) (value int) {
	c.Lock()

	if ee, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ee)
		counter := ee.Value.(*entry).value.(int)
		ee.Value.(*entry).value = counter + 1
		c.Unlock()
		return counter + 1
	}

	// 1st element
	ele := c.ll.PushFront(&entry{key, 1})
	c.cache[key] = ele
	if c.MaxEntries != 0 && c.ll.Len() > c.MaxEntries {
		// evict olded element
		c.removeOldest()
	}

	c.Unlock()
	return 1
}

func (c *LruCache) Del(key Key) {
	c.Lock()
	if ele, hit := c.cache[key]; hit {
		c.removeElement(ele)
	}
	c.Unlock()
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
	c.Lock()
	defer c.Unlock()

	if c.cache == nil {
		return 0
	}
	return c.ll.Len()
}
