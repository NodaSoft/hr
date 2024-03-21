package main

import (
	"errors"
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
// В конце должно выводить успешные таски и ошибки выполнения остальных тасков

const (
	TasksCount   = 10
	WorkersCount = 5
	TimeFormat   = "02.01.2006 15:04:05"
)

type Task struct {
	id             int
	creationTime   time.Time     // время создания
	completionTime time.Duration // время выполнения
	errorOccured   error
	taskResult     string
}

func main() {
	tasksChan := make(chan Task, 10)

	go func() {
		for i := 1; i <= TasksCount; i++ {
			creationTime := time.Now()
			var err error
			if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				err = errors.New("creation: some error occured")
			}
			tasksChan <- Task{id: i, creationTime: creationTime, errorOccured: err} // передаем таск на выполнение
		}
		close(tasksChan)
	}()

	taskWorker := func(t Task) Task {
		if !t.creationTime.After(time.Now().Add(-20 * time.Second)) {
			t.errorOccured = errors.New("processing: some error occured")
			return t
		}

		time.Sleep(time.Millisecond * 150) // что-то выполняется
		t.taskResult = "task has been successed"
		t.completionTime = time.Now().Sub(t.creationTime)
		return t
	}

	doneResults := []string{}
	failedErrors := []error{}
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}

	for i := 1; i <= WorkersCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// получение тасков
			for t := range tasksChan {
				t = taskWorker(t)
				mutex.Lock()
				if t.errorOccured == nil {
					doneResults = append(doneResults, fmt.Sprintf("Task ID: %d Time: %s Result: \"%s\" CompletionTime: %d milliseconds", t.id, t.creationTime.Format(TimeFormat), t.taskResult, t.completionTime.Milliseconds()))
				} else {
					failedErrors = append(failedErrors, fmt.Errorf("Task ID: %d Time: %s Error: \"%s\" Result: \"%s\"", t.id, t.creationTime.Format(TimeFormat), t.errorOccured.Error(), t.taskResult))
				}
				mutex.Unlock()
			}
		}()
	}

	wg.Wait()

	println("Done tasks:")
	for _, r := range doneResults {
		println(r)
	}

	println("Errors:")
	for _, e := range failedErrors {
		println(e.Error())
	}
}
