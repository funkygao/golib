package osutil

var netstatFn func() map[string]int64

func Netstat() map[string]int64 {
	if netstatFn != nil {
		return netstatFn()
	}
	return nil
}
