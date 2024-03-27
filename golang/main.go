package main

import (
	"fmt"
	"sync"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнения остальных тасков

const (
    errorMessage = "Error occurred"
    successMessage = "Task completed successfully" 
    numberOfTasks = 10
    numberOfWorkers = 5
)

type Task struct {
	ID         int
	CreatedAt  time.Time
	StartedAt  time.Time
	FinishedAt time.Time
	Result     string
}

func main() {
	taskCreator := func(out chan<- Task) {
		for {
			createdAt := time.Now()
			var result string
			if createdAt.Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				result = errorMessage
			} else {
				result = successMessage
			}
			task := Task{ID: int(createdAt.Unix()), CreatedAt: createdAt, Result: result}
			time.Sleep(time.Millisecond * 50) // имитация времени создания задачи
			select {
                case out <- task: // передаем таск на выполнение
                default:
                    return
                }
		}
	}

	taskWorker := func(in <-chan Task, out chan<- Task, wg *sync.WaitGroup) {
		defer wg.Done()
		for task := range in {
		    task.StartedAt = time.Now()
			time.Sleep(time.Millisecond * 150) // имитация времени обработки задачи
			task.FinishedAt = time.Now()
			out <- task
		}
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex
	doneTasks := make(map[int]Task)
	var errors []error

	tasks := make(chan Task, numberOfTasks)
	done := make(chan Task, numberOfTasks)

	go taskCreator(tasks)

	for i := 0; i < numberOfWorkers; i++ {
	    wg.Add(1)
		go taskWorker(tasks, done, &wg)
	}
	
	wg.Add(1)
	go func() {
	    defer wg.Done()
	    for task := range done {
    		if task.Result == successMessage {
    			mutex.Lock()
    			doneTasks[task.ID] = task
    			mutex.Unlock()
    		} else {
    			errors = append(errors, fmt.Errorf("Task ID %d created at %s: %s", task.ID, task.CreatedAt.Format(time.RFC3339), task.Result))
    		}
	    }
	}()
	
	go func() {
		wg.Wait()
		close(done)
	}()
    
    time.Sleep(time.Second * 3)
    close(tasks)

	fmt.Println("Errors:")
	for _, err := range errors {
		fmt.Println(err)
	}

	fmt.Println("\nDone tasks:")
	for id := range doneTasks {
		fmt.Println(id)
	}
}
