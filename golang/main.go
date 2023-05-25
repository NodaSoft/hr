package main

import (
	"context"
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

const (
	produceTime = time.Second * 3
)

// A Task represents a meaninglessness of our life
type Task struct {
	id         int
	createTime time.Time // время создания
	finishTime time.Time // время выполнения
	taskResult []byte
	err        error
}

type TaskStorage struct {
	mu    sync.Mutex
	tasks map[int]Task
}

func (ts *TaskStorage) add(task Task) {
	ts.mu.Lock()
	ts.tasks[task.id] = task
	ts.mu.Unlock()
}

func newTaskStorage() *TaskStorage {
	return &TaskStorage{
		mu:    sync.Mutex{},
		tasks: map[int]Task{},
	}
}

type ErrorStorage struct {
	mu     sync.Mutex
	errors []error
}

func (es *ErrorStorage) addFromTask(task Task) {
	err := fmt.Errorf("task id %d time %s, error %w", task.id, task.createTime, task.err)
	es.mu.Lock()
	es.errors = append(es.errors, err)
	es.mu.Unlock()
}

func newErrorStorage() *ErrorStorage {
	return &ErrorStorage{
		mu:     sync.Mutex{},
		errors: []error{},
	}
}

type addTask func(task Task)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), produceTime)
	defer cancel()

	doneTaskPipeline, unDoneTaskPipeline := taskSorter(taskConsumer(taskProducer(ctx.Done())))

	wg := sync.WaitGroup{}
	taskStorage := newTaskStorage()
	errStorage := newErrorStorage()

	wg.Add(1)
	go func() {
		defer wg.Done()
		taskAdd(doneTaskPipeline, taskStorage.add)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		taskAdd(unDoneTaskPipeline, errStorage.addFromTask)
	}()

	wg.Wait()
	println("Errors:")
	for _, err := range errStorage.errors {
		println(err.Error())
	}

	println("Done tasks:")
	for id := range taskStorage.tasks {
		println(id)
	}
}

func taskProducer(done <-chan struct{}) <-chan Task {
	var uniqTaskID int32
	taskCh := make(chan Task)
	go func() {
		defer close(taskCh)
		for {
			select {
			case <-done:
				return
			default:
				var err error
				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					err = errors.New("taskProducer: some error occurred")
				}

				taskCh <- Task{
					id:         int(atomic.AddInt32(&uniqTaskID, 1)),
					createTime: time.Now(),
					err:        err,
				} // передаем таск на выполнение
			}
		}
	}()
	return taskCh
}

func taskConsumer(taskCh <-chan Task) <-chan Task {
	resTaskCh := make(chan Task)
	go func() {
		defer close(resTaskCh)
		for task := range taskCh {
			if task.err == nil {
				if task.createTime.After(time.Now().Add(-20 * time.Second)) {
					task.taskResult = []byte("task has been succeeded")
				} else {
					task.taskResult = []byte("something went wrong")
					task.err = errors.New("taskConsumer: something went wrong")
				}
			} else {
				task.taskResult = []byte("something went wrong")
			}

			task.finishTime = time.Now()
			time.Sleep(time.Millisecond * 150)
			resTaskCh <- task
		}
	}()

	return resTaskCh
}

func taskSorter(taskCh <-chan Task) (<-chan Task, <-chan Task) {
	doneTasksCh := make(chan Task)
	unDoneTasksCh := make(chan Task)
	go func() {
		defer close(doneTasksCh)
		defer close(unDoneTasksCh)
		for task := range taskCh {
			if task.err == nil {
				doneTasksCh <- task
			} else {
				unDoneTasksCh <- task
			}
		}
	}()
	return doneTasksCh, unDoneTasksCh
}

func taskAdd(taskCh <-chan Task, add addTask) {
	wg := sync.WaitGroup{}
	for task := range taskCh {
		wg.Add(1)
		go func(task Task) {
			defer wg.Done()
			add(task)
		}(task)
	}
	wg.Wait()
}
