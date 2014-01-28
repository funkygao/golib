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

package bit_test

import (
	bit "."
	"fmt"
)

func ExampleMaxPos_allocation() {
	// Compute suitable allocation size (using MaxPos and MaxInt)
	// favoring powers of two and guaranteeing linear amortized cost
	// for a repeated number of allocations.
	newSize := func(want, had int) int {
		if n := NextPowerOfTwo(had); n > want {
			return n
		}
		return want
	}

	say := func(h, w, g int) {
		fmt.Printf("Had %#x, wanted %#x, got %#x.\n", h, w, g)
	}

	had := 0xff0
	want := had + 0x008       // shy of next power of two
	got := newSize(want, had) // got == NextPowerOfTwo(had)
	say(had, want, got)

	had = got
	want = had + 0x2000      // overshooting next power of two
	got = newSize(want, had) // got == want
	say(had, want, got)

	had = got
	want = had + 0x1000      // hitting next power of two
	got = newSize(want, had) // got == want
	say(had, want, got)

	// Output:
	// Had 0xff0, wanted 0xff8, got 0x1000.
	// Had 0x1000, wanted 0x3000, got 0x3000.
	// Had 0x3000, wanted 0x4000, got 0x4000.
}

// NextPowerOfTwo returns the smallest p = 1, 2, 4, ..., 1<<k
// such that p > n, or MaxInt if p > MaxInt.
func NextPowerOfTwo(n int) (p int) {
	if n <= 0 {
		return 1
	}

	if k := bit.MaxPos(uint64(n)) + 1; k < bit.BitsPerWord-1 {
		return 1 << uint(k)
	}
	return bit.MaxInt
}
