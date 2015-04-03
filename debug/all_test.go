package debug

import (
	"github.com/funkygao/assert"
	"testing"
)

func add(a, b int) int {
	return a + b
}

func TestTimeit(t *testing.T) {
	r, _, err := Timeit(add, 1, 2)
	assert.Equal(t, err, nil)
	assert.Equal(t, r[0].Int(), int64(3))
}

func TestTrace(t *testing.T) {
	EnableTrace()
	defer Un(Trace("test trace"))
	callA()
}

func callA() {
	defer Un(Trace("")) // if fn empty, auto get fn name
}
