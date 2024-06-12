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
type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

func taskCreator(a chan Ttype) {
	for {
		creationTime := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			creationTime = "Some error occured"
		}
		a <- Ttype{cT: creationTime, id: int(time.Now().Unix())} // передаем таск на выполнение
	}
}

func taskWorker(a Ttype) Ttype {
	tt, err := time.Parse(time.RFC3339, a.cT)
	if err != nil {
		a.taskRESULT = []byte("something went wrong")
	}
	if tt.After(time.Now().Add(-20 * time.Second)) {
		a.taskRESULT = []byte("task has been successed")
	} else {
		a.taskRESULT = []byte("something went wrong")
	}
	a.fT = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)

	return a
}

func taskSorter(t Ttype, doneTasks chan Ttype, undoneTasks chan error) {
	if string(t.taskRESULT[14:]) == "successed" {
		doneTasks <- t
	} else {
		undoneTasks <- fmt.Errorf("task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
	}
}

func main() {
	superChan := make(chan Ttype, 10)

	go taskCreator(superChan)

	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	go func() {
		// получение тасков
		for t := range superChan {
			t = taskWorker(t)
			go taskSorter(t, doneTasks, undoneTasks)
		}
	}()

	result := map[int]Ttype{}
	err := []error{}
	go func() {
		for {
			select {
			case r := <-doneTasks:
				result[r.id] = r
			case r := <-undoneTasks:
				err = append(err, r)
			}
		}
	}()

	quit := make(chan bool)
	go func() {
		<-time.After(10 * time.Second)
		quit <- true
	}()

	for {
		select {
		case <-quit:
			return
		case <-time.Tick(3 * time.Second):
			println("Errors:")
			for r := range err {
				println(r)
			}

			println("Done tasks:")
			for r := range result {
				println(r)
			}
		}
	}
}
