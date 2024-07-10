package main

import (
	"fmt"
	"time"
)

const (
	TaskChannelBufferSize   = 10
	TaskGenerationInterval  = 100 * time.Millisecond
	TaskProcessingDelay     = 150 * time.Millisecond
	PrintInterval           = 3 * time.Second
	TaskErrorMessage        = "Task processing failed"
	TaskSuccessMessage      = "Task has been successful"
	TaskErrorCreationString = "Task creation error occurred"
	ErrorCondition          = 2
	TimeThreshold           = 20 * time.Second
)

type Task struct {
	id                 int
	createdAt          time.Time
	processingDuration time.Time
	result             []byte
}

func main() {
	taskChannel := make(chan Task, TaskChannelBufferSize)
	doneTasks := make(chan Task)
	errorTasks := make(chan error)
	taskResults := make(map[int]Task)
	errors := []error{}

	go func() {
		for {
			createdAt := time.Now()
			if time.Now().Nanosecond()%ErrorCondition > 0 {
				createdAt = time.Time{}
			}
			taskChannel <- Task{id: int(time.Now().Unix()), createdAt: createdAt}
			time.Sleep(TaskGenerationInterval)
		}
	}()

	taskWorker := func(task Task) Task {
		if task.createdAt.IsZero() {
			task.result = []byte(TaskErrorMessage)
		} else if task.createdAt.After(time.Now().Add(-TimeThreshold)) {
			task.result = []byte(TaskErrorMessage)
		} else {
			task.result = []byte(TaskSuccessMessage)
		}
		task.processingDuration = time.Now()
		time.Sleep(TaskProcessingDelay)
		return task
	}

	taskSorter := func(task Task) {
		if string(task.result) == TaskSuccessMessage {
			doneTasks <- task
		} else if task.createdAt.IsZero() {
			errorTasks <- fmt.Errorf("Task id %d, error: %s. Details: Invalid creation time", task.id, task.result)
		} else {
			errorTasks <- fmt.Errorf("Task id %d, created at %s, error: %s. Details: Task did not pass the time check", task.id, task.createdAt.Format(time.RFC3339), task.result)
		}
	}

	go func() {
		for task := range taskChannel {
			processedTask := taskWorker(task)
			go taskSorter(processedTask)
		}
	}()

	go func() {
		for task := range doneTasks {
			taskResults[task.id] = task
		}
	}()

	go func() {
		for err := range errorTasks {
			errors = append(errors, err)
		}
	}()

	for {
		time.Sleep(PrintInterval)

		fmt.Println("Errors:")
		for _, err := range errors {
			fmt.Println(err)
		}

		fmt.Println("Done tasks:")
		for _, task := range taskResults {
			fmt.Printf("Task ID: %d, Created At: %s, Result: %s, Processed At: %s\n", task.id, task.createdAt.Format(time.RFC3339), task.result, task.processingDuration.Format(time.RFC3339))
		}
	}
}
