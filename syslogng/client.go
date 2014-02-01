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

func connIfNeccessary() {
    if conn != nil {
        return
    }
	var err error
	conn, err = net.Dial("unix", SYSLOGNG_SOCK)
	if err != nil {
		panic(err)
	}

}

func Printf(format string, args ...interface{}) (n int, err error) {
    connIfNeccessary()
	return fmt.Fprintf(conn, format, args...)
}

func Println(args ...interface{}) (n int, err error) {
    connIfNeccessary()
	return fmt.Fprintln(conn, args...)
}
