package sampling

import (
	"math/rand"
	"os"
	"time"
)

const (
	BASE = 10000
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano() + int64(os.Getpid()))
}

func SampleRateSatisfied(rate int) bool {
	return rand.Intn(BASE) <= rate
}
