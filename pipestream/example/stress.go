/*
stress test invoker
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/funkygao/golib/pipestream"
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type strslice []string

func (this *strslice) String() string {
	return ""
}

func (this *strslice) Set(value string) error {
	*this = append(*this, value)
	return nil
}

// globals
var (
	options struct {
		cmd  string
		args strslice

		n int
		c int
	}

	startedAt  time.Time
	succ, fail int64
	wg         sync.WaitGroup
)

func main() {
	startedAt = time.Now()

	log.Println("started")
	for i := 0; i < options.c; i++ {
		wg.Add(1)
		go playSession(i, options.n)
	}

	wg.Wait()
	showReport()
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	flag.IntVar(&options.n, "n", 100, "number of test rounds per session")
	flag.IntVar(&options.c, "c", 200, "number for concurrent game sessions")
	flag.StringVar(&options.cmd, "cmd", "", "run command")
	flag.Var(&options.args, "arg", "run command arg, multiple values accepted")

	flag.Parse()
	if options.cmd == "" {
		fmt.Fprintf(os.Stderr, "run command required\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

}

func playSession(seq int, rounds int) {
	defer func() {
		wg.Done()
		log.Printf("session[%d.%d] done", seq, rounds)
	}()

	for i := 0; i < rounds; i++ {
		testcase := pipestream.New(options.cmd, options.args...)
		err := testcase.Open()
		if err != nil {
			log.Printf("[%d.%d.%d] %s", seq, rounds, i, err.Error())
			continue
		}

		log.Printf("starting game session[%d.%d.%d]", seq, rounds, i)

		scanner := bufio.NewScanner(testcase.Reader())
		scanner.Split(bufio.ScanLines)
		var lastLine string
		for scanner.Scan() {
			lastLine = scanner.Text()
		}
		err = scanner.Err()
		if err != nil {
			log.Printf("[%d.%d.%d]: %s", seq, rounds, i, err)
		}
		testcase.Close()

		if strings.HasPrefix(lastLine, "OK") {
			atomic.AddInt64(&succ, 1)
		} else {
			atomic.AddInt64(&fail, 1)
		}

		log.Printf("finished game session[%d.%d.%d] %s", seq, rounds, i, lastLine)
	}

}

func showReport() {
	fmt.Println()
	fmt.Printf("elapsed: %s\n", time.Since(startedAt))
	fmt.Printf("total:%d, success:%d, fail:%d\n", succ+fail, succ, fail)
}
