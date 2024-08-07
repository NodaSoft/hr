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
	taskRESULT string // результат выполнения задачи
}

// taskCreator генерирует задачи и отправляет их в канал a
func taskCreator(a chan<- Ttype) {
	for {
		ft := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			ft = "Some error occurred"
		}
		a <- Ttype{cT: ft, id: int(time.Now().Unix())}
		time.Sleep(time.Millisecond * 500) // имитируем задержку создания задач
	}
}

// taskWorker обрабатывает задачу и возвращает её результат
func taskWorker(a Ttype) Ttype {
	tt, err := time.Parse(time.RFC3339, a.cT)
	if err != nil || tt.After(time.Now().Add(-20*time.Second)) {
		a.taskRESULT = "task has been succeeded"
	} else {
		a.taskRESULT = "something went wrong"
	}
	a.fT = time.Now().Format(time.RFC3339Nano)
	time.Sleep(time.Millisecond * 150)
	return a
}

// taskSorter распределяет задачи по каналам doneTasks и undoneTasks
func taskSorter(t Ttype, doneTasks chan<- Ttype, undoneTasks chan<- error) {
	if t.taskRESULT == "task has been succeeded" {
		doneTasks <- t
	} else {
		undoneTasks <- fmt.Errorf("task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
	}
}

// printResults периодически выводит результаты обработки задач
func printResults(errors []error, result map[int]Ttype, mutex *sync.Mutex) {
	for {
		time.Sleep(3 * time.Second)
		mutex.Lock()
		fmt.Println("Errors:")
		for _, r := range errors {
			fmt.Println(r)
		}

		fmt.Println("Done tasks:")
		for _, r := range result {
			fmt.Printf("Task ID: %d, Created: %s, Finished: %s, Result: %s\n", r.id, r.cT, r.fT, r.taskRESULT)
		}
		mutex.Unlock()
	}
}

func main() {
	var wg sync.WaitGroup
	var mutex sync.Mutex

	superChan := make(chan Ttype, 10)
	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	go taskCreator(superChan)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for t := range superChan {
			t = taskWorker(t)
			taskSorter(t, doneTasks, undoneTasks)
		}
	}()

	go func() {
		wg.Wait() // ожидание завершения всех задач
		close(doneTasks)
		close(undoneTasks)
	}()

	result := map[int]Ttype{}
	var errors []error

	wg.Add(1)
	go func() {
		defer wg.Done()
		for r := range doneTasks {
			mutex.Lock()
			result[r.id] = r
			mutex.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for r := range undoneTasks {
			mutex.Lock()
			errors = append(errors, r)
			mutex.Unlock()
		}
	}()

	go printResults(errors, result, &mutex)

	wg.Wait() // ожидаем завершения всех горутин
}
