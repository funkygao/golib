package terminal

import (
	"syscall"
	"unsafe"
)

var TIOCGWINSZ int

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func getWinsize() *winsize {
	ws := new(winsize)

	syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)),
	)

	return ws
}

func TerminalWidth() int {
	winsize := getWinsize()
	return int(winsize.Col)
}

func TerminalHeight() int {
	winsize := getWinsize()
	return int(winsize.Row)
}
