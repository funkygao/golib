package signal

import (
	"os"
	"os/signal"
	"sync"
)

type SignalHandler func(os.Signal)

var signals = struct {
	sync.Mutex // map in go is not thread/goroutine safe

	handlers map[os.Signal]SignalHandler
	ch       chan os.Signal
}{
	handlers: make(map[os.Signal]SignalHandler),
	ch:       make(chan os.Signal),
}

func init() {
	go func() {
		for sig := range signals.ch {
			signals.Lock() // map not goroutine safe in golang
			handler := signals.handlers[sig]
			signals.Unlock()
			if handler != nil {
				handler(sig)
			}
		}
	}()
}

func RegisterSignalHandler(sig os.Signal, handler SignalHandler) {
	signals.Lock()
	defer signals.Unlock()

	_, present := signals.handlers[sig]
	if !present {
		signals.handlers[sig] = handler
		signal.Notify(signals.ch, sig)
	}
}

func IgnoreSignal(sig ...os.Signal) {
	ignoreFunc := func(s os.Signal) {}
	for _, s := range sig {
		RegisterSignalHandler(s, ignoreFunc)
	}
}

// Send a signal to current running proc
func Kill(sig os.Signal) {
	go func() {
		signals.ch <- sig
	}()
}
