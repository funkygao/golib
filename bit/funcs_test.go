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
	"reflect"
	"testing"
)

func TestConstants(t *testing.T) {
	if BitsPerWord != 32 && BitsPerWord != 64 {
		t.Errorf("BitsPerWord = %v; want 32 or 64", BitsPerWord)
	}

	if BitsPerWord == 32 {
		if MaxUint != 1<<32-1 {
			t.Errorf("MaxUint = %#x; want 1<<32 - 1", uint64(MaxUint))
		}
		if MaxInt != 1<<31-1 {
			t.Errorf("MaxInt = %#x; want 1<<31 - 1", int64(MaxInt))
		}
		if MinInt != -1<<31 {
			t.Errorf("MaxUint = %#x; want -1 << 31", int64(MinInt))
		}
	}

	if BitsPerWord == 64 {
		if MaxUint != 1<<64-1 {
			t.Errorf("MaxUint = %#x; want 1<<64 - 1", uint64(MaxUint))
		}
		if MaxInt != 1<<63-1 {
			t.Errorf("MaxInt = %#x; want 1<<63 - 1", int64(MaxInt))
		}
		if MinInt != -1<<63 {
			t.Errorf("MaxUint = %#x; want -1 << 63", int64(MinInt))
		}
	}
}

// Tests MinPos, MaxPos, and Count on all words with one nonzero bit.
func TestWordOneBit(t *testing.T) {
	for i := 0; i < 64; i++ {
		var w uint64 = 1 << uint(i)

		min, max, count := MinPos(w), MaxPos(w), Count(w)
		if min != i {
			t.Errorf("MinPos(%#x) = %d; want %d", w, min, i)
		}
		if max != i {
			t.Errorf("MaxPos(%#x) = %d; want %d", w, max, i)
		}
		if count != 1 {
			t.Errorf("Count(%#x) = %d; want %d", w, count, 1)
		}
	}
}

func TestWordFuncs(t *testing.T) {
	for _, x := range []struct {
		w               uint64
		min, max, count int
	}{
		{0xa, 1, 3, 2},
		{0xffffffffffffffff, 0, 63, 64},
		{0x7ffffffffffffffe, 1, 62, 62},
		{0x5555555555555555, 0, 62, 32},
		{0xaaaaaaaaaaaaaaaa, 1, 63, 32},
	} {
		w := x.w

		min, max, count := MinPos(w), MaxPos(w), Count(w)
		if min != x.min {
			t.Errorf("MinPos(%#x) = %v; want %v", w, min, x.min)
		}
		if max != x.max {
			t.Errorf("MaxPos(%#x) = %v; want %v", w, max, x.max)
		}
		if count != x.count {
			t.Errorf("Count(%#x) = %v; want %v", w, count, x.count)
		}
	}

	if !Panics(MinPos, uint64(0)) {
		t.Errorf("MinPos(0) should panic.")
	}
	if !Panics(MaxPos, uint64(0)) {
		t.Errorf("MaxPos(0) should panic.")
	}
}

// Returns true if function f panics with parameters p.
func Panics(f interface{}, p ...interface{}) bool {
	fv := reflect.ValueOf(f)
	ft := reflect.TypeOf(f)

	if ft.NumIn() != len(p) {
		panic("wrong argument count")
	}

	pv := make([]reflect.Value, len(p))
	for i, v := range p {
		if reflect.TypeOf(v) != ft.In(i) {
			panic("wrong argument type")
		}
		pv[i] = reflect.ValueOf(v)
	}

	return call(fv, pv)
}

func call(fv reflect.Value, pv []reflect.Value) (b bool) {
	defer func() {
		if err := recover(); err != nil {
			b = true
		}
	}()

	fv.Call(pv)
	return
}
