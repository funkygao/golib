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

import (
	"strconv"
	"testing"
)

func TestNew(t *testing.T) {
	for _, S := range []*Set{
		New(),
		New(1),
		New(1, 1),
		New(65),
		New(1, 2, 3),
		New(100, 200, 300),
	} {
		CheckInvariants(t, "New", S)
	}
}

// Checks that the invariants for S.data and S.min hold.
func CheckInvariants(t *testing.T, msg string, S *Set) {
	len := len(S.data)
	cap := cap(S.data)
	data := S.data[:cap]
	min := S.min
	s := "Invariant for "

	if len > 0 && data[len-1] == 0 {
		t.Errorf("%s%s: data = %v, data[%d] = 0; want non-zero", s, msg, data, len-1)
	}

	for i := len; i < cap; i++ {
		if data[i] != 0 {
			t.Errorf("%s%s: data = %v, data[%d] = %#x; want 0", s, msg, data, i, data[i])
			break
		}
	}

	minExp := 0
	if len == 0 && min != minExp {
		t.Errorf("%s%s: S = %v, S.min = %d; want %d", s, msg, S, min, minExp)
	}
	if len > 0 {
		minExp = findMinFrom(0, data)
		if min != minExp {
			t.Errorf("%s%s: S = %v, S.min = %d; want %d", s, msg, S, min, minExp)
		}
	}
}

func TestCmp(t *testing.T) {
	Zero, One := New(), New(1)
	for _, x := range []struct {
		A, B                         *Set
		equals, intersects, subsetOf bool
	}{
		{Zero, Zero, true, false, true},
		{One, One, true, true, true},
		{New(), New(), true, false, true},
		{New(1), New(1), true, true, true},
		{New(64), New(64), true, true, true},
		{New(65), New(65), true, true, true},
		{New(1, 2, 3), New(1, 2, 3), true, true, true},
		{New(100, 200, 300), New(100, 200, 300), true, true, true},

		{New(), New(1), false, false, true},
		{New(1), New(), false, false, false},
		{New(1), New(2), false, false, false},
		{New(), New(65), false, false, true},
		{New(65), New(), false, false, false},
		{New(1), New(65), false, false, false},
		{New(1, 2, 3), New(100, 200, 300), false, false, false},

		{New(1), New(1, 2, 3), false, true, true},
		{New(1, 2, 3), New(1), false, true, false},
		{New(100), New(100, 200, 300), false, true, true},
		{New(100, 200, 300), New(100), false, true, false},
	} {
		A, B := x.A, x.B

		equals := A.Equals(B)
		if equals != x.equals {
			t.Errorf("%v.Equals(%v) = %t; want %t", A, B, equals, x.equals)
		}

		intersects := A.Intersects(B)
		if intersects != x.intersects {
			t.Errorf("%v.Intersects(%v) = %t; want %t", A, B, intersects, x.intersects)
		}

		subsetOf := A.SubsetOf(B)
		if subsetOf != x.subsetOf {
			t.Errorf("%v.SubsetOf(%v) = %t; want %t", A, B, subsetOf, x.subsetOf)
		}
	}
}

func TestMinMax(t *testing.T) {
	for _, x := range []struct {
		S        *Set
		min, max int
	}{
		{New(0), 0, 0},
		{New(65), 65, 65},
		{New(1, 2, 3), 1, 3},
		{New(100, 200, 300), 100, 300},
	} {
		S := x.S

		min, max := S.Min(), S.Max()
		if min != x.min {
			t.Errorf("%v.Min() = %d; want %d", S, min, x.min)
		}

		if max != x.max {
			t.Errorf("%v.Max() = %d; want %d", S, max, x.max)
		}
	}

	N := New()
	if !Panics((*Set).Min, N) {
		t.Errorf("S.Min() should panic for S = ∅.")
	}
	CheckInvariants(t, "S.Min() panic", N)

	N = New()
	if !Panics((*Set).Max, N) {
		t.Errorf("S.Max() should panic for S = ∅.")
	}
	CheckInvariants(t, "S.Max() panic", N)
}

