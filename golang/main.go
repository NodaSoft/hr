package main

import (
	"fmt"
	"sync"
	"time"
)

// Реализована синхронизация с помощью sync.WaitGroup для корректного завершения задач.
// Убраны лишние и неправильные go-routines, вызывающие гонки данных.
// Добавлен Ticker для регулярного вывода результатов каждые 3 секунды.
// Исправлена логика обработки успешных и ошибочных задач.
// Уменьшена частота генерации задач для лучшего контроля
// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

func taskCreator(a chan Ttype, wg *sync.WaitGroup) {
	defer wg.Done()
	for start := time.Now(); time.Since(start) < 10*time.Second; {
		ft := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 {
			ft = "Some error occured"
		}
		a <- Ttype{cT: ft, id: int(time.Now().UnixNano())}
		time.Sleep(100 * time.Millisecond) // чтобы избежать слишком быстрого создания задач
	}
	close(a)
}

func taskWorker(a Ttype) Ttype {
	tt, err := time.Parse(time.RFC3339, a.cT)
	if err == nil && tt.After(time.Now().Add(-20*time.Second)) {
		a.taskRESULT = []byte("task has been successed")
	} else {
		a.taskRESULT = []byte("something went wrong")
	}
	a.fT = time.Now().Format(time.RFC3339Nano)
	time.Sleep(time.Millisecond * 150)
	return a
}

func main() {
	superChan := make(chan Ttype, 10)
	doneTasks := make(chan Ttype, 10)
	errorTasks := make(chan error, 10)

	var wg sync.WaitGroup

	wg.Add(1)
	go taskCreator(superChan, &wg)

	go func() {
		for t := range superChan {
			go func(t Ttype) {
				t = taskWorker(t)
				if string(t.taskRESULT) == "task has been successed" {
					doneTasks <- t
				} else {
					errorTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
				}
			}(t)
		}
		close(doneTasks)
		close(errorTasks)
	}()

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			fmt.Println("Errors:")
			for len(errorTasks) > 0 {
				err := <-errorTasks
				fmt.Println(err)
			}

			fmt.Println("Done tasks:")
			for len(doneTasks) > 0 {
				task := <-doneTasks
				fmt.Printf("Task ID: %d, Created: %s, Finished: %s, Result: %s\n",
					task.id, task.cT, task.fT, task.taskRESULT)
			}
		}
	}()

	wg.Wait()
	time.Sleep(3 * time.Second) // ждем последний вывод
}
