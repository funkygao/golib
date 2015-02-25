// Null is zero memory overhead
// Classical usage:
// to design a Set datastructure, you can
// type Set map[interface{}]bool
// But bool is 8 bit, while Null is 0 byte
// So, type Set map[interface{}]NullStruct
package null
