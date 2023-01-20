package task

import (
	"fmt"
	"time"
)

const (
	executionTimeout         = 20 * time.Second
	taskCreationTimeErrorMsg = "task creation time error occurred"
	taskExecutionTimeoutMsg  = "task execution timeout"
	taskSucceededMsg         = "task has been succeeded"
)

type Task struct {
	ID         int64
	CreatedAt  time.Time
	FinishedAt time.Time
	Result     string
	Successful bool
}

func New() *Task {
	createdAt := time.Now()

	task := &Task{
		ID:        createdAt.Unix(),
		CreatedAt: createdAt,
	}

	return task
}

func (t *Task) Execute() {
	defer func() {
		t.FinishedAt = time.Now()
	}()

	if t.CreatedAt.Nanosecond()%2 > 0 {
		t.Result = taskCreationTimeErrorMsg
		return
	}

	if time.Since(t.CreatedAt) > executionTimeout {
		t.Result = taskExecutionTimeoutMsg
		return
	}

	t.Result = taskSucceededMsg
	t.Successful = true
}

func (t *Task) String() string {
	return fmt.Sprintf("Task ID: %d, time: %s, result: %s", t.ID, t.CreatedAt.Format(time.RFC3339), t.Result)
}
