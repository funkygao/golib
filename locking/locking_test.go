package locking

import (
	"github.com/funkygao/assert"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestFlockAndFunlock(t *testing.T) {
	filepath := ".lk"
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = Flock(f, time.Second)
	assert.Equal(t, nil, err)
	err = Flock(f, time.Second)
	assert.NotEqual(t, nil, err)

	assert.Equal(t, nil, Funlock(f))

	syscall.Unlink(filepath)
}
