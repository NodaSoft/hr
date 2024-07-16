package main

import (
	"fmt"
	"sync"
	"time"
)

type Ttype struct {
	id         int
	cT         string
	fT         string
	taskRESULT []byte
}

func main() {
	var wg sync.WaitGroup
	taskChan := make(chan Ttype)
	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	wg.Add(1)
	go func() {
		defer wg.Done()
		taskCreator(taskChan)
	}()

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			taskWorker(taskChan, doneTasks, undoneTasks)
		}()
	}

	go func() {
		wg.Wait()
		close(taskChan)
		close(doneTasks)
		close(undoneTasks)
	}()

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	result := make(map[int]Ttype)
	var errors []error

	for {
		select {
		case task := <-doneTasks:
			result[task.id] = task
		case err := <-undoneTasks:
			errors = append(errors, err)
		case <-ticker.C:
			printResults(result, errors)
		}

		if len(result)+len(errors) == 10 {
			break
		}
	}

	printResults(result, errors)
}

func taskCreator(taskChan chan<- Ttype) {
	for i := 0; i < 10; i++ {
		ft := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 {
			ft = "Some error occurred"
		}
		taskChan <- Ttype{cT: ft, id: int(time.Now().Unix())}
		time.Sleep(time.Second)
	}
}

func taskWorker(taskChan <-chan Ttype, doneTasks chan<- Ttype, undoneTasks chan<- error) {
	for task := range taskChan {
		tt, _ := time.Parse(time.RFC3339, task.cT)
		if tt.After(time.Now().Add(-20 * time.Second)) {
			task.taskRESULT = []byte("task has been successful")
		} else {
			task.taskRESULT = []byte("something went wrong")
		}
		task.fT = time.Now().Format(time.RFC3339Nano)
		time.Sleep(150 * time.Millisecond)

		if string(task.taskRESULT[14:]) == "successful" {
			doneTasks <- task
		} else {
			undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", task.id, task.cT, task.taskRESULT)
		}
	}
}

func printResults(result map[int]Ttype, errors []error) {
	fmt.Println("Errors:")
	for _, err := range errors {
		fmt.Println(err)
	}

	fmt.Println("Done tasks:")
	for _, task := range result {
		fmt.Printf("Task id: %d, creation time: %s, finish time: %s\n", task.id, task.cT, task.fT)
	}
	fmt.Println()
}