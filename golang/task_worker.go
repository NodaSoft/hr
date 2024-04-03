package main

import (
	"errors"
	"time"
)

var (
	taskProcessingTime    = time.Millisecond * 150
	taskValidityThreshold = 20 * time.Second

	ErrProcessingTask = errors.New("something went wrong")
)

func NewTaskWorker(puller Puller, pusher Pusher) *TaskProcessor {
	return &TaskProcessor{
		pull:    puller,
		process: processTask,
		push:    pusher,
	}
}

func processTask(task *Task) {
	if !task.createdAt.After(time.Now().Add(-taskValidityThreshold)) {
		task.err = ErrProcessingTask
	}

	time.Sleep(taskProcessingTime)

	var currentTime = time.Now()
	task.finishedAt = &currentTime
}
