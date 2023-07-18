package main

import (
	"context"
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

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskResult TaskResult
}

type TaskResult string

const (
	TaskResultSuccess TaskResult = "task has been success"
	TaskResultWrong   TaskResult = "something went wrong"
)

func taskCreturer(superChan chan<- Ttype) {
	defer close(superChan)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			ft := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				ft = "Some error occurred"
			}
			superChan <- Ttype{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
		}
	}
}

func taskWorker(a Ttype) Ttype {
	tt, _ := time.Parse(time.RFC3339, a.cT)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		a.taskResult = TaskResultSuccess
	} else {
		a.taskResult = TaskResultWrong
	}
	a.fT = time.Now().Format(time.RFC3339Nano)

	return a
}

func main() {
	superChan := make(chan Ttype, 10)

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		taskCreturer(superChan)
	}()

	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	go func() {
		// получение тасков
		defer wg.Done()

		defer close(doneTasks)
		defer close(undoneTasks)

		for t := range superChan {
			t = taskWorker(t)
			if t.taskResult == TaskResultSuccess {
				doneTasks <- t
			} else {
				undoneTasks <- fmt.Errorf("Task id %d time %s, error %s ", t.id, t.cT, t.taskResult)
			}
		}
	}()

	result := make(map[int]Ttype)
	errs := make([]error, 0)

	go func() {
		defer wg.Done()
		for r := range doneTasks {
			result[r.id] = r
		}
	}()
	go func() {
		defer wg.Done()
		for r := range undoneTasks {
			errs = append(errs, r)
		}
	}()

	wg.Wait()

	println("Errors:")
	for r := range errs {
		println(r)
	}

	println("Done tasks:")
	for r := range result {
		println(r)
	}
}
