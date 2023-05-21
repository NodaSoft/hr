package main

import (
	"fmt"
	"sync"
	"time"
)

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнения остальных тасков

// A Task represents a meaninglessness of our life
type Task struct {
	CreationTime string // время создания
	ID           int
	TaskResult   []byte
	FinishTime   string // время выполнения
}

func main() {
	taskCreator := func(taskChan chan Task) {
		go func() {
			for {
				formattedTime := time.Now().Format(time.RFC3339)
				if time.Now().Nanosecond()%2 > 0 {
					formattedTime = "Some error occurred"
				}
				taskChan <- Task{CreationTime: formattedTime, ID: int(time.Now().Unix())}
			}
		}()
	}

	taskChan := make(chan Task, 10)

	go taskCreator(taskChan)

	var wg sync.WaitGroup
	wg.Add(10)

	doneTasks := make(chan Task)
	undoneTasks := make(chan error)
	taskWorker := func(task Task, wg *sync.WaitGroup, doneTasks chan Task, undoneTasks chan error) {
		defer wg.Done()
		taskTime, _ := time.Parse(time.RFC3339, task.CreationTime)
		if taskTime.After(time.Now().Add(-20 * time.Second)) {
			task.TaskResult = []byte("task has been succeeded")
			select {
			case doneTasks <- task:
			default:
			}
		} else {
			task.TaskResult = []byte("something went wrong")
			select {
			case undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", task.ID, task.CreationTime, task.TaskResult):
			default:
			}
		}
		task.FinishTime = time.Now().Format(time.RFC3339Nano)
		time.Sleep(time.Millisecond * 150)
	}

	go func() {
		for task := range taskChan {
			go taskWorker(task, &wg, doneTasks, undoneTasks)
		}
		wg.Wait()
		close(doneTasks)
		close(undoneTasks)
	}()

	var mu sync.Mutex
	result := map[int]Task{}
	errors := []error{}

	go func() {
		for task := range doneTasks {
			mu.Lock()
			result[task.ID] = task
			mu.Unlock()
		}
	}()

	go func() {
		for err := range undoneTasks {
			mu.Lock()
			errors = append(errors, err)
			mu.Unlock()
		}
	}()

	time.Sleep(time.Second * 3)

	println("Errors:")
	for _, err := range errors {
		println(err)
	}

	println("Done tasks:")
	for id := range result {
		println(id)
	}
}
