package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, time.Second*5)

	tasks := make(chan Task, 10)
	doneTasks := make(chan Task)
	errs := make(chan error)

	go func() {
		taskGenerator(ctx, tasks)
		close(tasks)
	}()

	go taskProcessorGroup(ctx, tasks, doneTasks, errs, 10)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		fmt.Fprintln(os.Stderr, "Errors:")
		for err := range errs {
			fmt.Fprintln(os.Stderr, err)
		}
		wg.Done()
	}(wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		fmt.Fprintln(os.Stdout, "Done tasks:")
		for task := range doneTasks {
			fmt.Fprintln(os.Stdout, task.id)
		}
		wg.Done()
	}(wg)
	wg.Wait()
}

func NewTask() Task {
	return Task{
		// Не очень понятно насколько надо быть рандомными
		// Может вообще лучше не int использовать
		id:        rand.Int(),
		createdAt: time.Now(),
		result:    []byte{},
	}
}

type Task struct {
	id         int
	createdAt  time.Time
	finishedAt time.Time
	result     []byte
}

func (t *Task) DoneWithResult(result []byte) {
	t.finishedAt = time.Now()
	t.result = result
}

func (t *Task) IsFresh() bool {
	return t.createdAt.After(time.Now().Add(-20 * time.Second))
}

func (t *Task) IsSuccessed() bool {
	return strings.Contains(string(t.result), "successed")
}

var (
	taskSuccessedStatus = []byte("task has been successed")
	taskFailedStatus    = []byte("something went wrong")
	taskGenerationError = []byte("Some error occured")
)

func taskGenerator(ctx context.Context, tasks chan<- Task) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			task := NewTask()
			// if time.Now().UnixMilli()%2 > 0 { <- Может лучше так ?
			if time.Now().Nanosecond()%2 > 0 {
				task.DoneWithResult(taskGenerationError)
			}

			tasks <- task
		}
	}
}

func taskProcessorGroup(
	ctx context.Context,
	tasks <-chan Task,
	doneTasks chan<- Task,
	errs chan<- error,
	workerCount int,
) {
	wg := &sync.WaitGroup{}
	for ; workerCount != 0; workerCount-- {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			taskProcessor(ctx, tasks, doneTasks, errs)
			wg.Done()
		}(wg)
	}
	wg.Wait()
	close(doneTasks)
	close(errs)
}

func taskProcessor(
	ctx context.Context,
	tasks <-chan Task,
	doneTasks chan<- Task,
	errs chan<- error,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case task := <-tasks:
			result := processTask(task)
			if result.IsSuccessed() {
				doneTasks <- task
				continue
			}

			errs <- taskProcessError(task)
		}
	}
}

func taskProcessError(task Task) error {
	return fmt.Errorf(
		"Task id %d time %s, error %s",
		task.id,
		task.createdAt.Format(time.RFC3339),
		task.result,
	)
}

func processTask(task Task) Task {
	switch task.IsFresh() {
	case true:
		task.DoneWithResult(taskSuccessedStatus)
	case false:
		task.DoneWithResult(taskFailedStatus)
	}

	// Эмуляция долгой обработки ?
	time.Sleep(time.Millisecond * 150)

	return task
}
