package workerPool

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Task interface {
	Process()
	Result() string
	IsSuccess() bool
	Status() string
	Error() error
}

type WorkerPool struct {
	Wg           sync.WaitGroup
	workersCount int
	ctx          context.Context
	taskFactory  func() Task
	Tasks        chan Task
	successTasks chan Task
	errorTasks   chan error
}

func New(workerCount int, ctx context.Context) *WorkerPool {
	return &WorkerPool{
		Wg:           sync.WaitGroup{},
		workersCount: workerCount,
		ctx:          ctx,
		Tasks:        make(chan Task, 10),
		successTasks: make(chan Task, 10),
		errorTasks:   make(chan error, 10),
	}
}

func (wp *WorkerPool) StartFactory(errChan chan error) {
	if wp.taskFactory == nil {
		errChan <- fmt.Errorf("Task factory was not provided")
		return
	}

	defer close(wp.Tasks)
	defer wp.Wg.Done()
	for {
		select {
		case <-wp.ctx.Done():
			return
		default:
			wp.Tasks <- wp.taskFactory()
		}
	}
}

func (wp *WorkerPool) WithFactory(taskFactory func() Task) {
	wp.taskFactory = taskFactory
}

func (wp *WorkerPool) ProcessAndSortTask() {
	defer wp.Wg.Done()
	for task := range wp.Tasks {
		task.Process()

		wp.sortTask(task)

		time.Sleep(time.Millisecond * 150)
	}
}

func (wp *WorkerPool) sortTask(t Task) {
	if t.IsSuccess() {
		wp.successTasks <- t
	} else {
		wp.errorTasks <- t.Error()
	}
}

func (wp *WorkerPool) PrintResults(ticker *time.Ticker, stopChan chan struct{}) {
	for {
		select {
		case <-ticker.C:
			fmt.Println("Success tasks:")
			for i := 0; i <= len(wp.successTasks); i++ {
				task := <-wp.successTasks
				fmt.Println(task.Status())
			}

			fmt.Println("Errors:")
			for i := 0; i <= len(wp.successTasks); i++ {
				err := <-wp.errorTasks
				fmt.Println(err)
			}
		case <-stopChan:
			return
		}
	}
}

func (wp *WorkerPool) Shutdown() {
	close(wp.successTasks)
	close(wp.errorTasks)
}
