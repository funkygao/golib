package hack

// Append basically same as builtin append except that it controls the grow factor.
// Builtin append grow factor=cap*2.
func Append(s []interface{}, val interface{}, grow int) []interface{} {
	l := len(s)
	target := s
	if l == cap(target) {
		// will grow by specified size instead of default cap*2
		target = make([]interface{}, l, l+grow)
		copy(target, s[:])
	}

	target = append(target, val)
	return target
}
