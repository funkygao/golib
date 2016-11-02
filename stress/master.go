package stress

import (
	"fmt"
	"log"
	"net/http"
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

	if Flags.MasterAddr == "" {
		s := new(ReportService)
		rpc.Register(s)
		rpc.HandleHTTP()
		log.Println("Master report server ready on :10093")
		go http.ListenAndServe(":10093", nil) // TODO
	}
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
		log.Println(err)
	}
	var result ReportResult
	client.Call("ReportService.ReportStat", arg, &result)
	log.Printf("told %s {C:%-5d G:%-5d %+v}", Flags.MasterAddr, arg.C, arg.G, arg.Counter)
}
