package debug

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
)

const (
	depth_char = " "
	depth_step = 4
)

var (
	enabled bool
	depth   int32
	lock    *sync.Mutex
)

// Entering into a func.
func Trace(fn string) string {
	if !enabled {
		return fn
	}

	if fn == "" {
		fn = CallerFuncName(2)
	}

	space := strings.Repeat(depth_char, int(atomic.LoadInt32(&depth)))
	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	}
	fmt.Printf("%s %s %s\n", space, "Entering:", fn)

	atomic.AddInt32(&depth, depth_step)

	return fn
}

// Leaving from a func.
func Un(fn string) {
	if !enabled {
		return
	}

	if fn == "" {
		fn = CallerFuncName(2)
	}

	atomic.AddInt32(&depth, -depth_step)

	space := strings.Repeat(depth_char, int(atomic.LoadInt32(&depth)))
	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	}
	fmt.Printf("%s %s %s\n", space, "Leaving :", fn)
}

// Enable the trace output.
func EnableTrace() {
	enabled = true
}

// Disable the trace output.
func DisableTrace() {
	enabled = false
}

// Setup the global mutex.
func SetLock(l *sync.Mutex) {
	lock = l
}

// Caller func name with skip as the call stack level.
func CallerFuncName(skipFrames int) string {
	pc, fn, line, _ := runtime.Caller(skipFrames)
	f := runtime.FuncForPC(pc)
	return fmt.Sprintf("%s:%s:%d", f.Name(), fn, line) // the caller func name
}
