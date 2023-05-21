package main

import (
	"fmt"
	"time"
)

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// в конце должно выводить успешные таски и ошибки выполнения остальных тасков

// A Task represents a meaninglessness of our life
type Task struct {
	CreationTime string
	ID           int
	Result       []byte
	FinishTime   string
}

func main() {
	createTask := func(taskChan chan Task) {
		for {
			time.Sleep(1 * time.Microsecond)
			creationTime := time.Now().Format(time.RFC3339)
			id := int(time.Now().UnixMicro())
			if id%2 > 0 {
				creationTime = "Some error occurred"
			}
			taskChan <- Task{CreationTime: creationTime, ID: id, Result: nil, FinishTime: ""}
		}
	}

	processTask := func(task Task) Task {
		taskTime, err := time.Parse(time.RFC3339, task.CreationTime)
		if err == nil && taskTime.After(time.Now().Add(-20*time.Second)) {
			task.Result = []byte("task has been succeeded")
		} else {
			task.Result = []byte("something went wrong")
		}
		task.FinishTime = time.Now().Format(time.RFC3339Nano)
		time.Sleep(time.Millisecond * 150)
		return task
	}

	doneTasks := make(chan Task)
	defer close(doneTasks)

	undoneTasks := make(chan error)
	defer close(undoneTasks)

	sortTasks := func(task Task) {
		if string(task.Result[14:]) == "succeeded" {
			doneTasks <- task
		} else {
			undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", task.ID, task.CreationTime, task.Result)
		}
	}

	result := map[int]Task{}
	go func() {
		for task := range doneTasks {
			result[task.ID] = task
		}
	}()

	errors := []error{}
	go func() {
		for err := range undoneTasks {
			errors = append(errors, err)
		}
	}()

	taskChan := make(chan Task, 10)
	defer close(taskChan)

	go func() {
		for task := range taskChan {
			task = processTask(task)
			sortTasks(task)
		}
	}()

	go createTask(taskChan)

	time.Sleep(time.Second * 3)

	fmt.Println("Done tasks:")
	for id := range result {
		fmt.Println(id)
	}

	fmt.Println("Errors:")
	for _, err := range errors {
		fmt.Println(err)
	}
}
