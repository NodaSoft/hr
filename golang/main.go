package main

import (
	"context"
	"os"
	"time"
)

const (
	generateTime = 10 * time.Second
	sleepTime    = 3 * time.Second
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), generateTime)
	defer cancel()

	reporter := NewTaskReporter(os.Stdout)
	reporter.StartReporting(ctx, sleepTime)

	_ = StartPipeline(ctx,
		NewTaskGenerator(GeneratorDefaultChanCapacity),
		HandleTaskPipe,
		NewRecorderPipe(reporter),
		DiscardPipe[*Task],
	)

	<-ctx.Done()
}
