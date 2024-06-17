package app

import (
	"context"
	"taskConcurrency/internal/creator"
	"taskConcurrency/internal/domain/task"
	"taskConcurrency/internal/monitor"
	"taskConcurrency/internal/sorter"
	"taskConcurrency/internal/worker"
	"time"
)

type App struct {
}

// ctx - контекст для задания времени генерации тасков
func (a *App) Do(ctx context.Context) {
	superChan := make(chan task.Task, 10)

	creator := creator.Creator{}
	creator.Create(ctx, superChan)

	processed := make(chan task.Task)

	worker := worker.Worker{}
	worker.Work(superChan, processed)

	doneTasks := make(chan task.Task)
	undoneTasks := make(chan error)

	sorter := sorter.Sorter{}
	sorter.Sort(processed, doneTasks, undoneTasks)

	// Выводим и блокируемся
	monitor := monitor.Monitor{}
	monitor.PrintWithInterval(time.Second*3, doneTasks, undoneTasks)
}

func NewApp() *App {
	return &App{}
}
