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

func ExampleSet_Do() {
	A := bit.New(1, 2, 3, 4)
	sum := 0
	A.Do(func(n int) {
		sum += n
	})
	fmt.Printf("sum %v = %d\n", A, sum)
	// Output: sum {1..4} = 10
}

func ExampleSet_Next() {
	A := bit.New(2, 3, 5, 7, 11, 13)

	// Print all single digit numbers in A.
	for n, found := A.Next(-1); found && n < 10; n, found = A.Next(n) {
		fmt.Printf("%d ", n)
	}
	// Output: 2 3 5 7
}

func ExampleSet_Previous() {
	A := bit.New(2, 3, 5, 7)

	// Print all numbers in A in reverse order.
	if !A.IsEmpty() {
		for n, found := A.Max(), true; found; n, found = A.Previous(n) {
			fmt.Printf("%d ", n)
		}
	}
	// Output: 7 5 3 2
}

func ExampleSet_operators() {
	A := new(bit.Set).AddRange(0, 100)     // A = {0..99}
	B := bit.New(0, 200).AddRange(50, 150) // B = {0, 50..149, 200}
	S := A.Xor(B)                          // S = A ∆ B
	C := A.Or(B).AndNot(A.And(B))          // C = (A ∪ B) ∖ (A ∩ B)
	D := A.AndNot(B).Or(B.AndNot(A))       // D = (A ∖ B) ∪ (B ∖ A)

	if C.Equals(S) && D.Equals(S) {
		fmt.Printf("A ∆ B = %v\n", S)
	}
	// Output: A ∆ B = {1..49, 100..149, 200}
}

func ExampleSet_words() {
	const faraway = 46                         // billion light years
	Universe := bit.New().AddRange(0, faraway) // Universe = {0..faraway-1}
	Even := bit.New().SetWord(0, 1<<faraway/3) // Even = {0, 2, 4, ..., faraway-2}

	Odd := Universe.AndNot(Even)
	fmt.Printf("Odd = %v\n", Odd)

	Even.FlipRange(0, faraway)
	fmt.Printf("Even = %v\n", Even)
	// Output:
	// Odd = {1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23, 25, 27, 29, 31, 33, 35, 37, 39, 41, 43, 45}
	// Even = {1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23, 25, 27, 29, 31, 33, 35, 37, 39, 41, 43, 45}
}

// Create the set of all primes ≤ max using Sieve of Eratosthenes.
func ExampleSet_eratosthenes() {
	const max = 50
	sieve := bit.New().AddRange(2, max)
	for p, found := 2, true; found; p, found = sieve.Next(p) {
		for n := 2 * p; n <= max; n += p {
			sieve.Remove(n)
		}
	}
	fmt.Println(sieve)
	// Output: {2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47}
}
