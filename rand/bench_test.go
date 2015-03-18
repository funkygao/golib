package rand

import (
	"testing"
)

func BenchmarkRandomString20(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = RandomString(20)
	}
}

func BenchmarkSizedString20(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = SizedString(20)
	}
}

func BenchmarkNewPseudoSeed(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewPseudoSeed()
	}
}
