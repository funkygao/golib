package golib

import (
	"os"
	"os/signal"
	"sync"
)

type signalHandler func(os.Signal)

var signals = struct {
	sync.Mutex // map in go is not thread/goroutine safe

	handlers map[os.Signal]signalHandler
	ch       chan os.Signal
}{
	handlers: make(map[os.Signal]signalHandler),
	ch:       make(chan os.Signal, 20),
}

func init() {
	for sig := range signals.ch {
		signals.Lock() // map not goroutine safe in golang
		handler := signals.handlers[sig]
		signals.Unlock()
		if handler != nil {
			handler(sig)
		}
	}
}

func RegisterSignalHandler(sig os.Signal, handler signalHandler) {
	signals.Lock()
	defer signals.Unlock()

	if _, present := signals.handlers[sig]; !present {
		signals.handlers[sig] = handler
		signal.Notify(signals.ch, sig)
	}
}

func IgnoreSignal(sig os.Signal) {}
