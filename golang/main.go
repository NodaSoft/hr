package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	BufSize      = 10
	WorkersCount = 5
)

type Task struct {
	Id       int64
	Created  time.Time
	Finished time.Time
	Success  bool
	Err      error
}

func taskCreator(ch chan Task, wg *sync.WaitGroup) {
	var err error
	timeNow := time.Now()

	if timeNow.Nanosecond()%2 > 0 {
		err = errors.New("some error occurred")
	}

	ch <- Task{
		Id: time.Now().UnixNano(), Created: timeNow, Err: err}

	wg.Done()
}

func taskWorker(taskChan chan Task, doneTasks chan Task, undoneTasks chan Task, wg *sync.WaitGroup) {
	for task := range taskChan {
		if task.Err == nil {
			if task.Created.After(time.Now().Add(-20 * time.Second)) {
				task.Finished = time.Now()
				task.Success = true
			} else {
				task.Err = errors.New("something went wrong")
			}
		}

		taskSorter(task, doneTasks, undoneTasks)
	}
	wg.Done()
}

func taskSorter(task Task, doneTasks chan Task, undoneTasks chan Task) {
	if task.Err == nil {
		doneTasks <- task
	} else {
		undoneTasks <- task
	}
}

func main() {
	tasksChan := make(chan Task, WorkersCount)
	doneTasks := make(chan Task, BufSize)
	undoneTasks := make(chan Task, BufSize)
	var (
		wg  sync.WaitGroup
		wg2 sync.WaitGroup
	)

	for range WorkersCount {
		wg.Add(1)
		go taskWorker(tasksChan, doneTasks, undoneTasks, &wg)
	}

	startTime := time.Now()
	for time.Since(startTime).Seconds() < 3 {
		wg2.Add(1)
		go taskCreator(tasksChan, &wg2)
		time.Sleep(time.Millisecond * 150)
	}

	go func() {
		wg2.Wait()
		close(tasksChan)
	}()

	var (
		successTasks []Task
		errorTasks   []Task
	)

	go func() {
		for r := range doneTasks {
			successTasks = append(successTasks, r)
		}
	}()

	go func() {
		for r := range undoneTasks {
			errorTasks = append(errorTasks, r)
		}
	}()

	wg.Wait()
	close(doneTasks)
	close(undoneTasks)

	fmt.Println("Done tasks:")
	for i, j := range successTasks {
		fmt.Printf("%d | Task id %d | created at %s | finished at %s | success: %t \n",
			i, j.Id, j.Created, j.Finished, j.Success)
	}

	fmt.Println("Errors:")
	for i, j := range errorTasks {
		fmt.Printf("%d | Task id %d | created at %s | finished at %s | error: %s \n",
			i, j.Id, j.Created, j.Finished, j.Err)
	}
}
