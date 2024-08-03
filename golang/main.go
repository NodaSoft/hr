package main

import (
	"fmt"
	"sync"
	"time"
)

type Ttype struct {
	id         int
	cT         string // creation time
	fT         string // finish time
	taskRESULT string
}

func main() {
	taskCreator := func(taskChan chan Ttype) {
		for i := 0; i < 100; i++ {
			ft := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 {
				ft = "Some error occurred"
			}
			taskChan <- Ttype{cT: ft, id: int(time.Now().UnixNano())}
			time.Sleep(100 * time.Millisecond)
		}
		close(taskChan)
	}

	taskWorker := func(task Ttype) Ttype {
		tt, err := time.Parse(time.RFC3339, task.cT)
		if err != nil || tt.Before(time.Now().Add(-20*time.Second)) {
			task.taskRESULT = "something went wrong"
		} else {
			task.taskRESULT = "task has been succeeded"
		}
		task.fT = time.Now().Format(time.RFC3339Nano)
		time.Sleep(150 * time.Millisecond) // simulate work
		return task
	}

	taskSorter := func(task Ttype, doneTasks chan Ttype, errorTasks chan Ttype) {
		if task.taskRESULT == "task has been succeeded" {
			doneTasks <- task
		} else {
			errorTasks <- task
		}
	}

	taskChan := make(chan Ttype, 10)
	doneTasks := make(chan Ttype, 10)
	errorTasks := make(chan Ttype, 10)

	go taskCreator(taskChan)

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskChan {
				processedTask := taskWorker(task)
				taskSorter(processedTask, doneTasks, errorTasks)
			}
		}()
	}

	var mu sync.Mutex
	results := make(map[int]Ttype)
	var errors []Ttype

	go func() {
		for task := range doneTasks {
			mu.Lock()
			results[task.id] = task
			mu.Unlock()
		}
	}()

	go func() {
		for task := range errorTasks {
			mu.Lock()
			errors = append(errors, task)
			mu.Unlock()
		}
	}()

	go func() {
		wg.Wait()
		close(doneTasks)
		close(errorTasks)
	}()

	time.Sleep(10 * time.Second)

	mu.Lock()
	fmt.Println("Done tasks:")
	for _, task := range results {
		fmt.Printf("Task ID: %d, Creation Time: %s, Finish Time: %s, Result: %s\n", task.id, task.cT, task.fT, task.taskRESULT)
	}

	fmt.Println("\nErrors:")
	for _, task := range errors {
		fmt.Printf("Task ID: %d, Error: %s\n", task.id, task.cT)
	}
	mu.Unlock()
}
