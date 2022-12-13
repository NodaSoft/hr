package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	OK    = "ok"
	ERROR = "error"
)

type task struct {
	id         int
	createTime time.Time
	endTime    time.Time
	status     string
	result     []byte
}

func taskGenerator() chan task {
	timeout := time.After(2 * time.Second)

	tasksNewChan := make(chan task, 10)

	go func() {
		i := 0
		for {
			select {
			case <-timeout:
				fmt.Println("Generator is closed")
				close(tasksNewChan)

				return
			default:
				timeNow := time.Now()
				newTask := task{
					//id: int(timeNow.Unix()),
					id:         i,
					createTime: timeNow,
				}

				newTask.status = OK
				if timeNow.Nanosecond()%2 > 0 {
					newTask.status = ERROR
					newTask.result = []byte("some error occured")
				} else if timeNow.Second()%5 > 0 && i%2 > 0 { // Добавил для появления тасков у которых something went wrong
					newTask.createTime = timeNow.Add(-25 * time.Second)
				}

				tasksNewChan <- newTask
			}
			i++
		}
	}()

	return tasksNewChan
}

func taskProcess(tasksChan chan task) chan task {
	taskProcessedChan := make(chan task, 10)

	go func(taskProcessedChan chan task) {
		for t := range tasksChan {
			if t.status == OK {
				if t.createTime.After(time.Now().Add(-20 * time.Second)) {
					time.Sleep(time.Millisecond * 150)

					t.result = []byte("task has been successed")
					t.status = OK
				} else {
					t.result = []byte("something went wrong")
					t.status = ERROR
				}
			}

			t.endTime = time.Now()

			taskProcessedChan <- t
		}
		close(taskProcessedChan)
	}(taskProcessedChan)

	return taskProcessedChan
}

func taskSort(tasksChan chan task) (map[int]task, []error) {
	doneTasks := map[int]task{}
	errors := []error{}

	var wg sync.WaitGroup
	var mu sync.Mutex

	for t := range tasksChan {
		wg.Add(1)

		go func() {
			defer wg.Done()
			switch t.status {
			case OK:
				mu.Lock()
				doneTasks[t.id] = t
				mu.Unlock()
			case ERROR:
				mu.Lock()
				errors = append(errors, fmt.Errorf("task id %d time %s, error %s", t.id, t.createTime.Format(time.RFC3339), t.result))
				mu.Unlock()
			}
		}()

		wg.Wait()
	}

	return doneTasks, errors
}

func main() {

	tasksChan := taskProcess(taskGenerator())

	doneTasks, errors := taskSort(tasksChan)

	fmt.Println("Errors:")
	for _, err := range errors {
		fmt.Println(err)
	}

	fmt.Println("Done tasks:")
	for task := range doneTasks {
		fmt.Println(task)
	}
}
