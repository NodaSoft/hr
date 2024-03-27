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


const tasksAmount = 10

// Struct for a task represents when task was created, finished and whether any errors occurred.
type Task struct {
	id           int
	created      time.Time // время создания
	finished     time.Time // время выполнения
	errorOccured bool
}

// Я изменил структуру кода, сделав её чище
func main() {
	// Функция, которая отправляет в очередь n тасков
	taskCreator := func(queue chan<- Task, n int) {
		go func() {
			defer close(queue)

			for i := 0; i < n; i++ {
				currTime := time.Now()

				// Пример для проверки вывода ошибочных тасков
				// if i % 2 == 0 {
				// 	currTime = currTime.Add(time.Nanosecond)
				// }

				newTask := Task{created: currTime, id: int(currTime.Unix())}
				if currTime.Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					newTask.errorOccured = true
				}
				queue <- newTask
			}
		}()
	}

	// Функция выполняет работу над тасками и закрывает все каналы (в том числе done), когда работа завершена
	taskWorker := func(done chan<- struct{}, queue <-chan Task, succeeded chan<- Task, failed chan<- error) {
		work := func(task Task) {
			if !task.errorOccured {
				task.finished = time.Now()
				succeeded <- task
			} else {
				failed <- fmt.Errorf("Task: %x failed at %s", task.id, task.created)
			}
			time.Sleep(150 * time.Millisecond)
		}

		go func() {
			// Закрываем done в первую очередь, чтобы пришел сигнал о завершении работы над тасками
			defer close(done)
			defer close(succeeded)
			defer close(failed)

			for task := range queue {
				work(task)
			}
		}()
	}

	tasksQueue := make(chan Task, 10)
	succeededTasks := make(chan Task)
	failedTasks := make(chan error)
	done := make(chan struct{})

	taskCreator(tasksQueue, tasksAmount)
	taskWorker(done, tasksQueue, succeededTasks, failedTasks)

	// В условии сказано, что нужно выводить выполненые и ошибочные таски, а не записывать их в мапу
	go func() {
		for {
			select {
			case task := <-succeededTasks:
				fmt.Printf("Task: %x succeeded at %s\n", task.id, task.finished)
			case err := <-failedTasks:
				fmt.Println(err)
			case <-done:
				return
			}
		}
	}()

	<-done
}
