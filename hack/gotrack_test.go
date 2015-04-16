package hack

import (
	"fmt"
	"strings"
	"testing"
)

func TestGoroutineLock(t *testing.T) {
	DebugGoroutines = true
	g := newGoroutineLock()
	g.check()

	sawPanic := make(chan interface{})
	go func() {
		defer func() { sawPanic <- recover() }()
		g.check() // should panic
	}()
	e := <-sawPanic
	if e == nil {
		t.Fatal("did not see panic from check in other goroutine")
	}
	if !strings.Contains(fmt.Sprint(e), "wrong goroutine") {
		t.Errorf("expected on see panic about running on the wrong goroutine; got %v", e)
	}
}
