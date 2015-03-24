package locking

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"
	"time"
)

var (
	ErrTimeout = errors.New("timeout")
)

func InstanceLocked(lockfileName string) bool {
	_, err := os.Stat(lockfileName)
	return err == nil
}

func LockInstance(lockfileName string) {
	pid := fmt.Sprintf("%d", os.Getpid())
	if err := ioutil.WriteFile(lockfileName, []byte(pid), 0644); err != nil {
		panic(err)
	}
}

func UnlockInstance(lockfileName string) {
	syscall.Unlink(lockfileName)
}

func Flock(f *os.File, timeout time.Duration) error {
	var t time.Time
	for {
		if t.IsZero() {
			t = time.Now()
		} else if timeout > 0 && time.Since(t) > timeout {
			return ErrTimeout
		}

		err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		fmt.Println(err, int(f.Fd()))
		if err == nil {
			return nil
		} else if err != syscall.EWOULDBLOCK {
			return err
		}

		time.Sleep(50 * time.Millisecond)
	}
}

func Funlock(f *os.File) error {
	return syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
}
