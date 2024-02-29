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

// A Task represents a single task
type Task struct {
	id           int
	creationTime time.Time // время создания
	finishTime   time.Time // время выполнения
	taskResult   []byte
}

func main() {
	taskCreator := func(a chan Task) {
		go func() {
			for {
				ft := time.Now().Format(time.RFC3339)
				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					ft = "Some error occured"
				}
				a <- Task{creationTime: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
			}
		}()
	}

	superChan := make(chan Task, 10)

	go taskCreator(superChan)

	taskWorker := func(a Task) Task {
		tt, _ := time.Parse(time.RFC3339, a.creationTime)
		if tt.After(time.Now().Add(-20 * time.Second)) {
			a.taskResult = []byte("task has been successed")
		} else {
			a.taskResult = []byte("something went wrong")
		}
		a.finishTime = time.Now().Format(time.RFC3339Nano)

		time.Sleep(time.Millisecond * 150)

		return a
	}

	doneTasks := make(chan Task)
	undoneTasks := make(chan error)

	tasksorter := func(t Task) {
		if string(t.taskResult[14:]) == "successed" {
			doneTasks <- t
		} else {
			undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.creationTime, t.taskResult)
		}
	}

	go func() {
		// получение тасков
		for t := range superChan {
			t = taskWorker(t)
			go tasksorter(t)
		}
		close(superChan)
	}()

	result := map[int]Task{}
	err := []error{}
	go func() {
		for r := range doneTasks {
			go func() {
				result[r.id] = r
			}()
		}
		for r := range undoneTasks {
			go func() {
				err = append(err, r)
			}()
		}
		close(doneTasks)
		close(undoneTasks)
	}()

	time.Sleep(time.Second * 3)

	println("Errors:")
	for r := range err {
		println(r)
	}

	println("Done tasks:")
	for r := range result {
		println(r)
	}
}
