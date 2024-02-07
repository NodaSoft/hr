package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	id         int
	creationTime time.Time
	finishTime   time.Time
	result      string
}

func main() {
	taskCreator := func(taskChannel chan Task, wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			creationTime := time.Now()
			if creationTime.Nanosecond()%2 > 0 {
				taskChannel <- Task{id: int(creationTime.Unix()), creationTime: creationTime, result: "Some error occurred"}
			} else {
				taskChannel <- Task{id: int(creationTime.Unix()), creationTime: creationTime}
			}
			time.Sleep(time.Second) // Simulate task creation time
		}
	}

	taskWorker := func(task Task) Task {
		time.Sleep(150 * time.Millisecond) // Simulate task processing time
		task.finishTime = time.Now()
		if task.result == "" {
			task.result = "Task has been successfully completed"
		} else {
			task.result = "Something went wrong"
		}
		return task
	}

	taskSorter := func(task Task, wg *sync.WaitGroup, doneTasks chan Task, errorTasks chan Task) {
		defer wg.Done()
		if task.result == "Task has been successfully completed" {
			doneTasks <- task
		} else {
			errorTasks <- task
		}
	}

	taskChannel := make(chan Task, 10)
	doneTasks := make(chan Task)
	errorTasks := make(chan Task)
	var wg sync.WaitGroup

	wg.Add(1)
	go taskCreator(taskChannel, &wg)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskChannel {
				processedTask := taskWorker(task)
				wg.Add(1)
				go taskSorter(processedTask, &wg, doneTasks, errorTasks)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(doneTasks)
		close(errorTasks)
	}()

	var wg2 sync.WaitGroup
	wg2.Add(2)
	go func() {
		defer wg2.Done()
		for task := range doneTasks {
			fmt.Printf("Task ID %d: %s\n", task.id, task.result)
		}
	}()

	go func() {
		defer wg2.Done()
		for task := range errorTasks {
			fmt.Printf("Error in Task ID %d: %s\n", task.id, task.result)
		}
	}()

	wg2.Wait()
}
