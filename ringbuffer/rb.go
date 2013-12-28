package ringbuffer

import "fmt"
import "sync/atomic"
import "runtime"

const (
	queueSize uint64 = 4096
	indexMask uint64 = queueSize - 1
)

type ringBuffer struct {
	padding1           [8]uint64
	lastCommittedIndex uint64
	padding2           [8]uint64
	nextFreeIndex      uint64
	padding3           [8]uint64
	readerIndex        uint64
	padding4           [8]uint64
	contents           [queueSize]interface{}
	padding5           [8]uint64
}

func New() *ringBuffer {
	return &ringBuffer{lastCommittedIndex: 0, nextFreeIndex: 1, readerIndex: 1}
}

func (this *ringBuffer) Write(value interface{}) {
	var myIndex = atomic.AddUint64(&this.nextFreeIndex, 1) - 1
	// Wait for reader to catch up, so we don't clobber a slot which it is (or will be) reading
	for myIndex > (this.readerIndex + queueSize - 2) {
		runtime.Gosched()
	}

	// Write the item into it's slot
	this.contents[myIndex&indexMask] = value
	// Increment the lastCommittedIndex so the item is available for reading
	for !atomic.CompareAndSwapUint64(&this.lastCommittedIndex, myIndex-1, myIndex) {
		runtime.Gosched()
	}
}

func (this *ringBuffer) Read() interface{} {
	var myIndex = atomic.AddUint64(&this.readerIndex, 1) - 1
	// If reader has out-run writer, wait for a value to be committed
	for myIndex > this.lastCommittedIndex {
		runtime.Gosched()
	}

	return this.contents[myIndex&indexMask]
}

func (this *ringBuffer) PrintDump() {
	fmt.Printf("lastCommitted: %3d, nextFree: %3d, readerIndex: %3d, content:",
		this.lastCommittedIndex, this.nextFreeIndex, this.readerIndex)
	for idx, value := range this.contents {
		fmt.Printf("%5v : %5v", idx, value)
	}
	fmt.Print("\n")
}
