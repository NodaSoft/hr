package task

import (
	"errors"
	"time"
)

// Worker does the Task and sends result to the Logger
type Worker struct {
	NewTasks       chan Task
	CompletedTasks chan Task
}

// Start Worker for parallel doing Task
func (w Worker) Start() {
	for {
		task := <-w.NewTasks
		go w.do(task)
	}
}

func (w Worker) do(task Task) {
	if task.Start.Nanosecond()%2 > 0 {
		task.ErrorMessage = errors.New("invalid start time").Error()
	}

	if task.Start.After(time.Now().Add(MaxExecutionTime)) {
		task.ErrorMessage = errors.New("task execution timeout").Error()
	}

	task.Finish = time.Now()

	w.CompletedTasks <- task
}
