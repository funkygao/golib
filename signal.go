package main

import (
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

type signalHandler func(os.Signal)

var signals = struct {
	sync.Mutex // map in go is not thread/goroutine safe
	handlers   map[os.Signal]signalHandler
	C          chan os.Signal
}{
	handlers: make(map[os.Signal]signalHandler),
	C:        make(chan os.Signal, 20),
}

func trapSignals() {
	for sig := range signals.C {
		signals.Lock() // map not goroutine safe in golang
		handler := signals.handlers[sig]
		signals.Unlock()
		if handler != nil {
			handler(sig)
		}
	}
}

func registerSignalHandler(sig os.Signal, handler signalHandler) {
	signals.Lock()
	defer signals.Unlock()

	if _, present := signals.handlers[sig]; !present {
		signals.handlers[sig] = handler
		signal.Notify(signals.C, sig)
	}
}

func handleIgnore(sig os.Signal) {}

func showWorkers(sig os.Signal) {
	logger.Printf("workers: %+v\n", allWorkers)
	logger.Printf("parsers: %+v\n", parser.Parsers())
}

func handleInterrupt(sig os.Signal) {
	logger.Printf("got signal %s\n", strings.ToUpper(sig.String()))

	shutdown()
}

func setupSignals() {
	go trapSignals()

	//registerSignalHandler(syscall.SIGHUP, handleInterrupt)
	registerSignalHandler(syscall.SIGINT, handleInterrupt)
	registerSignalHandler(syscall.SIGTERM, handleInterrupt)
	registerSignalHandler(syscall.SIGUSR2, showWorkers)
}
