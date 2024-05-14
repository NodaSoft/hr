package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// После обработки тасков в течении 3 секунд приложение должно выводить накопленные к этому моменту успешные таски и отдельно ошибки обработки тасков.

// ЗАДАНИЕ: сделать из плохого кода хороший и рабочий - as best as you can.
// Важно сохранить логику появления ошибочных тасков.
// Важно оставить асинхронные генерацию и обработку тасков.
// Сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через pull-request в github
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id         int
	createdAt  string // время создания
	finishedAt string // время выполнения
	taskResult []byte
	err        error
}

const (
	taskChanSize = 10
	workersCount = 4
)

var (
	ErrSomethingWentWrong = errors.New("something went wrong")
	ErrTaskFailed         = errors.New("task failed")
)

func NewTask() Ttype {
	task := Ttype{createdAt: time.Now().Format(time.RFC3339), id: int(time.Now().Unix())}
	if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
		task.err = ErrTaskFailed
	}
	return task
}

func taskWorker(task Ttype) Ttype {
	createdAt, err := time.Parse(time.RFC3339, task.createdAt)
	if err != nil {
		task.err = ErrSomethingWentWrong
		return task
	}
	if createdAt.After(time.Now().Add(-20 * time.Second)) {
		time.Sleep(time.Millisecond * 150)
		task.taskResult = []byte("task has been successed")
	} else {
		task.err = ErrTaskFailed
	}
	task.finishedAt = time.Now().Format(time.RFC3339Nano)
	return task
}

func taskSorter(task Ttype, doneTasks chan Ttype, failedTasks chan error) {
	if task.err != nil {
		failedTasks <- fmt.Errorf("Task id %d time %s, error %s", task.id, task.createdAt, task.err)
		return
	}
	doneTasks <- task
}

func main() {
	taskChan := make(chan Ttype, taskChanSize)
	doneTasks := make(chan Ttype)
	failedTasks := make(chan error)

	go func() {
		for {
			task := NewTask()
			taskChan <- task // передаем таск на выполнение
		}
	}()

	var wg sync.WaitGroup
	for i := 0; i < workersCount; i++ { // Adjust the number of workers as needed
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskChan {
				task = taskWorker(task)
				go taskSorter(task, doneTasks, failedTasks)
			}
		}()
	}

	doneResults := map[int]Ttype{}
	errorResults := []error{}
	ticker := time.NewTicker(time.Second * 3)

	for {
		select {
		case result := <-doneTasks:
			doneResults[result.id] = result
		case result := <-failedTasks:
			errorResults = append(errorResults, result)
		case <-ticker.C:
			println("Errors:")
			for _, err := range errorResults {
				fmt.Println(err)
			}
			errorResults = []error{} // Reset results for next interval

			println("Done tasks:")
			for _, task := range doneResults {
				fmt.Printf("Task ID: %d, createdAt: %s, finishedAt: %s, Result: %s\n", task.id, task.createdAt, task.finishedAt, task.taskResult)
			}
			doneResults = map[int]Ttype{} // Reset results for next interval
		}
	}
}
