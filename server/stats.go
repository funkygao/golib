package server

import (
	"github.com/funkygao/golib/gofmt"
	log "github.com/funkygao/log4go"
	"runtime"
	"syscall"
	"time"
)

func RunSysStats(startedAt time.Time, interval time.Duration) {
	const nsInMs uint64 = 1000 * 1000

	ticker := time.NewTicker(interval)
	defer func() {
		ticker.Stop()
	}()

	var (
		ms           = new(runtime.MemStats)
		rusage       = &syscall.Rusage{}
		lastUserTime int64
		lastSysTime  int64
		userTime     int64
		sysTime      int64
		userCpuUtil  float64
		sysCpuUtil   float64
	)

	for _ = range ticker.C {
		runtime.ReadMemStats(ms)

		syscall.Getrusage(syscall.RUSAGE_SELF, rusage)
		syscall.Getrusage(syscall.RUSAGE_SELF, rusage)
		userTime = rusage.Utime.Sec*1000000000 + int64(rusage.Utime.Usec)
		sysTime = rusage.Stime.Sec*1000000000 + int64(rusage.Stime.Usec)
		userCpuUtil = float64(userTime-lastUserTime) * 100 / float64(interval)
		sysCpuUtil = float64(sysTime-lastSysTime) * 100 / float64(interval)

		lastUserTime = userTime
		lastSysTime = sysTime

		log.Info("ver:%s, since:%s, go:%d, gc:%dms/%d=%d, heap:{%s, %s, %s, %s} cpu:{%3.2f%%us, %3.2f%%sy}",
			BuildID,
			time.Since(startedAt),
			runtime.NumGoroutine(),
			ms.PauseTotalNs/nsInMs,
			ms.NumGC,
			ms.PauseTotalNs/(nsInMs*uint64(ms.NumGC))+1,
			gofmt.ByteSize(ms.HeapSys),      // bytes it has asked the operating system for
			gofmt.ByteSize(ms.HeapAlloc),    // bytes currently allocated in the heap
			gofmt.ByteSize(ms.HeapIdle),     // bytes in the heap that are unused
			gofmt.ByteSize(ms.HeapReleased), // bytes returned to the operating system, 5m for scavenger
			userCpuUtil,
			sysCpuUtil)

	}
}
