// Package osutil provides operating system-specific path information,
// and other utility functions.
package osutil

import (
	"errors"
	"os"
)

// ErrNotSupported is returned by functions (like Mkfifo and Mksocket)
// when the underlying operating system or environment doesn't support
// the operation.
var ErrNotSupported = errors.New("operation not supported")

// DirExists reports whether dir exists. Errors are ignored and are
// reported as false.
func DirExists(dir string) bool {
	fi, err := os.Stat(dir)
	return err == nil && fi.IsDir()
}
