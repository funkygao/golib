package locking

import (
	"fmt"
	"io/ioutil"
	"os"
	"syscall"
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
