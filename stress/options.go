package stress

import (
	"flag"
	"log"
	"os"
)

var flags struct {
	round int
	c1    int
	c2    int
	step  int
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	flag.IntVar(&flags.round, "round", 1, "Number of stress test rounds to run")
	flag.IntVar(&flags.c1, "c1", 1, "Number of low concurrency")
	flag.IntVar(&flags.c2, "c2", 1000, "Number of high concurrency")
	flag.IntVar(&flags.step, "step", 20, "Concurrency step between each round")
}
