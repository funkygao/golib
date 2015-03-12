package ringbuffer

import (
	"github.com/funkygao/assert"
	"testing"
)

func TestRingBuffer(t *testing.T) {
	rb := New()
	rb.Write("hello")
	rb.Write(189)
	v1 := rb.Read().(string)
	assert.Equal(t, "hello", v1)
	v2 := rb.Read().(int)
	assert.Equal(t, 189, v2)
}
