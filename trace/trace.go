package trace

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	DEPTH_CHAR = " "
	DEPTH_STEP = 4
)

var (
	enabled bool
	depth   int32
	lock    *sync.Mutex
)

// Entering into a func
func Trace(fn string) string {
	if !enabled {
		return fn
	}

	if fn == "" {
		fn = CallerFuncName(2)
	}

	space := strings.Repeat(DEPTH_CHAR, int(atomic.LoadInt32(&depth)))
	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	}
	fmt.Printf("%s %s %s\n", space, "Entering:", fn)

	atomic.AddInt32(&depth, DEPTH_STEP)

	return fn
}

// Leaving from a func
func Un(fn string) {
	if !enabled {
		return
	}

	if fn == "" {
		fn = CallerFuncName(2)
	}

	atomic.AddInt32(&depth, -DEPTH_STEP)

	space := strings.Repeat(DEPTH_CHAR, int(atomic.LoadInt32(&depth)))
	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	}
	fmt.Printf("%s %s %s\n", space, "Leaving :", fn)
}

// Enable the trace output
func Enable() {
	enabled = true
}

// Disable the trace output
func Disable() {
	enabled = false
}

// Setup the global mutex
func SetLock(l *sync.Mutex) {
	lock = l
}

// Measure how long it takes to run a func
func Timeit(fun interface{}, args ...interface{}) (result []reflect.Value, delta time.Duration, err error) {
	// assertions
	f := reflect.ValueOf(fun)
	if f.Kind().String() != "func" {
		err = errors.New("fun must be a func")
		return
	}
	if len(args) != f.Type().NumIn() {
		err = errors.New("input arguments not match")
		return
	}

	// prepare input arguments
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}

	defer func(start time.Time) {
		delta = time.Since(start)
	}(time.Now())
	result = f.Call(in) // call the universal func
	return
}

// Caller func name with skip as the call stack level
func CallerFuncName(calldepth int) string {
	pc, _, _, _ := runtime.Caller(calldepth)
	f := runtime.FuncForPC(pc)
	return f.Name() // the caller func name
}
