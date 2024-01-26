package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

var (
	InitialTaskError = errors.New("some error occurred")
	TimeoutTaskError = errors.New("timeout error")
)

// Task represents the complete meaningfulness of our lives
type Task struct {
	id         int64
	createTime time.Time
	finishTime time.Time
	err        error
	res        []byte
}

// SuccessTaskResult thread-safe map with ability to write
type SuccessTaskResult struct {
	mutex  sync.Mutex // Protects resMap.
	resMap map[int64]Task
}

func (res *SuccessTaskResult) Write(tsk Task) {
	res.mutex.Lock()
	res.resMap[tsk.id] = tsk
	res.mutex.Unlock()
}

type ErrorsResult struct {
	mutex  sync.Mutex // Protects result.
	result []error
}

func (res *ErrorsResult) Write(err error) {
	res.mutex.Lock()
	res.result = append(res.result, err)
	res.mutex.Unlock()
}

func main() {
	// general program context
	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	// create channels
	taskCh := make(chan Task)
	finishedTaskCh := make(chan Task) // pipeline channel
	successTaskCh := make(chan Task)
	errorsCh := make(chan error)

	// context with timeout for generating task
	taskGenCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// start generating task
	go generateTasks(taskGenCtx, taskCh)

	// process tasks
	go startWorker(ctx, 0, taskCh, finishedTaskCh)

	// sort tasks
	go sortFinishedTask(finishedTaskCh, successTaskCh, errorsCh)

	wg := sync.WaitGroup{}

	// write success result to map
	successResult := SuccessTaskResult{resMap: map[int64]Task{}, mutex: sync.Mutex{}}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for res := range successTaskCh {
			wg.Add(1)
			go func(tsk Task) {
				defer wg.Done()
				successResult.Write(tsk)
			}(res)
		}
	}()

	// write error result to slice
	errorsResult := ErrorsResult{result: make([]error, 0), mutex: sync.Mutex{}}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for err := range errorsCh {
			wg.Add(1)
			go func(err error) {
				defer wg.Done()
				errorsResult.Write(err)
			}(err)
		}
	}()

	wg.Wait()

	fmt.Println("Errors:")
	for _, err := range errorsResult.result {
		fmt.Println(err)
	}

	fmt.Println("Done tasks:")
	for k, v := range successResult.resMap {
		fmt.Println(k, v)
	}
}

func generateTasks(ctx context.Context, ch chan<- Task) {
	defer func() {
		close(ch)
		log.Printf("task channel was closed")
	}()

	taskCounter := 0
	for {
		tsk := createTask()

		select {
		case ch <- tsk:
			taskCounter++
		case <-ctx.Done():
			log.Printf("was added %d tasks", taskCounter)
			return
		}
	}
}

func createTask() Task {
	tsk := Task{
		id:         generateId(),
		createTime: time.Now(),
	}

	if isErrorOccurred(time.Now()) {
		tsk.err = InitialTaskError
	}

	return tsk
}

func generateId() int64 {
	return time.Now().UnixNano()
}

func isErrorOccurred(currentTime time.Time) bool {
	return currentTime.Nanosecond()%2 > 0
}

func startWorker(ctx context.Context, number int, tasks <-chan Task, finishedTasks chan<- Task) {
	log.Printf("task worker %d was started", number)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Context canceled. Task worker %d is ended", number)
			return
		case tsk, ok := <-tasks:
			// if channel was closed
			if !ok {
				close(finishedTasks)
				log.Printf("finished tasks channel was closed")

				log.Printf("Task channel was closed. Task worker %d is ended", number)
				return
			}

			finishedTasks <- processTask(tsk)
		}
	}
}

func processTask(tsk Task) Task {
	if errors.Is(tsk.err, InitialTaskError) {
		tsk.res = []byte("something went wrong")
		tsk.finishTime = time.Now()

		return tsk
	}

	if isTaskTimeout(tsk.createTime, 20*time.Millisecond) {
		tsk.res = []byte("something went wrong")
		tsk.err = TimeoutTaskError
	} else {
		time.Sleep(150 * time.Millisecond) // some process
		tsk.res = []byte("task has been succeeded")
	}

	tsk.finishTime = time.Now()

	return tsk
}

func isTaskTimeout(taskCreateTime time.Time, timeout time.Duration) bool {
	return time.Now().Sub(taskCreateTime) > timeout
}

func sortFinishedTask(finishedTasks <-chan Task, successTaskCh chan<- Task, errorsCh chan<- error) {
	defer func() {
		close(successTaskCh)
		log.Printf("success task channel was closed")

		close(errorsCh)
		log.Printf("errors channel was closed")
	}()

	for tsk := range finishedTasks {
		if tsk.err != nil {
			errorsCh <- fmt.Errorf("task: id=%d time=%s error: %w", tsk.id, tsk.createTime, tsk.err)
		} else {
			successTaskCh <- tsk
		}
	}
}