func TestNextPrev(t *testing.T) {
	for _, x := range []struct {
		S         *Set
		m         int
		nextN     int
		nextFound bool
		prevN     int
		prevFound bool
	}{
		{New(), 0, 0, false, 0, false},
		{New(), -1, 0, false, 0, false},

		{New(1), 0, 1, true, 0, false},
		{New(1), 1, 0, false, 0, false},
		{New(1), 2, 0, false, 1, true},

		{New(0, 2), -1, 0, true, 0, false},
		{New(0, 2), 0, 2, true, 0, false},
		{New(0, 2), 1, 2, true, 0, true},
		{New(0, 2), 2, 0, false, 0, true},
		{New(0, 2), 3, 0, false, 2, true},

		{New(63, 64), 63, 64, true, 0, false},
		{New(64, 63), 64, 0, false, 63, true},

		{New(100, 300), MinInt, 100, true, 0, false},
		{New(100, 300), 0, 100, true, 0, false},
		{New(100, 300), 100, 300, true, 0, false},
		{New(100, 300), 200, 300, true, 100, true},
		{New(100, 300), 300, 0, false, 100, true},
		{New(100, 300), 400, 0, false, 300, true},
		{New(100, 300), MaxInt, 0, false, 300, true},
	} {
		S := x.S
		m := x.m

		n, found := S.Next(m)
		if n != x.nextN || found != x.nextFound {
			t.Errorf("%v.Next(%d) = (%d, %v); want (%d, %v)", S, m, n, found, x.nextN, x.nextFound)
		}

		n, found = S.Previous(m)
		if n != x.prevN || found != x.prevFound {
			t.Errorf("%v.Previous(%d) = (%d, %v); want (%d, %v)", S, m, n, found, x.prevN, x.prevFound)
		}

	}
}

func TestRemoveMinMax(t *testing.T) {
	for _, x := range []struct {
		S        *Set
		min, max int
	}{
		{New(0), 0, 0},
		{New(65), 65, 65},
		{New(1, 2, 3), 1, 3},
		{New(100, 200, 300), 100, 300},
	} {
		S := new(Set)

		min := S.Set(x.S).RemoveMin()
		if min != x.min {
			t.Errorf("%v.RemoveMin() = %d; want %d", S, min, x.min)
		}
		CheckInvariants(t, "RemoveMin", S)

		max := S.Set(x.S).RemoveMax()
		if max != x.max {
			t.Errorf("%v.RemoveMax() = %d; want %d", S, max, x.max)
		}
		CheckInvariants(t, "RemoveMax", S)

		S.Set(x.S)
		T := new(Set)
		for !S.IsEmpty() {
			T.Add(S.RemoveMin())
		}
		if !T.Equals(x.S) {
			t.Errorf("RemoveMin, T = %v; want %v", T, x.S)
		}
		CheckInvariants(t, "Iterated RemoveMin", S)

		S.Set(x.S)
		T.Clear()
		for !S.IsEmpty() {
			T.Add(S.RemoveMax())
		}
		if !T.Equals(x.S) {
			t.Errorf("RemoveMax, T = %v; want %v", T, x.S)
		}
		CheckInvariants(t, "Iterated RemoveMax", S)
	}

	N := New()
	if !Panics((*Set).RemoveMin, N) {
		t.Errorf("S.RemoveMin() should panic for S = ∅.")
	}
	CheckInvariants(t, "S.RemoveMin() panic", N)

	N = New()
	if !Panics((*Set).RemoveMax, N) {
		t.Errorf("S.RemoveMax() should panic for S = ∅.")
	}
	CheckInvariants(t, "S.RemoveMax() panic", N)
}

