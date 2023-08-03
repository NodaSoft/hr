package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	Id              int
	GenerationTime  time.Time
	ExecutionTime   time.Time
	Error           error
}

func (task *Task) failWith(message string) {
	if (task.Error != nil) {
		return
	}
	
	task.Error = fmt.Errorf("Task id %d time %s, error: %s", 
				task.Id, 
				task.GenerationTime, 
				message)
}

func (task *Task) execute() {
	if task.GenerationTime.Before(time.Now().Add(-20 * time.Second)) {
		task.failWith("executed too late")
	}

	task.ExecutionTime = time.Now()

	time.Sleep(time.Millisecond * 150)
}

func generateTasks(tasks chan<- Task, taskAmount int) {
	for i := 0; i < taskAmount; i++ {
		task := Task{
			Id: i,
			GenerationTime: time.Now(),
		}

		if time.Now().Nanosecond()%2 > 0 { 
			task.failWith("wrong generation time")
		}

		tasks <- task
	}
	close(tasks)
}

// executeTasks generates taskAmount and executes them
func executeTasks(executedTasks chan<- Task, taskAmount int) {
	tasks := make(chan Task)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		generateTasks(tasks, taskAmount)
	}()

	for task := range tasks {
		wg.Add(1)
		go func(currentTask Task) {
			defer wg.Done()
			currentTask.execute()
			executedTasks <- currentTask
		}(task)
	}

	go func() {
		wg.Wait()
		close(executedTasks)	
	}()
}

// collectResults returns (okTasks, errorTasks), differentiation is based on Error field
func collectResults(executedTasks <-chan Task) ([]Task, []Task) {

	okTasks := []Task{}
	errorTasks := []Task{}

	for task := range executedTasks {
		currentTask := task
		if task.Error == nil {
			okTasks = append(okTasks, currentTask)
		} else {
			errorTasks = append(errorTasks, currentTask)
		}
	}

	return okTasks, errorTasks
} 

func main() {
	executedTasks := make(chan Task)

	executeTasks(executedTasks, 10)

	okTasks, errorTasks := collectResults(executedTasks)

	println("Done tasks:")
	for _, task := range okTasks {
		currentTask := task
		fmt.Println(currentTask)
	}

	println("Errors:")
	for _, task := range errorTasks {
		currentTask := task
		fmt.Println(currentTask.Error)
	}
}
