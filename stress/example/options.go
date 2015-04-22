package main

import (
	"flag"
)

var options struct {
	server         string
	jid            string
	passwd         string
	msgSize        int
	to             string
	loopPerSession int
	tick           int64
}

func init() {
	flag.StringVar(&options.server, "server", "10.77.140.98:10222", "ejabberd server host:port")
	flag.StringVar(&options.jid, "jid", "admin@localhost", "jid")
	flag.StringVar(&options.passwd, "pass", "password", "password")
	flag.StringVar(&options.to, "to", "admin@localhost", "sent messages to whom")
	flag.IntVar(&options.msgSize, "size", 100, "each message size")
	flag.Int64Var(&options.tick, "tick", 2, "ticker of console stats in seconds")
	flag.IntVar(&options.loopPerSession, "loop", 1000000, "number of messages to send per session")
	flag.Parse()

}
