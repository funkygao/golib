package hack

import (
	"testing"
)

func BenchmarkStringWithoutHack(b *testing.B) {
	b.ReportAllocs()
	ba := []byte{'h', 'e', 'l', 'l', 'o'}
	for i := 0; i < b.N; i++ {
		_ = string(ba)
	}
}

func BenchmarkStringWithHack(b *testing.B) {
	b.ReportAllocs()
	ba := []byte{'h', 'e', 'l', 'l', 'o'}
	for i := 0; i < b.N; i++ {
		_ = String(ba)
	}
}

func BenchmarkStringWithStringArena(b *testing.B) {
	b.ReportAllocs()
	ba := []byte{'h', 'e', 'l', 'l', 'o'}
	sa := NewStringArena(len(ba))
	for i := 0; i < b.N; i++ {
		sa.NewString(ba)
	}
}

func BenchmarkAppendVariant(b *testing.B) {
	var ba [100]byte
	b.ReportAllocs()
	y := []byte{'h', 'e', 'l', 'l', 'o'}
	for i := 0; i < b.N; i++ {
		x := ba[:0]
		x = append(x, y...)
		_ = x
	}
	b.SetBytes(int64(len(y)))
}

func BenchmarkByteWithoutHack(b *testing.B) {
	b.ReportAllocs()
	s := "hello"
	for i := 0; i < b.N; i++ {
		_ = []byte(s)
	}
}

func BenchmarkByteWithHack(b *testing.B) {
	b.ReportAllocs()
	s := "hello"
	for i := 0; i < b.N; i++ {
		_ = Byte(s)
	}
}
