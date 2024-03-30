package main

import (
	"context"
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
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

type Task struct {
	ID         uint64
	CreatedAt  time.Time
	FinishedAt time.Time
	Result     string
	Error      error
}

const layout = "2006-01-02 15:04:05"

var (
	TASK_CREATE_ERROR  = errors.New("create task error")
	TASK_PROCESS_ERROR = errors.New("something went wrong while processing")
)

func main() {
	task := &Task{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	tasks := make(chan Task, 10)
	results := make(chan Task, 10)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go task.taskProducer(ctx, tasks, &wg)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go task.taskWorker(ctx, tasks, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	if err := task.printResults(results); err != nil {
		fmt.Printf("Error printing results: %v\n", err)
	}
}

func (t *Task) taskProducer(ctx context.Context, tasks chan<- Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for id := 1; ; id++ {
		select {
		case <-ctx.Done():
			time.Sleep(time.Second)
			return
		case <-time.After(10 * time.Millisecond):
			var err error
			if time.Now().Nanosecond()%2 > 0 {
				err = TASK_CREATE_ERROR
			}
			tasks <- Task{
				ID:        uint64(id),
				CreatedAt: time.Now(),
				Error:     err,
			}
		}
	}
}

func (t *Task) taskWorker(ctx context.Context, tasks <-chan Task, results chan<- Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-tasks:
			if !ok {
				return
			}
			task.processTask()
			results <- task
		}
	}
}

func (t *Task) processTask() {
	time.Sleep(500 * time.Millisecond)
	if !t.CreatedAt.After(time.Now().Add(-time.Second)) {
		t.Error = TASK_PROCESS_ERROR
		t.Result = ""
	} else {
		t.Error = nil
		t.Result = "task has been succeeded"
	}
	t.FinishedAt = time.Now()
}

func (t *Task) printResults(results <-chan Task) error {
	var wg sync.WaitGroup
	wg.Add(3)

	successfulTasks := make(chan Task, 10)
	failedTasks := make(chan Task, 10)

	go func() {
		defer wg.Done()
		for task := range results {
			if task.Error != nil {
				failedTasks <- task
			} else {
				successfulTasks <- task
			}
		}
		close(successfulTasks)
		close(failedTasks)
	}()

	go func() {
		defer wg.Done()
		for task := range successfulTasks {
			fmt.Printf("Task ID: %d started at %s, finished at: %s, result: %s\n", task.ID, task.CreatedAt.Format(layout), task.FinishedAt.Format(layout), task.Result)
		}
	}()

	go func() {
		defer wg.Done()
		for task := range failedTasks {
			fmt.Printf("Task ID: %d failed with error: %s\n", task.ID, task.Error.Error())
		}
	}()

	wg.Wait()
	return nil
}
