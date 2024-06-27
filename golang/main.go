package main

import (
	"fmt"
	"time"
	"sync"
)

type Task struct {
	id int
	created time.Time
	finished time.Time
	failed bool
	result []byte
}

func taskCreator(taskC chan Task) {
	for {
		created := time.Now()
		failed := time.Now().Nanosecond() % 2 > 0
		taskC <- Task{
			id: int(time.Now().Unix()),
			created: created,
			failed: failed,
		}
	}
}

func taskWorker(task Task) Task {
	if task.created.After(time.Now().Add(-20 * time.Second)) && !task.failed {
		task.result = []byte("task finished successfully")
	} else {
		task.result = []byte("something went wrong")
	}
	task.finished = time.Now()
	time.Sleep(150 * time.Millisecond)
	return task
}

func main() {
	taskC := make(chan Task, 10)
	go taskCreator(taskC)
	
	mu := new(sync.Mutex)
	result := make(map[int]Task)
	errors := make([]error, 0)

	go func() {
		for task := range taskC{
			task := taskWorker(task)
			mu.Lock()
			if task.failed {
				err := fmt.Errorf("Task id=%d created=%s error=%s", task.id, task.created.Format(time.RFC3339Nano), string(task.result))
				errors = append(errors, err)
			} else {
				result[task.id] = task
			}
			mu.Unlock()
		}
	}()

	start := time.Now()
	for time.Now().Sub(start) < 10 * time.Second {
		mu.Lock()
		fmt.Println("Errors:")
		for _, err := range errors {
			fmt.Println(err)
		}
		fmt.Println("Done tasks:")
		for _, task := range result {
			fmt.Println("result", task)
		}
		mu.Unlock()
		time.Sleep(3 * time.Second)
	}
}