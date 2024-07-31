package main

import (
	"context"
	"time"
)

const (
	GeneratorDefaultChanCapacity = 10
)

// NewTaskGenerator returns a new Source that generates tasks.
func NewTaskGenerator(chanCapacity int) Source[*Task] {
	return func(ctx context.Context) <-chan *Task {
		out := make(chan *Task, chanCapacity)

		go func() {
			defer close(out)
			for {
				select {
				case <-ctx.Done():
					break
				default:
					creationTime := time.Now()
					task := NewTask(creationTime.Unix(), creationTime)
					out <- task
				}
			}
		}()

		return out
	}
}
