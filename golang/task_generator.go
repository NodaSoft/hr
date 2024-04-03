package main

import (
	"context"
	"log"
	"time"
)

var (
	retryMaxAttempts = 3
	retryInterval    = time.Millisecond * 100
)

type TaskGenerator struct {
	push Pusher
}

func NewTaskGenerator(push Pusher) *TaskGenerator {
	return &TaskGenerator{push: push}
}

func (gen *TaskGenerator) generate() *Task {
	createdAt := time.Now()
	if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
		createdAt = time.Time{} // что-то пошло не так
	}

	return &Task{
		id:        time.Now().Unix(),
		createdAt: createdAt,
	}
}

func (gen *TaskGenerator) pushhWithRetry(task *Task) {
	for attempt := 0; attempt < retryMaxAttempts; attempt++ {
		if err := gen.push(task); err != nil {
			log.Printf("failed to publish task %d (%d/%d)\n", task.id, attempt+1, retryMaxAttempts)

			if attempt+1 < retryMaxAttempts {
				time.Sleep(retryInterval)
			}
		} else {
			return
		}
	}

	log.Fatalf("reached max attempts while trying to publish task %d", task.id)
}

func (gen *TaskGenerator) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			task := gen.generate()
			gen.pushhWithRetry(task)
		}
	}
}
