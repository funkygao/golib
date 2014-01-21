package sortedmap

import (
	"github.com/funkygao/assert"
	"testing"
)

func TestSortedKeys(t *testing.T) {
	sm := NewSortedMap()
	sm.Set("a", 10)
	sm.Set("c", 43)
	sm.Set("d", 21)
	t.Logf("%v", sm.SortedKeys())
}

func TestBasic(t *testing.T) {
	sm := NewSortedMap()
	sm.Set("foo", 34)
	sm.Inc("foo", 3)
	assert.Equal(t, 37, sm.Get("foo"))
	sm.Inc("foo", -2)
	assert.Equal(t, 35, sm.Get("foo"))
	assert.Equal(t, -1, sm.Inc("non-exist", 5))
}
