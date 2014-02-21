package recycler

import (
	"math/rand"
	"testing"
)

func TestRecycler(t *testing.T) {
	get, give := New(0, func() interface{} {
		return make([]byte, rand.Intn(5000000)+5000000)
	})

	buf := <-get
	t.Logf("len=%v get=%+v give=%+v", len(buf.([]byte)), get, give)
}

func BenchmarkRecycler(b *testing.B) {
	b.ReportAllocs()
	get, give := New(20, func() interface{} {
		return make([]byte, rand.Intn(5000000)+5000000)
	})
	pool := make([][]byte, 20)

	for i := 0; i < b.N; i++ {
		buf := <-get
		i := rand.Intn(len(pool))
		if pool[i] != nil {
			give <- pool[i]
		}

		pool[i] = buf.([]byte)
	}

	b.Logf("makes: %v", makes)
}
