package mutexmap

import (
	"testing"
	"time"
)

func TestLockAndUnlock(t *testing.T) {
	const KEY = "10"
	m := New(800 << 10)
	m.Lock(KEY)
	val, _ := m.items.Get(KEY)
	t.Logf("items: %+v", *m.items)
	t.Logf("value of key[%v]: %+v", KEY, val)

	for i := 0; i < 100; i++ {
		go func() {
			m.Lock(KEY)
			m.Unlock(KEY)
		}()
	}

	time.Sleep(time.Second * 3)
	m.Unlock(KEY)
}

func BenchmarkLockAndUnlock(b *testing.B) {
	const KEY = "abc"

	b.ReportAllocs()
	m := New(800 << 10)
	for i := 0; i < b.N; i++ {
		m.Lock(KEY)
		m.Unlock(KEY)
	}
}
