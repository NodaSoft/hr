package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Ttype struct {
	id         int
	cT         time.Time
	fT         time.Time
	taskRESULT []byte
	error      bool
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	taskCreator := func(ctx context.Context, tasks chan<- Ttype) {
		for {
			select {
			case <-ctx.Done():
				close(tasks)
				return
			default:
				errOccurred := time.Now().Nanosecond()%2 > 0
				task := Ttype{
					cT:    time.Now(),
					id:    int(time.Now().Unix()),
					error: errOccurred,
				}
				tasks <- task
				time.Sleep(time.Millisecond * 100) // Generate tasks every 100 ms
			}
		}
	}

	taskProcessor := func(ctx context.Context, tasks <-chan Ttype, doneTasks chan<- Ttype, wg *sync.WaitGroup) {
		defer wg.Done()
		for task := range tasks {
			if task.error {
				task.taskRESULT = []byte("some error occurred")
			} else {
				task.taskRESULT = []byte("task has been succeeded")
				task.fT = time.Now()
			}
			doneTasks <- task
		}
	}

	displayResults := func(ctx context.Context, tasks <-chan Ttype) {
		ticker := time.NewTicker(3 * time.Second)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				var successCount, errorCount int
				for len(tasks) > 0 {
					task := <-tasks
					if task.error {
						errorCount++
					} else {
						successCount++
					}
				}
				fmt.Printf("Success: %d, Errors: %d\n", successCount, errorCount)
			}
		}
	}

	taskChan := make(chan Ttype, 100)
	doneTaskChan := make(chan Ttype, 100)
	var wg sync.WaitGroup

	wg.Add(1)
	go taskProcessor(ctx, taskChan, doneTaskChan, &wg)

	go taskCreator(ctx, taskChan)

	go displayResults(ctx, doneTaskChan)

	wg.Wait()
}
