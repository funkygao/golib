package peer

import (
	"sync/atomic"
)

type LamportClock struct {
	v uint64
}

type LamportTime uint64

func (l *LamportClock) Time() LamportTime {
	return LamportTime(atomic.LoadUint64(&l.v))
}

func (l *LamportClock) Inc() LamportTime {
	return LamportTime(atomic.AddUint64(&l.v, 1))
}

func (l *LamportClock) Witness(v LamportTime) {
	for {
		this := atomic.LoadUint64(&l.v)
		that := uint64(v)
		if that < this {
			// that value is older
			return
		}

		// ensure that our local clock is at least 1 ahead
		if atomic.CompareAndSwapUint64(&l.v, this, that+1) {
			break
		}
	}
}
