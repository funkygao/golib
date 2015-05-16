package vclock

import (
	"bytes"
	"testing"
)

func assert(expression bool) {
	if !expression {
		panic("assertion failed")
	}
}

func TestVectorClocks(t *testing.T) {
	// a round
	c1a := VectorClockWithValues(1, 0)
	c2a := VectorClockWithValues(0, 1)
	if c1a.Equals(c2a) {
		t.Error("should not be equal", c1a, c2a)
	}
	if !c1a.Equals(VectorClockWithValues(1, 0)) {
		t.Error("should be equal", c1a, VectorClockWithValues(1, 0))
	}

	// b round happens after a
	c1b := VectorClockWithValues(2, 1)
	c2b := VectorClockWithValues(1, 2)

	verifyHappensBefore := func(c1 VectorClock, c2 VectorClock) {
		if !c1.happensBefore(c2) {
			t.Error("Expected", c1, " <--happens before--", c2)
		}
		if c2.happensBefore(c1) {
			t.Error("Expected", c2, " NOT happens before", c1)
		}
	}

	// verify a round before b round
	verifyHappensBefore(c1a, c1b)
	verifyHappensBefore(c2a, c1b)
	verifyHappensBefore(c1a, c2b)
	verifyHappensBefore(c2a, c2b)

	// verify concurrent
	if !c1a.concurrentWith(c2a) {
		t.Error("concurrent", c1a, c2a)
	}
	if !c2a.concurrentWith(c1a) {
		t.Error("concurrent", c1a, c2a)
	}

	merged := c1b.merge(c2b)
	if !merged.Equals(VectorClockWithValues(2, 2)) {
		t.Error(merged)
	}
}

const example = `a 0, 1  
0, 2	  
	 c., !label	 1 100	  ` + "\n\n"

func TestParse(t *testing.T) {
	clocks, labels, err := Parse(bytes.NewBufferString(example))
	if err != nil {
		t.Error(err)
	}

	if len(clocks) != 3 {
		t.Fatal(clocks)
	}
	if !clocks[0].Equals(VectorClockWithValues(0, 1)) {
		t.Error(clocks[0])
	}
	if !clocks[1].Equals(VectorClockWithValues(0, 2)) {
		t.Error(clocks[1])
	}
	if !clocks[2].Equals(VectorClockWithValues(1, 100)) {
		t.Error(clocks[2])
	}

	if labels[VectorClockWithValues(0, 1).String()] != "a" {
		t.Error(labels[VectorClockWithValues(0, 1).String()])
	}
	if v, ok := labels[VectorClockWithValues(0, 2).String()]; ok {
		t.Error("label should not exist:", v, ok)
	}
	if labels[VectorClockWithValues(1, 100).String()] != "c., !label" {
		t.Error(labels[VectorClockWithValues(1, 100).String()])
	}
}
