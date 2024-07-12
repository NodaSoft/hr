package main

import (
	"fmt"
	"sync"
	"time"
)

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

func main() {
	superChan := make(chan Ttype, 10)   // Канал для передачи задач
	doneTasks := make(chan Ttype, 10)   // Канал для успешных задач
	undoneTasks := make(chan error, 10) // Канал для неуспешных задач
	var wg sync.WaitGroup               // WaitGroup для ожидания завершения всех горутин
	var mu sync.Mutex                   // Mutex для защиты общего ресурса

	// Функция генерации задач
	taskCreturer := func(a chan Ttype) {
		for i := 0; i < 10; i++ {
			time.Sleep(time.Second)
			ft := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 {
				ft = "Some error occured"
			}
			a <- Ttype{cT: ft, id: int(time.Now().Unix())}
		}
		close(a)
	}

	// Функция обработки задач
	task_worker := func(a Ttype) Ttype {
		tt, err := time.Parse(time.RFC3339, a.cT)
		if err != nil || tt.After(time.Now().Add(-20*time.Second)) {
			a.taskRESULT = []byte("something went wrong")
		} else {
			a.taskRESULT = []byte("task has been successed")
		}
		a.fT = time.Now().Format(time.RFC3339Nano)
		time.Sleep(time.Millisecond * 150) // Задержка для эмуляции обработки
		return a
	}

	// Функция сортировки задач
	tasksorter := func(t Ttype) {
		mu.Lock() // Блокировка для защиты общего ресурса
		defer mu.Unlock()
		if string(t.taskRESULT) == "task has been successed" {
			doneTasks <- t
		} else {
			undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
		}
	}

	// Горутина для обработки задач из канала superChan
	go func() {
		for t := range superChan {
			wg.Add(1)
			go func(t Ttype) {
				defer wg.Done()
				t = task_worker(t)
				tasksorter(t)
			}(t)
		}
		wg.Wait()
		close(doneTasks)
		close(undoneTasks)
	}()

	go taskCreturer(superChan)

	// Горутина для периодического вывода результатов
	go func() {
		for {
			time.Sleep(3 * time.Second)
			fmt.Println("Errors:")
			for len(undoneTasks) > 0 {
				err := <-undoneTasks
				fmt.Println(err)
			}
			fmt.Println("Done tasks:")
			for len(doneTasks) > 0 {
				task := <-doneTasks
				fmt.Println(task)
			}
		}
	}()

	time.Sleep(12 * time.Second) // Время выполнения программы
}
