package main

import (
	"github.com/bom-d-van/xmppclient"
	"io/ioutil"
	"log"
	"strings"
	"sync/atomic"
)

func runSession(seq int) {
	atomic.AddInt64(&concurrency, 1)
	defer atomic.AddInt64(&concurrency, -1)

	conn, err := xmppclient.Dial(
		options.server,
		options.jid,
		"localhost",
		"password",
		"", //let server generate the resource
		&xmppclient.Config{Log: ioutil.Discard},
	)
	if err != nil {
		log.Printf("[%d] %s", seq, err)
		return
	}

	msgBody := strings.Repeat("X", options.msgSize)

	defer conn.Close()

	//conn.JoinMUC("49qniykfbt9@conference.localhost", "y")
	//conn.SendGroupChatMessage("49qniykfbt9@conference.localhost", "I came from the darkness")

	//conn.SendMediatedMucInvitation("enn.raven-theplant@localhost", "49qniykfbt9@conference.localhost", "noreason")
	//conn.JoinMUC("bullshit@conference.localhost", "y")
	//conn.SendMediatedMucInvitation("enn.raven-theplant@localhost", "bullshit@conference.localhost", "noreason")
	//conn.SendDirectMucInvitation("enn.raven-theplant@localhost", "bullshit@conference.localhost", "noreason")

	conn.Handler = &xmppclient.BasicHandler{}
	go conn.Listen()
	for i := 0; i < options.loopPerSession; i++ {
		if false {
			if err := conn.SignalPresence("online"); err != nil {
				log.Printf("[%d] %s", seq, err)
				return
			}
			sentMsg()
		}

		if true {
			if err := conn.Send(options.to, msgBody); err != nil {
				log.Printf("[%d] %s", seq, err)
				return
			}

			sentMsg()
		}

		if false {
			if err := conn.RetrieveRoster(); err != nil {
				log.Printf("[%d] %s", seq, err)
				return
			}

			recvMsg()
		}
	}

}
