package main

import (
	"sync"
	"testing"
	"time"
)

func makeWorker() *Worker {
	return &Worker{
		doneTasks:   make(chan *Task, 10),
		undoneTasks: make(chan error, 10),
		taskWithErr: make(chan *Task, 10),
		destroy:     make(chan bool, 1),
		begin: func() {

		},
		done: func() {

		},
		once: sync.Once{},
	}
}

func TestWorkerSubscribe(t *testing.T) {
	doneTasks := make(chan *Task, 10)
	undoneTasks := make(chan error, 10)
	withErr := make(chan *Task, 10)
	w := &Worker{
		doneTasks:   doneTasks,
		undoneTasks: undoneTasks,
		taskWithErr: withErr,
		destroy:     make(chan bool, 1),
		begin: func() {

		},
		done: func() {

		},
		once: sync.Once{},
	}
	masterChan := make(chan *Task, 10)
	w.Subscribe(masterChan)
	//устаревшая таска без ошибки
	since := time.Date(2000, 1, 1, 0, 0, 0, 2, time.Local)
	task := NewTask(since)
	masterChan <- task
	//устаревшая таска с ошибкой
	since = time.Date(2000, 1, 1, 0, 0, 0, 1, time.Local)
	task = NewTask(since)
	masterChan <- task

	//свежая таска с ошибкой
	since = time.Date(3000, 1, 1, 0, 0, 0, 1, time.Local)
	task = NewTask(since)
	masterChan <- task

	//свежая таска без ошибки
	since = time.Date(3000, 1, 1, 0, 0, 0, 2, time.Local)
	task = NewTask(since)
	masterChan <- task
	//ждем обработки
	time.Sleep(3 * time.Second)

	//итого должно быть два сообщения в канале undone
	//одно сообщение в канале done
	//одно сообщение в канале withErr
	close(undoneTasks)
	close(doneTasks)
	close(withErr)
	i := 0
	for range undoneTasks {
		i++
	}
	if i != 2 {
		t.Fatalf("assertions don't eq %d !=%d", 2, i)
	}

	i = 0
	for range doneTasks {
		i++
	}

	if i != 1 {
		t.Fatalf("assertions don't eq %d !=%d", 1, i)

	}

	i = 0
	for range withErr {
		i++
	}

	if i != 1 {
		t.Fatalf("assertions don't eq %d !=%d", 1, i)

	}

}
