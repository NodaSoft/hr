package main

import (
	"context"
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


func taskCreator(ctx context.Context, wg *sync.WaitGroup, tasks chan<- Ttype) {
	defer wg.Done()
	taskID := 1

	for {
		select {
		case <-ctx.Done():
			return
		default:
			ft := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 {
				ft = "Some error occurred"
			}
			tasks <- Ttype{id: taskID, cT: ft}
			taskID++
			time.Sleep(500 * time.Millisecond) 
		}
	}
}

func taskWorker(wg *sync.WaitGroup, tasks <-chan Ttype, doneTasks chan<- Ttype, undoneTasks chan<- error) {
	defer wg.Done()

	for task := range tasks {
		// Simulate task processing
		task.fT = time.Now().Format(time.RFC3339Nano)
		parsedTime, err := time.Parse(time.RFC3339, task.cT)
		if err != nil || parsedTime.Before(time.Now().Add(-20*time.Second)) {
			task.taskRESULT = "something went wrong"
			undoneTasks <- fmt.Errorf("Task ID %d: %s", task.id, task.taskRESULT)
		} else {
			task.taskRESULT = "task has been successful"
			doneTasks <- task
		}
		time.Sleep(150 * time.Millisecond) 
	}
}

func resultPrinter(ctx context.Context, doneTasks <-chan Ttype, undoneTasks <-chan error) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fmt.Println("Processed Tasks:")
			for len(doneTasks) > 0 {
				task := <-doneTasks
				fmt.Printf("Success: Task ID %d, Created at %s, Finished at %s\n", task.id, task.cT, task.fT)
			}
			for len(undoneTasks) > 0 {
				err := <-undoneTasks
				fmt.Println("Error:", err)
			}
		}
	}
}

func main() {
	tasks := make(chan Ttype, 10)
	doneTasks := make(chan Ttype, 10)
	undoneTasks := make(chan error, 10)

	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	wg.Add(1)
	go taskCreator(ctx, &wg, tasks)

	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go taskWorker(&wg, tasks, doneTasks, undoneTasks)
	}

	go resultPrinter(ctx, doneTasks, undoneTasks)

	wg.Wait()
	close(tasks)
	close(doneTasks)
	close(undoneTasks)
}
