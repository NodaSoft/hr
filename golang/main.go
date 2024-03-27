package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

const (
	tasksCount   = 10
	workersCount = 5
)

type Task struct {
	ID        int
	CreatedAt time.Time     // время создания
	Dur       time.Duration // время выполнения
	Err       error
	Result    string
}

func taskProducer(count int) <-chan Task {
	ch := make(chan Task, count)

	go func(ch chan Task) {
		for i := 0; i < count; i++ {
			var err error
			if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				err = errors.New("error occured during creation")
			}
			ch <- Task{ID: i, CreatedAt: time.Now(), Err: err}
		}

		close(ch)
	}(ch)

	return ch
}

func taskWorker(t Task) Task {
	if !t.CreatedAt.After(time.Now().Add(-20 * time.Second)) {
		t.Err = errors.New("error occured during processing")
		return t
	}

	time.Sleep(time.Millisecond * 150)

	t.Result = "task has been successed"
	t.Dur = time.Now().Sub(t.CreatedAt)

	return t

}

func main() {
	tasksChan := taskProducer(tasksCount)

	doneTasksChan := make(chan Task)
	failedTasksChan := make(chan Task)

	wg := sync.WaitGroup{}
	for i := 0; i < workersCount; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for task := range tasksChan {
				t := taskWorker(task)

				if t.Err == nil {
					doneTasksChan <- t
				} else {
					failedTasksChan <- t
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(doneTasksChan)
		close(failedTasksChan)
	}()

	var results []string
	var errors []string

	for i := 0; i < tasksCount; i++ {
		select {
		case t := <-doneTasksChan:
			results = append(results, fmt.Sprintf("id: %d, duration: %s, result: %s", t.ID, t.Dur.String(), t.Result))
		case t := <-failedTasksChan:
			errors = append(errors, fmt.Sprintf("id: %d, duration: %s, error: %s", t.ID, t.Dur.String(), t.Err.Error()))
		}
	}

	fmt.Println("Failed tasks:")
	for _, v := range errors {
		fmt.Println(v)
	}

	fmt.Println("Done tasks:")
	for _, v := range results {
		fmt.Println(v)
	}
}

