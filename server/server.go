package server

import (
	"github.com/funkygao/golib/signal"
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
	"os"
	"runtime"
	"syscall"
	"time"
)

type Server struct {
	*conf.Conf

	Name      string
	StartedAt time.Time
	pid       int
	hostname  string
}

func NewServer(name string) (this *Server) {
	this = new(Server)
	this.Name = name

	return
}

func (this *Server) LoadConfig(fn string) *Server {
	log.Info("Server[%s %s@%s] loading config file: %s", this.Name, BuildId, Version, fn)

	var err error
	this.Conf, err = conf.Load(fn)
	if err != nil {
		panic(err)
	}

	return this
}

func (this *Server) Launch() {
	this.StartedAt = time.Now()
	this.hostname, _ = os.Hostname()
	this.pid = os.Getpid()
	signal.IgnoreSignal(syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGSTOP)

	runtime.GOMAXPROCS(this.Int("max_cpu", runtime.NumCPU()))
}
