package main

import (
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
	result       TaskError
}

func main() {
	// Atomic value for holding amount of successfully created tasks, from which id's is generated
	var successfulCounter atomic.Int32

	taskCreator := func(tasksChan chan Task) {
		go func() {
			for {
				// Set id to -1
				var id int32 = -1
				creationTime := time.Now()
				if !(creationTime.Nanosecond()%2 > 0) { // вот такое условие появления ошибочных тасков
					successfulCounter.Add(1)
					id = successfulCounter.Load()
				}
				tasksChan <- Task{creationTime: creationTime, id: id} // передаем таск на выполнение
			}
		}()
	}

	tasksChan := make(chan Task, 10)

	go taskCreator(tasksChan)

	taskWorker := func(task Task) Task {
		tt, _ := time.Parse(time.RFC3339, task.creationTime)
		if tt.After(time.Now().Add(-20 * time.Second)) {
			task.result = []byte("task has been successed")
		} else {
			task.result = []byte("something went wrong")
		}
		task.finishTime = time.Now().Format(time.RFC3339Nano)

		time.Sleep(time.Millisecond * 150)

		return task
	}

	doneTasks := make(chan Task)
	undoneTasks := make(chan error)

	tasksorter := func(t Task) {
		if string(t.result[14:]) == "successed" {
			doneTasks <- t
		} else {
			undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.creationTime, t.result)
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

	result := map[int]Task{}
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
