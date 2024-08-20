package main

import (
	"fmt"
	"sync"
	"time"
)

// Task represents a task for processing
type Task struct {
	id         int
	createTime string
	finishTime string
	taskResult []byte
}

func taskCreator(task chan<- Task) {
	for {
		start := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 { // условие появления ошибочных тасков
			start = "Some error occured"
		}
		task <- Task{createTime: start, id: int(time.Now().Unix())} // передаем таск на выполнение
	}
}

func worker(task Task) Task {
	tt, _ := time.Parse(time.RFC3339, task.createTime)

	if tt.After(time.Now().Add(-20 * time.Second)) {
		task.taskResult = []byte("task was completed successfully")
	} else {
		task.taskResult = []byte("something went wrong")
	}

	task.finishTime = time.Now().Format(time.RFC3339Nano)
	time.Sleep(time.Millisecond * 150)

	return task
}

func taskSorter(finishedTasks chan<- Task, failedTasks chan<- error, task Task) {
	if string(task.taskResult) == "task was completed successfully" {
		finishedTasks <- task
	} else {
		failedTasks <- fmt.Errorf("task_id: %d, time: %s, error: %s", task.id, task.createTime, task.taskResult)
	}
}

func printSortedTasks(ticker *time.Ticker, finishedTasks <-chan Task, failedTasks <-chan error) {
	var mx sync.RWMutex
	results := make(map[int]Task)
	var errors []error

	for {
		select {
		case task, ok := <-finishedTasks:
			if !ok {
				finishedTasks = nil
			} else {
				mx.Lock()
				results[task.id] = task
				mx.Unlock()
			}
		case err, ok := <-failedTasks:
			if !ok {
				failedTasks = nil
			} else {
				mx.Lock()
				errors = append(errors, err)
				mx.Unlock()
			}
		case <-ticker.C:
			mx.RLock()
			fmt.Println("Done tasks:")
			for _, task := range results {
				fmt.Printf("Task id: %d, finished at: %s\n", task.id, task.finishTime)
			}
			fmt.Println("Errors:")
			for _, err := range errors {
				fmt.Println(err)
			}
			mx.RUnlock()
		}
		if finishedTasks == nil && failedTasks == nil {
			break
		}
	}
}

func main() {
	finishedTasks := make(chan Task)
	failedTasks := make(chan error)
	taskChan := make(chan Task, 10)

	go taskCreator(taskChan)

	var wg sync.WaitGroup

	go func() {
		for t := range taskChan {
			wg.Add(1)
			go func(t Task) {
				defer wg.Done()
				processedTask := worker(t)
				taskSorter(finishedTasks, failedTasks, processedTask)
			}(t)
		}
		wg.Wait()
		close(finishedTasks)
		close(failedTasks)
	}()

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	go printSortedTasks(ticker, finishedTasks, failedTasks)

	time.Sleep(10 * time.Second)
}
