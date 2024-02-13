package Services

import (
	"context"
	"github.com/danyducky/go-abcp/Models"
	"time"
)

// A TaskProvider provides a channel of tasks.
type TaskProvider struct {
	Tasks chan *Models.Task // Contains incoming tasks.
}

// NewTaskProvider Create an instance of TaskProvider.
func NewTaskProvider(bandwidth int) TaskProvider {
	return TaskProvider{
		Tasks: make(chan *Models.Task, bandwidth),
	}
}

func (provider *TaskProvider) Initialize(ctx context.Context) {
	for {
		select {

		case <-time.After(500 * time.Millisecond):
			provider.Tasks <- provider.getTask()

		case <-ctx.Done():
			close(provider.Tasks)
			return
		}
	}
}

func (provider *TaskProvider) getTask() *Models.Task {
	var task = Models.NewTask()

	// Let's try % 3 to emulate failed tasks.
	if task.Id%3 > 0 {
		task.Update(Models.Error, "Something went wrong.")
	} else {
		task.Update(Models.Successful, "Task created successfully.")
	}

	return task
}
