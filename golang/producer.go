package main

import (
	"context"
	"time"
)

// taskProducer produce the tasks.
func taskProducer(
	ctx context.Context,
	queueCap int,
) <-chan *Task {

	out := make(chan *Task, queueCap)

	cnt := 0

	go func() {
	productionLoop:
		for {
			cnt++
			select {
			case <-ctx.Done():
				break productionLoop
			case out <- makeNewTask(cnt):
			}
		}
		close(out)
	}()

	return out
}

// makeNewTask returns one new task.
func makeNewTask(id int) *Task {

	task := &Task{
		ID:        id,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	// вот такое условие появления ошибочных тасков
	if time.Now().Nanosecond()/1000%2 > 0 {
		task.CreatedAt = "some error occurred"
	}

	return task
}
