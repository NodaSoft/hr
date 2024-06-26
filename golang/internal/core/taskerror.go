package core

import (
	"fmt"
	"strings"
)

type TaskError struct {
	msg  string
	task Task
	err  error
}

func (t *TaskError) Error() string {
	if t.msg == "" {
		return fmt.Sprintf("Task id: %d: Error: %s ", t.task.Id, t.err.Error())
	}
	return fmt.Sprintf("Task id: %d: Error: %s \tHint: %s", t.task.Id, t.err.Error(), t.msg)
}
func NewTaskError(err error, t Task, hints ...string) *TaskError {
	if len(hints) == 0 {
		return &TaskError{task: t, err: err}
	}
	return &TaskError{task: t, err: err, msg: strings.Join(hints, " ")}
}
