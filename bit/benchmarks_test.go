// Copyright 2012 Stefan Nilsson
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bit

import "testing"

func BenchmarkMinPos(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MinPos(0xcafecafecafecafe)
	}
}

func BenchmarkMaxPos(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MaxPos(0xcafecafecafecafe)
	}
}

func BenchmarkCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Count(0xcafecafecafecafe)
	}
}

// Number of words in test set.
const nw = 1 << 10

func BenchmarkSize(b *testing.B) {
	b.StopTimer()
	S := BuildTestSet(nw << 3) // Allocates nw<<3 bytes = nw words.
	b.StartTimer()

	for i := 0; i < b.N/nw; i++ { // Measure time per word.
		S.Size()
	}
}

func BenchmarkSetAnd(b *testing.B) {
	b.StopTimer()
	S := New().SetWord(nw-1, 1).Clear() // Allocates nw words.
	A := BuildTestSet(nw << 3)
	B := BuildTestSet(nw << 3)
	b.StartTimer()

	for i := 0; i < b.N/nw; i++ { // Measure time per word.
		S.SetAnd(A, B)
	}
}

func BenchmarkDo(b *testing.B) {
	b.StopTimer()
	S := BuildTestSet(b.N) // As Do is pretty fast, S can be pretty big.
	b.StartTimer()

	S.Do(func(i int) {})
}

func BenchmarkString(b *testing.B) {
	b.StopTimer()
	S := BuildTestSet(b.N) // As Do is pretty fast, S can be pretty big.
	b.StartTimer()

	S.String()
}

func BenchmarkNext(b *testing.B) {
	b.StopTimer()
	S := BuildTestSet(b.N)
	b.StartTimer()

	for n, found := S.Next(-1); found; {
		n, found = S.Next(n)
	}
}

func BenchmarkRemoveMin(b *testing.B) {
	b.StopTimer()
	S := BuildTestSet(b.N)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		S.RemoveMin()
	}
}

func BenchmarkRemoveMax(b *testing.B) {
	b.StopTimer()
	S := BuildTestSet(b.N)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		S.RemoveMax()
	}
}

// Quickly builds a set of n somewhat random elements from 0..8n-1.
func BuildTestSet(n int) *Set {
	S := New().Add(8*n - 1).Clear() // Allocates n bytes.

	lfsr := uint16(0xace1) // linear feedback shift register
	for i := 0; i < n; i++ {
		bit := (lfsr>>0 ^ lfsr>>2 ^ lfsr>>3 ^ lfsr>>5) & 1
		lfsr = lfsr>>1 | bit<<15
		e := i<<3 + int(lfsr&0x7)
		S.Add(e) // Add a number from 8i..8i+7.
	}
	return S
}
