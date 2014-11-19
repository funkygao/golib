package server

import (
	conf "github.com/funkygao/jsconf"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"
	"time"
)

type Server struct {
	*conf.Conf

	Name       string
	configFile string
	StartedAt  time.Time
	pid        int
	hostname   string
}

func NewServer(name string) (this *Server) {
	this = new(Server)
	this.Name = name

	return
}

func KillProcess(pidFile string) error {
	pidBody, err := ioutil.ReadFile(pidFile)
	if err != nil {
		return err
	}

	pid, err := strconv.Atoi(string(pidBody))
	if err != nil {
		return err
	}

	serverProcess, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	if err := serverProcess.Kill(); err != nil {
		return err
	}

	return syscall.Unlink(pidFile)
}
