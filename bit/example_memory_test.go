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
	"math/rand"
)

func ExampleSet_memory() {
	// Memory management
	S := bit.New(100, 100000) // Make a set that occupies a few kilobyes.
	S.RemoveMax()             // The big element is gone. Memory remains.
	S = bit.New().Set(S)      // Give excess capacity to garbage collector.

	// Memory hints and memory reuse
	IntegerRuns(10, 5)

	// Example output:
	//
	//      Max     Size     Runs
	//       10        9        2    {0..2, 4..9}
	//       20       18        2    {0, 1, 3..11, 13..19}
	//       30       24        4    {0..3, 6..8, 10..13, 15..27}
	//       40       32        5
	//       50       44        5
}

// IntegerRuns(s, n) generates random sets R(i), i = 1, 2, ..., n,
// with elements drawn from 0..i*s-1 and computes the number of runs
// in each set. A run is a sequence of at least three integers
// a, a+1, ..., b, such that {a..b} ⊆ S and {a-1, b+1} ∩ S = ∅.
func IntegerRuns(start, n int) {
	// Give a capacity hint.
	R := bit.New(n*start - 1).Clear()

	fmt.Printf("\n%8s %8s %8s\n", "Max", "Size", "Runs")
	for i := 1; i <= n; i++ {
		// Reuse memory from last iteration.
		R.Clear()

		// Create a random set with ~86% population.
		max := i * start
		for j := 2 * max; j > 0; j-- {
			R.Add(rand.Intn(max))
		}

		// Compute the number of runs.
		runs, length, prev := 0, 0, -2
		R.Do(func(i int) {
			if i == prev+1 { // Continue run.
				length++
			} else { // Start a new run.
				if length >= 3 {
					runs++
				}
				length = 1
			}
			prev = i
		})
		if length >= 3 {
			runs++
		}
		fmt.Printf("%8d %8d %8d", max, R.Size(), runs)
		if max <= 32 {
			fmt.Printf("%4s%v", "", R)
		}
		fmt.Println()
	}
	fmt.Println()
}
