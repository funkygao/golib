package idgen

import (
	"github.com/funkygao/assert"
	"testing"
)

func TestNextId(t *testing.T) {
	id := NewIdGenerator()
	for i := 0; i < 100; i++ {
		n, err := id.Next()
		assert.Equal(t, nil, err)
		n1, _ := id.Next()
		if n1 <= n {
			t.Errorf("expected %d > %d", n1, n)
		}
		t.Logf("n = %v", n)
		t.Logf("n1= %v", n1)
	}

}

func BenchmarkIdNext(b *testing.B) {
	id := NewIdGenerator()
	for i := 0; i < b.N; i++ {
		id.Next()
	}
}
