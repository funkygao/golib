package sampling

import (
	"testing"
)

func BenchmarkSampling(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		SampleRateSatisfied(3)
	}
}
