package main

import (
	"errors"
	"fmt"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// * обновленный код отправить через pull-request.

// Приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме.
// Должно выводить успешные таски и ошибки по мере выполнения.
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

// A Ttype represents a meaninglessness of our life
type Task struct {
	Id           int
	CreationTime string // время создания
	StartTime    string // время выполнения
	Error        error
	Result       []byte
}

func main() {
	taskCreator := func(out chan<- Task) {
		var id int = 0

		for {
			var err error = nil
			ft := time.Now().Format(time.RFC3339)

			// time.Now().Nanosecond()%2 > 0, у меня всегда возвращал false
			if (time.Now().Nanosecond()/1000)%2 > 0 {
				err = errors.New("Some error occured")
			}

			id += 1
			out <- Task{
				CreationTime: ft,
				Error:        err,
				Id:           id,
			}
		}
	}

	taskWorker := func(input <-chan Task, done chan<- Task, errChan chan<- error) {
		for {
			task := <-input

			if task.Error != nil {
				errChan <- fmt.Errorf("Task id [%d] time [%s], error [%s]", task.Id, task.CreationTime, task.Error.Error())
				continue
			}

			creationTime, err := time.Parse(time.RFC3339, task.CreationTime)
			if err != nil { // Немного бессмысленно, но пусть будет
				errChan <- fmt.Errorf("Task id [%d] time [%s], error [%s]", task.Id, task.CreationTime, task.Error.Error())
				continue
			}

			if creationTime.After(time.Now().Add(-20 * time.Second)) {
				task.Result = []byte("task has been successed")
			} else {
				task.Result = []byte("something went wrong")
			}
			task.StartTime = time.Now().Format(time.RFC3339Nano)
			time.Sleep(time.Millisecond * 150)

			done <- task
		}
	}

	taskChan := make(chan Task, 10)
	doneTasks := make(chan Task, 10)
	undoneTasks := make(chan error, 10)

	go taskCreator(taskChan)
	go taskWorker(taskChan, doneTasks, undoneTasks)

	fmt.Println("Here")
	for {
		select {
		case done := <-doneTasks:
			fmt.Printf("Task id [%d] time [%s], result: [%s]\n", done.Id, done.CreationTime, done.Result)
		case err := <-undoneTasks:
			fmt.Println(err)
		}
	}
}
