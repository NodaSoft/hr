package main

import (
	"errors"
	"fmt"
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
	return fmt.Sprintf("", e.reason)
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

func main() {
	// Atomic value for holding amount of successfully created tasks, from which id's is generated
	var taskCounter atomic.Int32

	taskCreator := func(tasksChan chan Task) {
		go func() {
			for {
				task := Task{
					id:           taskCounter.Add(1),
					creationTime: time.Now(),
				}

				if task.creationTime.Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					task.err = taskCreationError
				}
				tasksChan <- task // передаем таск на выполнение
			}
		}()
	}

	tasksChan := make(chan Task, 10)

	go taskCreator(tasksChan)

	taskWorker := func(task Task) Task {
		tt := task.creationTime
		if !tt.After(time.Now().Add(-20 * time.Second)) {
			task.err = taskProcessingError
		}
		task.finishTime = time.Now()

		time.Sleep(time.Millisecond * 150)

		return task
	}

	doneTasks := make(chan Task)
	undoneTasks := make(chan error)

	tasksorter := func(t Task) {
		if errors.As(t.err, &TaskError{}) {
			undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.creationTime, t.err)
		} else {
			doneTasks <- t
		}
	}

	go func() {
		// получение тасков
		for t := range tasksChan {
			t = taskWorker(t)
			go tasksorter(t)
		}
		close(tasksChan)
	}()

	result := map[int32]Task{}
	err := []error{}
	go func() {
		for r := range doneTasks {
			go func() {
				result[r.id] = r
			}()
		}
		for r := range undoneTasks {
			go func() {
				err = append(err, r)
			}()
		}
		close(doneTasks)
		close(undoneTasks)
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
