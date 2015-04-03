package conv

import (
	"strconv"
	"testing"
)

var v = "12232"

func BenchmarkParseInt(b *testing.B) {
	b.ReportAllocs()
	s := []byte(v)
	for i := 0; i < b.N; i++ {
		ParseInt(s)
	}
}

func BenchmarkStrconvParseInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.ParseInt(v, 10, 64)
	}
}

func BenchmarkClosestPow2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ClosestPow2(1234)
	}
}
