package endian

import (
	"github.com/funkygao/assert"
	"testing"
)

func TestEndian(t *testing.T) {
	if EndianIsLittle() {
		t.Log("current os is little endian")
	} else {
		t.Log("current os is big endian")
	}

	var val uint16 = 6 << 8
	low, high := SafeSplitUint16(val)
	assert.Equal(t, uint8(0), low)
	assert.Equal(t, uint8(0x6), high)
}
