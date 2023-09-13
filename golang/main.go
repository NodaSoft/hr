package main

import (
	"context"
	"time"
)

// app configuration
var (
	applicationDuration   = 3 * time.Second        // duration of application run
	queueCapacity         = 10                     // producer queue capacity
	parallelRun           = 2                      // max number of simultaneously processed tasks
	taskProcessingTimeout = 20 * time.Second       // threshold value for detect task processing error
	taskProcessingLatency = 150 * time.Millisecond // duration of task processing function run
)

// main func must be as short as possible
func main() {

	// make context
	ctx, cancel := context.WithTimeout(context.Background(), applicationDuration)
	defer cancel()

	// and run the app
	runTheApp(ctx)
}

// build and run your app (with context)
func runTheApp(ctx context.Context) {

	// make task producer
	taskChan := taskProducer(ctx, queueCapacity)

	// process all tasks
	processedTasks, errs := processTasks(ctx, taskChan, parallelRun)

	// show results
	handleProcessingResults(processedTasks, errs)
}
