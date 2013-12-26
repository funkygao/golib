package golib

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

func lockInstance(lockfileName string) {
	pid := fmt.Sprintf("%d", os.Getpid())
	if err := ioutil.WriteFile(lockfileName, []byte(pid), 0644); err != nil {
		panic(err)
	}
}

func unlockInstance(lockfileName string) {
	syscall.Unlink(lockfileName)
}
