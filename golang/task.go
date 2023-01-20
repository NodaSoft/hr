package main

import (
	"strconv"
	"time"
)

// A Task это задача которая проходит цикл исполнения
type Task struct {
	id     int64
	since  string // время создания
	till   string // время выполнения
	error  error
	result Result
}

func (t *Task) String() string {
	if t.error != nil {
		return t.error.Error()
	}

	return strconv.FormatInt(t.id, 10)
}

type Result string

const (
	Success Result = "task has been successed"
	Wrong   Result = "something went wrong"
)

func NewTask(since time.Time) *Task {
	var err error

	if since.Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
		err = SomeError
	}

	return &Task{
		id:    since.Unix(),
		since: since.Format(time.RFC3339),
		error: err,
	}
}
