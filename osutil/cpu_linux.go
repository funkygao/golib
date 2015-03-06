// +build linux,!appengine

package osutil

import (
	"syscall"
	"time"
)

func init() {
	cpuUsage = cpuLinux
}

func cpuLinux() time.Duration {
	var ru syscall.Rusage
	syscall.Getrusage(0, &ru)
	return time.Duration(ru.Utime.Nano())
}
