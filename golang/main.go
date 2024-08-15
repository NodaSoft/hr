package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	id          int
	createdTime string // время создания
	finishedTime string // время выполнения
	result      []byte
}

func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex

	taskChannel := make(chan Task, 10)
	successChannel := make(chan Task, 10)
	errorChannel := make(chan error, 10)

	taskResults := make(map[int]Task)
	taskErrors := []error{}

	// Генератор задач
	go func() {
		defer close(taskChannel)
		startTime := time.Now()
		for time.Since(startTime) < 10*time.Second {
			createdTime := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 { // условие появления ошибочных задач
				createdTime = "Some error occurred"
			}
			taskChannel <- Task{createdTime: createdTime, id: int(time.Now().UnixNano())}
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// Обработчик задач
	processTask := func(task Task) Task {
		parsedTime, err := time.Parse(time.RFC3339, task.createdTime)
		if err != nil || parsedTime.After(time.Now().Add(-20*time.Second)) {
			task.result = []byte("task has succeeded")
		} else {
			task.result = []byte("something went wrong")
		}
		task.finishedTime = time.Now().Format(time.RFC3339Nano)
		time.Sleep(150 * time.Millisecond) // имитация времени выполнения
		return task
	}

	// Сортировщик задач
	sortTask := func(task Task) {
		if string(task.result) == "task has succeeded" {
			successChannel <- task
		} else {
			errorChannel <- fmt.Errorf("Task ID %d time %s, error %s", task.id, task.createdTime, task.result)
		}
		wg.Done()
	}

	// Полученные задачи
	go func() {
		for task := range taskChannel {
			wg.Add(1)
			go func(t Task) {
				processedTask := processTask(t)
				sortTask(processedTask)
			}(task)
		}
	}()

	// Успешные задачи
	go func() {
		for successfulTask := range successChannel {
			mu.Lock()
			taskResults[successfulTask.id] = successfulTask
			mu.Unlock()
		}
	}()

	// Ошибочные задачи
	go func() {
		for taskError := range errorChannel {
			mu.Lock()
			taskErrors = append(taskErrors, taskError)
			mu.Unlock()
		}
	}()

	go func() {
		for {
			time.Sleep(5 * time.Second)
			mu.Lock()
			fmt.Println("Completed tasks:")
			for _, task := range taskResults {
				fmt.Printf("Task ID: %d, Created: %s, Finished: %s\n", task.id, task.createdTime, task.finishedTime)
			}
			fmt.Println("Errors:")
			for _, err := range taskErrors {
				fmt.Println(err)
			}
			mu.Unlock()
		}
	}()

	wg.Wait()
	close(successChannel)
	close(errorChannel)

	time.Sleep(5 * time.Second)
	fmt.Println("Final Completed tasks:")
	for _, task := range taskResults {
		fmt.Printf("Task ID: %d, Created: %s, Finished: %s\n", task.id, task.createdTime, task.finishedTime)
	}
	fmt.Println("Final Errors:")
	for _, err := range taskErrors {
		fmt.Println(err)
	}
}

