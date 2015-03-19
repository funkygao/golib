package osutil

import (
	"os/exec"
	"runtime"
)

func OpenURL(url string) error {
	if runtime.GOOS == "windows" {
		return exec.Command("cmd.exe", "/C", "start "+url).Run()
	}

	if runtime.GOOS == "darwin" {
		return exec.Command("open", url).Run()
	}

	return exec.Command("xdg-open", url).Run()
}
