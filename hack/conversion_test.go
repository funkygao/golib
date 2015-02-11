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

func BenchmarkString(b *testing.B) {
	b.ReportAllocs()
	ba := []byte{'h', 'e', 'l', 'l', 'o'}
	for i := 0; i < b.N; i++ {
		_ = String(ba)
	}
}

func BenchmarkString1(b *testing.B) {
	b.ReportAllocs()
	ba := []byte{'h', 'e', 'l', 'l', 'o'}
	for i := 0; i < b.N; i++ {
		_ = string(ba)
	}
}

func BenchmarkByte(b *testing.B) {
	b.ReportAllocs()
	s := "hello"
	for i := 0; i < b.N; i++ {
		_ = Byte(s)
	}
}

func BenchmarkByte1(b *testing.B) {
	b.ReportAllocs()
	s := "hello"
	for i := 0; i < b.N; i++ {
		_ = []byte(s)
	}
}
