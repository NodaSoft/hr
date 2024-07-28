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
	taskRESULT string
}


func taskCreator(a chan<- Ttype) {
	for i := 0; i < 100; i++ {
		ft := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 { 
			ft = "Some error occurred"
		}
		a <- Ttype{id: i, cT: ft} 
		time.Sleep(100 * time.Millisecond) 
	}
	close(a)
}


func taskWorker(a Ttype) Ttype {
	_, err := time.Parse(time.RFC3339, a.cT)
	if err == nil {
		a.taskRESULT = "task has been succeeded"
	} else {
		a.taskRESULT = "something went wrong"
	}
	a.fT = time.Now().Format(time.RFC3339Nano)
	time.Sleep(150 * time.Millisecond) // simulate task processing time
	return a
}


func resultCollector(doneTasks <-chan Ttype, undoneTasks <-chan Ttype, wg *sync.WaitGroup) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	results := map[int]Ttype{}
	errors := []Ttype{}

	for {
		select {
		case task, ok := <-doneTasks:
			if !ok {
				doneTasks = nil
			} else {
				results[task.id] = task
			}
		case task, ok := <-undoneTasks:
			if !ok {
				undoneTasks = nil
			} else {
				errors = append(errors, task)
			}
		case <-ticker.C:
			fmt.Println("Done tasks:")
			for _, task := range results {
				fmt.Printf("Task ID: %d, Creation Time: %s, Finish Time: %s, Result: %s\n", task.id, task.cT, task.fT, task.taskRESULT)
			}
			fmt.Println("Errors:")
			for _, task := range errors {
				fmt.Printf("Task ID: %d, Creation Time: %s, Finish Time: %s, Result: %s\n", task.id, task.cT, task.fT, task.taskRESULT)
			}
		}

		if doneTasks == nil && undoneTasks == nil {
			break
		}
	}

	wg.Done()
}

func main() {
	superChan := make(chan Ttype, 10)
	doneTasks := make(chan Ttype, 10)
	undoneTasks := make(chan Ttype, 10)

	var wg sync.WaitGroup
	wg.Add(1)

	go taskCreator(superChan)

	
	go func() {
		for t := range superChan {
			processedTask := taskWorker(t)
			if processedTask.taskRESULT == "task has been succeeded" {
				doneTasks <- processedTask
			} else {
				undoneTasks <- processedTask
			}
		}
		close(doneTasks)
		close(undoneTasks)
	}()

	go resultCollector(doneTasks, undoneTasks, &wg)

	
	time.Sleep(10 * time.Second)
	close(superChan)


	wg.Wait()
}
