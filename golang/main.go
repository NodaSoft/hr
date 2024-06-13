package main

import (
	"fmt"
	"log"
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

// taskCreator генерирует таски 10 сек.
func taskCreator(ch chan Ttype) {
	fmt.Println("=========================================================================================")
	log.Println("taskCreator started")
	fmt.Println("=========================================================================================")

	go func() {
		stop := time.After(10 * time.Second)
		for {
			select {
			case <-stop:
				close(ch)
				fmt.Println("=========================================================================================")
				log.Println("taskCreator stopped")
				fmt.Println("=========================================================================================")

				return
			default:
				ft := time.Now().Format(time.RFC3339)
				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					ft = "Some error occured"
				}
				ch <- Ttype{cT: ft, id: int(time.Now().UnixNano())} // передаем таск на выполнение, каждый таск должен иметь уникальный идентификатор
			}
		}
	}()
}

// taskWorker обрабатывает полученные таски
func taskWorker(t Ttype) Ttype {
	tt, _ := time.Parse(time.RFC3339, t.cT)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		t.taskRESULT = []byte("task has been successed")
	} else {
		t.taskRESULT = []byte("something went wrong")
	}
	t.fT = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)

	return t
}

func taskSorter(wg *sync.WaitGroup, t Ttype, doneTaskChan chan Ttype, undoneTaskChan chan error) {
	defer wg.Done()
	if string(t.taskRESULT[14:]) == "successed" {
		doneTaskChan <- t
	} else {
		undoneTaskChan <- fmt.Errorf("task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
	}

}

func main() {

	newTaskChan := make(chan Ttype, 10)

	// генерация тасков
	go taskCreator(newTaskChan)

	doneTaskChan := make(chan Ttype, 10)
	undoneTaskChan := make(chan error, 10)
	defer close(doneTaskChan)
	defer close(undoneTaskChan)

	var wg sync.WaitGroup
	wg.Add(1)

	go func(wg *sync.WaitGroup) {

		// получение и обработка тасков
		for t := range newTaskChan {
			t = taskWorker(t)

			wg.Add(1)
			go taskSorter(wg, t, doneTaskChan, undoneTaskChan)
		}
		defer wg.Done()
	}(&wg)

	// чтение  и вывод результатов
	go func() {
		time.Sleep(time.Second * 3)
		for {
			select {

			case r := <-doneTaskChan:
				fmt.Println("Succeed: task id", r.id)
			case r := <-undoneTaskChan:
				fmt.Println("Error: ", r)
			default:
				t := time.Now()
				time.Sleep(time.Second * 3)
				fmt.Println("")
				fmt.Printf("Paused: %v\n", time.Since(t))
				fmt.Println("")
			}
		}
	}()

	wg.Wait()
}
