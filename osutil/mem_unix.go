// +build linux,!appengine darwin freebsd

package osutil

import (
	"runtime"
	"syscall"
)

func init() {
	memUsage = memUnix
}

func memUnix() int64 {
	var ru syscall.Rusage
	syscall.Getrusage(0, &ru)
	if runtime.GOOS == "linux" {
		// in KB
		return int64(ru.Maxrss) << 10
	}
	// In bytes:
	return int64(ru.Maxrss)
}
