package main

import (
	"fmt"
	"time"
)

const (
	maxConcurrentTasks = 10
	workTime           = 30
)

type Task struct {
	id        int
	createdAt time.Time
}

func main() {
	startTime := time.Now()

	taskQueue := make(chan Task)

	go taskCreator(taskQueue, startTime)

	workerDone := make(chan bool)

	for i := 0; i < maxConcurrentTasks; i++ {
		go taskWorker(taskQueue, startTime, workerDone)
	}

	for i := 0; i < maxConcurrentTasks; i++ {
		<-workerDone
	}

	fmt.Println("Done.")
}

func taskCreator(c chan Task, startTime time.Time) {
	for time.Since(startTime).Seconds() < workTime {
		createdAt := time.Now()
		id := int(time.Now().Unix())
		task := Task{id: id, createdAt: createdAt}

		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			fmt.Printf("Task %d failed at creation time %s\n", task.id, task.createdAt)
		}

		c <- task
	}

	close(c)
}

func taskRunner(t Task, startTime time.Time) {
	if t.createdAt.Before(startTime.Add(20 * time.Second)) {
		fmt.Printf("Task %d (created at %s) is done at %s\n", t.id, t.createdAt, time.Now())
	} else {
		fmt.Printf("Task %d (created at %s) failed\n", t.id, t.createdAt)
	}
}

func taskWorker(taskQueue chan Task, startTime time.Time, workerDone chan bool) {
	for task := range taskQueue {
		taskRunner(task, startTime)
	}

	workerDone <- true
}
