// thread-safe free list
package pool

type Pool struct {
	factory func() interface{}
	objects chan interface{}
}

func New(size int, factory func() interface{}) (this *Pool) {
	this = new(Pool)
	this.objects = make(chan interface{}, size)
	this.factory = factory
	return
}

func (this *Pool) Get() interface{} {
	select {
	case obj := <-this.objects:
		return obj
	default:
	}

	// no free resouce, create new object
	return this.factory()
}

func (this *Pool) Recycle(obj interface{}) bool {
	select {
	case this.objects <- obj:
		return true
	default:
		// it's full, so will not recyle
		return false
	}
}
