package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	id         int
	createdAt  time.Time
	finishedAt time.Time
	result     string
}

const (
	workersCount = 5
	tasksCount   = 50
)

func worker(id int,
	tasksChannel <-chan Task,
	successResChannel chan<- Task,
	failureResChannel chan<- Task) {

	for t := range tasksChannel {
		t.finishedAt = time.Now()
		if t.finishedAt.Nanosecond()%2 > 0 {
			t.result = "failure"
			failureResChannel <- t
		} else {
			t.result = "success"
			successResChannel <- t
		}
	}
}

func main() {
	tasksChannel := make(chan Task, tasksCount)
	successResChannel := make(chan Task, tasksCount)
	failureResChannel := make(chan Task, tasksCount)

	for t := 1; t <= tasksCount; t++ {
		tasksChannel <- Task{id: t, createdAt: time.Now()}
	}
	close(tasksChannel)

	var wg sync.WaitGroup
	wg.Add(workersCount)

	for w := 0; w < workersCount; w++ {
		go func(id int) {
			defer wg.Done()
			worker(id, tasksChannel, successResChannel, failureResChannel)
		}(w)
	}

	wg.Wait()

	close(successResChannel)
	close(failureResChannel)

	fmt.Println("Done tasks:")
	for t := range successResChannel {
		fmt.Println(t.id)
	}

	fmt.Println("Errors:")
	for t := range failureResChannel {
		fmt.Printf("Task id %d time %s, error %s\n", t.id, t.createdAt, t.result)
	}
}
