package recycler

import (
	"container/list"
	"sync/atomic"
	"time"
)

var (
	makes int64
)

type queued struct {
	when time.Time
	data interface{}
}

func New(poolSize int, factory func() interface{}) (get, give chan interface{}) {
	get = make(chan interface{}, poolSize)
	give = make(chan interface{}, poolSize)

	go houseKeeping(get, give, factory)

	return
}

func houseKeeping(get, give chan interface{},
	factory func() interface{}) {
	q := new(list.List)
	for {
		if q.Len() == 0 {
			atomic.AddInt64(&makes, 1)
			q.PushFront(queued{when: time.Now(), data: factory()})
		}

		element := q.Front()
		timeout := time.NewTimer(time.Minute)
		select {
		case b := <-give:
			timeout.Stop()
			q.PushFront(queued{when: time.Now(), data: b})

		case get <- element.Value.(queued).data:
			timeout.Stop()
			q.Remove(element)

		case <-timeout.C:
			e := q.Front()
			for e != nil {
				n := e.Next()
				if time.Since(e.Value.(queued).when) > time.Minute {
					q.Remove(e)
					e.Value = nil
				}
				e = n
			}
		}
	}
}
