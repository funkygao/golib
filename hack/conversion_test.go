package hack

import (
	"github.com/funkygao/assert"
	"testing"
)

func TestString(t *testing.T) {
	b := []byte{'h', 'e', 'l', 'l', 'o'}
	assert.Equal(t, "hello", String(b))
}

func TestByte(t *testing.T) {
	s := "hello"
	assert.Equal(t, []byte{'h', 'e', 'l', 'l', 'o'}, Byte(s))
}
