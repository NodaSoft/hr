package Services

import (
	"context"
	"fmt"
	"github.com/danyducky/go-abcp/Models"
	"time"
)

// A TaskWorker represents some unit of work in our life.
type TaskWorker struct {
	Bandwidth int // TaskWorker bandwidth.

	DoneTasks   chan *Models.Task // Contains successful tasks.
	FailedTasks chan *Models.Task // Contains failed tasks.
}

// NewTaskWorker create an instance of a TaskWorker.
func NewTaskWorker(bandwidth int) *TaskWorker {
	return &TaskWorker{
		Bandwidth: bandwidth,

		DoneTasks:   make(chan *Models.Task, bandwidth),
		FailedTasks: make(chan *Models.Task, bandwidth),
	}
}

func (worker *TaskWorker) DoWork(duration time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	var provider = NewTaskProvider(worker.Bandwidth)

	go provider.Initialize(ctx)

	for task := range provider.Tasks {
		worker.ProcessTask(task)
	}

	close(worker.DoneTasks)
	close(worker.FailedTasks)
}

func (worker *TaskWorker) ReadTasks() {
	for task := range worker.DoneTasks {
		fmt.Println("Success task: ", task.Id)
	}

	for task := range worker.FailedTasks {
		fmt.Println("Failed task: ", task.Id)
	}
}

func (worker *TaskWorker) ProcessTask(task *Models.Task) {
	task.CompletedAt = time.Now()

	switch task.Status {
	case Models.Successful:
		worker.DoneTasks <- task
	case Models.Error:
		worker.FailedTasks <- task
	default:
		panic("Not supported task status.")
	}
}
