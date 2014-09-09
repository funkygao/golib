package pqueue

import (
	"container/heap"
)

type Item struct {
	Value    interface{}
	Priority int

	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func New() *PriorityQueue {
	return new(PriorityQueue)
}

func (this *PriorityQueue) PrioritySum() int {
	sum := 0
	for _, item := range *this {
		sum += item.Priority
	}

	return sum
}

func (this PriorityQueue) Len() int {
	return len(this)
}

func (this PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return this[i].Priority > this[j].Priority
}

func (this PriorityQueue) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
	this[i].index = i
	this[j].index = j
}

func (this *PriorityQueue) Push(x interface{}) {
	n := len(*this)
	item := x.(*Item)
	item.index = n
	*this = append(*this, item)
}

func (this *PriorityQueue) Pop() interface{} {
	old := *this
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*this = old[0 : n-1]
	return item
}

func (this *PriorityQueue) Peek() interface{} {
    // the highest priority item is at index zero
    return (*this)[0]
}

// update modifies the priority and value of an Item in the queue.
func (this *PriorityQueue) lupdate(item *Item, value interface{}, priority int) {
	heap.Remove(this, item.index)
	item.Value = value
	item.Priority = priority
	heap.Push(this, item)
}
