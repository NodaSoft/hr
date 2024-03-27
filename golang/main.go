package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных задач;
// * сделать правильную многопоточность обработки задач.
// Обновленный код отправить через merge-request.

// Приложение эмулирует получение и обработку задач - пытается это делать в многопоточном режиме.
// В конце должно выводить успешные таски и ошибки выполнения остальных задач.

const (
	taskCount   = 10
	workerCount = 3
	timeLayout  = time.RFC3339Nano
)

var (
	errCreateTask  = errors.New("error occurred during creation")
	errProcessTask = errors.New("error occurred during processing")
)

// Task олицетворяет бессмысленность нашей жизни.
type Task struct {
	ID         int
	CreatedAt  time.Time
	FinishedAt time.Time
	Err        error
	Result     string
}

func (t *Task) String() string {
	if t.Err != nil {
		return fmt.Sprintf("id: %d, error: %s",
			t.ID,
			t.Err.Error())
	}
	return fmt.Sprintf("id: %d, created: %s, finished: %s, result: %s",
		t.ID,
		t.CreatedAt.Format(timeLayout),
		t.FinishedAt.Format(timeLayout),
		t.Result)
}

func taskProduce(count int) <-chan Task {
	out := make(chan Task)

	go func() {
		for i := 0; i < count; i++ {
			var err error
			if time.Now().Nanosecond()%2 > 0 { // Условие появления ошибочных задач
				err = errCreateTask
			}
			out <- Task{ID: i, CreatedAt: time.Now(), Err: err}
		}

		close(out)
	}()

	return out
}

func taskProcess(t Task) Task {
	if t.Err != nil {
		return t
	}

	if !t.CreatedAt.After(time.Now().Add(-20 * time.Second)) {
		t.Err = errProcessTask
		return t
	}

	time.Sleep(time.Millisecond * 150) // Имитация времени обработки задачи

	t.Result = "task has been successfully done"
	t.FinishedAt = time.Now()

	return t
}

func main() {
	tasksCh := taskProduce(taskCount)

	doneTasksCh := make(chan Task)
	failedTasksCh := make(chan Task)

	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for task := range tasksCh {
				t := taskProcess(task)
				if t.Err != nil {
					failedTasksCh <- t
					continue
				}
				doneTasksCh <- t
			}
		}()
	}

	go func() {
		wg.Wait()
		close(doneTasksCh)
		close(failedTasksCh)
	}()

	doneTasks := make([]Task, 0, taskCount)
	failedTasks := make([]Task, 0, taskCount)

	for i := 0; i < taskCount; i++ {
		select {
		case t := <-doneTasksCh:
			doneTasks = append(doneTasks, t)
		case t := <-failedTasksCh:
			failedTasks = append(failedTasks, t)
		}
	}

	fmt.Println("Errors:")
	for _, t := range failedTasks {
		fmt.Println(t.String())
	}

	fmt.Println("Done tasks:")
	for _, t := range doneTasks {
		fmt.Println(t.String())
	}
}
