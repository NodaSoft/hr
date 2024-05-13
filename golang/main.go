package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Task struct {
	ID         int
	CreatedAt  time.Time
	FinishedAt time.Time
	Result     []byte
}

func main() {
	taskChan := make(chan Task, 10)
	doneTasks := make(chan Task)
	errChan := make(chan error)

	//создаем таску
	go func() {
		for {
			now := time.Now()
			var createdAt time.Time
			if now.Nanosecond()%2 > 0 {
				createdAt = time.Time{}
			} else {
				createdAt = now
			}
			task := Task{
				ID:        rand.Int(),
				CreatedAt: createdAt,
			}
			taskChan <- task
			time.Sleep(50 * time.Millisecond)
		}
	}()

	//иммитируем работу таски
	go func() {
		for task := range taskChan {
			task = processTask(task)
			if task.Result != nil {
				doneTasks <- task
			} else {
				errChan <- fmt.Errorf("Task id %d time %s, error", task.ID, task.CreatedAt)
			}
		}
		close(taskChan)
	}()

	//выводим резы
	go func() {
		for {
			select {
			case task := <-doneTasks:
				fmt.Printf("Task %d: Created at %s, Finished at %s, Result: %s\n", task.ID, task.CreatedAt.Format(time.RFC3339), task.FinishedAt.Format(time.RFC3339), string(task.Result))
			case err := <-errChan:
				fmt.Printf("Error: %s\n", err)
			}
		}
	}()

	time.Sleep(3 * time.Second)
}

func processTask(task Task) Task {
	if time.Since(task.CreatedAt) < 20*time.Second {
		task.Result = []byte("task has been successed")
	} else {
		task.Result = []byte("something went wrong")
	}
	task.FinishedAt = time.Now()
	time.Sleep(150 * time.Millisecond)
	return task
}
