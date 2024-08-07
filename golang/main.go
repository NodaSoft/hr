package main

import (
	"fmt"
	"strings"
	"sync"
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

const (
	tasksAmount          = 10
	timeFormat           = time.RFC3339
	executionTimeFormat  = time.RFC3339Nano
	taskHandlingDelay    = 150 * time.Millisecond
	maxTaskCreatingDelay = 20 * time.Second
)

// A Task represents a meaninglessness of our life
type Task struct {
	id            int
	createTime    string // время создания
	executionTime string // время выполнения
	isSuccess     bool
	isCompleted   bool
	logs          []string
}

func createTask() *Task {
	createTime := time.Now().Format(timeFormat)
	taskId := int(time.Now().Unix())
	if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
		return &Task{
			id:        taskId,
			isSuccess: false,
			logs:      []string{"Some error occurred"},
		}
	}
	return &Task{
		id:         taskId,
		createTime: createTime,
		isSuccess:  true,
	} // передаем таск на выполнение
}

func handleTask(a *Task) {
	createTime, _ := time.Parse(timeFormat, a.createTime)
	duration := createTime.Sub(time.Now())

	if duration <= maxTaskCreatingDelay {
		a.isCompleted = true
		a.logs = append(a.logs, "Task was completed successfully")
	} else {
		a.isCompleted = false
		a.logs = append(a.logs, "Task is working too long, something went wrong")
	}

	a.executionTime = time.Now().Format(executionTimeFormat)

	time.Sleep(taskHandlingDelay)
}

func getTasksChannel() <-chan *Task {
	tasksChan := make(chan *Task, tasksAmount)

	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(tasksAmount)

	for i := 0; i < tasksAmount; i++ {
		go func() {
			defer wg.Done()
			t := createTask()
			mu.Lock()
			handleTask(t)
			mu.Unlock()
			tasksChan <- t
		}()
	}
	wg.Wait()
	close(tasksChan)
	return tasksChan
}

func main() {
	tasksChan := getTasksChannel()

	result := map[int]*Task{}
	errors := make([]error, 0, tasksAmount)

	for r := range tasksChan {
		if !r.isCompleted {
			err := fmt.Errorf(
				"taskId: %d\t "+
					"createTime: %s\t "+
					"errors %s", r.id, r.createTime, strings.Join(r.logs, ", "))
			errors = append(errors, err)
		} else {
			result[r.id] = r
		}
	}

	println("Errors:")
	for r := range errors {
		println(r)
	}

	println("Done tasks:")
	for r := range result {
		println(r)
	}
}
