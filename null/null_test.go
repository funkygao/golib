package null

import (
	"github.com/funkygao/assert"
	"testing"
	"unsafe"
)

func TestNullIsSmallerThanBool(t *testing.T) {
	s1 := unsafe.Sizeof(Null)
	s2 := unsafe.Sizeof(true)
	assert.Equal(t, true, s2 > s1)
	assert.Equal(t, uintptr(0), s1)
}
