package main

import (
	"fmt"
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
type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

func taskCreturer() <-chan Ttype {
	tasksChan := make(chan Ttype, 10)

	go func() {
	loop:
		for timeout := time.After(10 * time.Second); ; {
			select {
			case <-timeout:
				break loop
			case <-time.Tick(100 * time.Millisecond):
				ct := time.Now().Format(time.RFC3339)

				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					ct = "Some error occured"
				}

				tasksChan <- Ttype{cT: ct, id: int(time.Now().Unix())} // передаем таск на выполнение
			}

		}
		close(tasksChan)
	}()

	return tasksChan
}

func taskWorker(tasksChan <-chan Ttype) (<-chan Ttype, <-chan error) {
	doneTasksChan := make(chan Ttype, 10)
	undoneTasksChan := make(chan error, 10)

	go func() {

		var taskWorkerWG sync.WaitGroup
		// получение тасков
		for t := range tasksChan {
			taskWorkerWG.Add(1)
			go func(t Ttype) {
				defer taskWorkerWG.Done()
				completed := false

				tt, err := time.Parse(time.RFC3339, t.cT)

				if err == nil && tt.After(time.Now().Add(-20*time.Second)) {
					completed = true
				}

				t.fT = time.Now().Format(time.RFC3339Nano)

				if completed {
					t.taskRESULT = []byte("task has been successed")
					doneTasksChan <- t
				} else {
					t.taskRESULT = []byte("something went wrong")
					undoneTasksChan <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
				}

				time.Sleep(time.Millisecond * 150)
			}(t)
		}
		taskWorkerWG.Wait()
		close(doneTasksChan)
		close(undoneTasksChan)
	}()

	return doneTasksChan, undoneTasksChan
}

func taskPrinter(doneTasksChan <-chan Ttype, undoneTasksChan <-chan error, workDoneChan chan bool) {

	result := map[int]Ttype{}
	errTasks := []error{}

	var taskSorterWG sync.WaitGroup
	sortDoneChan := make(chan bool)

	go func() {
		taskSorterWG.Add(1)
		go func() {
			defer taskSorterWG.Done()
			for r := range doneTasksChan {
				result[r.id] = r
			}
		}()

		taskSorterWG.Add(1)
		go func() {
			defer taskSorterWG.Done()
			for e := range undoneTasksChan {
				errTasks = append(errTasks, e)
			}
		}()

		taskSorterWG.Wait()
		sortDoneChan <- true
	}()

	var taskPrinterWG sync.WaitGroup
	taskPrinterWG.Add(1)
	go func() {
		defer taskPrinterWG.Done()
	loop:
		for {
			select {
			case <-sortDoneChan:
				break loop
			case <-time.Tick(3 * time.Second):
				println("Errors:")
				for _, r := range errTasks {
					println(r.Error())
				}

				println("Done tasks:")
				for r := range result {
					println(r)
				}
				println()
			}
		}
	}()

	taskPrinterWG.Wait()
	workDoneChan <- true
}

func main() {
	tasksChan := taskCreturer()
	doneTasksChan, undoneTasksChan := taskWorker(tasksChan)

	workDoneChan := make(chan bool)
	go taskPrinter(doneTasksChan, undoneTasksChan, workDoneChan)

	<-workDoneChan
}
