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

func TestStringArena(t *testing.T) {
	sa := NewStringArena(10)
	buf1 := []byte("01234")
	buf2 := []byte("5678")
	buf3 := []byte("ab")
	buf4 := []byte("9")

	s1 := sa.NewString(buf1)
	checkint(t, sa.SpaceLeft(), 10-len(buf1))
	checkstring(t, s1, "01234")
	checkint(t, sa.SpaceUsed(), 5)

	s2 := sa.NewString(buf2)
	checkint(t, sa.SpaceLeft(), 1)
	checkstring(t, s2, "5678")
	checkint(t, sa.SpaceUsed(), 9)

	// s3 will be allocated outside of sa, traditional go string
	s3 := sa.NewString(buf3)
	checkint(t, sa.SpaceLeft(), 1)
	checkstring(t, s3, "ab")
	checkint(t, sa.SpaceUsed(), 9)

	// s4 should still fit in sa
	s4 := sa.NewString(buf4)
	checkint(t, sa.SpaceLeft(), 0)
	checkstring(t, s4, "9")
	checkint(t, sa.SpaceUsed(), 10)
}

func checkstring(t *testing.T, actual, expected string) {
	if actual != expected {
		t.Errorf("received %s, expecting %s", actual, expected)
	}
}

func checkint(t *testing.T, actual, expected int) {
	if actual != expected {
		t.Errorf("received %d, expecting %d", actual, expected)
	}
}
