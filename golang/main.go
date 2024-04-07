package main

import (
	"log"
	"os"
	"runtime"
	"time"
)

const (
	WorkingTimeSeconds   float64 = 3
	SomeErrorMessage     string  = "some error occurred"
	WorkerExecutionError string  = "worker error occurred"
	InfoAbsenceMessage   string  = "none"
)

type Task struct {
	id         int64
	createTime time.Time
}

func main() {
	stdOut := log.New(os.Stdout, "DONE: ", log.LstdFlags)
	stdErr := log.New(os.Stderr, "FAIL: ", log.LstdFlags)

	workersAmount := runtime.NumCPU()
	tasksQueue := make(chan Task, workersAmount)

	go TaskCreator(tasksQueue, stdErr)

	for task := range tasksQueue {
		go RunWorker(task, stdOut, stdErr)
	}
}

func TaskCreator(tasksQueue chan Task, stdErr *log.Logger) {
	startTime := time.Now()
	for time.Since(startTime).Seconds() < WorkingTimeSeconds {
		timeNow := time.Now()
		newTask := Task{
			id:         timeNow.Unix(),
			createTime: timeNow,
		}
		if timeNow.Nanosecond()&1 == 1 {
			stdErr.Printf(
				"Task id: %d create time: %s, error: %s",
				newTask.id,
				InfoAbsenceMessage, SomeErrorMessage,
			)
		} else {
			tasksQueue <- newTask
		}
	}
	close(tasksQueue)
}

func RunWorker(task Task, stdOut, stdErr *log.Logger) {
	if task.createTime.After(time.Now().Add(-20 * time.Second)) {
		stdOut.Printf(
			"Task id: %d create time: %s, finish time: %s",
			task.id,
			task.createTime.Format(time.RFC3339Nano),
			time.Now().Format(time.RFC3339Nano),
		)
	} else {
		stdErr.Printf(
			"Task id: %d create time: %s, error: %s",
			task.id,
			task.createTime.Format(time.RFC3339Nano),
			WorkerExecutionError,
		)
	}
}
