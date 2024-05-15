package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

const numTasks = 10

var (
	errTaskWorker  = errors.New("task worker error")
	errTaskCreator = errors.New("task creator error")
)

type Task struct {
	Id           int64
	CreationTime time.Time
	FinishTime   time.Time
	Result       string
	Error        error
}

func taskCreator(ctx context.Context, taskChan chan<- Task) {
	for {
		select {
		case <-ctx.Done():
			close(taskChan)
			return
		default:
			var err error = nil

			now := time.Now()
			// if time.Now().Nanosecond()%2 > 0 {
			if rand.Int()%2 > 0 {
				err = errTaskCreator
			}

			task := Task{
				Id:           int64(time.Now().UnixNano()),
				CreationTime: now,
				Error:        err,
			}

			taskChan <- task
		}
	}
}

func taskWorker(task Task) Task {
	time.Sleep(time.Millisecond * 150)

	if time.Since(task.CreationTime) < 20*time.Second {
		task.Result = "task has been successed"
	} else {
		task.Error = errTaskWorker
	}

	task.FinishTime = time.Now()
	return task
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	taskChan := make(chan Task, numTasks)
	go taskCreator(ctx, taskChan)

	doneTasks := make(chan Task, numTasks)
	undoneTasks := make(chan error, numTasks)

	taskSorter := func(task Task) {
		if task.Error == nil {
			doneTasks <- task
		} else {
			undoneTasks <- fmt.Errorf("Task id %d time %s, error %w", task.Id, task.CreationTime.String(), task.Error)
		}
	}

	go func() {
		for task := range taskChan {
			task = taskWorker(task)
			go taskSorter(task)
		}
		close(taskChan)
	}()

	result := make(map[int64]Task)
	taskErrors := make([]error, 0)

	go func() {
		defer close(doneTasks)
		for task := range doneTasks {
			result[task.Id] = task
		}
	}()

	go func() {
		defer close(undoneTasks)
		for taskError := range undoneTasks {
			taskErrors = append(taskErrors, taskError)
		}
	}()

	time.Sleep(time.Second * 3)

	fmt.Println("Errors:")
	for _, taskError := range taskErrors {
		fmt.Println(taskError)
	}

	fmt.Println("Done tasks:")
	for id, task := range result {
		fmt.Println(id, " -> ", task)
	}
}
