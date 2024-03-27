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
// В конце должно выводить успешные таски и ошибки выполнения остальных тасков

var (
	TaskSuccessMessage = "task has been successed"
	TaskErrorMessage   = "something went wrong"
)

const workersCount = 5

// Task represents a task structure
type Task struct {
	ID         int
	CreatedAt  time.Time // время создания
	ExecutedAt time.Time // время выполнения
	HasError   bool
	Result     []byte
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tasks := make(chan Task)
	doneTasks := make(chan Task)
	undoneTasks := make(chan error)

	go taskCreator(ctx, tasks)

	for i := 0; i < workersCount; i++ {
		go taskWorker(ctx, tasks, doneTasks, undoneTasks)
	}

	var wg sync.WaitGroup

	results := map[int]Task{}
	errors := []error{}

	wg.Add(1)
	go func() {
		defer close(doneTasks)
		defer close(undoneTasks)
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case t, ok := <-doneTasks:
				if !ok {
					return
				}
				results[t.ID] = t
			case err, ok := <-undoneTasks:
				if !ok {
					return
				}
				errors = append(errors, err)
			}
		}
	}()
	wg.Wait()

	fmt.Println("Errors:")
	for _, err := range errors {
		fmt.Println(err)
	}

	fmt.Println("Done tasks:")
	for r := range results {
		fmt.Println(r)
	}
}

func taskCreator(ctx context.Context, out chan<- Task) {
	defer close(out)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			CreatedAt := time.Now()
			hasError := false

			if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				hasError = true
			}

			time.Sleep(50 * time.Millisecond)

			out <- Task{CreatedAt: CreatedAt, ID: int(time.Now().Unix()), HasError: hasError} // передаем таск на выполнение
		}
	}
}

func taskWorker(ctx context.Context, tasks <-chan Task, doneTasks chan<- Task, undoneTasks chan<- error) {
	for {
		select {
		case <-ctx.Done():
			return
		case t, ok := <-tasks:
			if !ok {
				return
			}
			t.ExecutedAt = time.Now()

			if t.CreatedAt.After(time.Now().Add(-20*time.Second)) && !t.HasError {
				t.Result = []byte(TaskSuccessMessage)
				t.ExecutedAt = time.Now()
				doneTasks <- t
			} else {
				t.Result = []byte(TaskErrorMessage)
				undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.ID, t.CreatedAt, t.Result)
			}

			time.Sleep(time.Millisecond * 150)
		}
	}
}
