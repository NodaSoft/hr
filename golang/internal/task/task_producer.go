package task

import (
	"context"
	"time"
)

type TaskProducer struct {
	taskExpirationDuration time.Duration
}

func NewTaskProducer(taskExpirationDuration time.Duration) *TaskProducer {
	return &TaskProducer{
		taskExpirationDuration: taskExpirationDuration,
	}
}

func (tp *TaskProducer) ProduceTasks(ctx context.Context) <-chan Task {
	tasks := make(chan Task) // buffering can be useful in some cases

	go func() {
		defer close(tasks)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				tasks <- Task{
					Id:             int(time.Now().Unix()),
					CreationTime:   time.Now(),
					ExpirationTime: time.Now().Add(tp.taskExpirationDuration),
				}
				time.Sleep(time.Millisecond * 500) // only for debugging to avoid a lot of output
			}
		}
	}()

	return tasks
}
