package main

import (
	"fmt"
	"time"
)

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// Приложение должно генерировать таски 10 сек. Каждые 3 секунды должно выводить в консоль результат всех обработанных к этому моменту тасков (отдельно успешные и отдельно с ошибками).

// ЗАДАНИЕ: сделать из плохого кода хороший и рабочий - as best as you can.
// Важно сохранить логику появления ошибочных тасков.
// Важно оставить асинхронные генерацию и обработку тасков.
// Сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через pull-request в github
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

// A Ttype represents a meaninglessness of our life
type TaskType struct {
	id             int
	createTime     time.Time // время создания
	executionTime  time.Time // время завершения
	result         string
	errorRecording bool
}

func tasksCreator(taskChannel chan TaskType) {
	ticker := time.NewTimer(10 * time.Second)
	defer ticker.Stop()
	defer close(taskChannel)

	for {
		select {
		case <-ticker.C:
			return
		default:
			now := time.Now()
			task := TaskType{
				id:             int(now.UnixNano()),
				createTime:     now,
				errorRecording: now.Nanosecond()%2 > 0,
			}
			taskChannel <- task
		}
	}
}

func taskWorker(task TaskType) TaskType {
	if time.Since(task.createTime) > 20*time.Second || task.errorRecording {
		task.errorRecording = true
		task.result = "something went wrong"
	} else {
		task.result = "task has been successed"
	}
	task.executionTime = time.Now()

	time.Sleep(time.Millisecond * 150)
	return task
}

func taskSorter(task TaskType, doneTask, undoneTask chan TaskType) {
	if task.errorRecording {
		undoneTask <- task
	} else {
		doneTask <- task
	}
}

func printing(tasks []TaskType, start *int) {
	end := len(tasks)
	for i := *start; i < end; i++ {
		if tasks[i].errorRecording {
			fmt.Printf("Task id %d time create %s, result %s\n",
				tasks[i].id,
				tasks[i].createTime.Format(time.RFC3339),
				tasks[i].result)
		} else {
			fmt.Printf("Task id %d time create %s, time ending %s, result %s\n",
				tasks[i].id,
				tasks[i].createTime.Format(time.RFC3339),
				tasks[i].executionTime.Format(time.RFC3339),
				tasks[i].result)
		}
	}
	*start = end
}

func main() {
	creatorChan := make(chan TaskType, 10)
	doneTask := make(chan TaskType)
	undoneTask := make(chan TaskType)
	defer close(doneTask)
	defer close(undoneTask)

	go tasksCreator(creatorChan)

	successResult := []TaskType{}
	errorsResult := []TaskType{}

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	startErr := 0
	startDone := 0
	for {
		select {
		case task, ok := <-creatorChan:
			if !ok {
				return
			}
			task = taskWorker(task)
			go taskSorter(task, doneTask, undoneTask)			
		case success := <-doneTask:
			successResult = append(successResult, success)
		case err := <-undoneTask:
			errorsResult = append(errorsResult, err)
		case <-ticker.C:
			fmt.Println("Done tasks:")
			printing(successResult, &startDone)

			fmt.Println("Errors:")
			printing(errorsResult, &startErr)
		}
	}
}
