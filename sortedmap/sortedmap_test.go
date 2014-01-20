package sortmap

import (
	"testing"
)

func TestSortedMap(t *testing.T) {
	sm := NewSortedMap()
	sm.Set("a", 10)
	sm.Set("c", 43)
	sm.Set("d", 21)
	t.Logf("%v", sm.SortedKeys())
}
