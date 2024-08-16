package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// type taskData информация о таске
type taskData struct {
	id            int
	createTime    string // время создания
	executionTime string // время выполнения
	taskResult    []byte
}

// type wPool все составляющие воркер пула + каналы для вывода результатов
type wPool struct {
	wg               sync.WaitGroup
	numWorkers       int
	createdTasksChan chan taskData
	doneTasks        chan taskData
	errTasks         chan error
}

// newWPool - билдер воркер пула
func newWPool() *wPool {
	return &wPool{
		wg:               sync.WaitGroup{},
		numWorkers:       5,
		createdTasksChan: make(chan taskData, 10),
		doneTasks:        make(chan taskData, 10),
		errTasks:         make(chan error, 10),
	}
}

// taskCreator - создание тасков
func (wp *wPool) taskCreator(ctx context.Context) {
	defer close(wp.createdTasksChan)
	defer wp.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			ft := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				ft = "Some error occured"
			}
			wp.createdTasksChan <- taskData{createTime: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
		}
	}
}

// taskWorkerAndSorter обработка задач и последующая сортировка их результатов по каналам
func (wp *wPool) taskWorkerAndSorter() {
	defer wp.wg.Done()
	for task := range wp.createdTasksChan {
		tt, err := time.Parse(time.RFC3339, task.createTime)
		var successFlag bool
		if err == nil && tt.After(time.Now().Add(-20*time.Second)) {
			task.taskResult = []byte("task has been successed")
			successFlag = true
		} else {
			task.taskResult = []byte("something went wrong")
			successFlag = false
		}
		task.executionTime = time.Now().Format(time.RFC3339Nano)
		if successFlag {
			wp.doneTasks <- task
		} else {
			wp.errTasks <- fmt.Errorf("Task id %d time %s, error %s", task.id, task.createTime, task.taskResult)
		}
		time.Sleep(time.Millisecond * 150)
	}
}

// printResult - вывод результатов выполнения задач
func (wp *wPool) printResult(ticker *time.Ticker, stopChan <-chan struct{}) {
	for {
		select {
		case <-ticker.C:
			fmt.Println("Done tasks: ")
			for len(wp.doneTasks) > 0 {
				task := <-wp.doneTasks
				fmt.Printf("Task ID: %d, Created: %s, Finished: %s\n", task.id, task.createTime, task.executionTime)

			}
			fmt.Println("Errors: ")
			for len(wp.errTasks) > 0 {
				err := <-wp.errTasks
				fmt.Println(err)
			}

		case <-stopChan:
			return
		}
	}
}

func main() {

	wp := newWPool()
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stopChan := make(chan struct{})

	wp.wg.Add(1)
	go wp.taskCreator(ctx)

	//создаем воркер пул
	for i := 1; i <= wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.taskWorkerAndSorter()
	}

	go wp.printResult(ticker, stopChan)

	wp.wg.Wait()
	go func() { // Закрываем каналы результатов
		close(wp.doneTasks)
		close(wp.errTasks)
		// Сигнализируем PrintResult завершить работу
		stopChan <- struct{}{}
		close(stopChan)
	}()

	fmt.Println("All tasks processed, exiting.")
}
