package main

import "errors"

// Let models be public

// Task represents a task (why not)
type Task struct {
	ID         int    // id, скорее всего тут должен быть string или uuid
	CreatedAt  string // время создания
	FinishedAt string // время завершения выполнения
	Result     string // результат выполнения
}

// ErrSomethingWentWrong represents some task processing error
var ErrSomethingWentWrong = errors.New("something went wrong")
