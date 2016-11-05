package stress

import (
	"fmt"
	"sync"
)

var (
	concurrency int32
	instanceId  string

	counterMutex   sync.RWMutex
	defaultCounter Counter = make(Counter)

	// master only
	activeSlaves int32
	globalStats  = make(map[string]*ReportArg)
)

const (
	MasterPort = 10093
)

func MasterAddr(ip string) string {
	return fmt.Sprintf("%s:%d", ip, MasterPort)
}
