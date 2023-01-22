package main

import (
	"context"
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

// A Task represents a meaninglessness of our life
type Task struct {
	successed    bool
	id           int64
	result       string
	creationTime time.Time // время создания
	finishedTime time.Time // время выполнения
}

// TODO: can store complete tasks and failed tasks
type taskManager struct {
	id int64
}

func NewTaskManager() *taskManager {
	return &taskManager{
		id: 0,
	}
}

func (t *taskManager) CreateTask(creationTime time.Time) Task {
	return Task{
		id:           atomic.AddInt64(&t.id, 1),
		creationTime: creationTime,
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.TODO())

	taskManager := NewTaskManager()

	taskWorkerChan := make(chan Task, 10)
	taskResultChan := make(chan Task)

	//creation workers pool
	for i := 0; i < cap(taskWorkerChan); i++ {
		go taskWorkers(ctx, taskWorkerChan, taskResultChan)
	}

	//start main work
	go func() {
		defer func() {
			close(taskWorkerChan)
			close(taskResultChan)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				creationTime := time.Now()
				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					creationTime = time.Time{}
				}

				taskWorkerChan <- taskManager.CreateTask(creationTime) // передаем таск на выполнение
			}
		}
	}()

	completeTasksChan := make(chan Task)
	errorsChan := make(chan error)

	//receiving result from workers pool
	go func() {
		for task := range taskResultChan {
			if task.successed {
				completeTasksChan <- task
			} else {
				errorsChan <- fmt.Errorf("Task id %d time %s, error %s", task.id, task.creationTime, task.result)
			}
		}
	}()

	completeTasks := make(map[int64]Task)
	errs := make([]error, 0)

	go func() {
		defer func() {
			close(completeTasksChan)
			close(errorsChan)
		}()

		for {
			select {
			case task := <-completeTasksChan:
				completeTasks[task.id] = task
			case err := <-errorsChan:
				errs = append(errs, err)
			case <-ctx.Done():
				return
			}
		}
	}()

	time.Sleep(time.Second * 3)
	cancel()

	fmt.Println("Errors:")
	for _, err := range errs {
		fmt.Println(err.Error())
	}

	fmt.Println("Done tasks:")
	for id := range completeTasks {
		fmt.Println(id)
	}
}

func taskWorkers(ctx context.Context, workChan <-chan Task, result chan<- Task) {
	for task := range workChan {
		select {
		case <-ctx.Done():
			return
		default:
			result <- taskWorker(task)
		}
	}
}

func taskWorker(task Task) Task {
	goal := time.Now().Add(-20 * time.Second)

	task.result = "something went wrong"
	if task.creationTime.After(goal) {
		task.result = "task has been successed"
		task.successed = true
	}

	task.finishedTime = time.Now()

	time.Sleep(time.Millisecond * 150)

	return task
}
