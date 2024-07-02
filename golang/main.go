package main

import (
	"fmt"
	"github.com/kxait/pterm"
	"sync"
	"time"
)

type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

func main() {
	superChan := make(chan Ttype, 10)
	var wg sync.WaitGroup

	// Генерация задач в течение 10 секунд
	go func() {
		endTime := time.Now().Add(10 * time.Second)
		for time.Now().Before(endTime) {
			ft := formatTime(time.Now())
			if time.Now().Nanosecond()%2 > 0 { // условие появления ошибочных тасков
				ft = "Some error occurred"
			}
			task := Ttype{cT: ft, id: int(time.Now().UnixNano())}
			pterm.Info.Printf("Received task ID %d: creation time %s\n", task.id, task.cT)
			superChan <- task // передаем таск на выполнение
			time.Sleep(500 * time.Millisecond)
		}
		close(superChan)
	}()

	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	// Обработка задач
	go func() {
		for task := range superChan {
			wg.Add(1)
			go func(t Ttype) {
				defer wg.Done()

				if t.cT == "Some error occurred" {
					t.taskRESULT = []byte("something went wrong")
				} else {
					t.taskRESULT = []byte("task completed successfully")
				}
				t.fT = formatTime(time.Now())
				time.Sleep(150 * time.Millisecond)

				if string(t.taskRESULT) == "task completed successfully" {
					pterm.Success.Printf("Completed task ID %d: creation time %s, completion time %s\n", t.id, t.cT, t.fT)
					doneTasks <- t
				} else {
					err := fmt.Errorf("Task ID %d creation time %s, error %s", t.id, t.cT, t.taskRESULT)
					pterm.Error.Println(err)
					undoneTasks <- err
				}
			}(task)
		}
		wg.Wait()
		close(doneTasks)
		close(undoneTasks)
	}()

	var mu sync.Mutex
	results := make(map[int]Ttype)
	errors := []error{}

	// Результаты
	go func() {
		for doneTask := range doneTasks {
			mu.Lock()
			results[doneTask.id] = doneTask
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

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		mu.Lock()
		pterm.Info.Println("Completed tasks:")
		for _, task := range results {
			pterm.Info.Printf("Task ID %d: creation time %s, completion time %s\n", task.id, task.cT, task.fT)
		}
		pterm.Error.Println("Errors:")
		for _, err := range errors {
			pterm.Error.Println(err)
		}
		mu.Unlock()
	}
}

// Форматирование времени
func formatTime(t time.Time) string {
	return t.Format("02.01.2006 15:04:05")
}
