package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type taskObserver struct {
	doneTaskCh    <-chan Task
	failedTaskCh  <-chan Task
	doneTasks     []Task
	failedTasks   []Task
	printInterval time.Duration
	wg            sync.WaitGroup
}

func newTaskObserver(
	doneTaskCh <-chan Task,
	failedTaskCh <-chan Task,
	printInterval time.Duration,
) *taskObserver {
	return &taskObserver{
		doneTaskCh:    doneTaskCh,
		failedTaskCh:  failedTaskCh,
		doneTasks:     make([]Task, 0),
		failedTasks:   make([]Task, 0),
		printInterval: printInterval,
		wg:            sync.WaitGroup{},
	}
}

func (to *taskObserver) Start(ctx context.Context) {
	to.wg.Add(1)

	go func() {
		mt := sync.Mutex{}
		defer to.wg.Done()

		for {
			select {
			case task := <-to.doneTaskCh:
				mt.Lock()
				to.doneTasks = append(to.doneTasks, task)
				mt.Unlock()
			case task := <-to.failedTaskCh:
				mt.Lock()
				to.failedTasks = append(to.failedTasks, task)
				mt.Unlock()
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (to *taskObserver) PrintResultsPeriodically(ctx context.Context) {
	to.wg.Add(1)

	go func() {
		defer to.wg.Done()
		ticker := time.NewTicker(to.printInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				to.printResults()
			case <-ctx.Done():
				fmt.Println("Last PrintResults after cancel context")
				// Необходим еще один вывод на тот случай если таймт аут не кратен периоду вывода, и будет еще часть данные которые успели сгенерироваться но не попли в интервал
				to.printResults()
				return
			}
		}
	}()
}

func (to *taskObserver) Stop() {
	to.wg.Wait()
}

// Если нужно можно доработать для вывода любых данных по задаче. Не стал нагромождать логи
// Так же по хорошему использовать нормальный логгер.
func (to *taskObserver) printResults() {
	mt := sync.RWMutex{}

	fmt.Println("Success tasks:")
	mt.RLock()
	for _, doneTask := range to.doneTasks {
		fmt.Println(fmt.Sprintf("task id=%d,  statusInfo=%s", doneTask.id, doneTask.statusInfo))
	}
	mt.RUnlock()

	fmt.Println("Failed tasks:")
	mt.RLock()
	for _, failedTask := range to.failedTasks {
		fmt.Println(fmt.Errorf("task id=%d, errorInfo=%v", failedTask.id, failedTask.statusInfo))
	}
	mt.RUnlock()
}
