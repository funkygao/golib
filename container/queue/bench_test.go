package queue

import (
	"testing"
)

func BenchmarkPush(b *testing.B) {
	b.ReportAllocs()
	queue := New()
	for i := 0; i < b.N; i++ {
		queue.Push(i)
	}
}

func BenchmarkPop(b *testing.B) {
	b.ReportAllocs()
	queue := New()
	for i := 0; i < b.N; i++ {
		queue.Push(i)
	}
	b.ResetTimer()
	for !queue.Empty() {
		queue.Pop()
	}
}
