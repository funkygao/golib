package idgen

import (
	"errors"
	"os"
	"sync"
	"time"
)

var (
	ErrorClockBackwards = errors.New("Clock backwards")
)

const (
	workerIdBits       = uint64(5)
	datacenterIdBits   = uint64(5)
	sequenceBits       = uint64(12)
	workerIdShift      = sequenceBits
	datacenterIdShift  = sequenceBits + workerIdBits
	timestampLeftShift = sequenceBits + workerIdBits + datacenterIdBits
	sequenceMask       = int64(-1) ^ (int64(-1) << sequenceBits)

	// Tue, 21 Mar 2006 20:50:14.000 GMT
	twepoch = int64(1288834974657)
)

// throughput of 5Million/s
type IdGenerator struct {
	mutex         sync.Mutex
	cookie        uint32 // random number to mitigate brute force lookups TODO
	did           int64  // data center id
	wid           int64  // worker id
	seq           int64
	lastTimestamp int64
}

func NewIdGenerator() (this *IdGenerator) {
	this = new(IdGenerator)
	this.wid = int64(os.Getpid())
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
		this.seq = (this.seq + 1) & sequenceMask
		if this.seq == 0 {
			for ts <= this.lastTimestamp {
				ts = this.milliseconds()
			}
		}
	} else {
		this.seq = 0
	}

	this.lastTimestamp = ts

	r := ((ts - twepoch) << timestampLeftShift) |
		(this.did << datacenterIdShift) |
		(this.wid << workerIdShift) |
		this.seq
	return r, nil
}
