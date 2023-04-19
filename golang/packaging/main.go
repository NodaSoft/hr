package main

import (
	"fmt"
	"sync"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

const (
	_TotalRunTime = time.Second * 3

	_CreatorTime    = 4 * time.Millisecond
	_CreatorThreads = 3
)

func main() {
	superChan := make(chan any, 10)
	resChan := make(chan any)
	errChan := make(chan error)

	canalize := NewCanalizer(superChan, resChan, errChan)

	// Запуск эмулятора тасок в несколько потоков
	go func() {
		canalize.RunCreators(_CreatorThreads, _CreatorTime, NewTask)
	}()

	// Запуск обработчика тасок, каждая в своем потоке
	go func() {
		canalize.RunWorker(ProcessTask)
	}()

	// сбор результатов обработки.
	var results map[int]any
	var errs []error
	var wgResult sync.WaitGroup
	wgResult.Add(1)
	go func() {
		results, errs = canalize.RunResults(GuidTask)
		wgResult.Done()
	}()

	time.Sleep(_TotalRunTime)
	wgResult.Wait()

	fmt.Printf("\nErrors: %d\n", len(errs))
	for _, er := range errs {
		fmt.Println(er.Error())
	}

	fmt.Printf("\nDone tasks: %d", len(results))
	for id, res := range results {
		fmt.Printf("\n%d: %s", id, res.(*Task).taskRESULT)
	}
}
