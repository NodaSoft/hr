package worker

import (
	"sync"
	"taskConcurrency/internal/domain/task"
	"time"
)

type Worker struct{}

func (w *Worker) Work(tasks <-chan task.Task, processed chan<- task.Task) {
	go func() {
		wg := &sync.WaitGroup{}
		for task := range tasks {
            wg.Add(1)
            go func() {
                w.workOneTask(task, processed)
                wg.Done()
            }()
        }
		wg.Wait()
		close(processed)
	}()
}

func (w *Worker) workOneTask(task task.Task, processed chan<- task.Task) {
	creationTime, err := time.Parse(time.RFC3339, task.CreationTime)
	if err != nil {
		task.TaskResult = []byte("something went wrong")
		processed <- task
	}
	if creationTime.After(time.Now().Add(-20 * time.Second)) {
		task.TaskResult = []byte("task has been successed")
	} else {
		task.TaskResult = []byte("something went wrong")
	}
	task.Executiontime = time.Now().Format(time.RFC3339Nano)

	time.Sleep(150 * time.Millisecond)

	processed <- task
}
