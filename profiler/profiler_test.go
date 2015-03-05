package profiler

func ExampleStart() {
	// start a simple CPU profile and register
	// a defer to Stop (flush) the profiling data.
	defer Start(CPUProfile).Stop()
}
