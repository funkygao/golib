package debug

import (
	"errors"
	"reflect"
	"time"
)

// Measure how long it takes to run a func.
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
