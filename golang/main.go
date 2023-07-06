package main

import (
	"fmt"
	"math/rand"
	"sync"
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

type ResultMessage = string

var (
	TaskSuccessed ResultMessage = "task has been successed"
	TaskFailed    ResultMessage = "something went wrong"
)

type Task struct {
	id            int
	createdAt     string
	finishedAt    string
	resultMessage ResultMessage
}

func (t *Task) IsAlive() bool {
	tt, _ := time.Parse(time.RFC3339, t.createdAt)
	return tt.After(time.Now().Add(-20 * time.Second))
}

func taskCreator() <-chan Task {
	ch := make(chan Task, 10)
	go func() {
		for i := 0; i < 10; i++ {
			ft := time.Now().Format(time.RFC3339)
			// UnixNano() round to 1 microsecond, so the condition is always false
			if time.Now().UnixMicro()%2 > 0 { // вот такое условие появления ошибочных тасков
				ft = "Some error occured"
			}
			ch <- Task{createdAt: ft, id: int(time.Now().UnixMicro())} // передаем таск на выполнение

			// add some delay to avoid conflicts with ids
			time.Sleep(time.Microsecond * 1)
		}
		close(ch)
	}()
	return ch
}

type Worker struct {
	doneTasks  chan<- Task
	failedTask chan<- error
}

func (w *Worker) Listen(taskCh <-chan Task) {
	// outer goroutine to avoid blocking from caller side
	go func() {
		wg := sync.WaitGroup{}

		for t := range taskCh {
			wg.Add(1)
			go func(t Task) {
				t = w.doHardWork(t)
				w.demultiplex(t)
				wg.Done()
			}(t)
		}

		wg.Wait()
		close(w.doneTasks)
		close(w.failedTask)
	}()
}

// This method determines apporopriate out channel
func (w *Worker) demultiplex(t Task) {
	if t.resultMessage == TaskSuccessed {
		w.doneTasks <- t
	} else {
		w.failedTask <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.createdAt, t.resultMessage)
	}
}

func (w *Worker) doHardWork(task Task) Task {
	if task.IsAlive() {
		task.resultMessage = TaskSuccessed
	} else {
		task.resultMessage = TaskFailed
	}

	// add some fluctuations to delay
	time.Sleep(time.Millisecond*150 + time.Duration(rand.Float32()*100)*time.Millisecond)

	task.finishedAt = time.Now().Format(time.RFC3339Nano)
	return task
}

func main() {
	tasks := taskCreator()

	doneTasks := make(chan Task)
	failedTasks := make(chan error)

	worker := Worker{doneTasks: doneTasks, failedTask: failedTasks}
	worker.Listen(tasks)

	result := map[int]Task{}
	errResult := []error{}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		for r := range doneTasks {
			result[r.id] = r
		}
		wg.Done()
	}()
	go func() {
		for r := range failedTasks {
			errResult = append(errResult, r)
		}
		wg.Done()
	}()

	wg.Wait()

	println("Errors:")
	for _, errorMess := range errResult {
		fmt.Println(errorMess)
	}
	println("Done tasks:")
	for _, r := range result {
		fmt.Println(r.id)
	}
}
