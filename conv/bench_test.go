package conv

import (
	"testing"
)

func BenchmarkParseInt(b *testing.B) {
	s := []byte("12232")
	for i := 0; i < b.N; i++ {
		ParseInt(s)
	}
}
