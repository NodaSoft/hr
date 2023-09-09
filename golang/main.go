package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	// Constant error message for task creation error
	ErrTaskCreation = "Some error occurred"

	// Rate at which tasks are generated
	TaskCreationRate = 150 * time.Millisecond

	// Max age for a task to be processed
	MaxTaskAge = 2 * time.Second
)

// Task struct represents a task with its metadata.
type Task struct {
	ID          int       // Unique identifier for the task
	CreateTime  time.Time // Timestamp when the task was created
	FinishTime  time.Time // Timestamp when the task was completed
	TaskOutcome string    // Outcome message for the task
}

func main() {
	// Channel to produce new tasks
	taskChan := make(chan Task, 10)

	// Channel for successfully completed tasks
	doneTasks := make(chan Task, 10)

	// Channel for tasks that faced errors
	undoneTasks := make(chan error, 10)

	var wg sync.WaitGroup

	// Start task producer in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		taskProducer(taskChan)
	}()

	// Start multiple consumers to process tasks in parallel
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			taskConsumer(taskChan, doneTasks, undoneTasks)
		}()
	}

	// Goroutine to close result channels after all tasks are processed
	go func() {
		wg.Wait()
		close(doneTasks)
		close(undoneTasks)
	}()

	// Collect results and print
	taskResults := collectTasks(doneTasks)
	errorResults := collectErrors(undoneTasks)
	printResults(taskResults, errorResults)
}

// taskProducer simulates the creation of tasks at a specified rate.
func taskProducer(ch chan<- Task) {
	ticker := time.NewTicker(TaskCreationRate)
	defer ticker.Stop()

	// Timeout after 3 seconds to stop task production
	timeout := time.After(3 * time.Second)

	// Random generator for introducing random delays
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		select {
		case <-ticker.C:
			ct := time.Now()
			// Introducing a random delay half of the time
			if rnd.Intn(2) == 0 {
				randomDelay := time.Duration(rnd.Intn(5000)) * time.Millisecond
				delayed := ct.Add(-randomDelay)
				ch <- Task{ID: int(ct.Unix()), CreateTime: delayed, TaskOutcome: ErrTaskCreation}
			} else {
				ch <- Task{ID: int(ct.Unix()), CreateTime: ct}
			}
		case <-timeout:
			// Stop producing tasks after the timeout
			close(ch)
			return
		}
	}
}

// taskConsumer processes tasks and categorizes them as done or error.
func taskConsumer(taskChan <-chan Task, doneChan chan<- Task, errorChan chan<- error) {
	for task := range taskChan {
		task = processTask(task)
		if task.TaskOutcome == ErrTaskCreation {
			errorChan <- fmt.Errorf("Task id %d created at %s, error: %s", task.ID, task.CreateTime, task.TaskOutcome)
		} else {
			doneChan <- task
		}
	}
}

// processTask processes a task and updates its outcome and finish time.
func processTask(t Task) Task {
	// Checking if task is too old to be processed
	if time.Since(t.CreateTime) <= MaxTaskAge {
		t.TaskOutcome = "Task completed successfully"
	} else {
		t.TaskOutcome = ErrTaskCreation
	}
	t.FinishTime = time.Now()

	// Simulating task processing time
	time.Sleep(TaskCreationRate)
	return t
}

// collectTasks collects processed tasks from the done channel.
func collectTasks(done <-chan Task) map[int]Task {
	result := make(map[int]Task)
	for t := range done {
		result[t.ID] = t
	}
	return result
}

// collectErrors collects errors from the error channel.
func collectErrors(errors <-chan error) []error {
	var errs []error
	for e := range errors {
		errs = append(errs, e)
	}
	return errs
}

// printResults outputs the results of the processed tasks.
func printResults(tasks map[int]Task, errors []error) {
	fmt.Println("Done tasks:")
	for _, task := range tasks {
		fmt.Printf("Task ID: %d, Finish Time: %s, Outcome: %s\n", task.ID, task.FinishTime, task.TaskOutcome)
	}
	fmt.Println("\nErrors:")
	for _, err := range errors {
		fmt.Println(err)
	}
}
