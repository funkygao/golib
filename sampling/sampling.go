package servant

import (
	"math/rand"
)

func SampleRateSatisfied(rate int) bool {
	return rand.Intn(100) <= rate
}
