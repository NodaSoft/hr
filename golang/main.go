package main

import (
	"fmt"
	"sync"
	"time"
)

type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex

	tasks := make(chan Ttype, 10)
	doneTasks := make([]Ttype, 0)
	errorTasks := make([]error, 0)

	taskCreator := func(tasks chan<- Ttype) {
		defer close(tasks)
		for i := 0; i < 10; i++ {
			ft := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 {
				ft = "Some error occurred"
			}
			tasks <- Ttype{cT: ft, id: int(time.Now().UnixNano())}
			time.Sleep(time.Second)
		}
	}

	taskWorker := func(task Ttype) {
		defer wg.Done()
		if task.cT == "Some error occurred" {
			task.taskRESULT = []byte("something went wrong")
		} else {
			task.taskRESULT = []byte("task has been successed")
		}
		task.fT = time.Now().Format(time.RFC3339Nano)

		mu.Lock()
		if string(task.taskRESULT) == "task has been successed" {
			doneTasks = append(doneTasks, task)
		} else {
			errorTasks = append(errorTasks, fmt.Errorf("Task id %d time %s, error %s", task.id, task.cT, task.taskRESULT))
		}
		mu.Unlock()
	}

	go taskCreator(tasks)

	for task := range tasks {
		wg.Add(1)
		go taskWorker(task)
	}

	ticker := time.NewTicker(3 * time.Second)
	go func() {
		for range ticker.C {
			printResults(&mu, doneTasks, errorTasks)
		}
	}()

	wg.Wait()
	ticker.Stop()

	printResults(&mu, doneTasks, errorTasks)
}

func printResults(mu *sync.Mutex, doneTasks []Ttype, errorTasks []error) {
	mu.Lock()
	defer mu.Unlock()

	fmt.Println("Done tasks:")
	for _, task := range doneTasks {
		fmt.Printf("ID: %d, Created: %s, Finished: %s, Result: %s\n", task.id, task.cT, task.fT, task.taskRESULT)
	}

	fmt.Println("Errors:")
	for _, err := range errorTasks {
		fmt.Println(err)
	}
}
