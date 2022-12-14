package main

import (
	"fmt"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки

const (
	SuccessTaskResult = "task has been succeeded"
	FailTaskResult    = "something went wrong"
)

// A Task represents a meaninglessness of our life
type Task struct {
	id       int
	created  string // время создания
	finished string // время выполнения
	result   TaskResult
}

type TaskResult string

func main() {
	tasksCount := 10

	tasks := make(chan Task)
	done := make(chan struct{})

	successesTasks := make([]any, 0)
	failedTasks := make([]any, 0)

	go taskCreator(tasks, tasksCount)
	go taskWorker(tasks, done, &successesTasks, &failedTasks)

	<-done

	fmt.Println("Errors:", failedTasks)
	fmt.Println("Done tasks:", successesTasks)
}

func taskCreator(tasks chan<- Task, count int) {
	for i := 0; i < count; i++ {
		ft := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			ft = "Some error occurred"
		}
		tasks <- Task{created: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
	}
	close(tasks)
}

func taskExecutor(task Task) Task {
	tt, _ := time.Parse(time.RFC3339, task.created)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		task.result = SuccessTaskResult
	} else {
		task.result = FailTaskResult
	}
	task.finished = time.Now().Format(time.RFC3339Nano)
	time.Sleep(time.Millisecond * 150)

	return task
}

func taskWorker(tasks <-chan Task, done chan<- struct{}, st, ft *[]any) {
	for {
		task, ok := <-tasks
		if !ok {
			done <- struct{}{}
			return
		} else {
			t := taskExecutor(task)
			if t.result == SuccessTaskResult {
				*st = append(*st, t.id)
			} else {
				*ft = append(*ft, fmt.Errorf("task id %d time %s, error %s", t.id, t.created, t.result))
			}
		}
	}
}
