package dlog

import (
	"github.com/funkygao/golib/syslogng"
	"time"
)

// Writes to syslog-ng for a logging event.
//
// ident is used to persist you different log events to different files.
// tag is can be used for your own purpose.
func Dlog(ident string, tag string, jsonStr string) (err error) {
	_, err = syslogng.Printf(":%s,%s,%d,%s\n",
		ident,
		tag,
		time.Now().UTC().Unix(),
		jsonStr)
	return
}

// Specifies local syslog-ng unix domain socket path.
func SetSyslogngSocketPath(sock string) {
	syslogng.SetSocketPath(sock)
}
