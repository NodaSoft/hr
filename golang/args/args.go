package args

import (
	"flag"
	"os"
	"runtime"
	"time"

	"rsc.io/getopt"
)

// Must be executed at the top of main function
func Get() (concurrency, capacity int, timeToRun time.Duration) {
	var help bool

	flag.IntVar(&concurrency, "concurrency", runtime.NumCPU(), "Number of workers created. Default to runtime.NumCPU()")
	flag.IntVar(&capacity, "queue", 100, "Capacity of task and result channels.")
	flag.DurationVar(&timeToRun, "time", 2*time.Second, "Time to run process.")
	flag.BoolVar(&help, "help", false, "Prints available arguments")

	getopt.Aliases(
		"c", "concurrency",
		"t", "time",
		"q", "queue",
		"h", "help",
	)
	getopt.Parse()

	if help {
		println("List of available arguments:")
		getopt.PrintDefaults()
		os.Exit(0)
	}

	return
}
