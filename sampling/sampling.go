package sampling

import (
	"math/rand"
	"time"
)

const (
	BASE = 10000
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func SampleRateSatisfied(rate int) bool {
	return rand.Intn(BASE) <= rate
}
