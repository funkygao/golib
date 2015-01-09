package pool

import (
	"fmt"
	"github.com/funkygao/golib/sync2"
	log "github.com/funkygao/log4go"
	"time"
)

// ResourcePool allows you to use a pool of resources.
type ResourcePool struct {
	name         string
	resourcePool chan resourceWrapper
	factory      Factory
	capacity     sync2.AtomicInt64
	idleTimeout  sync2.AtomicDuration

	// stats
	waitCount sync2.AtomicInt64
	waitTime  sync2.AtomicDuration

	diagnosticTracker *DiagnosticTracker
}

type resourceWrapper struct {
	resource Resource
	timeUsed time.Time
}

// NewResourcePool creates a new ResourcePool pool.
// capacity is the initial capacity of the pool.
// maxCap is the maximum capacity of the pool.
// If a resource is unused beyond idleTimeout, it's discarded.
// An idleTimeout of 0 means that there is no timeout.
func NewResourcePool(name string, factory Factory, capacity, maxCap int,
	idleTimeout time.Duration, diagnosticInterval time.Duration,
	borrowMaxSeconds int) *ResourcePool {
	if capacity <= 0 || maxCap <= 0 || capacity > maxCap {
		panic("Invalid/out of range capacity")
	}

	this := &ResourcePool{
		name:         name,
		resourcePool: make(chan resourceWrapper, maxCap),
		factory:      factory,
		capacity:     sync2.AtomicInt64(capacity),
		idleTimeout:  sync2.AtomicDuration(idleTimeout),
	}
	this.diagnosticTracker = NewDiagnosticTracker(this)

	// put initial empty resources to pool
	for i := 0; i < capacity; i++ {
		this.resourcePool <- resourceWrapper{}
	}

	go this.diagnosticTracker.Run(diagnosticInterval, borrowMaxSeconds)

	return this
}

// Close empties the pool calling Close on all its resources.
// You can call Close while there are outstanding resources.
// It waits for all resources to be returned (Put).
// After a Close, Get and TryGet are not allowed.
func (this *ResourcePool) Close() {
	this.SetCapacity(0)
}

func (this *ResourcePool) IsClosed() (closed bool) {
	if this == nil {
		return true
	}
	return this.capacity.Get() == 0
}

// Get will return the next available resource. If capacity
// has not been reached, it will create a new one using the factory.
// Otherwise, it will indefinitely wait till the next resource becomes available.
func (this *ResourcePool) Get() (resource Resource, err error) {
	return this.get(true)
}

// TryGet will return the next available resource.
// If none is available, and capacity has not been reached, it
// will create a new one using the factory.
// Otherwise, it will return nil with no error.
func (this *ResourcePool) TryGet() (resource Resource, err error) {
	return this.get(false)
}

func (this *ResourcePool) get(wait bool) (resource Resource, err error) {
	if this == nil || this.IsClosed() {
		return nil, CLOSED_ERR
	}

	var (
		wrapper resourceWrapper
		ok      bool
	)
	select {
	case wrapper, ok = <-this.resourcePool:
		this.waitCount.Set(0) // reset
		if wrapper.resource != nil {
			this.diagnosticTracker.BorrowResource(wrapper.resource)
		}

	default:
		if !wait {
			return nil, nil
		}

		this.waitCount.Add(1)
		log.Warn("ResourcePool[%s] empty ready, pending: %d",
			this.name, this.WaitCount())

		t1 := time.Now()
		wrapper, ok = <-this.resourcePool
		this.waitTime.Add(time.Now().Sub(t1))
	}

	if !ok {
		return nil, CLOSED_ERR
	}

	// Close the aged idle resource
	timeout := this.idleTimeout.Get()
	if wrapper.resource != nil && timeout > 0 &&
		wrapper.timeUsed.Add(timeout).Sub(time.Now()) < 0 {
		log.Warn("ResourcePool[%s] resource:%d idle too long: closed", this.name,
			wrapper.resource.Id())
		wrapper.resource.Close()
		wrapper.resource = nil
	}

	if wrapper.resource == nil {
		wrapper.resource, err = this.factory()
		if err != nil {
			this.resourcePool <- resourceWrapper{}
		} else {
			this.diagnosticTracker.BorrowResource(wrapper.resource)
		}
	}

	return wrapper.resource, err
}

func (this *ResourcePool) Kill(resource Resource) {
	this.diagnosticTracker.ReturnResource(resource)
}

// Put will return a resource to the pool. For every successful Get,
// a corresponding Put is required. If you no longer need a resource,
// you will need to call Put(nil) instead of returning the closed resource.
// The will eventually cause a new resource to be created in its place.
func (this *ResourcePool) Put(resource Resource) {
	if this == nil {
		panic(CLOSED_ERR)
	}

	var wrapper resourceWrapper
	if resource != nil {
		wrapper = resourceWrapper{resource: resource, timeUsed: time.Now()}
	}

	select {
	case this.resourcePool <- wrapper:
		if wrapper.resource != nil {
			this.diagnosticTracker.ReturnResource(wrapper.resource)
		}

	default:
		if wrapper.resource != nil {
			wrapper.resource.Close() // FIXME should close it before discard it?
			log.Warn("ResourcePool[%s] full, resource:%d closed", this.name,
				wrapper.resource.Id())
		}
	}
}

// SetCapacity changes the capacity of the pool.
// You can use it to shrink or expand, but not beyond
// the max capacity. If the change requires the pool
// to be shrunk, SetCapacity waits till the necessary
// number of resources are returned to the pool.
// A SetCapacity of 0 is equivalent to closing the ResourcePool.
func (this *ResourcePool) SetCapacity(capacity int) error {
	if this == nil || capacity < 0 || capacity > cap(this.resourcePool) {
		return fmt.Errorf("capacity %d is out of range", capacity)
	}

	// Atomically swap new capacity with old, but only
	// if old capacity is non-zero.
	var oldcap int
	for {
		oldcap = int(this.capacity.Get())
		if oldcap == 0 {
			return CLOSED_ERR
		}
		if oldcap == capacity {
			return nil
		}
		if this.capacity.CompareAndSwap(int64(oldcap), int64(capacity)) {
			break
		}
	}

	if capacity < oldcap {
		for i := 0; i < oldcap-capacity; i++ {
			wrapper := <-this.resourcePool
			if wrapper.resource != nil {
				wrapper.resource.Close()
			}
		}
	} else {
		for i := 0; i < capacity-oldcap; i++ {
			this.resourcePool <- resourceWrapper{}
		}
	}

	if capacity == 0 {
		close(this.resourcePool)
	}

	return nil
}

func (this *ResourcePool) IdleTimeout() time.Duration {
	if this == nil {
		return 0
	}
	return this.idleTimeout.Get()
}

func (this *ResourcePool) SetIdleTimeout(idleTimeout time.Duration) {
	if this == nil {
		return
	}
	this.idleTimeout.Set(idleTimeout)
}

func (this *ResourcePool) Capacity() int64 {
	if this == nil {
		return 0
	}
	return this.capacity.Get()
}

func (this *ResourcePool) MaxCapacity() int64 {
	if this == nil {
		return 0
	}
	return int64(cap(this.resourcePool))
}

func (this *ResourcePool) Available() int64 {
	if this == nil {
		return 0
	}
	return int64(len(this.resourcePool))
}

func (this *ResourcePool) WaitCount() int64 {
	if this == nil {
		return 0
	}
	return this.waitCount.Get()
}

func (this *ResourcePool) WaitTime() time.Duration {
	if this == nil {
		return 0
	}
	return this.waitTime.Get()
}
