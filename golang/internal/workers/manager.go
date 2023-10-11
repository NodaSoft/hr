package workers

import (
	"fmt"
	"task_service/internal/domain"
)

type WorkerManager struct {
	superCh      <-chan domain.Task
	taskCh       chan<- domain.Task
	failedTaskCh chan<- domain.Task
	//pool      map[string]domain.Worker
	//mu *sync.Mutex
}

func NewWorkerManager(superCh, taskCh, failedTaskCh chan domain.Task) WorkerManager {
	return WorkerManager{
		superCh:      superCh,
		taskCh:       taskCh,
		failedTaskCh: failedTaskCh,
	}
}

func (m WorkerManager) Run(worker domain.Worker) {
	fmt.Println("start worker")
	for task := range m.superCh {
		task, err := worker.Handle(task)
		if err != nil {
			m.failedTaskCh <- task
			continue
		}
		m.taskCh <- task
	}
}
