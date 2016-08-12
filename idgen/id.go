// Package idgen is a twitter snowflake implementation in golang.
package idgen

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

var (
	ErrorClockBackwards  = errors.New("Clock backwards")
	ErrorInvalidWorkerId = errors.New("Too big worker id")
	ErrorInvalidTag      = errors.New("Too big tag")
)

const (
	WorkerIdBits = uint64(5)  // max 31
	TagIdBits    = uint64(5)  // max 31
	SequenceBits = uint64(12) // max 4095, limit 3M/s

	WorkerIdShift      = SequenceBits
	TagIdShift         = SequenceBits + WorkerIdBits
	TimestampLeftShift = SequenceBits + WorkerIdBits + TagIdBits

	SequenceMask = int64(-1) ^ (int64(-1) << SequenceBits)

	MaxTagId    = (1 << TagIdBits) - 1
	MaxWorkerId = (1 << WorkerIdBits) - 1

	// Sat Jan  4 19:29:34 2014
	twepoch = int64(1388834974657)
)

type IdGenerator struct {
	mutex         sync.Mutex
	cookie        uint32 // random number to mitigate brute force lookups TODO
	wid           int64  // worker id
	seq           int64
	lastTimestamp int64
	randStartSeq  bool
}

func NewIdGenerator(wid int) (this *IdGenerator, err error) {
	this = new(IdGenerator)
	this.wid = int64(wid)
	if wid > MaxWorkerId {
		return nil, ErrorInvalidWorkerId
	}
	return
}

// WithRandStartSeq returns an IdGenerator whose starting sequence at
// the millionsecond switching point be randomized, instead of always
// being zero.
// In this way, the generated id's ending digit will be more evenly
// distributed. Otherwise, most id's ending digit will be 0.
func (this *IdGenerator) WithRandStartSeq(r bool) *IdGenerator {
	this.randStartSeq = r
	return this
}

func (this *IdGenerator) milliseconds() int64 {
	return time.Now().UnixNano() / 1e6
}

func (this *IdGenerator) Next() (int64, error) {
	return this.nextId(0)
}

func (this *IdGenerator) NextWithTag(tag int16) (int64, error) {
	return this.nextId(tag)
}

func (this *IdGenerator) nextId(tag int16) (int64, error) {
	if tag > MaxTagId {
		return 0, ErrorInvalidTag
	}

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
		if this.randStartSeq {
			this.seq = rand.Int63n(10)
		} else {
			this.seq = 0
		}
	}

	this.lastTimestamp = ts

	r := ((ts - twepoch) << TimestampLeftShift) |
		(int64(tag) << TagIdShift) |
		(this.wid << WorkerIdShift) |
		this.seq
	return r, nil
}
