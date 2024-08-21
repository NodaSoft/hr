package main

import (
	"fmt"
	"sync"
	"time"
)

// Task представляет задачу, которую нужно выполнить
type Task struct {
	ID             int64     // Уникальный идентификатор задачи
	CreationTime   time.Time // Время создания задачи
	CompletionTime time.Time // Время завершения задачи
	Result         string    // Результат выполнения задачи
}

// Функция для создания задач
func taskCreator(tasks chan<- Task) {
	defer close(tasks)
	for {
		select {
		case <-time.After(10 * time.Second):
			return
		default:
			creationTime := time.Now()
			if time.Now().Nanosecond()%2 > 0 { // Условие появления ошибочных задач
				creationTime = time.Time{} // Пустое время для обозначения ошибки
			}
			task := Task{
				ID:           time.Now().UnixNano(), // Уникальный ID на основе времени
				CreationTime: creationTime,
			}
			tasks <- task
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// Функция для обработки задач
func taskWorker(task Task) Task {
	if !task.CreationTime.IsZero() {
		task.Result = "Task has been successful"
	} else {
		task.Result = "Something went wrong"
	}
	task.CompletionTime = time.Now()
	time.Sleep(150 * time.Millisecond) // Симуляция времени обработки задачи
	return task
}

// Функция для сортировки задач на успешные и неуспешные
func taskSorter(task Task, doneTasks chan<- Task, errorTasks chan<- Task) {
	if task.Result == "Task has been successful" {
		doneTasks <- task
	} else {
		errorTasks <- task
	}
}

// Функция для периодической печати результатов
func resultPrinter(doneTasks <-chan Task, errorTasks <-chan Task, done <-chan struct{}) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("=== Checking progress ===")
			printResults(doneTasks, errorTasks)
		case <-done:
			fmt.Println("=== Final results ===")
			printResults(doneTasks, errorTasks)
			return
		}
	}
}

// Функция для печати результатов
func printResults(doneTasks <-chan Task, errorTasks <-chan Task) {
	fmt.Println("Done tasks:")
	for task := range doneTasks {
		fmt.Printf("Task ID: %d, Created: %s, Finished: %s\n", task.ID, task.CreationTime.Format(time.RFC3339), task.CompletionTime.Format(time.RFC3339Nano))
	}

	fmt.Println("Errors:")
	for task := range errorTasks {
		fmt.Printf("Task ID: %d, Error: %s\n", task.ID, task.Result)
	}
}

func main() {
	tasks := make(chan Task, 10)
	doneTasks := make(chan Task, 10)
	errorTasks := make(chan Task, 10)
	done := make(chan struct{})

	var wg sync.WaitGroup

	// Запуск горутины для создания задач
	go taskCreator(tasks)

	// Запуск нескольких горутин для обработки задач
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasks {
				processedTask := taskWorker(task)
				taskSorter(processedTask, doneTasks, errorTasks)
			}
		}()
	}

	// Запуск горутины для периодической печати результатов
	go resultPrinter(doneTasks, errorTasks, done)

	// Ждем завершения всех горутин
	wg.Wait()

	// Закрытие каналов задач после завершения обработки
	close(doneTasks)
	close(errorTasks)

	// Сообщаем resultPrinter о завершении работы
	close(done)
}
