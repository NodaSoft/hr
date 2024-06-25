package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// A struct intended to bring some meaning to our lives
type Task struct {
	id         int
	createdAt  time.Time
	finishedAt time.Time
	res        string
	isError    bool
}

func taskCreator(ctx context.Context, taskChan chan<- Task) {
	defer close(taskChan)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			createdAt := time.Now()
			// Keeping the error condition unchanged. It grows with numWorkers
			// within the same timeframe (10 sec in our case)
			isError := createdAt.Nanosecond()%2 > 0
			res := ""
			if isError {
				res = "err: nanosecond is odd"
			}
			taskChan <- Task{id: int(createdAt.Unix()), createdAt: createdAt, isError: isError, res: res}
		}
	}
}

func taskWorker(ctx context.Context, taskChan <-chan Task, resultChan chan<- Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range taskChan {
		select {
		case <-ctx.Done():
			return
		default:
			if time.Since(task.createdAt) > 20*time.Microsecond {
				task.res = "something went wrong"
				task.isError = true
			} else if !task.isError {
				task.res = "task has completed successfully"
			}
			task.finishedAt = time.Now()
			time.Sleep(time.Millisecond * 150)
			resultChan <- task
		}
	}
}

func collectResults(ctx context.Context, resultChan <-chan Task, successTasks *[]Task, errorTasks *[]Task, mtx *sync.Mutex) {
	for task := range resultChan {
		select {
		case <-ctx.Done():
			return
		default:
			mtx.Lock()
			if task.isError {
				*errorTasks = append(*errorTasks, task)
			} else {
				*successTasks = append(*successTasks, task)
			}
			mtx.Unlock()
		}
	}
}

func printPeriodicResults(ctx context.Context, successTasks, errorTasks *[]Task, mtx *sync.Mutex) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			mtx.Lock()
			if len(*successTasks) > 0 || len(*errorTasks) > 0 {
				fmt.Println("===SUCCESSFUL TASKS:")
				for _, task := range *successTasks {
					fmt.Printf("Task id %d created at %s, finished at %s: %s\n", task.id, task.createdAt, task.finishedAt, task.res)
				}
				fmt.Println("===FAILED TASKS:")
				for _, task := range *errorTasks {
					fmt.Printf("Task id %d created at %s, error %s\n", task.id, task.createdAt, task.res)
				}
				fmt.Println("---")
			}
			mtx.Unlock()
		}
	}
}

func main() {
	taskChan := make(chan Task)
	resultChan := make(chan Task)

	var mtx sync.Mutex
	var wg sync.WaitGroup

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		taskCreator(ctx, taskChan)
	}()

	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go taskWorker(ctx, taskChan, resultChan, &wg)
	}

	var successTasks []Task
	var errorTasks []Task
	go func() {
		wg.Wait() // Waiting for all workers to finish before closing resultChan
		close(resultChan)
	}()

	go collectResults(ctx, resultChan, &successTasks, &errorTasks, &mtx)

	go printPeriodicResults(ctx, &successTasks, &errorTasks, &mtx)

	<-ctx.Done()

	mtx.Lock()
	fmt.Printf("Total tasks processed: %d\n", len(successTasks)+len(errorTasks))
	mtx.Unlock()
}
