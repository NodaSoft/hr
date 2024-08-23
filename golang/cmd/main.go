package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"task/internal/report"
	"task/internal/task"
	"time"
)

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// Приложение должно генерировать таски 10 сек. Каждые 3 секунды должно выводить в консоль результат всех обработанных к этому моменту тасков (отдельно успешные и отдельно с ошибками).

// ЗАДАНИЕ: сделать из плохого кода хороший и рабочий - as best as you can.
// Важно сохранить логику появления ошибочных тасков.
// Важно оставить асинхронные генерацию и обработку тасков.
// Сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через pull-request в github
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

// A Task represents a meaninglessness of our life
const (
	ttl           = 5 * time.Second
	printInterval = 1 * time.Second
)

func main() {

	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), ttl)
	defer func() {
		timeoutCancel()
		log.Println("main wg done")

	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		timeoutCancel()
	}()

	pendingQueue := make(chan *task.Task)
	completedQueue := make(chan *task.Task)
	doneTasks := make(chan *task.Task)
	errTasks := make(chan *task.Task)

	b := report.NewBuilder()

	wg := &sync.WaitGroup{}
	wg.Add(5)

	go task.Generate(timeoutCtx, wg, pendingQueue)
	go task.Process(wg, pendingQueue, completedQueue)
	go task.Sort(wg, completedQueue, doneTasks, errTasks)

	buildCtx, buildStop := context.WithCancel(context.Background())

	go report.BuildReport(wg, b, doneTasks, errTasks, buildStop)
	go report.SchedulePrint(buildCtx, wg, b, printInterval)

	wg.Wait()
}
