package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Ttype struct {
	id         int
	cT         string // Creation time
	fT         string // Finish time
	taskRESULT []byte
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	taskChan := make(chan Ttype)
	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	// Task creator
	go taskCreator(ctx, taskChan)

	// Workers for processing tasks
	const workerCount = 4
	for i := 0; i < workerCount; i++ {
		go taskWorker(ctx, taskChan, doneTasks, undoneTasks)
	}

	// Collect results
	var wg sync.WaitGroup
	results := make(map[int]Ttype)
	var errors []error

	wg.Add(2)
	go collectDoneTasks(ctx, &wg, doneTasks, results)
	go collectErrors(ctx, &wg, undoneTasks, &errors)
	wg.Wait()

	// Output results
	fmt.Println("Errors:")
	for _, err := range errors {
		fmt.Println(err)
	}

	fmt.Println("Done tasks:")
	for id := range results {
		fmt.Println(id)
	}
}

func taskCreator(ctx context.Context, taskChan chan<- Ttype) {
	defer close(taskChan)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			ft := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 {
				ft = "Some error occurred"
			}
			taskChan <- Ttype{cT: ft, id: int(time.Now().Unix())}
			time.Sleep(50 * time.Millisecond) // Throttle task creation
		}
	}
}

func taskWorker(ctx context.Context, taskChan <-chan Ttype, doneTasks chan<- Ttype, undoneTasks chan<- error) {
	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-taskChan:
			if !ok {
				return
			}
			processTask(task, doneTasks, undoneTasks)
		}
	}
}

func processTask(task Ttype, doneTasks chan<- Ttype, undoneTasks chan<- error) {
	// Simulate task processing time
	time.Sleep(150 * time.Millisecond)

	tt, err := time.Parse(time.RFC3339, task.cT)
	if err == nil && tt.After(time.Now().Add(-20*time.Second)) {
		task.taskRESULT = []byte("task has been succeeded")
	} else {
		task.taskRESULT = []byte("something went wrong")
	}
	task.fT = time.Now().Format(time.RFC3339Nano)

	if string(task.taskRESULT) == "task has been succeeded" {
		doneTasks <- task
	} else {
		undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", task.id, task.cT, task.taskRESULT)
	}
}

func collectDoneTasks(ctx context.Context, wg *sync.WaitGroup, doneTasks <-chan Ttype, results map[int]Ttype) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-doneTasks:
			if !ok {
				return
			}
			results[task.id] = task
		}
	}
}

func collectErrors(ctx context.Context, wg *sync.WaitGroup, undoneTasks <-chan error, errors *[]error) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case err, ok := <-undoneTasks:
			if !ok {
				return
			}
			*errors = append(*errors, err)
		}
	}
}
