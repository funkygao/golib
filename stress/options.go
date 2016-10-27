package stress

import (
	"flag"
	"log"
	"os"
)

var Flags struct {
	Round int
	C1    int
	C2    int
	Step  int
	Tick  int64
	Neat  bool

	MasterAddr string
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	flag.IntVar(&Flags.Round, "round", 1, "Number of stress test rounds to run")
	flag.IntVar(&Flags.C1, "c1", 1, "Number of low concurrency")
	flag.IntVar(&Flags.C2, "c2", 1000, "Number of high concurrency")
	flag.IntVar(&Flags.Step, "step", 20, "Concurrency step between each round")
	flag.Int64Var(&Flags.Tick, "tick", 2, "Console stats runner ticker in seconds")
	flag.BoolVar(&Flags.Neat, "neat", false, "Display in neat mode, with less output")
	flag.StringVar(&Flags.MasterAddr, "master", "", "Master stress address. If empty, works as master")
}
