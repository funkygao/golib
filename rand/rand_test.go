package rand

import (
	"github.com/funkygao/assert"
	"testing"
)

func TestRandomBytes(t *testing.T) {
	b := RandomByteSlice(20)
	t.Logf("%v", b)
	assert.Equal(t, 20, len(b))
}

func TestNewPseudoSeed(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Logf("%d\n", NewPseudoSeed())
	}
}

func TestShuffleInts(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 9, 7}
	t.Logf("%+v", ShuffleInts(a))
}

func TestWeightedRand(t *testing.T) {
	weights := map[string]int{
		"A20": 20,
		"B40": 40,
		"C10": 10,
		"D80": 80,
	}
	for i := 0; i < 20; i++ {
		t.Logf("%s", WeightedRand(weights))
	}
}