func TestContains(t *testing.T) {
	for _, x := range []struct {
		A        *Set
		n        int
		contains bool
	}{
		{New(), 1, false},
		{New(), 100, false},
		{New(1), 0, false},
		{New(1), 1, true},
		{New(1), 100, false},
		{New(65), 0, false},
		{New(65), 1, false},
		{New(65), 65, true},
		{New(65), 100, false},

		{New(1, 2, 3), 0, false},
		{New(1, 2, 3), 1, true},
		{New(1, 2, 3), 2, true},
		{New(1, 2, 3), 3, true},
		{New(1, 2, 3), 4, false},

		{New(100, 200, 300), 0, false},
		{New(100, 200, 300), 100, true},
		{New(100, 200, 300), 200, true},
		{New(100, 200, 300), 300, true},
		{New(100, 200, 300), 400, false},
	} {
		A, n := x.A, x.n

		contains := A.Contains(n)
		if contains != x.contains {
			t.Errorf("%v.Contains(%d) = %t; want %t", A, n, contains, x.contains)
		}
	}

	N := New()
	if !Panics((*Set).Contains, N, -1) {
		t.Errorf("S.Contains(-1) should panic.")
	}
	CheckInvariants(t, "S.Contains(-1) panic", N)
}

func TestSize(t *testing.T) {
	for _, x := range []struct {
		A    *Set
		size int
	}{
		{New(), 0},
		{New(1), 1},
		{New(64), 1},
		{New(65), 1},
		{New(1, 2, 3), 3},
		{New(100, 200, 300), 3},
		{New().AddRange(0, 64), 64},
		{New().AddRange(1, 64), 63},
		{New().AddRange(0, 63), 63},
	} {
		A := x.A

		size := A.Size()
		if size != x.size {
			t.Errorf("%v.Size() = %d; want %d", A, size, x.size)
		}
	}
}

func TestSet(t *testing.T) {
	for _, x := range []struct {
		S, A *Set
	}{
		{New(), New()},
		{New(), New(1)},
		{New(), New(65)},
		{New(), New(1, 2, 3)},
		{New(), New(100, 200, 300)},

		{New(1, 2, 3), New()},
		{New(1, 2, 3), New(1)},
		{New(1, 2, 3), New(65)},
		{New(1, 2, 3), New(1, 2, 3)},
		{New(1, 2, 3), New(100, 200, 300)},

		{New(100, 200, 300), New()},
		{New(100, 200, 300), New(1)},
		{New(100, 200, 300), New(65)},
		{New(100, 200, 300), New(1, 2, 3)},
		{New(100, 300, 300), New(100, 200, 300)},
	} {
		S := x.S

		T := S.Set(x.A)
		if T != S {
			t.Errorf("&(S.Set(%v)) = %p, &S = %p; want same", x.A, T, S)
		}
		if !T.Equals(x.A) {
			t.Errorf("S.Set(%v) = %v; want %v", x.A, T, x.A)
		}
		CheckInvariants(t, "Set", T)
	}
}

func TestSetWord(t *testing.T) {
	for _, x := range []struct {
		S *Set
		i int
		w uint64
	}{
		{New(), 0, 0x0},
		{New(), 0, 0x55},
		{New(), 2, 0x55},
		{New(1), 0, 0x0},
		{New(1), 0, 0x55},
		{New(1), 2, 0x55},
		{New(1, 2, 3), 0, 0x0},
		{New(1, 2, 3), 0, 0x55},
		{New(1, 2, 3), 2, 0x55},
		{New(100), 1, 0x00},
		{New(100, 200, 300), 0, 0x0},
		{New(100, 200, 300), 1, 0x0},
		{New(100, 200, 300), 2, 0x0},
		{New(100, 200, 300), 5, 0x0},
		{New(100, 200, 300), 0, 0x55},
		{New(100, 200, 300), 1, 0x55},
		{New(100, 200, 300), 2, 0x55},
		{New(100, 200, 300), 5, 0x55},
	} {
		S, T := new(Set), new(Set)
		S.Set(x.S)
		T.Set(x.S)
		i, w := x.i, x.w

		S.SetWord(i, w)
		T.RemoveRange(64*i, 64*(i+1))
		for n := 64 * i; n < 64*(i+1); n++ {
			if w&1 != 0 {
				T.Add(n)
			}
			w >>= 1
		}
		if !S.Equals(T) {
			t.Errorf("%v.SetWord(%d, %#x) = %v; want %v\n", x.S, x.i, x.w, S, T)
		}
		CheckInvariants(t, "SetWord", S)
	}

	for _, i := range []int{-1, MaxInt/64 + 1} {
		N := New()
		if !Panics((*Set).SetWord, N, i, uint64(0)) {
			t.Errorf("S.SetWord(%#x, %#x) should panic.", i, 0)
		}
		CheckInvariants(t, "S.SetWord(x, 0) panic", N)
	}
}

