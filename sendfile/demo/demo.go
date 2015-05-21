package main

import (
	"net"
	"os"

	"github.com/funkygao/golib/sendfile"
)

func dieIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	l, err := net.Listen("tcp", ":10111")
	dieIfError(err)

	for {
		conn, err := l.Accept()
		dieIfError(err)

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	f, err := os.Open("demo.go")
	dieIfError(err)

	sendfile.Sendfile(conn.(*net.TCPConn), f, 100)
}
