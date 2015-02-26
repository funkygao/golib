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
	sa := NewStringArena(len(ba) + 1)
	for i := 0; i < b.N; i++ {
		sa.NewString(ba)
	}
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
