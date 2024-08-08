package main

import (
	"context"
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

// A Ttype represents a meaninglessness of our life
type Task struct {
	ID         int
	CreateTime string // время создания
	RunTime    string // время выполнения
	Result     []byte
}

func StartTasksGeneration(timeout time.Duration, c chan Task, wg *sync.WaitGroup) {
	defer close(c)
	defer wg.Done()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			task := Task{
				CreateTime: time.Now().Format(time.RFC3339),
				ID:         int(time.Now().Unix()),
			}

			if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				task.CreateTime = "Some error occured"
			}
			// передаем таск на выполнение
			c <- task
		}
	}
}

func StartTaskHandler(genTasksChan chan Task, doneTasksChan chan Task, undoneTasksChan chan error, wg *sync.WaitGroup) {
	defer close(doneTasksChan)
	defer close(undoneTasksChan)

	defer wg.Done()

	for task := range genTasksChan {
		task = composeTaskResult(task)

		time.Sleep(time.Millisecond * 150)

		taskSortInChan(doneTasksChan, undoneTasksChan, task)
	}
}

func taskSortInChan(doneTasksChan chan Task, undoneTasksChan chan error, t Task) {
	if strings.Contains(string(t.Result), "successed") {
		doneTasksChan <- t
	} else {
		undoneTasksChan <- fmt.Errorf("task id %d time %s, error %s", t.ID, t.CreateTime, t.Result)
	}
}

func composeTaskResult(t Task) Task {
	tt, _ := time.Parse(time.RFC3339, t.CreateTime)

	t.Result = []byte("task has been successed")
	if !tt.After(time.Now().Add(-20 * time.Second)) {
		t.Result = []byte("something went wrong")
	}

	t.RunTime = time.Now().Format(time.RFC3339Nano)
	return t
}

func StartTaskReader(doneTasksChan chan Task, undoneTasksChan chan error, wg *sync.WaitGroup) {
	taskResults := map[int]Task{}
	taskErrors := []error{}

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	defer wg.Done()

	for {
		select {
		case t, ok := <-doneTasksChan:
			if !ok {
				return
			}
			taskResults[t.ID] = t
		case t, ok := <-undoneTasksChan:
			if !ok {
				return
			}
			taskErrors = append(taskErrors, t)
		case <-ticker.C:
			go WriteTasksResult(taskResults, taskErrors)
		}
	}
}

func WriteTasksResult(taskResults map[int]Task, taskErrors []error) {
	println("Errors:")
	for t := range taskErrors {
		fmt.Println(t)
	}

	println("Done tasks:")
	for id, task := range taskResults {
		fmt.Printf("task id %d, create time %s, run time %s, result %s\n", id, task.CreateTime, task.RunTime, task.Result)
	}
}

func main() {
	genTasksChan := make(chan Task, 10)
	doneTasksChan := make(chan Task)
	undoneTasksChan := make(chan error)

	var wg sync.WaitGroup

	wg.Add(3)
	// Запуск генератора задач в фоне с указанием времени работы в указанное время, а именно 10 сек
	go StartTasksGeneration(10*time.Second, genTasksChan, &wg)
	// Запуск обработчика задач в фоне
	go StartTaskHandler(genTasksChan, doneTasksChan, undoneTasksChan, &wg)
	// Запуск ридера задач
	go StartTaskReader(doneTasksChan, undoneTasksChan, &wg)
	wg.Wait()
}
