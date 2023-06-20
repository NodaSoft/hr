package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id         int
	createTime string // время создания
	finishTime string // время выполнения

	err error
}

var taskId = &atomic.Int64{}

func newTask() *Ttype {
	now := time.Now()
	ft := now.Format(time.RFC3339)

	if now.Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
		ft = "Some error occured"
	}
	return &Ttype{createTime: ft, id: int(taskId.Add(1))}
}

func taskSpawner(done <-chan struct{}) <-chan *Ttype {
	ch := make(chan *Ttype)

	go func() {
		defer close(ch)

		for {
			select {
			case <-done:
				return
			case ch <- newTask(): // передаем таск на выполнение
			}
		}
	}()

	return ch
}

func processTask(task *Ttype) {
	if tt, err := time.Parse(time.RFC3339, task.createTime); err != nil {
		task.err = fmt.Errorf("time parse error: %w", err)

		return
	} else if !tt.After(time.Now().Add(-20 * time.Second)) {
		task.err = fmt.Errorf("something went wrong")

		return
	}

	task.finishTime = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150) // work emulation
}

func taskWorker(tasks <-chan *Ttype) <-chan *Ttype {
	ch := make(chan *Ttype)

	go func() {
		defer close(ch)

		for task := range tasks {
			processTask(task)

			ch <- task
		}
	}()

	return ch
}

func fanIn[T any](ch1, ch2 <-chan T) <-chan T {
	ch := make(chan T)

	go func() {
		defer close(ch)

		for {
			if ch1 == nil && ch2 == nil {
				return
			}

			select {
			case val, ok := <-ch1:
				if !ok {
					ch1 = nil
					break
				}

				ch <- val
			case val, ok := <-ch2:
				if !ok {
					ch2 = nil
					break
				}

				ch <- val
			}
		}
	}()

	return ch
}

func taskSorter(tasks <-chan *Ttype) (<-chan *Ttype, <-chan *Ttype) {
	doneTasks, undoneTasks := make(chan *Ttype), make(chan *Ttype)

	go func() {
		defer func() {
			close(doneTasks)
			close(undoneTasks)
		}()

		for task := range tasks {
			if task.err != nil {
				undoneTasks <- task
			} else {
				doneTasks <- task
			}
		}
	}()

	return doneTasks, undoneTasks
}

const workersCnt = 5

func main() {
	done := make(chan struct{})

	tasks := taskSpawner(done)

	var taskResults <-chan *Ttype = nil
	for i := 0; i < workersCnt; i++ {
		taskResults = fanIn(taskResults, taskWorker(tasks))
	}

	doneTasks, undoneTasks := taskSorter(taskResults)

	time.Sleep(time.Second * 3)

	close(done)

	errs := []error{}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	// if we don't read one channel and sorter is trying to write to that inactive channel we will block

	go func() {
		for r := range undoneTasks {
			errs = append(errs, r.err)
		}
		wg.Done()
	}()

	fmt.Println("Done tasks:")
	for r := range doneTasks {
		fmt.Println(r.id)
	}

	wg.Wait()

	fmt.Println("Errors:")
	for _, err := range errs {
		fmt.Println(err)
	}
}
