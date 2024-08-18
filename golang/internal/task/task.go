package task

import (
	"fmt"
	"time"
)

type Result struct {
	msg string
	err error
}

func (r Result) Error() error {
	return r.err
}

type Task struct {
	Id       int
	InputErr bool
	Created  time.Time     // время создания
	Duration time.Duration // время выполнения
	Result   Result
}

func NewTask() *Task {
	// if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков - c таким условием ошибок не появляется

	isErrInput := time.Now().UnixMicro()%2 > 0
	return &Task{
		Created:  time.Now(),
		Id:       int(time.Now().Nanosecond()),
		InputErr: isErrInput,
	}
}

func (t *Task) String() string {
	if t.Result.err != nil {
		return fmt.Sprintf("task id: %d, error: %s, created: %s, duration: %s;\n",
			t.Id, t.Result.err.Error(), t.Created, t.Duration)
	}
	return fmt.Sprintf("task id: %d, result: %s, created: %s, duration: %s;\n",
		t.Id, t.Result.msg, t.Created, t.Duration)
}
