package idgen

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrorClockBackwards  = errors.New("Clock backwards")
	ErrorInvalidWorkerId = errors.New("Too big worker id")
)

const (
	WorkerIdBits       = uint64(5) // max 31
	SequenceBits       = uint64(12)
	WorkerIdShift      = SequenceBits
	TimestampLeftShift = SequenceBits + WorkerIdBits
	SequenceMask       = int64(-1) ^ (int64(-1) << SequenceBits)

	// Sat Jan  4 19:29:34 2014
	twepoch = int64(1388834974657)
)

// throughput of 5Million/s
type IdGenerator struct {
	mutex         sync.Mutex
	cookie        uint32 // random number to mitigate brute force lookups TODO
	wid           int64  // worker id
	seq           int64
	lastTimestamp int64
}

func NewIdGenerator(wid int) (this *IdGenerator, err error) {
	this = new(IdGenerator)
	this.wid = int64(wid)
	if wid > (1<<WorkerIdBits)-1 {
		return nil, ErrorInvalidWorkerId
	}
	return
}

func (this *IdGenerator) milliseconds() int64 {
	return time.Now().UnixNano() / 1e6
}

func (this *IdGenerator) Next() (int64, error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	ts := this.milliseconds()
	if ts < this.lastTimestamp {
		return 0, ErrorClockBackwards
	}

	if this.lastTimestamp == ts {
		this.seq = (this.seq + 1) & SequenceMask
		if this.seq == 0 {
			for ts <= this.lastTimestamp {
				ts = this.milliseconds()
			}
		}
	} else {
		this.seq = 0
	}

	this.lastTimestamp = ts

	r := ((ts - twepoch) << TimestampLeftShift) |
		(this.wid << WorkerIdShift) |
		this.seq
	return r, nil
}
