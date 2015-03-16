package pool_test

import (
	"fmt"

	"github.com/funkygao/golib/pool/queue"
)

// Simple usage example that inserts the numbers 0, 1, 2 into a queue and then
// removes them one by one, printing them to the standard output.
func Example_usage() {
	// Create a queue an push some data in
	q := queue.New()
	for i := 0; i < 3; i++ {
		q.Push(i)
	}
	// Pop out the queue contents and display them
	for !q.Empty() {
		fmt.Println(q.Pop())
	}
	// Output:
	// 0
	// 1
	// 2
}
