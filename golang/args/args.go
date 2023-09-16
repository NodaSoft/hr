package args

import (
	"flag"
	"os"
	"runtime"
	"time"

	"rsc.io/getopt"
)

// Must be executed at the top of main function
func Get() (concurrency int, timeToRun time.Duration) {
	var help bool

	flag.IntVar(&concurrency, "concurrency", runtime.NumCPU(), "Number of workers created. Default to runtime.NumCPU()")
	flag.DurationVar(&timeToRun, "time", 2*time.Second, "Time to run process.")
	flag.BoolVar(&help, "help", false, "Prints available arguments")

	getopt.Aliases(
		"c", "concurrency",
		"t", "time",
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
