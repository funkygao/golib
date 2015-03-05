package server

import (
	"fmt"
	"os"
	"runtime"
)

var (
	BuildId = "unknown"
	Version = "unknown"
)

func ShowVersionAndExit() {
	fmt.Fprintf(os.Stderr, "%s %s (build: %s)\n", os.Args[0], Version, BuildId)
	fmt.Fprintf(os.Stderr, "Built with %s %s for %s %s\n", runtime.Compiler,
		runtime.Version(), runtime.GOOS, runtime.GOARCH)
	os.Exit(0)
}