func TestWord(t *testing.T) {
	for _, x := range []struct {
		S *Set
		i int
		w uint64
	}{
		{New(), 0, 0x0},
		{New(), 10, 0x0},
		{New(0), 0, 0x1},
		{New(0), 1, 0x0},
		{New(1, 2, 3), 0, 0x0e},
		{New(1, 2, 3), 10, 0x0},
		{New(64), 0, 0x0},
		{New(64), 1, 0x1},
		{New(64), 2, 0x0},
	} {
		w := x.S.Word(x.i)
		if w != x.w {
			t.Errorf("%v.Word(%d) = %#x; want %#x", x.S, x.i, w, x.w)
		}
	}

	N := New()
	if !Panics((*Set).Word, N, -1) {
		t.Errorf("S.Word(-1) should panic.")
	}
	CheckInvariants(t, "S.Word(-1) panic", N)
}

type rangeFunc struct {
	fInt   func(S *Set, n int) *Set
	fRange func(S *Set, m, n int) *Set
	name   string
}

func TestRange(t *testing.T) {
	rangeFuncs := []rangeFunc{
		{(*Set).Add, (*Set).AddRange, "AddRange"},
		{(*Set).Flip, (*Set).FlipRange, "FlipRange"},
		{(*Set).Remove, (*Set).RemoveRange, "RemoveRange"},
	}

	for _, x := range []struct {
		S    *Set
		m, n int
	}{
		{New(), 0, 0},
		{New(), 1, 10},
		{New(), 64, 66},
		{New(), 1, 1000},

		{New(1, 2, 3), 0, 1},
		{New(1, 2, 3), 0, 2},
		{New(1, 2, 3), 0, 3},
		{New(1, 2, 3), 0, 4},
		{New(1, 2, 3), 1, 2},
		{New(1, 2, 3), 1, 4},
		{New(1, 2, 3), 1, 5},
		{New(1, 2, 3), 1, 1000},

		{New(100, 200, 300), 50, 100},
		{New(100, 200, 300), 50, 101},
		{New(100, 200, 300), 50, 250},
		{New(100, 200, 300), 50, 350},
		{New(100, 200, 300), 250, 350},
		{New(100, 200, 300), 300, 350},
		{New(100, 200, 300), 350, 400},
		{New(100, 200, 300), 1, 1000},
	} {
		for _, o := range rangeFuncs {
			fInt, fRange, name := o.fInt, o.fRange, o.name
			S := x.S
			m, n := x.m, x.n

			res := fRange(new(Set).Set(S), m, n)
			exp := new(Set).Set(S)
			for i := m; i < n; i++ {
				fInt(exp, i)
			}
			if !res.Equals(exp) {
				t.Errorf("%v.%v(%d, %d) = %v; want %v", S, name, m, n, res, exp)
			}
			CheckInvariants(t, name, res)
		}
	}

	for _, x := range []struct {
		m, n int
	}{
		{2, 1},
		{-2, -1},
		{-1, 0},
		{-1, -1},
	} {
		for _, o := range rangeFuncs {
			N := New()
			if !Panics(o.fRange, N, x.m, x.n) {
				t.Errorf("S.%v(%d, %d) should panic.", o.name, x.m, x.n)
			}
			CheckInvariants(t, "S."+o.name+"(x, x) panic", N)
		}
	}
}

