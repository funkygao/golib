package idgen

import (
	"github.com/funkygao/assert"
	"testing"
)

func TestNextId(t *testing.T) {
	id, _ := NewIdGenerator(1)
	for i := 0; i < 100; i++ {
		n, err := id.Next()
		assert.Equal(t, nil, err)
		n1, _ := id.Next()
		if n1 <= n {
			t.Errorf("expected %d > %d", n1, n)
		}
		t.Logf("n=%d, n1=%d", n, n1)
	}

}

func TestNextIdWithTag(t *testing.T) {
	id, _ := NewIdGenerator(1)
	for i := 0; i < 100; i++ {
		n, err := id.NextWithTag(5)
		assert.Equal(t, nil, err)
		n1, _ := id.NextWithTag(5)
		if n1 <= n {
			t.Errorf("expected %d > %d", n1, n)
		}
		t.Logf("n=%d, n1=%d", n, n1)
	}

}

func TestNextIdWithTagError(t *testing.T) {
	_, err := NewIdGenerator(32)
	assert.Equal(t, ErrorInvalidWorkerId, err)
	id, err := NewIdGenerator(31)
	assert.Equal(t, nil, err)
	val, err := id.NextWithTag(32)
	assert.Equal(t, int64(0), val)
	assert.Equal(t, ErrorInvalidTag, err)
}

func BenchmarkIdNext(b *testing.B) {
	id, _ := NewIdGenerator(3)
	for i := 0; i < b.N; i++ {
		id.Next()
	}
}
