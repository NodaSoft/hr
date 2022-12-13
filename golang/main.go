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

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

// A Task represents the meaning of our life as a path, not as a goal
type Task struct {
	id       int
	created  string // время создания
	finished string // время выполнения
	result   TaskResult
}

type TaskResult string

// the less hardcode the better
const (
	SuccessTaskResult = "task has been successed"
	FailTaskResult    = "something went wrong"
)

func main() {
	// send me feedback - t.me/kingarthurfish

	tasksCount := 10 // in the best world we get it from config or something
	tasks := make(chan Task)
	done := make(chan struct{})
	successedTasks := make([]any, 0)
	failedTasks := make([]any, 0)

	go taskCreator(tasks, tasksCount)
	go taskWorker(tasks, done, &successedTasks, &failedTasks)

	<-done // maybe better way is to use Context, but KISS is win there

	fmt.Println("Errors:", failedTasks)
	fmt.Println("Done tasks:", successedTasks)
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

// in fact I have never got an error here lol
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
