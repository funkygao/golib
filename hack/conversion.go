package hack

import (
	"reflect"
	"unsafe"
)

// String force casts a []byte to a string.
// USE AT YOUR OWN RISK
func String(b []byte) (s string) {
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pstring.Data = pbytes.Data
	pstring.Len = pbytes.Len
	return
}

// Caution, returned byte array is read only.
// If you change its content:
// panic: runtime error: invalid memory address or nil pointer dereference.
func Byte(s string) []byte {
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	pbytes.Cap = len(s)
	pbytes.Len = pbytes.Cap
	return *(*[]byte)(unsafe.Pointer(pbytes))
}
