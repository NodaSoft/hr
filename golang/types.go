package main

import "time"

type Puller func() *Task

type Pusher func(task *Task) error

type Processor func(task *Task)

type Task struct {
	id int64

	createdAt  time.Time  // время создания
	finishedAt *time.Time // время выполнения

	err error
}
