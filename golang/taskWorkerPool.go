package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var validTaskPeriod = -20 * time.Second

type taskWorkerPool struct {
	newTasksCh   <-chan Task
	doneTaskCh   chan<- Task
	failedTaskCh chan<- Task
	numWorkers   int
	wg           sync.WaitGroup
}

func newTaskWorkerPool(
	newTasksCh <-chan Task,
	doneTaskCh chan<- Task,
	failedTaskCh chan<- Task,
	numWorkers int,
) *taskWorkerPool {
	return &taskWorkerPool{
		newTasksCh:   newTasksCh,
		doneTaskCh:   doneTaskCh,
		failedTaskCh: failedTaskCh,
		numWorkers:   numWorkers,
	}
}

func (wp *taskWorkerPool) Start(ctx context.Context) {
	wp.wg.Add(wp.numWorkers)
	for i := 0; i < wp.numWorkers; i++ {
		go wp.taskWorker(ctx)
	}
}

func (wp *taskWorkerPool) taskWorker(ctx context.Context) {
	defer wp.wg.Done()
	for {
		select {
		case task, ok := <-wp.newTasksCh:
			if !ok {
				return // Если канал закрыт и больше нет задач, выходим
			}
			if task.status != ErrorStatus {

				time.Sleep(time.Millisecond * 150) // Эмуляция работы
				task.executeTime = time.Now()

				if task.createTime.After(time.Now().Add(validTaskPeriod)) {
					task.status = Success
					task.statusInfo = "Task completed successfully"

					wp.doneTaskCh <- task
				} else {
					task.status = ErrorStatus
					task.statusInfo = fmt.Sprintf("create time not after period=%v", validTaskPeriod)

					wp.failedTaskCh <- task
				}

			} else {
				wp.failedTaskCh <- task
			}
		case <-ctx.Done():
			return // выход через явную отмену
		}
	}
}

func (wp *taskWorkerPool) Stop() {
	wp.wg.Wait()
	close(wp.doneTaskCh)
	close(wp.failedTaskCh)
}
