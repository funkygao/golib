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

	configFile         string
	configFileLastStat os.FileInfo
}

func NewServer(name string) (this *Server) {
	this = new(Server)
	this.Name = name

	return
}

func (this *Server) LoadConfig(fn string) *Server {
	log.Info("Server[%s %s@%s] loading config file: %s", this.Name, BuildId, Version, fn)
	this.configFile = fn

	var err error
	this.Conf, err = conf.Load(fn)
	if err != nil {
		panic(err)
	}

	return this
}

// Hot reload of config file
func (this *Server) WatchConfig(interval time.Duration, ch chan *conf.Conf) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for _ = range ticker.C {
		stat, _ := os.Stat(this.configFile)
		if stat.ModTime() != this.configFileLastStat.ModTime() {
			this.configFileLastStat = stat

			cf := this.LoadConfig(this.configFile)
			log.Info("config[%s] reloaded", this.configFile)
			ch <- cf.Conf
		}
	}
}

func (this *Server) Launch() {
	this.StartedAt = time.Now()
	this.hostname, _ = os.Hostname()
	this.pid = os.Getpid()
	signal.IgnoreSignal(syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGSTOP)

	runtime.GOMAXPROCS(this.Int("max_cpu", runtime.NumCPU()))
}
