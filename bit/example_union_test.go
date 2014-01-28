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

// Variadic Union function efficiently implemented with SetOr.
func Union(A ...*bit.Set) *bit.Set {
	// Optimization: Allocate empty set with adequate capacity.
	max := 0
	for _, X := range A {
		if e := X.Max(); e > max {
			max = e
		}
	}
	S := bit.New(max).Clear()

	for _, X := range A {
		S.SetOr(S, X)
	}
	return S
}

func ExampleSet_union() {
	A, B, C := bit.New(1, 2), bit.New(2, 3), bit.New(5)
	fmt.Println(Union(A, B, C))
	// Output: {1..3, 5}
}
