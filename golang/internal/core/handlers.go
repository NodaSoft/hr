package core

import (
	"context"
	"fmt"
	"sync"
	log "taskhandler/internal/logger"
)

// Fill channel with tasks
// Anything that implements TaskFactory may be passed as factory
func FillChannel(ctx context.Context, tch chan Task, factory TaskFactory) chan Task {
	tch = make(chan Task, 10)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(tch)
				log.Info("FillChannel. Reason: context.Done(), channel closed")
				return
			case tch <- factory.MakeTask():
			}
		}
	}()
	return tch
}

type TaskWorker func(t Task) Task

// Passes Task to Arbitrary worker, writes modified tasks to channel
func HandleTasks(ctx context.Context, tasks chan Task, worker TaskWorker) chan Task {
	tch := make(chan Task, 10)

	go func(ctx context.Context) {
		wg := sync.WaitGroup{}
		for {
			select {
			case <-ctx.Done():
				wg.Wait()
				close(tch)
				log.Info("Pileline closed. Reason: context.Done(), channels closed")
				return
			case t := <-tasks:
				wg.Add(1)
				go func(t Task) {
					defer wg.Done()
					tch <- worker(t)
				}(t)
			}
		}
	}(ctx)
	return tch
}

// Closes both channels. Be careful with separator channel
// Writes separated tasks to channel, gets it in arguments
func SeparateBrokenTasks(ctx context.Context, tch chan Task, separator func(Task) int8, separated chan error) chan Task {
	correct := make(chan Task, 10)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(correct)
				close(separated)
				log.Info("Separator closed. Reason: context.Done(), channels closed")
				return
			case t := <-tch:
				res := separator(t)

				// In case separator logic cant handle broken task
				// We 'll give separator an option to return -1
				// and mark it in logs by __Error

				if res < 0 {
					log.Info("__Error in separator, cant handle broken task. Task id: ", t.Id)
					separated <- NewTaskError(fmt.Errorf("given separator cant handle broken task"), t)
				}

				// Common cases for separator, not exactly error
				if res > 0 {
					// Do not uncomment, will spam to the logs
					// log.Debug("Broken task")
					separated <- NewTaskError(fmt.Errorf("common broken task, separated as usual"), t)
					continue
				}

				// If separator returns 0 , it found no bugs, we'll pass it to the next stage
				correct <- t
			}
		}
	}()
	return correct
}
