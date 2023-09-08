package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	ID           int
	CreatedTime  string
	FinishedTime string
	Result       string
}

func main() {
	taskCreator := func(tasks chan<- Task) {
		for {
			createdTime := time.Now().Format(time.RFC3339)
			var result string
			if time.Now().Nanosecond()%2 > 0 {
				result = "error occurred"
			} else {
				result = "task has been succeeded"
			}
			tasks <- Task{ID: int(time.Now().Unix()), CreatedTime: createdTime, Result: result}
			time.Sleep(time.Second) // Задержка между созданием задач
		}
	}

	taskWorker := func(task Task, wg *sync.WaitGroup, successTasks chan<- Task, errorTasks chan<- error) {
		defer wg.Done()

		createdTime, err := time.Parse(time.RFC3339, task.CreatedTime)
		if err != nil {
			errorTasks <- fmt.Errorf("Task ID %d time %s, error: %s", task.ID, task.CreatedTime, err.Error())
			return
		}

		if time.Now().Sub(createdTime) > 20*time.Second {
			task.FinishedTime = time.Now().Format(time.RFC3339Nano)
			successTasks <- task
		} else {
			errorTasks <- fmt.Errorf("Task ID %d time %s, error: something went wrong", task.ID, task.CreatedTime)
		}
	}

	// Используем WaitGroup для ожидания завершения всех задач
	var wg sync.WaitGroup

	successTasks := make(chan Task, 10)
	errorTasks := make(chan error, 10)

	go func() {
		// Обработка успешных задач
		for task := range successTasks {
			fmt.Printf("Task ID %d finished successfully\n", task.ID)
		}
	}()

	go func() {
		// Обработка ошибок
		for err := range errorTasks {
			fmt.Println(err)
		}
	}()

	// Создаем и обрабатываем задачи в пуле горутин
	poolSize := 5 // Размер пула горутин

	for i := 0; i < poolSize; i++ {
		go func() {
			for {
				task := <-tasksChan
				wg.Add(1)
				go taskWorker(task, &wg, successTasks, errorTasks)
			}
		}()
	}

	tasksChan := make(chan Task, 10)

	go taskCreator(tasksChan)

	// Ожидаем завершения всех задач
	wg.Wait()
}
