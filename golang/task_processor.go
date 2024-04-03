package main

import (
	"context"
)

type TaskProcessor struct {
	pull    Puller
	process Processor
	push    Pusher
}

func (processor *TaskProcessor) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			task := processor.pull()
			processor.process(task)
			processor.push(task)
		}
	}
}