func TestModify(t *testing.T) {
	for _, x := range []struct {
		S       *Set
		n       int
		A, F, R *Set // Expected results for Add, Flip, Remove
	}{
		{New(), 0, New(0), New(0), New()},
		{New(), 65, New(65), New(65), New()},
		{New(), 1000, New(1000), New(1000), New()},

		{New(0), 0, New(0), New(), New()},
		{New(0), 65, New(0, 65), New(0, 65), New(0)},
		{New(0), 1000, New(0, 1000), New(0, 1000), New(0)},

		{New(65), 0, New(0, 65), New(0, 65), New(65)},
		{New(65), 65, New(65), New(), New()},
		{New(65), 1000, New(65, 1000), New(65, 1000), New(65)},

		{New(1, 2, 3), 0, New(0, 1, 2, 3), New(0, 1, 2, 3), New(1, 2, 3)},
		{New(1, 2, 3), 1, New(1, 2, 3), New(2, 3), New(2, 3)},
		{New(1, 2, 3), 2, New(1, 2, 3), New(1, 3), New(1, 3)},
		{New(1, 2, 3), 3, New(1, 2, 3), New(1, 2), New(1, 2)},
		{New(1, 2, 3), 4, New(1, 2, 3, 4), New(1, 2, 3, 4), New(1, 2, 3)},

		{New(100, 200, 300), 0, New(0, 100, 200, 300), New(0, 100, 200, 300), New(100, 200, 300)},
		{New(100, 200, 300), 100, New(100, 200, 300), New(200, 300), New(200, 300)},
		{New(100, 200, 300), 200, New(100, 200, 300), New(100, 300), New(100, 300)},
		{New(100, 200, 300), 300, New(100, 200, 300), New(100, 200), New(100, 200)},
		{New(100, 200, 300), 400, New(100, 200, 300, 400), New(100, 200, 300, 400), New(100, 200, 300)},
	} {
		S := x.S
		n := x.n

		A := New().Set(S).Add(n)
		F := New().Set(S).Flip(n)
		R := New().Set(S).Remove(n)

		if !A.Equals(x.A) {
			t.Errorf("%v.Add(%d) = %v; want %v", S, n, A, x.A)
		}
		if !F.Equals(x.F) {
			t.Errorf("%v.Flip(%d) = %v; want %v", S, n, F, x.F)
		}
		if !R.Equals(x.R) {
			t.Errorf("%v.Remove(%d) = %v; want %v", S, n, R, x.R)
		}
		CheckInvariants(t, "Add", A)
		CheckInvariants(t, "Flip", F)
		CheckInvariants(t, "Remove", R)
	}

	N := New()
	if !Panics((*Set).Add, N, -1) {
		t.Errorf("S.Add(-1) should panic.")
	}
	CheckInvariants(t, "S.Add(-1) panic", N)
	N = New()
	if !Panics((*Set).Flip, New(), -1) {
		t.Errorf("S.Flip(-1) should panic.")
	}
	CheckInvariants(t, "S.Flip(-1) panic", N)
	N = New()
	if !Panics((*Set).Remove, New(), -1) {
		t.Errorf("S.Remove(-1) should panic.")
	}
	CheckInvariants(t, "S.Remove(-1) panic", N)
}

type binOp struct {
	f    func(S *Set, A, B *Set) *Set
	name string
}

