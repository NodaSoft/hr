package main

import (
	"context"
	"time"
)

// 1) TODO Вопрос по условию нужно ли выовдить каждый раз Все отработанные за период между выводами, или все отработанные за все время работы приложения?
// Реализовал вариант где каждый раз выводятся вообще все отработанные за все время работы приложения, без вычисление уже выведенных ранее.

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// Приложение должно генерировать таски 10 сек. Каждые 3 секунды должно выводить в консоль результат всех обработанных к этому моменту тасков (отдельно успешные и отдельно с ошибками).

// ЗАДАНИЕ: сделать из плохого кода хороший и рабочий - as best as you can.
// Важно сохранить логику появления ошибочных тасков.
// Важно оставить асинхронные генерацию и обработку тасков.
// Сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через pull-request в github
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

func main() {

	// должны получать при запуске программы например из конфига.
	generateTimeout := 10 * time.Second
	printInterval := 3 * time.Second
	numWorkers := 5
	bufferSize := 10

	newTasksCh := make(chan Task, bufferSize)
	doneTasksCh := make(chan Task, bufferSize)
	failedTasksCh := make(chan Task, bufferSize)

	// Как вариант можно заменить на stop канал и сделать его частью структур taskWorkerPool и
	ctx, cancel := context.WithCancel(context.Background())

	go taskGenerator(ctx, generateTimeout, newTasksCh)

	wp := newTaskWorkerPool(newTasksCh, doneTasksCh, failedTasksCh, numWorkers)
	wp.Start(ctx)

	to := newTaskObserver(doneTasksCh, failedTasksCh, printInterval)
	to.Start(ctx)
	to.PrintResultsPeriodically(ctx)

	wp.Stop()
	cancel()
	to.Stop()
}
