package main

import (
	"fmt"
	"sync"
	"time"
)

type task struct {
	id         int
	createdAt  string // время создания
	finishedAt string // время выполнения
	flag       bool
	result     []byte
}

// Функция, которая генерирует таски
func taskCreturer(ch chan<- task, stop <-chan struct{}) {
	for {
		select {
		case <-stop:
			close(ch)
			return
		default:
			flag := true
			start := time.Now().Format(time.RFC3339)

			if time.Now().Nanosecond()%2 > 0 { // условие появления ошибочных тасков
				start = "Some error occured"
				flag = false
			}
			task := task{
				createdAt: start,
				//id:        int(time.Now().Unix()), //секнда - слишком долго
				id:   int(time.Now().Nanosecond()),
				flag: flag,
			}
			ch <- task
		}
	}

}

// Функция, которая обрабатывает таски
func taskWorker(in <-chan task, out chan<- task) {
	for task := range in {
		if task.flag {
			creationTime, _ := time.Parse(time.RFC3339, task.createdAt)
			if creationTime.After(time.Now().Add(-20 * time.Second)) {
				task.result = []byte("task has been successed")
			} else {
				task.result = []byte("something went wrong")
				task.flag = false
			}
		} else {
			task.result = []byte("something went wrong")
		}

		task.finishedAt = time.Now().Format(time.RFC3339Nano)

		time.Sleep(150 * time.Millisecond) // Пауза для имитации обработки
		out <- task
	}
	close(out)
	return
}

// Функция, которая распределяет таски по каналам
func distributeTask(num task, doneChan, errChan chan<- task) {
	if num.flag {
		doneChan <- num
	} else {
		errChan <- num
	}
}

// Функция для записи тасков из каналов в мапы
func recordTask(doneChan, errChan <-chan task, mapsDone, mapsErr *sync.Map, stop <-chan struct{}) {
	for {
		select {
		case num := <-doneChan:
			mapsDone.Store(num.id, num)
		case num := <-errChan:
			mapsErr.Store(num.id, num)
		case <-stop:
			return
		}
	}
}

// Функция для вывода содержимого мап
func printTask(mapsDone, mapsErr *sync.Map) {
	fmt.Println("Done tasks:")
	mapsDone.Range(func(key, value interface{}) bool {
		task := value.(task)
		fmt.Printf("Task ID: %d, createdAt: %s, finishedAt: %s, Result: %s\n", key, task.createdAt, task.finishedAt, task.result)
		return true
	})

	fmt.Println("Errors:")
	mapsErr.Range(func(key, value interface{}) bool {
		task := value.(task)
		fmt.Printf("Task ID: %d, error: %s\n", key, task.result)
		return true
	})
}

func main() {

	taskChan := make(chan task, 10)
	processedChan := make(chan task)
	doneChan := make(chan task)
	errChan := make(chan task)
	stop := make(chan struct{})

	mapsDone := &sync.Map{}
	mapsErr := &sync.Map{}

	go taskCreturer(taskChan, stop)
	go taskWorker(taskChan, processedChan)

	go func() {
		for num := range processedChan {
			distributeTask(num, doneChan, errChan)
		}
		close(doneChan)
		close(errChan)
		return
	}()

	go recordTask(doneChan, errChan, mapsDone, mapsErr, stop)

	timer := time.NewTimer(10 * time.Second)
	ticker := time.NewTicker(3 * time.Second)

	for {
		select {
		case <-timer.C:
			printTask(mapsDone, mapsErr)
			stop <- struct{}{}
			stop <- struct{}{}
			close(stop)
			fmt.Println("FINISH")
			return
		case <-ticker.C:
			printTask(mapsDone, mapsErr)
		}
	}
}
