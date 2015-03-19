package osutil

var memUsage func() int64

// MemUsage returns the number of bytes used by the process.
// On unsupported operating systems, it returns zero.
func MemUsage() int64 {
	if f := memUsage; f != nil {
		return f()
	}
	return 0
}
