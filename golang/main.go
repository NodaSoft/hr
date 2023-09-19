package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	timeOut     = 20 * time.Second
	idleTimeOut = 150 * time.Millisecond
)

// Task - стурктура задачи
type Task struct {
	ID       int
	CreateAt string
	FinishAt string
	Result   string
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	runApp(ctx)
}

// runApp - запуск программы
func runApp(ctx context.Context) {
	jobs := make(chan *Task, 10)
	doneTasks := make(chan *Task)
	undoneTasks := make(chan error)

	result := []*Task{}
	errors := []error{}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		worker(ctx, jobs, doneTasks, undoneTasks)
	}()

	go func() {
	exit:
		for {
			select {
			case <-ctx.Done():
				break exit
			case jobs <- makeNewTask(int(time.Now().Unix())):
			}
		}
		close(jobs)
	}()

	go func() {
		for r := range doneTasks {
			result = append(result, r)
		}
		close(doneTasks)
	}()

	go func() {
		for r := range undoneTasks {
			errors = append(errors, r)
		}
		close(undoneTasks)
	}()

	wg.Wait()
	fmt.Println("Errors:")
	for _, err := range errors {
		fmt.Println(err.Error())
	}

	fmt.Println("Done tasks:")
	for _, value := range result {
		fmt.Printf("Task id %d time %s, result: %s\n", value.ID, value.CreateAt, value.Result)
	}
}

func makeNewTask(id int) *Task {
	task := &Task{
		ID:       id,
		CreateAt: time.Now().Format(time.RFC3339Nano),
	}

	// Условие появления ошибочных задач
	if time.Now().Nanosecond()%2 > 0 {
		task.CreateAt = "Some error occured"
	}
	return task
}

func worker(ctx context.Context, jobs <-chan *Task, doneTasks chan<- *Task, undoneTasks chan<- error) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(idleTimeOut):
			return
		case task, ok := <-jobs:
			if !ok {
				return
			}
			err := processTask(ctx, task)
			if err != nil {
				undoneTasks <- err
			} else {
				doneTasks <- task
			}
		}
	}
}

func processTask(ctx context.Context, task *Task) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(idleTimeOut):
	}

	taskCreatedAt, err := time.Parse(time.RFC3339Nano, task.CreateAt)
	if err != nil {
		return fmt.Errorf("Task id %d, result: %s", task.ID, "something went wrong")
	}

	if time.Since(taskCreatedAt) > timeOut {
		return fmt.Errorf("Task id %d, result: %s", task.ID, "something went wrong")
	}
	task.Result = "task has been successed"
	task.FinishAt = time.Now().Format(time.RFC3339Nano)

	return nil
}
