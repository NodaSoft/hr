package main

import (
	"fmt"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

type Task struct {
	ID         int
	CreatedAt  time.Time
	FinishedAt time.Time
	Result     string
}

func main() {
	taskCreator := func(tasks chan Task) {
		go func() {
			for {
				createdAt := time.Now()
				var finishedAt time.Time
				if time.Now().Nanosecond()%2 > 0 {
					finishedAt = time.Time{}
				} else {
					finishedAt = createdAt.Add(20 * time.Second)
				}
				tasks <- Task{ID: int(createdAt.Unix()), CreatedAt: createdAt, FinishedAt: finishedAt}
				time.Sleep(500 * time.Millisecond)
			}
		}()
	}

	taskWorker := func(task Task) Task {
		if task.FinishedAt.IsZero() || task.FinishedAt.Before(time.Now()) {
			task.Result = "something went wrong"
		} else {
			task.Result = "task has been succeeded"
		}
		task.FinishedAt = time.Now()
		return task
	}

	doneTasks := make(chan Task)
	undoneTasks := make(chan Task)

	taskSorter := func(task Task) {
		if task.Result == "task has been succeeded" {
			doneTasks <- task
		} else {
			undoneTasks <- task
		}
	}

	go func() {
		taskCreatorChan := make(chan Task)
		go taskCreator(taskCreatorChan)

		for task := range taskCreatorChan {
			go func(t Task) {
				task := taskWorker(t)
				taskSorter(task)
			}(task)
		}
	}()

	results := make(map[int]Task)
	var errors []Task

	go func() {
		for task := range doneTasks {
			results[task.ID] = task
		}
	}()

	go func() {
		for task := range undoneTasks {
			errors = append(errors, task)
		}
	}()

	time.Sleep(5 * time.Second)

	fmt.Println("Errors:")
	for _, err := range errors {
		fmt.Printf("Task ID: %d, Created At: %s, Result: %s\n", err.ID, err.CreatedAt, err.Result)
	}

	fmt.Println("\nDone Tasks:")
	for id, result := range results {
		fmt.Printf("Task ID: %d, Created At: %s, Result: %s\n", id, result.CreatedAt, result.Result)
	}
}

