package main

import (
	"context"
	"fmt"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
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

type Task struct {
	id          string
	createdAt   string
	processedAt string
	success     bool
	payload     []byte
}

func main() {
	context, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	tasks := taskWorker(taskGenerator(context))

	doneTasks := []*Task{}
	failedTasks := []*Task{}

	for t := range tasks {
		if t.success {
			doneTasks = append(doneTasks, t)
		} else {
			failedTasks = append(failedTasks, t)
		}
	}

	fmt.Println("Errors:")
	for _, t := range failedTasks {
		fmt.Printf("Task id %s time %s, error %s\n", t.id, t.createdAt, t.payload)
	}

	fmt.Println("Done tasks:")
	for _, t := range doneTasks {
		fmt.Printf("Task id %s time %s, completed\n", t.id, t.createdAt)
	}

}

func taskGenerator(ctx context.Context) <-chan *Task {
	fmt.Println("Task generator started...")

	tasks := make(chan *Task)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(tasks)
				fmt.Println("Task generator finished")
				return
			default:
				id := uuid.New().String()
				createdAt := time.Now().Format(time.RFC3339)

				if time.Now().Nanosecond()/1000%2 > 0 {
					createdAt = "Some error occured"
				}

				tasks <- &Task{
					id:        id,
					createdAt: createdAt,
				}
			}
		}
	}()

	return tasks
}

func processTask(task Task) Task {
	task.processedAt = time.Now().Format(time.RFC3339Nano)
	tt, err := time.Parse(time.RFC3339, task.createdAt)

	if err != nil {
		task.success = false
		task.payload = []byte("something went wrong")
		return task
	}

	if tt.After(time.Now().Add(-20 * time.Second)) {
		task.success = true
		task.payload = []byte("task has been successed")
	} else {
		task.success = false
		task.payload = []byte("something went wrong")
	}

	time.Sleep(time.Millisecond * 150)

	return task
}

func taskWorker(tasksIn <-chan *Task) <-chan *Task {
	fmt.Println("Task worker started...")

	tasksOut := make(chan *Task)
	maxWorkersCount := runtime.GOMAXPROCS(0)
	workersLimit := make(chan struct{}, maxWorkersCount)
	wg := sync.WaitGroup{}

	go func() {
		for task := range tasksIn {
			t := task

			workersLimit <- struct{}{}
			wg.Add(1)

			go func() {
				tt := processTask(*t)
				tasksOut <- &tt

				<-workersLimit
				wg.Done()
			}()
		}

		wg.Wait()
		close(tasksOut)
		close(workersLimit)
		fmt.Println("Task worker finished")
	}()

	return tasksOut
}
