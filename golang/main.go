package main

import (
	"fmt"
	"time"
)

type Task struct {
	ID         int
	CreatedAt  time.Time
	FinishedAt time.Time
	Result     string
}

func main() {
	taskCreator := func(out chan<- Task) {
		defer close(out)
		for {
			createdAt := time.Now()
			var result string
			if createdAt.Nanosecond()%2 > 0 {
				result = "Что-то пошло не так("
			} else {
				result = "Задача завершена"
			}
			out <- Task{ID: int(createdAt.Unix()), CreatedAt: createdAt, Result: result}
		}
	}

	taskWorker := func(in <-chan Task, out chan<- Task, errCh chan<- error) {
		for task := range in {
			task.FinishedAt = time.Now()
			if task.Result == "Что-то пошло не так(" {
				errCh <- fmt.Errorf("Задача с id %d создана в %s, ошибка: %s", task.ID, task.CreatedAt, task.Result)
			} else {
				out <- task
			}
			time.Sleep(time.Millisecond * 150)
		}
	}

	doneTasks := make(chan Task)
	errTasks := make(chan error)

	go func() {
		for {
			select {
			case task, ok := <-doneTasks:
				if !ok {
					doneTasks = nil
					continue
				}
				fmt.Printf("Done Task: %v\n", task)
			case err, ok := <-errTasks:
				if !ok {
					errTasks = nil
					continue
				}
				fmt.Printf("Error: %v\n", err)
			}
		}
	}()

	tasks := make(chan Task, 10)
	go taskCreator(tasks)

	workers := 3 
	for i := 0; i < workers; i++ {
		go taskWorker(tasks, doneTasks, errTasks)
	}

	time.Sleep(time.Second * 3)
}
