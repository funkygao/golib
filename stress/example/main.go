package main

import (
	"github.com/funkygao/golib/stress"
)

// An example of using stress lib to load test xmpp server.
func main() {
	go runStats()
	stress.RunStress(runSession)
}
