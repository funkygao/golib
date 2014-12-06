package mutexmap

import (
	"github.com/funkygao/golib/cache"
	"sync"
)

type MutexMap struct {
	mu    sync.Mutex
	items *cache.LruCache
}

func New(maxEntries int) *MutexMap {
	this := &MutexMap{}
	this.items = cache.NewLruCache(maxEntries)
	return this
}

// block till acquire the lock
func (this *MutexMap) Lock(key cache.Key) {
	this.mu.Lock()

	value, present := this.items.Get(key)
	if !present {
		lock := &sync.Mutex{}
		lock.Lock()
		this.items.Set(key, lock)

		this.mu.Unlock()
		return
	}

	// this key already exists in items
	lock := value.(*sync.Mutex)
	lock.Lock()

	this.mu.Unlock()
}

func (this *MutexMap) Unlock(key cache.Key) {
	value, present := this.items.Get(key) // must be always present
	if !present {
		return
	}

	lock := value.(*sync.Mutex)
	lock.Unlock()
}
