// Package queue implements a FIFO (first in first out) data structure supporting
// arbitrary types (even a mixture).
//
// Internally it uses a dynamically growing circular slice of blocks, resulting
// in faster resizes than a simple dynamic array/slice would allow.
package queue

// The size of a block of data
const blockSize = 4096

// FIFO data structure not goroutine safe.
type Queue struct {
	tailIndex  int
	headIndex  int
	tailOffset int
	headOffset int

	blocks [][]interface{}
	head   []interface{}
	tail   []interface{}
}

// Creates a new, empty queue.
func New() *Queue {
	result := new(Queue)
	result.blocks = [][]interface{}{make([]interface{}, blockSize)}
	result.head = result.blocks[0]
	result.tail = result.blocks[0]
	return result
}

// Pushes a new element into the queue, expanding it if necessary.
func (q *Queue) Push(data interface{}) {
	q.tail[q.tailOffset] = data
	q.tailOffset++
	if q.tailOffset == blockSize {
		q.tailOffset = 0
		q.tailIndex = (q.tailIndex + 1) % len(q.blocks)

		// If we wrapped over to the end, insert a new block and update indices
		if q.tailIndex == q.headIndex {
			buffer := make([][]interface{}, len(q.blocks)+1)
			copy(buffer[:q.tailIndex], q.blocks[:q.tailIndex])
			buffer[q.tailIndex] = make([]interface{}, blockSize)
			copy(buffer[q.tailIndex+1:], q.blocks[q.tailIndex:])
			q.blocks = buffer
			q.headIndex++
			q.head = q.blocks[q.headIndex]
		}
		q.tail = q.blocks[q.tailIndex]
	}
}

// Pops out an element from the queue. Note, no bounds checking are done.
func (q *Queue) Pop() (res interface{}) {
	res, q.head[q.headOffset] = q.head[q.headOffset], nil
	q.headOffset++
	if q.headOffset == blockSize {
		q.headOffset = 0
		q.headIndex = (q.headIndex + 1) % len(q.blocks)
		q.head = q.blocks[q.headIndex]
	}
	return
}

// Returns the first element in the queue. Note, no bounds checking are done.
func (q *Queue) Front() interface{} {
	return q.head[q.headOffset]
}

// Checks whether the queue is empty.
func (q *Queue) Empty() bool {
	return q.headIndex == q.tailIndex && q.headOffset == q.tailOffset
}

// Returns the number of elements in the queue.
func (q *Queue) Size() int {
	if q.tailIndex > q.headIndex {
		return (q.tailIndex-q.headIndex)*blockSize - q.headOffset + q.tailOffset
	} else if q.tailIndex < q.headIndex {
		return (len(q.blocks)-q.headIndex+q.tailIndex)*blockSize - q.headOffset + q.tailOffset
	} else {
		return q.tailOffset - q.headOffset
	}
}

// Clears out the contents of the queue.
func (q *Queue) Reset() {
	// Rewind the queue indices
	q.headIndex = 0
	q.tailIndex = 0
	q.headOffset = 0
	q.tailOffset = 0

	// Reset the active blocks
	q.head = q.blocks[0]
	q.tail = q.blocks[0]

	// Set all elements to nil to allow garbage collection
	for _, block := range q.blocks {
		for i := 0; i < len(block); i++ {
			block[i] = nil
		}
	}
}
