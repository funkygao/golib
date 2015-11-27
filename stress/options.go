package stress

import (
	"flag"
	"log"
	"os"
	//"runtime"
)

var flags struct {
	round int
	c1    int
	c2    int
	step  int
	tick  int64
	neat  bool
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	flag.IntVar(&flags.round, "round", 1, "Number of stress test rounds to run")
	flag.IntVar(&flags.c1, "c1", 1, "Number of low concurrency")
	flag.IntVar(&flags.c2, "c2", 1000, "Number of high concurrency")
	flag.IntVar(&flags.step, "step", 20, "Concurrency step between each round")
	flag.Int64Var(&flags.tick, "tick", 2, "Console stats runner ticker in seconds")
	flag.BoolVar(&flags.neat, "neat", false, "Display in neat mode, with less output")

	//runtime.GOMAXPROCS(runtime.NumCPU())
}
