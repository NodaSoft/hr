package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	tasksChanSize             = 10
	taskExecutorsLimit        = 10
	taskGenerationDuration    = time.Second * 10
	taskResultReportingPeriod = time.Second * 3
)

func main() {
	tasks := make(chan Task, tasksChanSize)
	ctx, cancel := context.WithTimeout(context.Background(), taskGenerationDuration)
	defer cancel()
	go createTasks(ctx, tasks)

	doneTasks := make(chan Task)
	go executeTasks(tasks, doneTasks)

	reportTaskResults(doneTasks)
}

type Task struct {
	id             int
	creationTime   time.Time
	completionTime *time.Time
	result         []byte
	error          error
}

func (task Task) String() string {
	return fmt.Sprintf(`Task{id: %d, creationTime: %v, completionTime: %v, result: %s, error: %v}`,
		task.id, task.creationTime.Format(time.RFC3339),
		task.completionTime.Format(time.RFC3339Nano), string(task.result),
		task.error)
}

func createTasks(ctx context.Context, tasks chan<- Task) {
	for {
		select {
		case <-ctx.Done():
			close(tasks)
			return
		default:
			tasks <- Task{
				id:           int(time.Now().Unix()),
				creationTime: time.Now(),
			}
			time.Sleep(time.Millisecond * 500) // only for debugging to avoid a lot of output
		}
	}
}

func executeTasks(tasks <-chan Task, doneTasks chan<- Task) {
	var taskExecutionWG sync.WaitGroup
	executorsLimiter := make(chan struct{}, taskExecutorsLimit)
	for task := range tasks {
		executorsLimiter <- struct{}{}
		taskExecutionWG.Add(1)
		go func(task Task) {
			doneTasks <- executeTask(task)
			taskExecutionWG.Done()
			<-executorsLimiter
		}(task)
	}
	taskExecutionWG.Wait()

	close(doneTasks)
}

func executeTask(task Task) Task {
	if task.creationTime.Nanosecond()%2 == 1 {
		task.result = []byte("task has been executed successfuly")
	} else {
		task.error = fmt.Errorf("something went wrong")
	}
	task.completionTime = new(time.Time)
	*task.completionTime = time.Now()

	time.Sleep(time.Millisecond * 150)

	return task
}

func reportTaskResults(doneTasks <-chan Task) {
	successfulTasks := []Task{}
	failedTasks := []Task{}
	var successfulTasksMutex sync.Mutex
	var failedTasksMutex sync.Mutex

	var readingWG sync.WaitGroup
	readingWG.Add(1)
	go func() {
		for task := range doneTasks {
			if task.error == nil {
				successfulTasksMutex.Lock()
				successfulTasks = append(successfulTasks, task)
				successfulTasksMutex.Unlock()
			} else {
				failedTasksMutex.Lock()
				failedTasks = append(failedTasks, task)
				failedTasksMutex.Unlock()
			}
		}
		readingWG.Done()
	}()
	doneReading := make(chan struct{})
	go func() {
		readingWG.Wait()
		close(doneReading)
	}()

	printReport := func() {
		fmt.Println("Successful tasks:")
		successfulTasksMutex.Lock()
		for _, task := range successfulTasks {
			fmt.Println(task)
		}
		successfulTasksMutex.Unlock()

		fmt.Println("Failed tasks:")
		failedTasksMutex.Lock()
		for _, err := range failedTasks {
			fmt.Println(err)
		}
		failedTasksMutex.Unlock()
	}

	timeToReportTicker := time.NewTicker(taskResultReportingPeriod)
	for {
		select {
		case <-timeToReportTicker.C:
			printReport()
		case <-doneReading:
			printReport()
			timeToReportTicker.Stop()
			return
		}
	}
}
