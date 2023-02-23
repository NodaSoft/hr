package main

import (
	"context"
	"errors"
	"fmt"
	"runtime"
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

type Task struct {
	id             int64
	creationTime   time.Time // время создания
	processingTime time.Time // время выполнения
}

func (t *Task) String() string {
	return fmt.Sprintf(
		"Id: %d. Created at: %s. Processing At: %s",
		t.id,
		t.creationTime.Format(time.RFC3339),
		t.processingTime.Format(time.RFC3339))
}

type FailedTask struct {
	task Task
	err  error
}

func (f *FailedTask) String() string {
	return fmt.Sprintf("%s Err: %s", f.task.String(), f.err.Error())
}

func taskCreator(ctx context.Context, inputTasks chan Task) {
	for {
		select {
		case <-ctx.Done():
			close(inputTasks)
			return
		default:
			creationTime := time.Now()
			task := Task{
				id: creationTime.UnixMilli(),
			}
			/// Nanoseconds will always return true, my suggestion is use Milliseconds.
			/// May be this is a mistake in your code(processor take a delay in Milliseconds)
			if creationTime.Nanosecond()%2 == 0 {
				task.creationTime = creationTime
			}

			inputTasks <- task
		}
	}
}

func taskProcessor(input chan Task, successTasks chan Task, failureTasks chan FailedTask, wg *sync.WaitGroup) {
	for t := range input {
		t.processingTime = time.Now()
		if t.creationTime.After(time.Now().Add(-20 * time.Second)) {
			successTasks <- t
		} else {
			failureTasks <- FailedTask{
				task: t,
				err:  errors.New("some processing error"),
			}
		}

		time.Sleep(150 * time.Millisecond)
	}
	wg.Done()
}

type SuccessStorage struct {
	mtx     *sync.Mutex
	success []Task
}

func (ss *SuccessStorage) Collect(tasks chan Task, wg *sync.WaitGroup) {
	for task := range tasks {
		ss.mtx.Lock()
		ss.success = append(ss.success, task)
		ss.mtx.Unlock()
	}
	wg.Done()
}

func (ss *SuccessStorage) Print() {
	fmt.Println("Successes: ")
	for _, s := range ss.success {
		fmt.Println(s.String())
	}
}

type FailureStorage struct {
	mtx     *sync.Mutex
	failure []FailedTask
}

func (fs *FailureStorage) Collect(tasks chan FailedTask, wg *sync.WaitGroup) {
	for task := range tasks {
		fs.mtx.Lock()
		fs.failure = append(fs.failure, task)
		fs.mtx.Unlock()
	}
	wg.Done()
}

func (fs *FailureStorage) Print() {
	fmt.Println("Errors: ")
	for _, f := range fs.failure {
		fmt.Println(f.String())
	}
}

func main() {
	ctx, stop := context.WithCancel(context.Background())
	taskGetterCh := make(chan Task, 10)
	go taskCreator(ctx, taskGetterCh)

	successTasksCh := make(chan Task, 10)
	failureTasksCh := make(chan FailedTask, 10)
	resultsWg := &sync.WaitGroup{}

	failures := FailureStorage{
		mtx:     &sync.Mutex{},
		failure: make([]FailedTask, 0),
	}
	successes := SuccessStorage{
		mtx:     &sync.Mutex{},
		success: make([]Task, 0),
	}
	resultsWg.Add(2)
	go successes.Collect(successTasksCh, resultsWg)
	go failures.Collect(failureTasksCh, resultsWg)

	processorWg := &sync.WaitGroup{}
	processorWg.Add(runtime.NumCPU())
	/// Split by CPU num for example
	for c := 0; c < runtime.NumCPU(); c++ {
		go taskProcessor(taskGetterCh, successTasksCh, failureTasksCh, processorWg)
	}

	time.Sleep(time.Second * 3)
	stop()
	processorWg.Wait()
	close(successTasksCh)
	close(failureTasksCh)
	resultsWg.Wait()

	failures.Print()
	successes.Print()
}

