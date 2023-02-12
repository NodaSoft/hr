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

// A Task represents a meaninglessness of our life
type Task struct {
	id         int64
	createdAt  time.Time // время создания
	executedAt time.Time // время выполнения
	isFail     bool
	result     string
}

func taskCreator(ctx context.Context, taskChan chan<- Task) {
	for {
		now := time.Now()

		isFail := false
		if now.Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			isFail = true
		}

		// передаем таск на выполнение
		select {
		case <-ctx.Done():
			close(taskChan)
			return
		case taskChan <- Task{
			id:        now.Unix(),
			createdAt: now,
			isFail:    isFail,
		}:
		}
	}
}

func taskWorker(wg *sync.WaitGroup, taskChan <-chan Task, doneChan chan<- Task, errChan chan<- error) {
	for task := range taskChan {
		if !task.isFail && task.createdAt.After(time.Now().Add(-20*time.Second)) {
			task.result = "task has been successed"
			doneChan <- task
		} else {
			task.result = "something went wrong"
			errChan <- fmt.Errorf("Task id %d time %s, error %s", task.id, task.createdAt, task.result)
		}

		time.Sleep(time.Millisecond * 150)
	}

	wg.Done()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	var (
		workerCount = 100
		tasks       = make(chan Task, workerCount)
		doneTasks   = make(chan Task)
		errorTasks  = make(chan error)
		wg          = sync.WaitGroup{}
		results     = make(map[int64]Task)
		errors      = make([]error, 0)
	)

	go taskCreator(ctx, tasks)

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go taskWorker(&wg, tasks, doneTasks, errorTasks)
	}

	go func() {
		go func() {
			for task := range doneTasks {
				results[task.id] = task
			}
		}()
		go func() {
			for err := range errorTasks {
				errors = append(errors, err)
			}
		}()
	}()

	time.Sleep(time.Second * 3)
	cancel()
	wg.Wait()
	close(doneTasks)
	close(errorTasks)

	println("Errors:")
	for _, err := range errors {
		println(err)
	}

	println("Done tasks:")
	for result := range results {
		println(result)
	}
}
