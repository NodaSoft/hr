package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

// A Task represents a meaninglessness of our life
type Task struct {
	id         uuid.UUID
	createTime string // время создания
	finishTime string // время выполнения
	status     Status
}

type Status uint8

const (
	Success = iota
	Error
)

func NewTask() Task {
	return Task{
		id:         uuid.New(),
		createTime: time.Now().Format(time.RFC3339),
		status:     generateStatus(),
	}
}

func generateStatus() Status {
	if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
		return Error
	}
	return Success
}

func worker(wg *sync.WaitGroup, tasksChan, done chan Task, fail chan error) {
	defer wg.Done()

	for t := range tasksChan {
		switch t.status {
		case Success:
			done <- t
		case Error:
			fail <- fmt.Errorf("task id %s time %s, error: something went wrong", t.id.String(), t.createTime)
		}
		t.finishTime = time.Now().Format(time.RFC3339)
	}
}

func taskCreator(t chan Task, numTasks int) {
	go func() {
		defer close(t)

		for i := 0; i < numTasks; i++ {
			t <- NewTask() // передаем таск на выполнение
		}
	}()
}

func main() {
	numTasks := 10
	numWorkers := 5
	tasksChan := make(chan Task, numTasks)

	go taskCreator(tasksChan, numTasks)

	doneTasks := make(chan Task, numTasks)
	failTasks := make(chan error, numTasks)

	wg := &sync.WaitGroup{}
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker(wg, tasksChan, doneTasks, failTasks)
	}
	wg.Wait()

	close(failTasks)
	close(doneTasks)

	println("Errors:")
	for t := range failTasks {
		println(t.Error())
	}

	println("Done tasks:")
	for t := range doneTasks {
		println(t.id.String())
	}
}
