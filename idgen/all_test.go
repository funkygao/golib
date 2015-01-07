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

func BenchmarkIdNext(b *testing.B) {
	id, _ := NewIdGenerator(3)
	for i := 0; i < b.N; i++ {
		id.Next()
	}
}
