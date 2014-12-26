package server

import (
	"github.com/funkygao/golib/signal"
	"os"
	"runtime"
	"syscall"
	"time"
)

func (this *Server) Launch() {
	this.StartedAt = time.Now()
	this.hostname, _ = os.Hostname()
	this.pid = os.Getpid()
	signal.IgnoreSignal(syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGSTOP)

	runtime.GOMAXPROCS(this.Int("max_cpu", runtime.NumCPU()))
}
