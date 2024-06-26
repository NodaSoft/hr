package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	tasksChanSize             = 10
	taskGenerationDuration    = time.Second * 10
	taskResultReportingPeriod = time.Second * 3
)

func main() {
	tasks := make(chan Task, tasksChanSize)
	ctx, _ := context.WithTimeout(context.Background(), taskGenerationDuration)
	go createTasks(ctx, tasks)

	doneTasks := make(chan Task)
	taskErrors := make(chan error)
	go executeTasks(tasks, doneTasks, taskErrors)

	reportTaskResults(doneTasks, taskErrors)
}

type Task struct {
	id             int
	creationTime   string
	completionTime string
	result         []byte
}

func (t Task) IsSuccessful() bool {
	if t.completionTime == "" || string(t.result[14:]) != "successed" {
		return false
	}
	return true
}

func createTasks(ctx context.Context, tasks chan<- Task) {
	for {
		select {
		case <-ctx.Done():
			close(tasks)
			return
		default:
			creationTime := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 == 1 {
				creationTime = "Some error occured"
			}
			tasks <- Task{
				id:           int(time.Now().Unix()),
				creationTime: creationTime,
			}
			time.Sleep(time.Millisecond * 500) // чтобы не было слишком много значений в выводе во время проверки
		}
	}
}

func executeTasks(tasks <-chan Task, doneTasks chan<- Task, taskErrors chan<- error) {
	var taskExecutionWG sync.WaitGroup
	for task := range tasks {
		taskExecutionWG.Add(1)
		go func(task Task) {
			task = executeTask(task)
			if task.IsSuccessful() {
				doneTasks <- task
			} else {
				taskErrors <- fmt.Errorf("Task id %d time %s, error %s", task.id, task.creationTime, task.result)
			}
			taskExecutionWG.Done()
		}(task)
	}
	taskExecutionWG.Wait()

	close(doneTasks)
	close(taskErrors)
}

func executeTask(task Task) Task {
	tt, _ := time.Parse(time.RFC3339, task.creationTime)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		task.result = []byte("task has been successed")
	} else {
		task.result = []byte("something went wrong")
	}
	task.completionTime = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)

	return task
}

func reportTaskResults(doneTasksChan <-chan Task, taskErrorsChan <-chan error) {
	doneTasks := []Task{}
	taskErrors := []error{}
	var doneTasksMutex sync.Mutex
	var taskErrorsMutex sync.Mutex

	var readingWG sync.WaitGroup
	readingWG.Add(2)
	go func() {
		for task := range doneTasksChan {
			doneTasksMutex.Lock()
			doneTasks = append(doneTasks, task)
			doneTasksMutex.Unlock()
		}
		readingWG.Done()
	}()
	go func() {
		for err := range taskErrorsChan {
			taskErrorsMutex.Lock()
			taskErrors = append(taskErrors, err)
			taskErrorsMutex.Unlock()
		}
		readingWG.Done()
	}()
	doneReading := make(chan struct{})
	go func() {
		readingWG.Wait()
		close(doneReading)
	}()

	printReport := func() {
		fmt.Println("Done tasks:")
		doneTasksMutex.Lock()
		for r := range doneTasks {
			fmt.Println(r)
		}
		doneTasksMutex.Unlock()

		println("Errors:")
		taskErrorsMutex.Lock()
		for r := range taskErrors {
			println(r)
		}
		taskErrorsMutex.Unlock()
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
