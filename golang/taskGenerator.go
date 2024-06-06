package main

import (
	"context"
	"time"
)

type Status int8

const (
	New         Status = 0
	Success     Status = 1
	ErrorStatus Status = 2
)

// A Task represents a meaninglessness of our life
type Task struct {
	id          int
	createTime  time.Time
	executeTime time.Time
	status      Status
	statusInfo  string
}

func taskGenerator(ctx context.Context, generateTimeout time.Duration, newTaskCh chan<- Task) {
	defer close(newTaskCh)
	id := 0
	timer := time.NewTimer(generateTimeout)

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			return
		default:
			task := Task{
				createTime: time.Now(),
				id:         id, // так читабельнее в логах поэтому оставил.
				//id:         int(time.Now().UnixNano()), // Заменил unix на unixnano иначе генерились таски с одинаковыми id
				status: New,
			}
			id++

			// При использовании nano последние 3 знака всегда нули, и failed таски не генерятся.
			if time.Now().UnixMilli()%2 > 0 { // вот такое условие появления ошибочных тасков
				task.status = ErrorStatus
				task.statusInfo = "Failed task generated"
			}
			newTaskCh <- task
		}
	}
}
