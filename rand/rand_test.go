package rand

import (
	"github.com/funkygao/assert"
	"testing"
)

func TestRandomBytes(t *testing.T) {
	b := RandomByteSlice(20)
	t.Logf("%v", b)
	assert.Equal(t, 20, len(b))
}

func BenchmarkRandomString20(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = RandomString(20)
	}
}

func TestNewPseudoSeed(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Logf("%d\n", NewPseudoSeed())
	}
}

func BenchmarkNewPseudoSeed(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewPseudoSeed()
	}
}
