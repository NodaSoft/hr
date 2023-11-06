package main

import (
	"context"
	"fmt"
	"strings"
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

const TASK_LIMIT = 10 // parallel tasks

type TaskStatus int

const (
	TASK_CREATED TaskStatus = iota
	TASK_FAILED
	TASK_FINISHED
)

const SUCESS_STATUS_STR = "task has been succeeded"
const FAILED_STATUS_STR = "something went wrong"
const SUCCESS_STR = "succeeded"

type Task struct {
	id         uint64
	cT         string // время создания
	fT         string // время выполнения
	taskResult string
	status     TaskStatus
}

func taskProducer(ctx context.Context) <-chan Task {
	dst := make(chan Task)

	go func() {
		defer close(dst)

		var taskIdCnt uint64 = 0
		for {
			select {
			case <-ctx.Done():
				return
			default:
				creationTime := time.Now().Format(time.RFC3339)
				status := TASK_CREATED
				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					status = TASK_FAILED
				}
				dst <- Task{cT: creationTime, id: taskIdCnt, status: status} // передаем таск на выполнение
				taskIdCnt++
			}
		}
	}()

	return dst
}

func executeTasks(ctx context.Context, tasksCh <-chan Task, taskLimit int) <-chan Task {
	tasksResults := make(chan Task)

	limitCh := make(chan struct{}, taskLimit)
	work := func(task Task) {
		tasksResults <- taskWorker(task)
		select {
		case <-ctx.Done():
		case <-limitCh:
		}
	}

	go func() {
		defer close(tasksResults)
		for {
			select {
			case <-ctx.Done():
				return
			case task, ok := <-tasksCh:
				if !ok {
					return
				}
				select {
				case <-ctx.Done():
					return
				case limitCh <- struct{}{}:
					go work(task)
				}
			}
		}
	}()

	return tasksResults
}

func taskWorker(task Task) Task {
	if task.status == TASK_CREATED {
		task.taskResult = SUCESS_STATUS_STR
	} else {
		task.taskResult = FAILED_STATUS_STR
	}

	task.fT = time.Now().Format(time.RFC3339Nano)
	time.Sleep(time.Millisecond * 150)
	task.status = TASK_FINISHED

	return task
}

func countTasks(ctx context.Context, tasksCh <-chan Task, tasksOk *atomic.Uint64, tasksFailed *atomic.Uint64, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-tasksCh:
			if !ok {
				return
			}
			if task.status == TASK_FINISHED && strings.Contains(task.taskResult, SUCCESS_STR) {
				tasksOk.Add(1)
			} else {
				tasksFailed.Add(1)
			}
		}
	}
}

func main() {

	timeout := 3 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	taskChan := taskProducer(ctx)
	finishedTasks := executeTasks(ctx, taskChan, TASK_LIMIT)

	var tasksOk atomic.Uint64
	var tasksFailed atomic.Uint64

	wg := sync.WaitGroup{}
	wg.Add(1)
	go countTasks(ctx, finishedTasks, &tasksOk, &tasksFailed, &wg)
	wg.Wait()

	fmt.Printf("Errors: %v\n", tasksFailed.Load())
	fmt.Printf("Done tasks: %v\n", tasksOk.Load())
}
