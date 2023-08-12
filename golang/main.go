package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

// A Task struct represents a task data or job for a worker
type Task struct {
	id           int
	createdTime  string // время создания
	finishedTime string // время выполнения
	value        []byte
	errors       error
}

// generateTasks function generates a new tasks and returns a Task channel
func generateTasks(totalTasks int) <-chan Task {

	tasks := make(chan Task, totalTasks)

	go func() {

		for t := 1; t <= totalTasks; t++ {
			createdTime := time.Now().Format(time.RFC3339Nano)

			var err error
			if time.Now().Nanosecond()%2 > 0 {
				err = errors.New("Some error occurred")
			}

			task := Task{
				id:          t,
				createdTime: createdTime,
				errors:      err,
			}
			tasks <- task
		}
		close(tasks)
	}()

	return tasks
}

// worker function range over tasks and runs a work
func worker(id int, tasks <-chan Task, doneTasks chan<- Task, failedTasks chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range tasks {
		process(task, doneTasks, failedTasks)
	}
}

// process function checks condition of successfull or failed tasks and save a task value
func process(task Task, doneTasks chan<- Task, failedTasks chan<- error) {

	tt, _ := time.Parse(time.RFC3339Nano, task.createdTime)

	const timeDiff = -20

	if tt.After(time.Now().Add(timeDiff * time.Second)) {
		task.value = []byte("task has been succeeded")
	} else {
		task.value = []byte("something went wrong")
	}
	task.finishedTime = time.Now().Format(time.RFC3339Nano)

	// simulate work
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)))

	filterTask(task, doneTasks, failedTasks)
}

// filterTask function filter successfull or failed tasks
func filterTask(task Task, doneTasks chan<- Task, failedTasks chan<- error) {

	const startIndex = 14
	const taskSubstring = "succeeded"

	if string(task.value[startIndex:]) == taskSubstring {
		doneTasks <- task
	} else {
		failedTasks <- fmt.Errorf("Task id: %d, Created time: %s, Finished time: %s, Value: %s, Error: %s", task.id, task.createdTime, task.finishedTime, task.value, task.errors)
	}
}

func main() {

	const totalTasks = 500
	const numWorkers = 20 // or runtime.NumCPU()

	tasks := generateTasks(totalTasks)

	doneTasks := make(chan Task, totalTasks)
	failedTasks := make(chan error, totalTasks)

	wg := &sync.WaitGroup{}

	wg.Add(numWorkers)
	for i := 1; i <= numWorkers; i++ {
		go worker(i, tasks, doneTasks, failedTasks, wg)
	}

	// for non-blocking read from doneTasks
	go func() {
		wg.Wait()
		close(doneTasks)
		close(failedTasks)
	}()

	// for non-blocking read from failedTasks
	go func() {
		for errorTask := range failedTasks {
			fmt.Println("Errors:", errorTask)
		}
	}()

	for doneTask := range doneTasks {
		fmt.Printf("Done tasks: Task id: %d, Created time: %s, Finished time: %s, Value: %s, Error: %s\n", doneTask.id, doneTask.createdTime, doneTask.finishedTime, doneTask.value, doneTask.errors)
	}
}
