// Package signal provides helpers to handle OS signals.
package signal

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

type SignalHandler func(os.Signal)

var (
	signals = struct {
		sync.Mutex // map in go is not thread/goroutine safe
		handlers   map[os.Signal]SignalHandler
		ch         chan os.Signal
	}{
		handlers: make(map[os.Signal]SignalHandler),
		ch:       make(chan os.Signal, syscall.SIGUSR2), // SIGUSR2 is max
	}
)

func init() {
	go func() {
		for sig := range signals.ch {
			signals.Lock()
			sigHandler := signals.handlers[sig]
			signals.Unlock()
			if sigHandler != nil {
				sigHandler(sig)
			}
		}
	}()
}

func RegisterSignalHandler(sig os.Signal, handler SignalHandler) {
	signals.Lock()
	defer signals.Unlock()

	if _, present := signals.handlers[sig]; !present {
		signals.handlers[sig] = handler
		signal.Notify(signals.ch, sig)
	}
}

// Let current process ignore some os signals.
func IgnoreSignal(sig ...os.Signal) {
	ignoreFunc := func(s os.Signal) {}
	for _, s := range sig {
		RegisterSignalHandler(s, ignoreFunc)
	}
}

// Send a signal to current running process.
func Kill(sigs ...os.Signal) error {
	for _, sig := range sigs {
		select {
		case signals.ch <- sig:
		default:
			return errors.New(fmt.Sprintf("signal:%v discarded", sig))
		}
	}

	return nil
}

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

// Send a signal to a process by pid.
func SignalProcess(pid int, sig os.Signal) error {
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	return p.Signal(sig)
}

func SignalProcessByPidFile(pidFile string, sig os.Signal) error {
	p, err := findProcess(pidFile)
	if err != nil {
		return err
	}

	return p.Signal(sig)
}
