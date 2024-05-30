package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	periodTasksGeneration   = 100 * time.Millisecond
	durationTasksGeneration = 10 * time.Second
	thresholdDurationTask   = 20 * time.Second
	delayDurationTaskWork   = 150 * time.Millisecond
)

var (
	errInvalidTask = errors.New("some error occurred")
	errTaskTimeout = errors.New("something went wrong")
)

type Task struct {
	Id          int
	CreatedAt   time.Time
	FinalyzedAt time.Time
	Result      []byte
	Err         error
}

func (t *Task) work() {
	if t.Err == nil {
		if t.CreatedAt.After(time.Now().Add(-thresholdDurationTask)) {
			t.Result = []byte("task has been successed")
		} else {
			t.Err = errTaskTimeout
		}
	}

	t.FinalyzedAt = time.Now()
	time.Sleep(delayDurationTaskWork)
}

func tasksCreator(outCh chan Task) {
	defer close(outCh)
	timeoutCh := time.After(durationTasksGeneration)

	for {
		createTime := time.Now()

		// Id generation algorithm has a collision risk
		createdTask := Task{
			Id:        int(createTime.UnixMicro()),
			CreatedAt: createTime,
		}

		// Task validity condition
		// Nanoseconds were replaced by microseconds, because otherwise the condition has no meaning
		if (createTime.UnixMicro() % 2) > 0 {
			createdTask.Err = errInvalidTask
		}

		select {
		case <-timeoutCh:
			return
		case outCh <- createdTask:
			// Delay to avoid unnecessary CPU load
			time.Sleep(periodTasksGeneration)
		}
	}
}

func main() {
	tasksCh := make(chan Task, 10)
	doneTasksCh := make(chan Task)

	// Creator launch
	go tasksCreator(tasksCh)

	// Worker launch
	wg := sync.WaitGroup{}
	go func() {
		for t := range tasksCh {
			wg.Add(1)
			go func() {
				defer wg.Done()
				t.work()
				doneTasksCh <- t
			}()
		}

		wg.Wait()
		close(doneTasksCh)
	}()

	// Sorting launch
	mtx := sync.Mutex{}
	successDoneTasks := []Task{}
	errorDoneTasks := []Task{}
	exitSignal := make(chan struct{})
	go func() {
		for t := range doneTasksCh {
			mtx.Lock()
			if t.Err == nil {
				successDoneTasks = append(successDoneTasks, t)
			} else {
				errorDoneTasks = append(errorDoneTasks, t)
			}
			mtx.Unlock()
		}

		// Send a signal about the end of task processing
		exitSignal <- struct{}{}
	}()

	// Launching a terminal output goroutine
	wgEnd := sync.WaitGroup{}
	wgEnd.Add(1)
	go func() {
		defer wgEnd.Done()
		for {
			time.Sleep(time.Second * 3)

			mtx.Lock()

			fmt.Println("Done tasks:")
			for _, t := range successDoneTasks {
				fmt.Printf("Task id: %d time: %s, result: %s\n", t.Id, t.CreatedAt.Format(time.RFC3339), t.Result)
			}

			fmt.Println("Error tasks:")
			for _, t := range errorDoneTasks {
				fmt.Printf("Task id: %d time: %s, error: %s\n", t.Id, t.CreatedAt.Format(time.RFC3339), t.Err)
			}

			// Clear out slices
			successDoneTasks = successDoneTasks[:0]
			errorDoneTasks = errorDoneTasks[:0]

			mtx.Unlock()

			select {
			case <-exitSignal:
				return
			default:
			}
		}
	}()
	wgEnd.Wait()
}
