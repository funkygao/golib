package rand

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"math/rand"
)

const (
	minChar = '!'
	maxChar = '~'
)

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
