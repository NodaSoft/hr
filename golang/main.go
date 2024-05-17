package main

import (
	"fmt"
	"time"
)

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// После обработки тасков в течении 3 секунд приложение должно выводить накопленные к этому моменту успешные таски и отдельно ошибки обработки тасков.

// ЗАДАНИЕ: сделать из плохого кода хороший и рабочий - as best as you can.
// Важно сохранить логику появления ошибочных тасков.
// Важно оставить асинхронные генерацию и обработку тасков.
// Сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через pull-request в github
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

type Task struct {
	id             int
	creationTime   string
	completionTime string
	result         []byte
}

func main() {
	scheduleTasks := func(ch chan Task) {
		go func() {
			for {
				nowStr := time.Now().Format(time.RFC3339)
				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					nowStr = "Some error occured"
				}
				taskID := int(time.Now().UnixNano())         // NOTE: но лучше вообще генерировать случайный ID (например, GUID)
				ch <- Task{creationTime: nowStr, id: taskID} // передаем таск на выполнение
			}
		}()
	}

	superChan := make(chan Task, 10)

	go scheduleTasks(superChan)

	processTask := func(task Task) Task {
		parsedCreationTime, _ := time.Parse(time.RFC3339, task.creationTime)
		if parsedCreationTime.After(time.Now().Add(-20 * time.Second)) {
			task.result = []byte("task has been successed")
		} else {
			task.result = []byte("something went wrong")
		}
		task.completionTime = time.Now().Format(time.RFC3339Nano)

		time.Sleep(time.Millisecond * 150)

		return task
	}

	doneTasksChan := make(chan Task)
	taskErrorsChan := make(chan error)

	sortTask := func(task Task) {
		if string(task.result[14:]) == "successed" {
			doneTasksChan <- task
		} else {
			taskErrorsChan <- fmt.Errorf("Task id %d time %s, error %s", task.id, task.creationTime, task.result)
		}
	}

	go func() {
		// получение тасков
		for task := range superChan {
			task = processTask(task)
			go sortTask(task)
		}
		close(superChan)
	}()

	doneTasks := map[int]Task{}
	taskErrors := []error{}
	go func() {
		for doneTask := range doneTasksChan {
			go func() {
				doneTasks[doneTask.id] = doneTask
			}()
		}
		for taskErr := range taskErrorsChan {
			go func() {
				taskErrors = append(taskErrors, taskErr)
			}()
		}
		close(doneTasksChan)
		close(taskErrorsChan)
	}()

	time.Sleep(time.Second * 3)

	println("Errors:")
	for r := range taskErrors {
		println(r)
	}

	println("Done tasks:")
	for r := range doneTasks {
		println(r)
	}
}
