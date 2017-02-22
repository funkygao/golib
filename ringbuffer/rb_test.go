package ringbuffer

import (
	"testing"
	"time"

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

func TestReadTimeout(t *testing.T) {
	rb, _ := New(16)
	t0 := time.Now()
	r, ok := rb.ReadTimeout(time.Second)
	assert.Equal(t, nil, r)
	assert.Equal(t, false, ok)
	assert.Equal(t, true, time.Since(t0) > time.Second)
}

func TestRinbBufferAdvanceAndRewind(t *testing.T) {
	rb, _ := New(8)
	go func() {
		for i := 0; i < 30; i++ {
			rb.Write(i + 1)

			if i < 3 {
				rb.Advance()
			}
		}
	}()

	receivedInts := make(map[int]struct{})
	rewinded := false
	for {
		r, ok := rb.ReadTimeout(time.Second)
		if !ok {
			break
		}

		v := r.(int)
		receivedInts[v] = struct{}{}
		t.Logf("<- %d", v)

		if !rewinded && v == 5 {
			rb.Rewind()
			rewinded = true
		}
	}

	assert.Equal(t, 30, len(receivedInts))
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
