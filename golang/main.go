package main

import (
	"fmt"
	"sync"
	"time"
)

// Ttype представляет собой структуру задачи
type Ttype struct {
	id         int
	cT         time.Time // время создания
	fT         time.Time // время выполнения
	taskResult string
}

func main() {
	var wg sync.WaitGroup
	tasks := make(chan Ttype, 10)
	doneTasks := make(chan Ttype, 10)
	undoneTasks := make(chan Ttype, 10)

	// Генератор задач
	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(1 * time.Second) // Генерация задачи каждую секунду
			task := Ttype{
				id: i,
				cT: time.Now(),
			}
			if time.Now().Nanosecond()%2 > 0 { // Условие для ошибочных задач
				task.taskResult = "error"
			} else {
				task.taskResult = "success"
			}
			tasks <- task
		}
		close(tasks)
	}()

	// Рабочий обрабатывает задачи
	wg.Add(1)
	go func() {
		defer wg.Done()
		for task := range tasks {
			task.fT = time.Now()
			if task.taskResult == "success" {
				doneTasks <- task
			} else {
				undoneTasks <- task
			}
		}
	}()

	// Вывод результатов каждые 3 секунды
	go func() {
		for {
			select {
			case task := <-doneTasks:
				fmt.Printf("Успешная задача: %+v\n", task)
			case task := <-undoneTasks:
				fmt.Printf("Задача с ошибкой: %+v\n", task)
			case <-time.After(10 * time.Second):
				close(doneTasks)
				close(undoneTasks)
				return
			}
		}
	}()

	wg.Wait()
}