func TestBinOp(t *testing.T) {
	And := binOp{(*Set).SetAnd, "SetAnd"}
	AndNot := binOp{(*Set).SetAndNot, "SetAndNot"}
	Or := binOp{(*Set).SetOr, "SetOr"}
	Xor := binOp{(*Set).SetXor, "SetXor"}
	for _, x := range []struct {
		op   binOp
		A, B *Set
		exp  *Set
	}{
		{And, New(), New(), New()},
		{And, New(1), New(), New()},
		{And, New(), New(1), New()},
		{And, New(1), New(1), New(1)},
		{And, New(1), New(2), New()},
		{And, New(1), New(1, 2), New(1)},
		{And, New(1, 2), New(2, 3), New(2)},
		{And, New(100), New(), New()},
		{And, New(), New(100), New()},
		{And, New(100), New(100), New(100)},
		{And, New(100), New(100, 200), New(100)},
		{And, New(200), New(100, 200), New(200)},
		{And, New(100, 200), New(200, 300), New(200)},

		{AndNot, New(), New(), New()},
		{AndNot, New(1), New(), New(1)},
		{AndNot, New(), New(1), New()},
		{AndNot, New(1), New(1), New()},
		{AndNot, New(1), New(2), New(1)},
		{AndNot, New(1), New(1, 2), New()},
		{AndNot, New(1, 2), New(2, 3), New(1)},
		{AndNot, New(100), New(), New(100)},
		{AndNot, New(), New(100), New()},
		{AndNot, New(100), New(100), New()},
		{AndNot, New(100), New(100, 200), New()},
		{AndNot, New(200), New(100, 200), New()},
		{AndNot, New(100, 200), New(200, 300), New(100)},

		{Or, New(), New(), New()},
		{Or, New(), New(1), New(1)},
		{Or, New(1), New(), New(1)},
		{Or, New(1), New(1), New(1)},
		{Or, New(1), New(2), New(1, 2)},
		{Or, New(1), New(1, 2), New(1, 2)},
		{Or, New(1, 2), New(2, 3), New(1, 2, 3)},
		{Or, New(100), New(), New(100)},
		{Or, New(), New(100), New(100)},
		{Or, New(100), New(100), New(100)},
		{Or, New(100), New(100, 200), New(100, 200)},
		{Or, New(200), New(100, 200), New(100, 200)},
		{Or, New(100, 200), New(200, 300), New(100, 200, 300)},

		{Xor, New(), New(), New()},
		{Xor, New(1), New(), New(1)},
		{Xor, New(), New(1), New(1)},
		{Xor, New(1), New(1), New()},
		{Xor, New(1), New(2), New(1, 2)},
		{Xor, New(1), New(1, 2), New(2)},
		{Xor, New(1, 2), New(2, 3), New(1, 3)},
		{Xor, New(100), New(), New(100)},
		{Xor, New(), New(100), New(100)},
		{Xor, New(100), New(100), New()},
		{Xor, New(100), New(100, 200), New(200)},
		{Xor, New(200), New(100, 200), New(100)},
		{Xor, New(100, 200), New(200, 300), New(100, 300)},
	} {
		op, name := x.op.f, x.op.name
		A, B := New().Set(x.A), New().Set(x.B)
		S := New()

		res := op(S, A, B)
		exp := x.exp
		if S != res {
			t.Errorf("&(S.%s(%v, %v)) = %p &S = %p; want same", name, A, B, S, res)
		}
		if !res.Equals(exp) {
			t.Errorf("S.%s(%v, %v) = %v; want %v", name, x.A, x.B, res, exp)
		}
		CheckInvariants(t, name, res)

		A.Set(x.A)
		B.Set(x.B)
		S = A
		res = op(S, A, B)
		if !res.Equals(exp) {
			t.Errorf("S.%s(%v, %v) = %v; want %v", name, x.A, x.B, res, exp)
		}
		CheckInvariants(t, name, res)

		A.Set(x.A)
		B.Set(x.B)
		S = B
		res = op(S, A, B)
		if !res.Equals(exp) {
			t.Errorf("S.%s(%v, %v) = %v; want %v", name, x.A, x.B, res, exp)
		}
		CheckInvariants(t, name, res)

		A.Set(x.A)
		B.Set(x.B)
		S = New().AddRange(150, 250)
		res = op(S, A, B)
		if !res.Equals(exp) {
			t.Errorf("S.%s(%v, %v) = %v; want %v", name, x.A, x.B, res, exp)
		}
		CheckInvariants(t, name, res)
	}
}

