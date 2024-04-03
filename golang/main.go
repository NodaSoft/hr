package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	ID         int
	CreateTime time.Time
	FinishTime time.Time
	Result     string
}

func main() {
	taskCreator := func(taskChan chan Task, wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			createTime := time.Now()
			var finishTime time.Time
			if time.Now().Nanosecond()%2 > 0 {
				// имитация ошибочного задания
				finishTime = time.Time{}
			} else {
				finishTime = createTime.Add(20 * time.Second)
			}
			task := Task{
				ID:         int(createTime.Unix()),
				CreateTime: createTime,
				FinishTime: finishTime,
			}
			taskChan <- task
			time.Sleep(time.Second) // задержка между созданием заданий
		}
	}

	taskWorker := func(task Task, wg *sync.WaitGroup) {
		defer wg.Done()
		if task.FinishTime.IsZero() || time.Now().After(task.FinishTime) {
			task.Result = "error"
		} else {
			task.Result = "success"
		}
		task.FinishTime = time.Now()
		time.Sleep(time.Millisecond * 150)
	}

	var wg sync.WaitGroup
	taskChan := make(chan Task)
	doneTasks := make(chan Task)
	errors := make(chan error)

	// Запуск горутины для создания задач
	wg.Add(1)
	go taskCreator(taskChan, &wg)

	// Запуск горутин для обработки задач
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskChan {
				taskWorker(task, &wg)
				if task.Result == "success" {
					doneTasks <- task
				} else {
					errors <- fmt.Errorf("Task ID %d created at %s: Error occurred", task.ID, task.CreateTime.Format(time.RFC3339))
				}
			}
		}()
	}

	// Закрытие каналов после завершения работы всех горутин
	go func() {
		wg.Wait()
		close(taskChan)
		close(doneTasks)
		close(errors)
	}()

	// Сбор результатов и вывод
	var (
		successfulTasks []Task
		errorMessages   []error
	)

	for task := range doneTasks {
		successfulTasks = append(successfulTasks, task)
	}

	for err := range errors {
		errorMessages = append(errorMessages, err)
	}

	fmt.Println("Errors:")
	for _, errMsg := range errorMessages {
		fmt.Println(errMsg)
	}

	fmt.Println("Successful tasks:")
	for _, task := range successfulTasks {
		fmt.Printf("ID: %d, Create Time: %s, Finish Time: %s, Result: %s\n", task.ID, task.CreateTime.Format(time.RFC3339), task.FinishTime.Format(time.RFC3339), task.Result)
	}
}
