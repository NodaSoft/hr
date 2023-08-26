package main

import (
	"fmt"
	"slices"
	"sync"
	"time"
)

const CH_BUFSIZE = 10

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id           int
	creationTime string // время создания
	finishTime   string // время выполнения
	taskResult   []byte
}

func taskCreturer(a chan Ttype) {
	go func() {
		for {
			ft := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 {
				ft = "Some error occured"
			}
			a <- Ttype{creationTime: ft, id: int(time.Now().Unix())}
		}
	}()
}

func handleTask(task Ttype) (Ttype, error) {
	if task.creationTime == "Some error occured" {
		task.taskResult = []byte("something went wrong")
	} else {
		var (
			createdTime time.Time
			err         error
		)

		if createdTime, err = time.Parse(time.RFC3339, task.creationTime); err != nil {
			return task, fmt.Errorf("Parsing error: %s", err)
		}
		if time.Now().Sub(createdTime) < (20 * time.Second) {
			task.taskResult = []byte("task has been successed")
		} else {
			task.taskResult = []byte("something went wrong")
		}
	}

	task.finishTime = time.Now().Format(time.RFC3339)
	time.Sleep(time.Millisecond * 150)
	return task, nil
}

func getTask(mainCh chan Ttype, doneTasks chan Ttype, undoneTasks chan error, doneCh chan struct{}) {
	go func() {
		for t := range mainCh {
			t, err := handleTask(t)
			if err != nil {
				undoneTasks <- fmt.Errorf("Error handling task: %s", err)
			}

			taskSorter(t, doneTasks, undoneTasks)
		}
	}()

	<-doneCh
	close(undoneTasks)
	close(doneTasks)
}

func taskSorter(task Ttype, doneTasks chan Ttype, undoneTasks chan error) {
	if slices.Equal(task.taskResult, []byte("something went wrong")) {
		doneTasks <- task
		return
	} else {
		undoneTasks <- fmt.Errorf(
			"Task id %d time %s, error %s",
			task.id,
			task.creationTime,
			task.taskResult)
	}
}

func main() {
	superChan := make(chan Ttype, CH_BUFSIZE)
	doneHandling := make(chan struct{})

	go taskCreturer(superChan)

	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	go getTask(superChan, doneTasks, undoneTasks, doneHandling)

	result := make([]Ttype, 0)
	errors := make([]error, 0)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for r := range doneTasks {
			result = append(result, r)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range undoneTasks {
			errors = append(errors, err)
		}
	}()

	time.Sleep(time.Second * 3)
	doneHandling <- struct{}{}

	wg.Wait()

	println("Errors:")
	for _, e := range errors {
		fmt.Println(e)
	}

	println("Done tasks:")
	for _, r := range result {
		fmt.Println(r.id)
	}
}
