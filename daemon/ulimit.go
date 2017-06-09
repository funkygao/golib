package daemon

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func mustMaxOpenFile(min int) {
	ulimitN, err := exec.Command("/bin/sh", "-c", "ulimit -n").Output()
	if err != nil {
		panic(err)
	}

	n, err := strconv.Atoi(strings.TrimSpace(string(ulimitN)))
	if err != nil || n < min {
		panic(fmt.Sprintf("ulimit too small: %d, should be at least %d", n, min))
	}
}
