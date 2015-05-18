// Package sendfile implements a handy zero-copy sendfile in golang.
package sendfile

import (
	"io"
	"net"
	"os"
)

// Sendfile sends count bytes from f to remote a TCP connection.
// f offset is always relative to the current offset.
func Sendfile(conn *net.TCPConn, f *os.File, count int64) (n int64, err error) {
	lr := &io.LimitedReader{N: count, R: f}
	n, err = conn.ReadFrom(lr)
	return
}