func TestClear(t *testing.T) {
	for _, S := range []*Set{
		New(),
		New(1),
		New(65),
		New(1, 2, 3),
		New(100, 200, 300),
	} {
		A := New().Set(S)
		empty := New()

		T := S.Clear()
		if T != S {
			t.Errorf("&(S.Clear()) = %p, &S = %p; want same", T, S)
		}
		if !T.Equals(empty) {
			t.Errorf("%v.Clear() = %v; want %v", A, T, empty)
		}
		CheckInvariants(t, "Clear", T)
	}
}

func TestIsempty(t *testing.T) {
	for _, x := range []struct {
		S     *Set
		empty bool
	}{
		{New(), true},
		{New(1), false},
		{New(65), false},
		{New(1, 2, 3), false},
		{New(100, 200, 300), false},
	} {
		S := x.S

		empty := S.IsEmpty()
		if empty != x.empty {
			t.Errorf("%v.IsEmpty() = %v; want %v", S, empty, x.empty)
		}
	}
}

func TestDo(t *testing.T) {
	for _, x := range []struct {
		S   *Set
		res string
	}{
		{New(), ""},
		{New(0), "0"},
		{New(1, 2, 3, 62, 63, 64), "123626364"},
		{New(1, 22, 333, 4444), "1223334444"},
	} {
		S := x.S
		res := ""

		S.Do(func(n int) {
			res += strconv.Itoa(n)
		})
		if res != x.res {
			t.Errorf("%v.Do(func(n int) { s += Itoa(n) }) -> s=%q; want %q", S, res, x.res)
		}

		S = x.S
		res = ""
		S.Do(func(n int) {
			S.RemoveRange(0, n+1)
			res += strconv.Itoa(n)
		})
		if res != x.res {
			t.Errorf("%v.Do(func(n int) { S.RemoveRange(0, n+1); s += Itoa(n) }) -> s=%q; want %q", S, res, x.res)
		}
	}
}

func TestString(t *testing.T) {
	for _, x := range []struct {
		S   *Set
		res string
	}{
		{New(), "{}"},
		{New(1), "{1}"},
		{New(1, 2), "{1, 2}"},
		{New(1, 3), "{1, 3}"},
		{New(0, 2, 3), "{0, 2, 3}"},
		{New(0, 1, 3), "{0, 1, 3}"},
		{New(0, 2, 3, 5), "{0, 2, 3, 5}"},
		{New(0, 1, 2, 4, 5), "{0..2, 4, 5}"},
		{New(0, 1, 2, 3, 5, 7, 8, 9), "{0..3, 5, 7..9}"},
		{New(65), "{65}"},
		{New(100, 200, 300), "{100, 200, 300}"},
	} {
		res := x.S.String()
		if res != x.res {
			t.Errorf("S.String() = %q; want %q", res, x.res)
		}
	}
}

func TestNextPow2(t *testing.T) {
	for _, x := range []struct {
		n, p int
	}{
		{MinInt, 1},
		{-1, 1},
		{0, 1},
		{1, 2},
		{2, 4},
		{3, 4},
		{4, 8},
		{1<<19 - 1, 1 << 19},
		{1 << 19, 1 << 20},
		{MaxInt >> 1, MaxInt>>1 + 1},
		{MaxInt>>1 + 1, MaxInt},
		{MaxInt - 1, MaxInt},
		{MaxInt, MaxInt},
	} {
		n := x.n

		p := nextPow2(n)
		if p != x.p {
			t.Errorf("nextPow2(%#x) = %#x; want %#x", n, p, x.p)
		}
	}
}
