package syslogng

import (
	"fmt"
	"net"
)

var (
	conn     net.Conn
	sockPath = "/tmp/als.sock"
)

func connIfNeccessary() error {
	if conn != nil {
		return nil
	}

	var err error
	conn, err = net.Dial("unix", sockPath)
	if err != nil {
		return err
	}

	return nil
}

func SetSocketPath(path string) {
	sockPath = path
}

func Printf(format string, args ...interface{}) (n int, err error) {
	err = connIfNeccessary()
	if err != nil {
		return 0, err
	}
	return fmt.Fprintf(conn, format, args...)
}

func Println(args ...interface{}) (n int, err error) {
	err = connIfNeccessary()
	if err != nil {
		return 0, err
	}
	return fmt.Fprintln(conn, args...)
}
