package main

import (
	"fmt"
	"sync"
	"time"
)

// ЗАДАНИЕ:
// + сделать из плохого кода хороший;
// + важно сохранить логику появления ошибочных тасков;
// + сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

// A Ttype represents a meaninglessness of our life :)
type Task struct {
	ID         int
	CreatedAt  string // время создания
	FinishedAt string // время выполнения
	Result     []byte
	IsDone     bool
}

func main() {
	numTasks := 10 // Set tasks number
	superChan := make(chan Task, numTasks)
	doneTasks := make(chan Task)
	undoneTasks := make(chan Task)

	wg := sync.WaitGroup{}
	wg.Add(1)

	// Create a tasks
	go taskCreator(superChan, &wg, numTasks)

	go func() {
		// Get a tasks
		for t := range superChan {
			wg.Add(1)
			go func(task Task) {
				defer wg.Done()
				t = taskWorker(t)
				go taskSorter(t, doneTasks, undoneTasks)
			}(t)
		}
		wg.Wait()
		close(doneTasks)
		close(undoneTasks)
	}()

	result := map[int]Task{}
	errors := map[int]Task{}
	mu := sync.Mutex{}

	// Get results (success and fails)
	go func() {
		for r := range doneTasks {
			mu.Lock()
			result[r.ID] = r
			mu.Unlock()
		}
	}()

	go func() {
		for r := range undoneTasks {
			mu.Lock()
			errors[r.ID] = r
			mu.Unlock()
		}
	}()

	time.Sleep(time.Second * 3)

	fmt.Println("Errors:")
	mu.Lock()
	for _, err := range errors {
		fmt.Printf("ID: %d, Result: %s\n", err.ID, err.Result)
	}
	mu.Unlock()

	fmt.Println("Done tasks:")
	mu.Lock()
	for _, res := range result {
		fmt.Printf("ID: %d, Result: %s\n", res.ID, res.Result)
	}
	mu.Unlock()
}

func taskCreator(tasks chan Task, wg *sync.WaitGroup, numTasks int) {
	defer wg.Done()
	for i := 0; i < numTasks; i++ {
		crunchTime := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			crunchTime = "Some error occurred"
		}
		tasks <- Task{CreatedAt: crunchTime, ID: i + 1} // передаем таск на выполнение
	}
}

func taskWorker(task Task) Task {
	createTime, _ := time.Parse(time.RFC3339, task.CreatedAt)
	if createTime.After(time.Now().Add(-20 * time.Second)) {
		task.Result = []byte("task has been successed")
		task.IsDone = true
	} else {
		task.Result = []byte("something went wrong")
		task.IsDone = false
	}
	task.FinishedAt = time.Now().Format(time.RFC3339Nano)
	time.Sleep(time.Millisecond * 150)
	return task
}

func taskSorter(task Task, done chan Task, undone chan Task) {
	if task.IsDone {
		done <- task
	} else {
		undone <- task
	}
}
