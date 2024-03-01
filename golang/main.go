package main

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

// Custom TaskError for clarifying type of errors in Task
type TaskError struct {
	reason string
}

func NewTaskError(reason string) TaskError {
	return TaskError{reason: reason}
}

func (e TaskError) Error() string {
	return fmt.Sprintf("%s", e.reason)
}

var (
	taskCreationError   = NewTaskError("Failed to create task")
	taskProcessingError = NewTaskError("Failed to proccess task")
)

// A Task represents a single task
type Task struct {
	id           int32
	creationTime time.Time // время создания
	finishTime   time.Time // время выполнения
	err          TaskError
}

var taskCounter atomic.Int32

// Producer function for creating tasks
func taskCreator(tasksChan chan Task, numberOfTasks int32) {

	for {
		if taskCounter.Load() == numberOfTasks {
			break
		}
		task := Task{
			id:           taskCounter.Add(1),
			creationTime: time.Now(),
		}
		if task.creationTime.Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			println("taskProcessingError")

			task.err = taskCreationError
		}
		tasksChan <- task // передаем таск на выполнение
	}
	close(tasksChan)
}

// Consumer function for processing tasks
func tasksWorker(in chan Task, out chan Task) {
	wg := sync.WaitGroup{}

	for task := range in {
		wg.Add(1)
		go func(task Task) {
			defer wg.Done()
			tt := task.creationTime
			if !tt.After(time.Now().Add(-20 * time.Second)) {
				task.err = taskProcessingError
			}
			task.finishTime = time.Now()

			time.Sleep(time.Millisecond * 150)
			out <- task
		}(task)
	}
	wg.Wait()
	close(out)
}

// Tasks result sorter
func tasksSorter(in chan Task, taskErrors chan error, doneTasks chan Task) {
	wg := sync.WaitGroup{}
	for t := range in {
		wg.Add(1)
		go func(t Task) {
			defer wg.Done()
			if errors.Is(t.err, &TaskError{}) {
				taskErrors <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.creationTime, t.err)
			} else {
				doneTasks <- t
			}
		}(t)
	}
	wg.Wait()
	close(doneTasks)
	close(taskErrors)
}

func main() {
	// Atomic value for holding amount of successfully created tasks, from which id's is generated

	const numberOfTasks = 200
	const lenOfBuf = 2

	tasksChan := make(chan Task, lenOfBuf)

	go taskCreator(tasksChan, numberOfTasks)

	processedTasks := make(chan Task, lenOfBuf)

	go tasksWorker(tasksChan, processedTasks)

	doneTasks := make(chan Task, lenOfBuf)
	taskErrors := make(chan error, lenOfBuf)
	go tasksSorter(processedTasks, taskErrors, doneTasks)

	result := map[int32]Task{}
	err := []error{}

	go func() {
		for r := range doneTasks {
			go func() {
				result[r.id] = r
			}()
		}
		for r := range taskErrors {
			go func() {
				err = append(err, r)
			}()
		}
		close(doneTasks)
		close(taskErrors)
	}()

	time.Sleep(time.Second * 3)

	println("Errors:")
	for r := range err {
		println(r)
	}

	println("Done tasks:")
	for r := range result {
		println(r)
	}
}
