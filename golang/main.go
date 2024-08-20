package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// A TType represents a meaninglessness of our life
type TType struct {
	id            int
	createTime    time.Time // Время создания.
	executionTime time.Time // Время выполнения.
	taskRESULT    bool
}

var (
	mu     sync.Mutex
	lastID int
)

func main() {
	// Канал для получаения задач.
	superChan := make(chan TType, 10)
	// Канал для успешных задач.
	doneTasks := make(chan TType, 10)
	// Канал для задач с ошибкой
	undoneTasks := make(chan error, 10)
	// Контекст который ограничивает время выполнения программы.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var wg sync.WaitGroup

	go taskCreate(ctx, superChan)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go taskWorker(ctx, &wg, superChan, doneTasks, undoneTasks)
	}

	// Функция для вывода информации каждын 3 секунды.
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				printResults(doneTasks, undoneTasks)
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		wg.Wait()
		close(doneTasks)
		close(superChan)
	}()

	<-ctx.Done()
}

// Функция создания таски
func taskCreate(ctx context.Context, superChan chan TType) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			ft := time.Now()
			id := int(ft.UnixNano())
			mu.Lock()
			if id <= lastID {
				mu.Unlock()
				continue
			}
			lastID = id
			mu.Unlock()

			task := TType{
				id:         id,
				createTime: ft,
				taskRESULT: time.Now().Nanosecond()%2 > 0, // Вот такое условие появления ошибочных тасков.
			}
			superChan <- task
		}
	}
}

// Функция в которой таски работают.
func taskWorker(ctx context.Context, wg *sync.WaitGroup, taskChan <-chan TType, doneTasks chan<- TType, undoneTasks chan<- error) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-taskChan:
			if !ok {
				return
			}
			task.executionTime = time.Now()
			if task.taskRESULT {
				if task.createTime.After(time.Now().Add(-20 * time.Second)) {
					doneTasks <- task
				} else {
					undoneTasks <- fmt.Errorf("Task id %d time %s failed", task.id, task.createTime)
				}
			} else {
				undoneTasks <- fmt.Errorf("Task id %d time %s failed", task.id, task.createTime)
			}
			time.Sleep(150 * time.Millisecond)
		}
	}
}

// Функция для печати.
func printResults(doneTasks <-chan TType, undoneTasks <-chan error) {
	fmt.Println("Result:")
	doneCount := 0
	undoneCount := 0
	fmt.Println("Done:")
	for len(doneTasks) > 0 {
		task := <-doneTasks
		fmt.Println(task)
		doneCount++
	}
	fmt.Println("Undone:")
	for len(undoneTasks) > 0 {
		err := <-undoneTasks
		fmt.Println(err)
		undoneCount++
	}
	fmt.Printf("Done: %d, Undone: %d\n", doneCount, undoneCount)
}
