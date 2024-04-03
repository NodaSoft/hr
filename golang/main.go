package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Task struct {
	Id         int
	CreatedAt  time.Time
	FinishedAt time.Time
	Result     string
	Error      error
}

const (
	workersAmount = 5
)

var (
	ErrorCreate  = errors.New("create task error")
	ErrorProcess = errors.New("something went wrong while processing")
)

func main() {
	task := &Task{}
	tasks := make(chan Task, 10)
	results := make(chan Task, 10)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go taskProducer(ctx, tasks, &wg)

	for i := 0; i < workersAmount; i++ {
		wg.Add(1)
		go taskWorker(ctx, tasks, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	if err := task.printResults(results); err != nil {
		fmt.Printf("Error while printing results: %v\n", err)
	}
}

func taskProducer(ctx context.Context, tasks chan<- Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for id := 1; ; id++ {
		select {
		case <-ctx.Done():
			close(tasks)
			return
		default:
			var err error
			if time.Now().Nanosecond()%2 > 0 {
				err = ErrorCreate
			}

			tasks <- Task{
				Id:        id,
				CreatedAt: time.Now(),
				Error:     err,
			}
		}
	}
}

func taskWorker(ctx context.Context, tasks <-chan Task, results chan<- Task, wg *sync.WaitGroup) {
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
		t.Error = ErrorProcess
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

	doneTasks := make(chan Task)
	undoneTasks := make(chan Task)

	go func() {
		defer wg.Done()

		for task := range results {
			if task.Error != nil {
				undoneTasks <- task
			} else {
				doneTasks <- task
			}
		}

		close(doneTasks)
		close(undoneTasks)
	}()

	go func() {
		defer wg.Done()
		for task := range doneTasks {
			fmt.Printf("Task with id=%d is done\n", task.Id)
		}
	}()

	go func() {
		defer wg.Done()
		for task := range undoneTasks {
			fmt.Printf("Task with id=%d is undone\n", task.Id)
		}
	}()

	wg.Wait()
	return nil
}
