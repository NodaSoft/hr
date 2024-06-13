package main

import (
	"fmt"
	"time"
)

type Type struct {
	ID            int
	CreatedAt     string
	ExecutionTime string
	TaskResult    []byte
}

func main() {
	tasksCreator := func(taskChannel chan Type) {
		go func() {
			ticker := time.NewTicker(10 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					now := time.Now()
					task := Type{
						ID:        int(now.Unix()),
						CreatedAt: now.Format(time.RFC3339),
					}
					if now.Nanosecond()%2 > 0 {
						task.ExecutionTime = "Some error occurred"
						task.TaskResult = []byte("something went wrong")
					} else {
						task.ExecutionTime = now.Format(time.RFC3339Nano)
						task.TaskResult = []byte("task has been succeeded")
					}
					taskChannel <- task
				}
			}
		}()
	}

	superChan := make(chan Type, 10)
	doneTasks := make(chan Type)
	undoneTasks := make(chan Type)

	go tasksCreator(superChan)

	taskWorker := func(task Type) (Type, error) {
		parsedTime, err := time.Parse(time.RFC3339, task.CreatedAt)
		if err != nil || parsedTime.Before(time.Now().Add(-20*time.Second)) {
			return Type{}, fmt.Errorf("Task ID %d failed to process", task.ID)
		}
		task.ExecutionTime = time.Now().Format(time.RFC3339Nano)
		return task, nil
	}

	go func() {
		for task := range superChan {
			task, err := taskWorker(task)
			if err != nil {
				undoneTasks <- task
			} else {
				doneTasks <- task
			}
		}
		close(doneTasks)
		close(undoneTasks)
	}()

	results := make(map[int]Type)
	var errors []error

	go func() {
		for {
			select {
			case r := <-doneTasks:
				results[r.ID] = r
			case r := <-undoneTasks:
				errors = append(errors, fmt.Errorf("task id %d time %s, error %s", r.ID, r.CreatedAt, r.TaskResult))
			}
		}
	}()

	stop := time.After(30 * time.Second)
	<-stop

	fmt.Println("Errors:")
	for _, e := range errors {
		fmt.Println(e)
	}

	fmt.Println("Done tasks:")
	for id, task := range results {
		fmt.Printf("ID: %d, Created At: %s, Execution Time: %s\n", id, task.CreatedAt, task.ExecutionTime)
	}
}
