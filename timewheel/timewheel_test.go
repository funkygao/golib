package timewheel

import (
	"testing"
	"time"
)

func TestTimeWheel(t *testing.T) {
	w := NewTimeWheel(100*time.Millisecond, 10)

	for {
		select {
		case <-w.After(200 * time.Millisecond):
			return
		}
	}
}
