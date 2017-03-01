package ringbuffer

import (
	"testing"

	"github.com/funkygao/assert"
)

func TestRingBufferBasics(t *testing.T) {
	rb, _ := New(16)
	rb.Write("hello")
	rb.Write(189)
	v1 := rb.Read().(string)
	assert.Equal(t, "hello", v1)
	v2 := rb.Read().(int)
	assert.Equal(t, 189, v2)
	if testing.Verbose() {
		t.Logf("%+v", rb)
	}
}

func TestRingBufferAdvanced(t *testing.T) {
	rb, _ := New(128)
	go func() {
		for i := 0; i < 1<<10; i++ {
			rb.Write(i + 1)
		}
	}()

	var last int
	for i := 0; i < 1<<10; i++ {
		r := rb.Read().(int)
		if r-last != 1 {
			t.Fatalf("%d %d", r, last)
		}

		last = r
	}
}

func TestRingBufferRewind(t *testing.T) {
	rb, _ := New(128)
	for i := 0; i < 1<<20; i++ {
		rb.Write(i)
		assert.Equal(t, i, rb.Read().(int))
	}
	if testing.Verbose() {
		t.Logf("%+v", rb)
	}
}

func TestNewWithError(t *testing.T) {
	for _, sz := range []uint64{1, 3, 5, 10, 12, 28, 4098} {
		_, err := New(sz)
		if err != ErrInvalidQueueSize {
			t.Fatal("should return err")
		}
	}

	_, err := New(2 << 12)
	if err != nil {
		t.Fatal("should not return err")
	}
}

func BenchmarkRingBuffer(b *testing.B) {
	rb, _ := New(4096)
	data := "value"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		rb.Write(data)
		rb.Read()
	}
}
