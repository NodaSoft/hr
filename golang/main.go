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

const SuccessfulResult = "successed"

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskResult []byte
}

func (t *Ttype) Work() {
	tt, _ := time.Parse(time.RFC3339, t.cT)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		t.taskResult = []byte("task has been successed")
	} else {
		t.taskResult = []byte("something went wrong")
	}
	t.fT = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	superChan := make(chan Ttype, 10)

	generateTasks(ctx, superChan)

	doneTasks := make(chan Ttype)
	failedTasks := make(chan error)

	receiveTasks(superChan, doneTasks, failedTasks)

	result := map[int]Ttype{}
	var errors []error

	resultMx := sync.Mutex{}
	errorsMx := sync.Mutex{}

	receiveDoneTasks(doneTasks, &resultMx, result)
	receiveErrors(failedTasks, &errorsMx, &errors)

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.NewTicker(time.Second * 3).C:
			showFailedTasks(errors, &errorsMx)
			showDoneTasks(result, &resultMx)
		}
	}
}

func showDoneTasks(result map[int]Ttype, mx *sync.Mutex) {
	mx.Lock()
	println("Done tasks:")
	for id, _ := range result {
		println(id)
	}
	mx.Unlock()
}

func showFailedTasks(errors []error, mx *sync.Mutex) {
	mx.Lock()
	println("Errors:")
	for _, r := range errors {
		println(r.Error())
	}
	mx.Unlock()
}

func receiveDoneTasks(tasks <-chan Ttype, mx *sync.Mutex, result map[int]Ttype) {
	go func() {
		for doneTask := range tasks {
			go func(task Ttype) {
				mx.Lock()
				result[task.id] = task
				mx.Unlock()
			}(doneTask)
		}
	}()
}

func receiveErrors(failedTasks <-chan error, mx *sync.Mutex, errors *[]error) {
	go func() {
		for failedTaskError := range failedTasks {
			go func(error error) {
				mx.Lock()
				*errors = append(*errors, failedTaskError)
				mx.Unlock()
			}(failedTaskError)
		}
	}()
}

func receiveTasks(superChan chan Ttype, doneTasks chan Ttype, undoneTasks chan error) {
	go func() {
		// получение тасков
		for task := range superChan {
			task.Work()
			go func(task Ttype) {
				if strings.Contains(string(task.taskResult), SuccessfulResult) {
					doneTasks <- task
				} else {
					undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", task.id, task.cT, task.taskResult)
				}
			}(task)
		}
	}()
}

func generateTasks(ctx context.Context, superChan chan Ttype) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(superChan)
				return
			default:
				ft := time.Now().Format(time.RFC3339)
				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					ft = "Some error occured"
				}
				superChan <- Ttype{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
			}
		}
	}()
}
