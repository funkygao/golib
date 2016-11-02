package stress

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"github.com/pborman/uuid"
)

type ReportArg struct {
	Id      string
	Counter Counter
	T       time.Time
	C       int32 // concurrency
	G       int   // goroutines
}

type ReportResult struct{}

type ReportService struct{}

func (r *ReportService) ReportStat(arg *ReportArg, result *ReportResult) error {
	counterMutex.Lock()
	defer counterMutex.Unlock()

	globalStats[arg.Id] = arg
	atomic.AddInt32(&activeSlaves, 1)

	return nil
}

func init() {
	host, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	instanceId = fmt.Sprintf("%s-%s", host, strings.Replace(uuid.New(), "-", "", -1))
}

func reportToMaster() {
	counterMutex.RLock()
	counterClone := make(Counter, len(defaultCounter))
	for k, v := range defaultCounter {
		counterClone[k] = v
	}
	counterMutex.RUnlock()

	arg := &ReportArg{
		Id:      instanceId,
		Counter: counterClone,
		G:       runtime.NumGoroutine(),
		C:       atomic.LoadInt32(&concurrency),
		T:       time.Now(),
	}
	var client, err = rpc.DialHTTP("tcp", Flags.MasterAddr)
	if err != nil {
		log.Printf("report to master: %v", err)
		return
	}
	var result ReportResult
	err = client.Call("ReportService.ReportStat", arg, &result)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("told %s {C:%-5d G:%-5d %+v}", Flags.MasterAddr, arg.C, arg.G, arg.Counter)
	}
}
