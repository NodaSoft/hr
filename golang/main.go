package main

import (
	"context"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновлённый код отправить через merge-request.

// приложение эмулирует получение и обработку тасков: пытается одновременно получать и обрабатывать их в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнения остальных тасков

var (
	bufferSize = 1024
)

func main() {
	unprocessedTasks := make(chan *Task, bufferSize)
	processedTasks := make(chan *Task, bufferSize)

	taskGenerator := NewTaskGenerator(func(task *Task) error {
		unprocessedTasks <- task
		return nil
	})

	taskWorker := NewTaskWorker(func() *Task {
		task := <-unprocessedTasks
		return task
	}, func(task *Task) error {
		processedTasks <- task
		return nil
	})

	resultCollector := NewResultCollector(func() *Task {
		task := <-processedTasks
		return task
	})

	ctx, cancelContext := context.WithCancel(context.Background())

	go taskGenerator.Run(ctx)
	go taskWorker.Run(ctx)
	go resultCollector.Run(ctx)

	time.Sleep(time.Second * 3)

	cancelContext()

	resultCollector.printErrors()
	resultCollector.printDone()
}
