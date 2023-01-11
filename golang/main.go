package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type TaskStatus string
type TaskFailedReason string

const (
	TaskStatusSuccess TaskStatus = "success"
	TaskStatusFailed  TaskStatus = "failed"
)

const (
	Timeout     = 3 * time.Second
	TickerCycle = 150 * time.Millisecond
)

const NullDivision TaskFailedReason = "nullDivision"
const NumWorkers = 10

type Task struct {
	id           int
	createdAt    string
	doneAt       string
	result       TaskStatus
	failedReason TaskFailedReason
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	taskChannel := make(chan Task, NumWorkers)
	go TaskCreator(ctx, taskChannel)

	successChannels := make([]<-chan Task, NumWorkers)
	errorChannels := make([]<-chan Task, NumWorkers)

	for i := 0; i < NumWorkers; i++ {
		successChannel, errorsChannel := Worker(ctx, taskChannel)
		successChannels[i] = successChannel
		errorChannels[i] = errorsChannel
	}

	successChannel := MergeChannels(successChannels)
	errorsChannel := MergeChannels(errorChannels)

	successTasks, errorTasks := TasksProcessing(successChannel, errorsChannel)

	fmt.Println("Errors:")
	for _, errorTask := range errorTasks {
		fmt.Println(errorTask)
	}

	fmt.Println("Done tasks:")
	for _, successTask := range successTasks {
		fmt.Println(successTask.failedReason)
	}
}

func TaskCreator(ctx context.Context, taskChannel chan Task) {
	defer close(taskChannel)
	for {
		start := time.Now()
		currentTime := start.Format(time.RFC3339Nano)
		task := Task{createdAt: currentTime, id: int(time.Now().Unix())}
		if time.Since(start).Nanoseconds()%2 > 0 {
			task.failedReason = NullDivision
		}

		select {
		case taskChannel <- task:
		case <-ctx.Done():
			return
		}
	}
}

func Worker(ctx context.Context, taskChannel <-chan Task) (<-chan Task, <-chan Task) {
	doneTasks := make(chan Task)
	errorTasks := make(chan Task)

	go func() {
		ticker := time.NewTicker(TickerCycle)
		defer close(doneTasks)
		defer close(errorTasks)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				task, _ := <-taskChannel
				if task.failedReason != NullDivision {
					task.result = TaskStatusSuccess
					doneTasks <- task
				} else {
					task.result = TaskStatusFailed
					errorTasks <- task
				}
				task.doneAt = time.Now().Format(time.RFC3339Nano)
			}
		}
	}()

	return doneTasks, errorTasks
}

func MergeChannels(channels []<-chan Task) <-chan Task {
	wg := sync.WaitGroup{}
	outer := make(chan Task)
	for _, channel := range channels {
		wg.Add(1)
		go func(inner <-chan Task) {
			for {
				task, ok := <-inner
				if ok == false {
					wg.Done()
					break
				}
				outer <- task
			}
		}(channel)
	}

	go func() {
		wg.Wait()
		close(outer)
	}()

	return outer
}

func TasksProcessing(successChannel <-chan Task, errorsChannel <-chan Task) ([]Task, []Task) {
	successTasks := make([]Task, 0)
	errorTasks := make([]Task, 0)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		for task := range successChannel {
			successTasks = append(successTasks, task)
		}
		wg.Done()
	}()

	go func() {
		for task := range errorsChannel {
			errorTasks = append(errorTasks, task)
		}
		wg.Done()
	}()

	wg.Wait()
	return successTasks, errorTasks
}
