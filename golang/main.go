package main

import (
	"errors"
	"fmt"
	"hr/model"
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
	PrintInterval           = 3 * time.Second
	TaskErrorCreationString = "Task creation error occurred"
	ErrorCondition          = 2
	TaskGenerateDuration    = 10 * time.Second
)

func generateTasks(superChan chan model.Task, wg *sync.WaitGroup) {
	startTime := time.Now()
	elapsed := time.Duration(0)

	for elapsed < TaskGenerateDuration {
		currentTime := time.Now()
		tt, err := time.Parse(time.RFC3339, currentTime.Format(time.RFC3339))

		if time.Now().Nanosecond()%ErrorCondition > 0 && err == nil {
			err = errors.New(TaskErrorCreationString)
		}
		superChan <- model.Task{
			CreateTime: tt,
			Id:         int(currentTime.Unix()),
			Error:      err,
		}
		elapsed = time.Since(startTime)
	}
	wg.Done()
}

func parseTasks(superChan chan model.Task, successTasks chan<- model.Task, failTasks chan error, wg *sync.WaitGroup) {
	for {
		task, ok := <-superChan
		if !ok {
			break
		} else {
			wg.Add(1)
			go func(val *model.Task) {
				val.Check()
				val.Sort(successTasks, failTasks)
				wg.Done()
			}(&task)
		}
	}
	wg.Done()
}

func parseSuccess(mu *sync.Mutex, successTasks <-chan model.Task, result map[int]model.Task, wg *sync.WaitGroup) {
	for r := range successTasks {
		wg.Add(1)
		go func(r model.Task) {
			mu.Lock()
			result[r.Id] = r
			mu.Unlock()
			wg.Done()
		}(r)
	}
	wg.Done()
}

func parseFailure(mu *sync.Mutex, failTasks <-chan error, err *[]error, wg *sync.WaitGroup) {
	for e := range failTasks {
		wg.Add(1)
		go func(e error) {
			mu.Lock()
			*err = append(*err, e)
			mu.Unlock()
			wg.Done()
		}(e)
	}
	wg.Done()
}

func printer(result map[int]model.Task, err *[]error, wg *sync.WaitGroup) {
	ticker := time.NewTicker(PrintInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			wg.Add(1)
			go func() {
				for e := range *err {
					fmt.Printf("error:%v\n", e)
				}

				for _, r := range result {
					fmt.Println(r)
				}
				wg.Done()
			}()
		}
		wg.Done()
	}
}

func main() {
	superChan := make(chan model.Task, 10)
	successTasks := make(chan model.Task)
	failTasks := make(chan error)
	defer close(successTasks)
	defer close(failTasks)
	defer close(superChan)
	var wg sync.WaitGroup
	var mu sync.Mutex
	result := map[int]model.Task{}
	var err []error

	wg.Add(1)
	go generateTasks(superChan, &wg)

	wg.Add(1)
	go parseTasks(superChan, successTasks, failTasks, &wg)

	wg.Add(1)
	go parseSuccess(&mu, successTasks, result, &wg)

	wg.Add(1)
	go parseFailure(&mu, failTasks, &err, &wg)

	wg.Add(1)
	go printer(result, &err, &wg)

	wg.Wait()
}
