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

// Привет! Перепишем код, учитывая следующие изменения:
//
// 1 - Переименуем некоторые переменные и функции, чтобы сделать код более понятным и легче читаемым.
// 2 - Добавим многопоточность для обработки задач, чтобы ускорить выполнение и увеличить производительность.
// 3 - Изменим логику обработки ошибок в задачах, чтобы более явно отслеживать ошибки.
// 4 - Оптимизируем работу горутин, чтобы избежать блокировок и зависаний.

type Task struct {
	ID          int
	CreatedTime time.Time
	Result      string
	Err         error
}

func main() {
	taskCreator := func(taskChan chan Task) {
		for {
			createdTime := time.Now()
			var result string
			if createdTime.Nanosecond()%2 > 0 {
				result = "Something went wrong"
			} else {
				result = "Task has been succeeded"
			}
			taskChan <- Task{
				ID:          int(createdTime.Unix()),
				CreatedTime: createdTime,
				Result:      result,
			}
		}
	}

	taskWorker := func(task Task) Task {
		time.Sleep(time.Millisecond * 150)
		if task.Result == "Something went wrong" {
			task.Err = fmt.Errorf("Task ID %d failed at %s", task.ID, task.CreatedTime.Format(time.RFC3339))
		}
		return task
	}

	tasks := make(chan Task, 10)
	doneTasks := make(chan Task)
	failedTasks := make(chan Task)

	go taskCreator(tasks)

	for i := 0; i < 5; i++ {
		go func() {
			for task := range tasks {
				processedTask := taskWorker(task)
				if processedTask.Err != nil {
					failedTasks <- processedTask
				} else {
					doneTasks <- processedTask
				}
			}
		}()
	}

	go func() {
		for doneTask := range doneTasks {
			fmt.Printf("Task ID %d succeeded at %s\n", doneTask.ID, doneTask.CreatedTime.Format(time.RFC3339))
		}
	}()

	go func() {
		for failedTask := range failedTasks {
			fmt.Printf("%s\n", failedTask.Err)
		}
	}()

	time.Sleep(time.Second * 3)

	close(tasks)
	close(doneTasks)
	close(failedTasks)
}
