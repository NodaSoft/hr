package config

import (
	"flag"
	"time"
)

var LogToStdout = flag.Bool("v", false, "Should app write log to stdout?")
var TaskTimeout = flag.Int("tt", 0, "Time between each task emission (in milliseconds)")
var WorkerTimeout = flag.Int("wt", 500, "Time worker sleeps after processing a task (in milliseconds)")
var WorkerLimit = flag.Int("w", 2, "Maximum simultaneous workers count")
var PrintResults = flag.Bool("pr", true, "Should app print results?")
var RunTime = flag.Int("rt", 3, "Time for app to work (seconds)")

func Init() {
	flag.Parse()
}

// GetTaskTimeout
// отдаёт отформатированный таймаут для эмиттера тасок
func GetTaskTimeout() *time.Duration {
	if *TaskTimeout > 0 {
		tempTime := time.Millisecond * time.Duration(*TaskTimeout)
		return &tempTime
	} else {
		return nil
	}
}

// GetWorkerTimeout
// отдаёт отформатированный таймаут для воркеров
func GetWorkerTimeout() *time.Duration {
	if *WorkerTimeout > 0 {
		tempTime := time.Millisecond * time.Duration(*WorkerTimeout)
		return &tempTime
	} else {
		return nil
	}
}
