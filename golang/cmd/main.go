package main

import (
	"context"
	"fmt"
	"os"
	meaninglessTask "test_task/internal/tasks/meaningless"
	"test_task/pkg/workerPool"
	"time"
)

const (
	MAX_WORKER_COUNT      = 10
	LOG_INTERVAL          = 3 * time.Second
	TASK_CREATION_TIMEOUT = 10 * time.Second
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), TASK_CREATION_TIMEOUT)
	defer cancel()

	wp := workerPool.New(MAX_WORKER_COUNT, ctx)
	wp.WithFactory(func() workerPool.Task {
		return meaninglessTask.New()
	})

	wp.Wg.Add(1)

	errChan := make(chan error, 1)
	go wp.StartFactory(errChan)

	go func() {
		if err := <-errChan; err != nil {
			fmt.Printf("Error starting task factory: %v\n", err)
			os.Exit(1)
		}
	}()

	for i := 0; i <= MAX_WORKER_COUNT-1; i++ {
		wp.Wg.Add(1)
		go wp.ProcessAndSortTask()
	}

	ticker := time.NewTicker(LOG_INTERVAL)
	defer ticker.Stop()

	stopChan := make(chan struct{})
	go wp.PrintResults(ticker, stopChan)

	wp.Wg.Wait()

	fmt.Println("exiting...")
	go func() {
		stopChan <- struct{}{}
		close(stopChan)
		wp.Shutdown()
	}()
}
