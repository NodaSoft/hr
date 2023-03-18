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

// A Task represents a meaninglessness of our life
type Task struct {
	id            int
	createTime    string
	executionTime string
	taskResult    []byte
}

func main() {
	taskCreator := func(a chan Task) {
		go func() {
			for {
				createTime := time.Now().Format(time.RFC3339)
				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					createTime = "Some error occured"
				}
				a <- Task{id: int(time.Now().Unix()), createTime: createTime} // передаем таск на выполнение
			}
		}()
	}

	taskChan := make(chan Task, 10)

	go taskCreator(taskChan)

	taskWorker := func(a Task) Task {
		tt, _ := time.Parse(time.RFC3339, a.createTime)
		if tt.After(time.Now().Add(-20 * time.Second)) {
			a.taskResult = []byte("task has been successed")
		} else {
			a.taskResult = []byte("something went wrong")
		}
		a.executionTime = time.Now().Format(time.RFC3339Nano)

		//time.Sleep(time.Millisecond * 150)

		return a
	}

	doneTasks := make(chan Task)
	undoneTasks := make(chan error)

	taskSorter := func(t Task) {
		if string(t.taskResult[14:]) == "successed" {
			doneTasks <- t
		} else {
			undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.createTime, t.taskResult)
		}
	}

	go func() {
		// получение тасков
		for task := range taskChan {
			task = taskWorker(task)
			go taskSorter(task)
		}
		close(taskChan)
	}()

	result := map[int]Task{}
	err := []error{}
	go func() {
		for doneTask := range doneTasks {
			go func() {
				result[doneTask.id] = doneTask
			}()
		}
		for undoneTask := range undoneTasks {
			go func() {
				err = append(err, undoneTask)
			}()
		}
		close(doneTasks)
		close(undoneTasks)
	}()

	//time.Sleep(time.Second * 3)

	println("Errors:")
	for r := range err {
		println(r)
	}

	println("Done tasks:")
	for r := range result {
		println(r)
	}
}
