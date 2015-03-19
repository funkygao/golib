package osutil

import (
	"syscall"
	"time"
)

func init() {
	cpuUsage = cpuFreeBSD
}

func cpuFreeBSD() time.Duration {
	var ru syscall.Rusage
	syscall.Getrusage(0, &ru)
	return time.Duration(ru.Utime.Nano())
}
