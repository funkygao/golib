package signal

import (
	"github.com/funkygao/assert"
	"os"
	"syscall"
	"testing"
)

func TestAll(t *testing.T) {
	var gotSignals []os.Signal = make([]os.Signal, 0, 10)
	RegisterSignalHandler(syscall.SIGINT, func(sig os.Signal) {
		gotSignals = append(gotSignals, sig)
		t.Logf("got sig: %v, all:%+v", sig, gotSignals)
		assert.Equal(t, []os.Signal{syscall.SIGINT}, gotSignals)
	})
	RegisterSignalHandler(syscall.SIGHUP, func(sig os.Signal) {
		gotSignals = append(gotSignals, sig)
		t.Logf("got sig: %v, all:%+v", sig, gotSignals)
		assert.Equal(t, []os.Signal{syscall.SIGINT, syscall.SIGHUP}, gotSignals)
	})

	IgnoreSignal(syscall.SIGALRM, syscall.SIGABRT)

	assert.Equal(t, nil, Kill(syscall.SIGALRM,
		syscall.SIGINT, syscall.SIGHUP, syscall.SIGABRT))

	pid := os.Getpid()
	assert.Equal(t, nil, SignalProcess(pid, syscall.SIGALRM))
	assert.Equal(t, nil, SignalProcess(pid, syscall.SIGINT))
	assert.Equal(t, nil, SignalProcess(pid, syscall.SIGHUP))
	assert.Equal(t, nil, SignalProcess(pid, syscall.SIGABRT))
}
