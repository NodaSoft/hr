package task

import (
	"errors"
	"time"
)

func Worker(t *Task) *Task {
	time.Sleep(time.Millisecond * 150)

	if t.InputErr {
		t.Result.err = errors.New("task input error")
		t.Duration = 0
		return t
	}
	if t.Created.After(time.Now().Add(-20 * time.Second)) { // ?
		t.Result.msg = "task process successed"
	} else {
		t.Result.err = errors.New("something went wrong")
	}

	t.Duration = time.Since(t.Created)

	return t
}
