package task

import (
	"fmt"
	"sync"
	"testtask/internal/config"
	"time"
)

type TaskExecutor struct {
	taskExecutorsLimit int
	doneTasksChanSize  int
}

func NewTaskExecutor(
	taskExecutorsLimit int,
	doneTasksBufferSize int,
) *TaskExecutor {
	return &TaskExecutor{
		taskExecutorsLimit: taskExecutorsLimit,
		doneTasksChanSize:  doneTasksBufferSize,
	}
}

func (te *TaskExecutor) ExecuteTasks(toDoTasks <-chan Task) <-chan Task {
	doneTasks := make(chan Task, te.doneTasksChanSize)

	go func() {
		defer close(doneTasks)

		var taskExecutionWG sync.WaitGroup
		executorsLimiter := make(chan struct{}, config.MustNew().TaskExecutorsLimit)

		for task := range toDoTasks {
			executorsLimiter <- struct{}{}
			taskExecutionWG.Add(1)
			go func(task Task) {
				doneTasks <- ExecuteTask(task)
				taskExecutionWG.Done()
				<-executorsLimiter
			}(task)
		}

		taskExecutionWG.Wait()
	}()

	return doneTasks
}

func ExecuteTask(task Task) Task {
	if task.CreationTime.Nanosecond()%2 == 0 || task.IsExpired() {
		task.Error = fmt.Errorf("something went wrong")
	} else {
		task.Result = []byte("task has been executed successfuly")
	}
	task.CompletionTime = new(time.Time)
	*task.CompletionTime = time.Now()

	time.Sleep(time.Millisecond * 150)

	return task
}
