package syslogng

import (
	"fmt"
	"net"
)

const (
	SYSLOGNG_SOCK = "/tmp/als.sock"
)

var (
	conn net.Conn
)

func init() {
	var err error
	conn, err = net.Dial("unix", SYSLOGNG_SOCK)
	if err != nil {
		panic(err)
	}

}

func Printf(format string, args ...interface{}) (n int, err error) {
	return fmt.Fprintf(conn, format, args...)
}

func Println(args ...interface{}) (n int, err error) {
	return fmt.Fprintln(conn, args...)
}
