package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// processTasks read tasks from taskChan and process them concurrently.
func processTasks(
	ctx context.Context,
	taskChan <-chan *Task,
	parallelRun int,
) (
	processed []*Task,
	errs []error,
) {

	mx := new(sync.Mutex)
	semaphore := make(chan struct{}, parallelRun)

	for task := range taskChan {

		semaphore <- struct{}{}

		go func(task *Task) {

			defer func() {
				select {
				case <-semaphore:
				default:
				}
			}()

			err := processOneTask(ctx, task)

			// store processed tasks
			mx.Lock()
			defer mx.Unlock()

			if err != nil {
				err = fmt.Errorf(
					"task_id: %3d, created_at: %25s, error: %s",
					task.ID, task.CreatedAt, err,
				)
				errs = append(errs, err)
				return
			}

			processed = append(processed, task)

		}(task)
	}

	return
}

// one task processing.
func processOneTask(
	ctx context.Context,
	task *Task,
) error {

	// skip task processing in case ctx done
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(taskProcessingLatency):
	}

	task.FinishedAt = time.Now().Format(time.RFC3339Nano)

	taskCreatedAt, err := time.Parse(time.RFC3339, task.CreatedAt)
	if err != nil {
		return fmt.Errorf("task time parsing err: %s", err)
	}

	if time.Since(taskCreatedAt) > taskProcessingTimeout {
		return ErrSomethingWentWrong
	}

	task.Result = "task was completed successfully"

	return nil
}
