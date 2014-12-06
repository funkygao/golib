package server

import (
	"io/ioutil"
	"os"
	"strconv"
	"syscall"
)

func findProcess(pidFile string) (p *os.Process, err error) {
	pidBody, err := ioutil.ReadFile(pidFile)
	if err != nil {
		return nil, err
	}

	pid, err := strconv.Atoi(string(pidBody))
	if err != nil {
		return nil, err
	}

	return os.FindProcess(pid)
}

func KillProcess(pidFile string) error {
	process, err := findProcess(pidFile)
	if err != nil {
		return err
	}

	if err := process.Kill(); err != nil {
		return err
	}

	syscall.Unlink(pidFile)
	return nil
}

func SignalProcess(pidFile string, sig os.Signal) error {
	process, err := findProcess(pidFile)
	if err != nil {
		return err
	}

	return process.Signal(sig)
}
