package main

import (
	"sync"
	"time"

	"github.com/NodaSoft/hr/args"
	"github.com/NodaSoft/hr/jobs"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

// go run . --help
func main() {
	concurrency, capacity, waitDuration := args.Get()

	var workerWg, loggerWg sync.WaitGroup

	taskChannel := make(chan *jobs.Task, capacity)
	resultChannel := make(chan *jobs.TaskResult, capacity)

	sp := jobs.NewTaskSpawner()
	sp.Start(taskChannel)

	for i := 0; i < concurrency; i++ {
		jobs.StartWorker(taskChannel, resultChannel, &workerWg)
		jobs.StartLogger(resultChannel, &loggerWg)
	}

	time.Sleep(waitDuration)

	// по порядку закрываем каналы и останавливаем воркеров
	sp.Stop()

	close(taskChannel)
	workerWg.Wait()

	close(resultChannel)
	loggerWg.Wait()
}
