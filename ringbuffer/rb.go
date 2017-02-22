package ringbuffer

import (
	"sync/atomic"
	"time"
)

const backoff = time.Millisecond * 5

type RingBuffer struct {
	queueSize uint64
	indexMask uint64

	padding1           [6]uint64
	lastCommittedIndex uint64 //
	padding2           [8]uint64
	nextFreeIndex      uint64 //
	padding3           [8]uint64
	readerIndex        uint64 //
	padding4           [8]uint64
	highWatermark      uint64 //
	padding5           [8]uint64
	contents           []interface{}
}

// New creates a ring buffer.
func New(queueSize uint64) (*RingBuffer, error) {
	if queueSize == 1 || queueSize&(queueSize-1) != 0 {
		return nil, ErrInvalidQueueSize
	}

	return &RingBuffer{
		queueSize:          queueSize,
		indexMask:          queueSize - 1,
		lastCommittedIndex: 0,
		nextFreeIndex:      1,
		readerIndex:        1,
		contents:           make([]interface{}, queueSize),
	}, nil
}

func (rb *RingBuffer) Write(value interface{}) {
	var myIndex = atomic.AddUint64(&rb.nextFreeIndex, 1) - 1
	// Wait for reader to catch up, so we don't clobber a slot which it is (or will be) reading
	for myIndex > (atomic.LoadUint64(&rb.readerIndex) + rb.queueSize - 2) {
		time.Sleep(backoff)
	}

	// Write the item into it's slot
	rb.contents[myIndex&rb.indexMask] = value
	// Increment the lastCommittedIndex so the item is available for reading
	for !atomic.CompareAndSwapUint64(&rb.lastCommittedIndex, myIndex-1, myIndex) {
		time.Sleep(backoff)
	}
}

func (rb *RingBuffer) Read() interface{} {
	var myIndex = atomic.AddUint64(&rb.readerIndex, 1) - 1
	// If reader has out-run writer, wait for a value to be committed
	for myIndex > atomic.LoadUint64(&rb.lastCommittedIndex) {
		time.Sleep(backoff)
	}

	return rb.contents[myIndex&rb.indexMask]
}

func (rb *RingBuffer) ReadTimeout(timeout time.Duration) (interface{}, bool) {
	var (
		t0      = time.Now()
		myIndex = atomic.AddUint64(&rb.readerIndex, 1) - 1
	)
	for myIndex > atomic.LoadUint64(&rb.lastCommittedIndex) {
		if time.Since(t0) < timeout {
			time.Sleep(backoff)
		} else {
			return nil, false
		}
	}

	return rb.contents[myIndex&rb.indexMask], true
}

// Rewind will rewind contents to last advanced position for reading.
func (rb *RingBuffer) Rewind() {
	atomic.StoreUint64(&rb.readerIndex, atomic.LoadUint64(&rb.highWatermark))
}

func (rb *RingBuffer) Advance() {
	atomic.AddUint64(&rb.highWatermark, 1)
}
