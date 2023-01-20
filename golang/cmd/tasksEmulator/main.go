package main

import (
	"context"
	"time"

	"github.com/Quantum12k/hr/golang/internal/creator"
	"github.com/Quantum12k/hr/golang/internal/worker"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

const (
	appLifetime = 5 * time.Second
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), appLifetime)
	defer cancel()

	tasksCreator := creator.New(ctx)
	tasksWorker := worker.New(ctx, tasksCreator.NewTasksCh)

	failed := make([]string, 0)
	success := make(map[int64]struct{})

	for task := range tasksWorker.DoneTasksCh {
		if task.Successful {
			success[task.ID] = struct{}{}
			continue
		}

		failed = append(failed, task.String())
	}

	println("Errors:")
	for _, msg := range failed {
		println(msg)
	}

	println("Done tasks:")
	for id := range success {
		println(id)
	}
}
