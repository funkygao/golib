package osutil

import (
	"time"
)

var cpuUsage func() time.Duration

// CPUUsage returns how much cumulative user CPU time the process has
// used. On unsupported operating systems, it returns zero.
func CPUUsage() time.Duration {
	if f := cpuUsage; f != nil {
		return f()
	}
	return 0
}
