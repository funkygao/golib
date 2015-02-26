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

// StringArena lets you consolidate allocations for a group of strings
// that have similar life length.
type StringArena struct {
	buf []byte
	str string
}

func NewStringArena(size int) *StringArena {
	this := &StringArena{buf: make([]byte, 0, size)}
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&this.buf))
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&this.str))
	pstring.Data = pbytes.Data
	pstring.Len = pbytes.Cap
	return this
}

// NewString copies a byte slice into the arena and returns it as a string.
// If the arena is full, it returns a traditional go string.
func (this *StringArena) NewString(b []byte) string {
	if len(this.buf)+len(b) > cap(this.buf) {
		return string(b)
	}

	start := len(this.buf)
	this.buf = append(this.buf, b...)
	return this.str[start : start+len(b)]
}

func (this *StringArena) SpaceLeft() int {
	return cap(this.buf) - len(this.buf)
}

func (this *StringArena) SpaceUsed() int {
	return len(this.buf)
}
