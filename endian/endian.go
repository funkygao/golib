package endian

import (
	"unsafe"
)

func EndianIsLittle() bool {
	var word uint16 = 1
	littlePtr := (*uint8)(unsafe.Pointer(&word))
	return *littlePtr == 1
}

func EndianIsBig() bool {
	return !EndianIsLittle()
}

func SafeSplitUint16(val uint16) (leastSignificant, mostSignificant uint8) {
	bytes := (*[2]uint8)(unsafe.Pointer(&val))
	if EndianIsLittle() {
		return bytes[0], bytes[1]
	}

	return bytes[1], bytes[0]
}
