package main

import (
	"fmt"
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

// Assumptions:
// 1) Task is created in the different microservice and
//    is passed to us through some event bus.
//	  But here we simulate it via TaskCreator
// 2) Task result must have some sort of helpful message

type ResultType int64

const (
	Success ResultType = iota
	Failure
)

type TaskResult struct {
	status  ResultType
	message []byte
}

// Task represents a task in our system
type Task struct {
	ID           int
	CreationTime string // время создания
	FinishTime   string // время выполнения
	TaskResult   TaskResult
}

// CreationTime and FinishTime are assumed to be from an external event (1)

func main() {
	tasks := TaskCreator(3 * time.Second)
	doneTasks, undoneTasks := Work(tasks)
	results, err := CollectResults(doneTasks, undoneTasks)

	PrintResults(err, results)
}

func TaskCreator(duration time.Duration) chan Task {
	tasks := make(chan Task, 10)
	end := time.Now().Add(duration)
	go func() {
		for {
			now := time.Now()
			if now.After(end) {
				close(tasks)
				break
			}
			ft := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				// Nanoseconds usually are errorous because
				// system clock doesn't send interrupts on nanosecond level.
				// To imitate the nanoseconds systems usually try to estimate them.
				// For example in Linux it is done using cpu ticks
				ft = "Some error occured"
			}

			tasks <- Task{CreationTime: ft, ID: int(time.Now().Unix())} // передаем таск на выполнение

			// This line was added for manual testing. It can be uncomment if you don't want your terminal being loaded with messages
			// time.Sleep(1 * time.Millisecond)
		}
	}()
	return tasks
}

func TaskWorker(task Task) Task {
	_, err := time.Parse(time.RFC3339, task.CreationTime)
	if err != nil {
		task.TaskResult = TaskResult{status: Failure, message: []byte("something went wrong")}
	} else {
		task.TaskResult = TaskResult{status: Success, message: []byte("task has been successed")}
	}
	task.FinishTime = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)

	return task
}

func TaskSorter(task Task, doneTasks chan<- Task, undoneTasks chan<- error) {
	switch task.TaskResult.status {
	case Success:
		doneTasks <- task
	case Failure:
		undoneTasks <- fmt.Errorf("task id %d time %s, error %s", task.ID, task.CreationTime, task.TaskResult.message)
	default:
		return
	}
}

func Work(tasks <-chan Task) (<-chan Task, <-chan error) {
	doneTasks := make(chan Task)
	undoneTasks := make(chan error)

	go func() {
		var wg sync.WaitGroup
		defer close(doneTasks)
		defer close(undoneTasks)
		for task := range tasks {
			wg.Add(1)
			go func(task Task) {
				defer wg.Done()
				task = TaskWorker(task)
				TaskSorter(task, doneTasks, undoneTasks)
			}(task)
		}
		wg.Wait()
	}()
	return doneTasks, undoneTasks
}

func CollectResults(doneTasks <-chan Task, undoneTasks <-chan error) (map[int]Task, []error) {
	results := map[int]Task{}
	err := []error{}
	for {
		select {
		case result, ok := <-doneTasks:
			if !ok {
				doneTasks = nil
			} else {
				results[result.ID] = result
			}
		case result, ok := <-undoneTasks:
			if !ok {
				undoneTasks = nil
			} else {
				err = append(err, result)
			}
		}
		if doneTasks == nil && undoneTasks == nil {
			break
		}
	}

	return results, err
}

func PrintResults(err []error, result map[int]Task) {
	fmt.Println("Errors:")
	for _, r := range err {
		fmt.Println(r)
	}

	fmt.Println("Done tasks:")
	for _, r := range result {
		fmt.Println(r)
	}
}
