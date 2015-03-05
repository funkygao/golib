package fixture

import (
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
