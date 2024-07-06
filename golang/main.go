package main

import (
	"context"
	"fmt"
	"hh-task/task"
	"runtime"
	"sync"
	"time"
)

func taskCreator(taskChannel chan task.Task) {
	for {
		taskChannel <- task.Create() // передаем таск на выполнение
	}
}

func taskManager(taskChannel chan task.Task, workedTask chan task.Task, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	for task := range taskChannel {
		// В данном случае не требуется заставлять выполнять работу таска в горутине
		// Т.к. уже создано оптимальное количество обработчиков
		task.Work()

		select {
		case <-ctx.Done():
			return
		default:
			workedTask <- task
		}
	}
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	taskCh := make(chan task.Task, 10)
	taskResultCh := make(chan task.Task, 10)

	wg := sync.WaitGroup{}
	// Создаем оптимальное количество обработчиков тасков
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go taskManager(taskCh, taskResultCh, &wg, ctx)
	}

	go taskCreator(taskCh)

	// Мап не нужен
	var successfulTasks []task.Task
	var erroredTasks []task.Task

	timer := time.NewTimer(3 * time.Second)

	for {
		select {
		case <-timer.C:

			fmt.Println("Done Tasks:")
			for _, task := range successfulTasks {
				fmt.Println(task.String())
			}

			fmt.Println("Failed Tasks:")
			for _, task := range erroredTasks {
				fmt.Println(task.String())
			}

			cancel()
		case <-ctx.Done():
			wg.Wait()
			return
		default:
			// Тут используется default т.к. может быть ситуация где таймер кончился, а необработанные таски все еще есть
			// В таком случае т.к. селект случайно выбирает, после окончания таймера могут залететь еще таски
			// При дефолте такого произойти не может
			task := <-taskResultCh

			if task.IsCompleted() {
				successfulTasks = append(successfulTasks, task)
			} else {
				erroredTasks = append(erroredTasks, task)
			}
		}
	}

}
