package sampling

import (
	"math/rand"
)

// Base is 1000
func SampleRateSatisfied(rate int) bool {
	return rand.Intn(1000) <= rate
}
