package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

// A TaskType represents a meaninglessness of our life
type TaskType struct {
	id         int
	tsCreate   string // время создания
	tsExecuted string // время выполнения
	result     []byte
}

type Errs struct {
	mx     sync.RWMutex
	values []error
}

func (e *Errs) Append(err error) {
	e.mx.Lock()
	defer e.mx.Unlock()

	e.values = append(e.values, err)
}

func (e *Errs) Values() []error {
	e.mx.Lock()
	defer e.mx.Unlock()

	return e.values
}

func createTasks(stopChan chan os.Signal, out chan TaskType) {
	wt := &sync.WaitGroup{}

	for {

		// For kill creating task
		select {
		case <-stopChan:
			log.Println("Stop signal")
			wt.Wait()
			return
		default:
		}

		wt.Add(1)
		go func() {
			defer wt.Done()
			taskTime := time.Now().Format(time.RFC3339)

			if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				taskTime = "Some error occured"
			}

			out <- TaskType{
				id:       int(time.Now().Unix()),
				tsCreate: taskTime,
			} // передаем таск на выполнение
		}()
	}
}

func taskExecute(a *TaskType) *TaskType {
	tt, _ := time.Parse(time.RFC3339, a.tsCreate)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		a.result = []byte("task has been successed")
	} else {
		a.result = []byte("something went wrong")
	}
	a.tsExecuted = time.Now().Format(time.RFC3339Nano)
	time.Sleep(time.Millisecond * 150)
	return a
}

func sortTask(doneTasks chan TaskType, undoneTasks chan error, t *TaskType) {
	if len(t.result) > 0 {
		if string(t.result[14:]) == "successed" {
			doneTasks <- *t
			return
		}
	}

	undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.tsCreate, t.result)
}

func main() {
	undoneTasks := make(chan error)
	superChan := make(chan TaskType, 10)
	doneTasks := make(chan TaskType)
	result := map[int]TaskType{}
	wg := sync.WaitGroup{}
	errs := Errs{}

	// Check stop from OS
	var stopChan = make(chan os.Signal, 2)
	signal.Notify(
		stopChan,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGKILL,
	)

	// Init tasks
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(superChan)

		createTasks(stopChan, superChan)
	}()

	// Handling task
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(doneTasks)
		defer close(undoneTasks)

		wt := &sync.WaitGroup{}

		// Get tasks
		for t := range superChan {
			wt.Add(1)
			t := t

			go func() {
				sortTask(
					doneTasks,
					undoneTasks,
					taskExecute(&t),
				)
			}()
		}
	}()

	// Check success handling a tasks
	wg.Add(1)
	go func() {
		defer wg.Done()

		for r := range doneTasks {
			wg.Add(1)
			r := r

			go func() {
				defer wg.Done()
				result[r.id] = r
			}()
		}
	}()

	// Check filed handling a tasks
	wg.Add(1)
	go func() {
		defer wg.Done()

		for r := range undoneTasks {
			wg.Add(1)
			r := r

			go func() {
				defer wg.Done()
				errs.Append(r)
			}()
		}
	}()

	// End all a gorutines
	wg.Wait()

	println("Done tasks:")
	for r := range result {
		println(r)
	}

	println("Errors:")
	for r := range errs.Values() {
		println(r)
	}
}
