package fixture

import (
	"math/rand"
)

func RandomString(l int) string {
	return string(RandomByteSlice(l))
}

func RandomByteSlice(l int) []byte {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(65 + rand.Intn(90-65))
	}
	return bytes
}
