package trace

import (
    "github.com/bmizerany/assert"
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
