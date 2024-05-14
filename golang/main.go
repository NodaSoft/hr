package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Maximum allowed execution time for a task.
const taskMaxLifeTime = 20 * time.Second

// Task represents some work to be executed.
type Task struct {
	id         int
	startedAt  time.Time
	finishedAt *time.Time
	mustFail   bool
}

// NewTask returns a new Task instance.
func NewTask() Task {
	now := time.Now()

	task := Task{
		id:        int(now.Unix()) + now.Nanosecond(),
		startedAt: now,
	}

	if now.Nanosecond()%2 > 0 {
		task.mustFail = true
	}

	return task
}

// Execute attempts to run the task and returns error if fail.
func (t *Task) Execute() error {
	if t.mustFail {
		return fmt.Errorf("task %d failed: some error occured", t.id)
	}

	if t.startedAt.Add(taskMaxLifeTime).Before(time.Now()) {
		return fmt.Errorf("task %d failed: timeout", t.id)
	}

	finishedAt := time.Now()
	t.finishedAt = &finishedAt

	return nil
}

// Result contains the outcome of a task execution.
type Result struct {
	Task  Task
	Error error
}

// runCreator continuously generates new tasks and sends them on the channel.
func runCreator() <-chan Task {
	tasks := make(chan Task)

	go func() {
		for {
			tasks <- NewTask()
		}
	}()

	return tasks
}

// runWorker receives tasks from the channel and executes them, sending results back.
func runWorker(tasks <-chan Task) <-chan Result {
	results := make(chan Result)

	go func() {
		for task := range tasks {
			var result Result
			err := task.Execute()

			if err != nil {
				result = Result{
					Task:  task,
					Error: err,
				}
			} else {
				result = Result{Task: task}
			}

			results <- result

			time.Sleep(time.Millisecond * 150)
		}
	}()

	return results
}

// aggregateResults continuously collects and prints task execution results.
func aggregateResults(results <-chan Result) {
	var (
		failed []error
		done   []Task
		mu     sync.Mutex
	)

	go func() {
		for {
			time.Sleep(3 * time.Second)

			fmt.Println("Errors:")
			for _, f := range failed {
				fmt.Println(f.Error())
			}

			fmt.Println("Done tasks:")
			for _, d := range done {
				fmt.Println(d.id)
			}

			mu.Lock()
			failed, done = []error{}, []Task{}
			mu.Unlock()
		}
	}()

	for result := range results {
		if result.Error != nil {
			mu.Lock()
			failed = append(failed, result.Error)
			mu.Unlock()
		} else {
			mu.Lock()
			done = append(done, result.Task)
			mu.Unlock()
		}
	}
}

// joinChanels merges multiple result channels into a single channel.
func joinChanels(channels ...<-chan Result) <-chan Result {
	merged := make(chan Result)

	go func() {
		defer close(merged)

		wg := sync.WaitGroup{}

		wg.Add(len(channels))
		for _, ch := range channels {
			go func(ch <-chan Result) {
				for result := range ch {
					merged <- result
				}

				wg.Done()
			}(ch)
		}

		wg.Wait()
	}()

	return merged
}

func main() {
	tasks := runCreator()

	results := []<-chan Result{}

	for i := 0; i < runtime.NumCPU(); i++ {
		results = append(results, runWorker(tasks))
	}

	muxResults := joinChanels(results...)
	aggregateResults(muxResults)
}
