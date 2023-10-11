package receivers

import (
	"fmt"
	"task_service/internal/domain"
)

type TaskReceiver struct {
	taskCh    chan domain.Task
	errorCh   chan domain.Task
	taskRepo  domain.TaskRepository
	errorRepo domain.ErrorRepository
}

func New(taskCh, errorCh chan domain.Task, taskRepo domain.TaskRepository, errorRepo domain.ErrorRepository) TaskReceiver {
	return TaskReceiver{
		taskCh:    taskCh,
		errorCh:   errorCh,
		taskRepo:  taskRepo,
		errorRepo: errorRepo,
	}
}

func (r TaskReceiver) Run() {
	for {
		select {
		case task := <-r.taskCh:
			r.taskRepo.Add(task)
		case errorTask := <-r.errorCh:
			err := fmt.Errorf("task id %s time %s, error %s", errorTask.ID, errorTask.CreationTime, errorTask.Result)
			r.errorRepo.Add(errorTask.ID, err)
		}
	}
}
