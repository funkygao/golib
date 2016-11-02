package stress

import (
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
