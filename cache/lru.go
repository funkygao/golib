package cache

import (
	"container/list"
	"sync"
)

// goroutine safe LRU cache implementation.
type LruCache struct {
	Cacheable
	HasLength

	lock sync.RWMutex // not embedded because lock is transparent for caller

	// maxItems is the maximum number of cache entries before
	// an item is evicted. Zero means no limit.
	maxItems    int
	initialSize int

	// OnEvicted optionally specificies a callback function to be
	// executed when an entry is purged from the cache.
	OnEvicted func(key Key, value interface{})

	ll    *list.List // double linked list
	items map[interface{}]*list.Element
}

// New creates a new LruCache.
// If maxItems is zero, the cache has no limit and it's assumed
// that eviction is done by the caller.
func NewLruCache(maxItems int) *LruCache {
	const M = 1 << 20
	var sz = maxItems
	if maxItems > M {
		sz = M
	}
	return &LruCache{
		maxItems:    maxItems,
		ll:          list.New(),
		items:       make(map[interface{}]*list.Element, sz),
		initialSize: sz,
	}
}

func (c *LruCache) Purge() {
	c.lock.Lock()
	c.ll = list.New()
	c.items = make(map[interface{}]*list.Element, c.initialSize)
	c.lock.Unlock()
}

func (c *LruCache) MaxItems() int {
	return c.maxItems
}

// Set adds a value to the cache.
// If key already exists, its value gets overwritten.
func (c *LruCache) Set(key Key, value interface{}) {
	c.lock.Lock()

	if item, ok := c.items[key]; ok {
		c.ll.MoveToFront(item)
		item.Value.(*entry).value = value
		c.lock.Unlock()
		return
	}

	item := c.ll.PushFront(&entry{key, value})
	c.items[key] = item
	if c.maxItems != 0 && c.ll.Len() > c.maxItems {
		// evict olded element
		c.removeOldest()
	}

	c.lock.Unlock()
}

// Get looks up a key's value from the cache.
func (c *LruCache) Get(key Key) (value interface{}, ok bool) {
	c.lock.RLock()

	if item, hit := c.items[key]; hit {
		c.ll.MoveToFront(item)
		c.lock.RUnlock()
		return item.Value.(*entry).value, true
	}

	c.lock.RUnlock()
	return
}

// Keys return active keys in the cache.
// Order is not garranteed.
func (c *LruCache) Keys() []interface{} {
	c.lock.RLock()

	keys := make([]interface{}, len(c.items))
	i := 0
	for k, _ := range c.items {
		keys[i] = k
		i++
	}

	c.lock.RUnlock()
	return keys
}

func (c *LruCache) Decr(key Key) (value int) {
	c.lock.Lock()

	if item, ok := c.items[key]; ok {
		c.ll.MoveToFront(item)
		counter := item.Value.(*entry).value.(int)
		item.Value.(*entry).value = counter - 1
		c.lock.Unlock()
		return counter - 1
	}

	// 1st element
	item := c.ll.PushFront(&entry{key, 0})
	c.items[key] = item
	if c.maxItems != 0 && c.ll.Len() > c.maxItems {
		// evict olded element
		c.removeOldest()
	}

	c.lock.Unlock()
	return 0
}

func (c *LruCache) Inc(key Key) (value int) {
	c.lock.Lock()

	if item, ok := c.items[key]; ok {
		c.ll.MoveToFront(item)
		counter := item.Value.(*entry).value.(int)
		item.Value.(*entry).value = counter + 1
		c.lock.Unlock()
		return counter + 1
	}

	// 1st element
	item := c.ll.PushFront(&entry{key, 1})
	c.items[key] = item
	if c.maxItems != 0 && c.ll.Len() > c.maxItems {
		// evict olded element
		c.removeOldest()
	}

	c.lock.Unlock()
	return 1
}

func (c *LruCache) Del(key Key) {
	c.lock.Lock()
	if item, hit := c.items[key]; hit {
		c.removeElement(item)
	}
	c.lock.Unlock()
}

// Len returns the number of items in the cache.
func (c *LruCache) Len() int {
	c.lock.RLock()

	if c.items == nil {
		c.lock.RUnlock()
		return 0
	}

	c.lock.RUnlock()
	return c.ll.Len()
}

// RemoveOldest removes the oldest item from the cache.
func (c *LruCache) removeOldest() {
	if c.items == nil {
		return
	}
	if item := c.ll.Back(); item != nil {
		c.removeElement(item)
	}
}

func (c *LruCache) removeElement(e *list.Element) {
	c.ll.Remove(e)
	kv := e.Value.(*entry)
	delete(c.items, kv.key)
	if c.OnEvicted != nil {
		c.OnEvicted(kv.key, kv.value)
	}
}
