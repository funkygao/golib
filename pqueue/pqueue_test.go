package pq

import (
	"container/heap"
	"github.com/funkygao/assert"
	"testing"
)

func TestPriorityQueue(t *testing.T) {
	pq := New()
	heap.Init(pq)
	heap.Push(pq, &Item{Value: "hello3", Priority: 3})
	heap.Push(pq, &Item{Value: "hello1", Priority: 1})
	heap.Push(pq, &Item{Value: "hello8", Priority: 8})

	item := heap.Pop(pq)
	assert.Equal(t, "hello8", item.(*Item).Value.(string))

	item = heap.Pop(pq)
	assert.Equal(t, "hello3", item.(*Item).Value.(string))

	item = heap.Pop(pq)
	assert.Equal(t, "hello1", item.(*Item).Value.(string))

	assert.Equal(t, 0, pq.Len())
}
