package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

const TaskCount = 10

const (
	Success = 1 << iota
	Failed
)

// A Task represents a meaninglessness of our life
type Task struct {
	id       uuid.UUID
	statTime string // время создания
	endTime  string // время выполнения
	status   uint8
}

func taskProducer(c chan Task) {
	for i := 0; i < TaskCount; i++ {
		formattedTime := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			formattedTime = "Some error occured"
		}
		c <- Task{statTime: formattedTime, id: uuid.New()} // передаем таск на выполнение
	}
}

func taskWorker(t Task) Task {
	tt, err := time.Parse(time.RFC3339, t.statTime)
	if err == nil && tt.After(time.Now().Add(-20*time.Second)) {
		t.status = Success
	} else {
		t.status = Failed
	}
	t.endTime = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)

	return t
}

func taskSorter(t Task, done chan Task, failed chan error) {
	if t.status == Success {
		done <- t
	} else {
		failed <- fmt.Errorf("Task id %d time %s failed", t.id, t.statTime)
	}
}

func main() {
	tasks := make(chan Task, TaskCount)

	go taskProducer(tasks)

	doneTasks := make(chan Task)
	failedTasks := make(chan error)

	go func() {
		// получение тасков
		for t := range tasks {
			t = taskWorker(t)
			go taskSorter(t, doneTasks, failedTasks)
		}
		close(tasks)
	}()

	result := map[uuid.UUID]Task{}
	errors := []error{}
	go func() {
		for task := range doneTasks {
			go func(t Task) {
				result[t.id] = t
			}(task)
		}
		for err := range failedTasks {
			go func(e error) {
				errors = append(errors, e)
			}(err)
		}
		close(doneTasks)
		close(failedTasks)
	}()

	time.Sleep(time.Second * 3)

	println("Errors:")
	for err := range errors {
		fmt.Println(err)
	}

	println("Done tasks:")
	for _, v := range result {
		fmt.Printf("%#v\n", v)
	}
}
