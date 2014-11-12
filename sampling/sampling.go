package sampling

import (
	"math/rand"
    "time"
)

func init() {
    rand.Seed(time.Now().UTC().UnixNano())
}

// Base is 1000
func SampleRateSatisfied(rate int) bool {
	return rand.Intn(1000) <= rate
}
