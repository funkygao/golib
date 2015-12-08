package timewheel

import (
	"testing"
	"time"
)

func TestTimeWheel(t *testing.T) {
	w := NewTimeWheel(time.Second, 10)
	defer w.Stop()

	t1 := time.Now()

	go func() {
		select {
		case <-w.After(time.Second):
			t.Logf("expected 1s, got %s", time.Since(t1))
		}
	}()

	go func() {
		select {
		case <-w.After(2 * time.Second):
			t.Logf("expected 2s, got %s", time.Since(t1))
		}
	}()

	for {
		select {
		case <-w.After(3 * time.Second):
			t.Logf("expected 3s, got %s", time.Since(t1))
			return
		}
	}
}
