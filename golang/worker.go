package main

import (
	"context"
	"fmt"
	"time"
)

// HandleTask finishes the Task t. It returns non-nil error if the context is cancelled.
func HandleTask(ctx context.Context, t *Task) error {
	const op = "app.task.HandleTask"

	if t.Finished() {
		return nil
	}
	if t.IsCorrupted() {
		err := t.Error()
		if err == nil {
			err = fmt.Errorf("%s: task is corrupted but no error", op)
		}
		t.Finish(NewResult(false).WithMessages(err.Error()))
		return nil
	}

	var result Result
	deadline := time.Now().Add(time.Duration(-20) * time.Second)
	if t.CreationTime().After(deadline) {
		result = NewResult(true).WithMessages("task has been successfully finished")
	} else {
		result = NewResult(false).WithMessages("task deadline exceeded")
	}

	if err := ctx.Err(); err == nil {
		time.Sleep(150 * time.Millisecond)
		t.Finish(result)
	}

	if err := ctx.Err(); err == nil {
		return nil
	}

	return fmt.Errorf("%s: %w", op, ctx.Err())
}

// HandleTaskPipe is a Pipe that runs HandleTask function on each input Task.
func HandleTaskPipe(ctx context.Context, in <-chan *Task) <-chan *Task {
	out := make(chan *Task, cap(in))

	go func() {
		defer close(out)

		for t := range in {
			if err := ctx.Err(); err != nil {
				break
			}

			//go func() {
			err := HandleTask(ctx, t)
			if err != nil {
				return
			}
			out <- t
			//}()
		}
	}()

	return out
}
