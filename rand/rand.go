package rand

import (
	crand "crypto/rand"
	crypto_rand "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"io"
	"math/rand"
	"time"
)

const (
	minChar = '!'
	maxChar = '~'
)

func RandSeedWithTime() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func RandomString(size int) string {
	return string(RandomByteSlice(size))
}

func RandomByteSlice(size int) []byte {
	bytes := make([]byte, size)
	for i := 0; i < size; i++ {
		bytes[i] = minChar + byte(rand.Intn(int(maxChar-minChar)))
	}
	return bytes
}

// NewPseudoSeed generates a seed from crypto/rand.
func NewPseudoSeed() (seed int64) {
	binary.Read(crypto_rand.Reader, binary.LittleEndian, &seed)
	return
}

func SizedString(size int) string {
	u := make([]byte, size)
	io.ReadFull(crand.Reader, u)
	return hex.EncodeToString(u)
}

func WeightedRand(weights map[string]int) string {
	sum := 0
	for _, v := range weights {
		sum += v
	}

	r := rand.Intn(sum)
	for k, v := range weights {
		r -= v
		if r <= 0 {
			return k
		}
	}

	// should never reach here
	return ""
}
