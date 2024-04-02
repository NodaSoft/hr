package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const countTasks = 10

type Task struct {
	id               int64
	time_of_creation time.Time
	time_completed   time.Time
	task_result      string
	task_error       error
}

func main() {

	tasks := make(chan Task, countTasks)
	doneTasks := make(chan Task, countTasks)
	undoneTasks := make(chan error, countTasks)

	for i := 0; i < countTasks; i++ {
		go taskCreturer(tasks)
	}

	time.Sleep(time.Second * 2)
	close(tasks)

	wg := sync.WaitGroup{}

	for t := range tasks {
		wg.Add(1)
		go task_worker(t, doneTasks, undoneTasks, &wg)
	}

	wg.Wait()
	close(doneTasks)
	close(undoneTasks)

	println("Done tasks:")
	for r := range doneTasks {
		fmt.Printf("task id %d  task_result: %s\n", r.id, r.task_result)
	}

	println("Errors:")
	for r := range undoneTasks {
		println(r.Error())
	}
}

func taskCreturer(a chan Task) {

	var err error
	tn := time.Now()
	if tn.Nanosecond()%2 > 0 {
		err = errors.New("some error occured")
	}
	a <- Task{
		time_of_creation: tn,
		id:               time.Now().UnixNano(),
		task_error:       err}

}

func task_worker(t Task, doneTasks chan Task, undoneTasks chan error, wg *sync.WaitGroup) {

	defer wg.Done()

	if t.time_of_creation.After(time.Now().Add(-20 * time.Second)) {
		t.task_result = "task has been successed"
	} else if t.task_error == nil {
		t.task_error = errors.New("something went wrong")
	}

	t.time_completed = time.Now()

	if t.task_error != nil {
		undoneTasks <- fmt.Errorf("task id %d time %s, error %w", t.id, t.time_of_creation, t.task_error)

	} else {
		doneTasks <- t
	}
}