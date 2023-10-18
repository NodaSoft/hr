package workers

import (
	"errors"
	"time"

	"github.com/NodaSoft/hr/golang/internal/tasks"
)

func (p *Pool) launchWorker(task *tasks.Task) {
	defer func() { <-p.workers }()

	time.Sleep(time.Millisecond * 150) // workload simulation

	createdAt, _ := time.Parse(time.RFC3339, task.CreatedAt()) // this is part of the task (error simulation)

	if createdAt.After(time.Now().Add(-20 * time.Second)) {
		task.MarkAsCompleted(nil)
	} else {
		task.MarkAsCompleted(errors.New("something went wrong"))
	}

	go p.submitProcessed(task)
}

func (p *Pool) submitProcessed(task *tasks.Task) {
	p.processed <- task
}
