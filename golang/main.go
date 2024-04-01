package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	ID         int
	CreatedAt  string
	FinishedAt string
	Result     []byte
}

func main() {
	const maxTasks = 10
	taskChan := make(chan Task, maxTasks)
	doneTasks := make(chan Task, maxTasks)
	errorTasks := make(chan error, maxTasks)
	var wg sync.WaitGroup

	// Task creator
	go createTasks(taskChan)

	// Task processor
	for i := 0; i < maxTasks; i++ {
		wg.Add(1)
		go processTasks(taskChan, doneTasks, errorTasks, &wg)
	}

	// Collect results for a fixed time
	time.Sleep(3 * time.Second)
	close(taskChan)
	wg.Wait()
	close(doneTasks)
	close(errorTasks)

	displayResults(doneTasks, errorTasks)
}

func createTasks(taskChan chan<- Task) {
	for {
		task := Task{
			ID:        int(time.Now().Unix()),
			CreatedAt: time.Now().Format(time.RFC3339),
		}
		if time.Now().Nanosecond()%2 > 0 {
			task.CreatedAt = "Some error occurred"
		}
		taskChan <- task
		time.Sleep(150 * time.Millisecond)
	}
}

func processTasks(taskChan <-chan Task, doneTasks chan<- Task, errorTasks chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range taskChan {
		t, err := time.Parse(time.RFC3339, task.CreatedAt)
		if err != nil || t.Before(time.Now().Add(-20*time.Second)) {
			task.Result = []byte("something went wrong")
			errorTasks <- fmt.Errorf("Task id %d, error: %s", task.ID, task.Result)
		} else {
			task.Result = []byte("task has been succeeded")
			doneTasks <- task
		}
		task.FinishedAt = time.Now().Format(time.RFC3339Nano)
	}
}

func displayResults(doneTasks <-chan Task, errorTasks <-chan error) {
	fmt.Println("Done tasks:")
	for task := range doneTasks {
		fmt.Printf("Task ID: %d, Finished At: %s\n", task.ID, task.FinishedAt)
	}

	fmt.Println("\nErrors:")
	for err := range errorTasks {
		fmt.Println(err)
	}
}
