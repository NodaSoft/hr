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

func main() {
	taskCreturer := func(a chan Ttype) {
		go func() {
			for {
				ft := time.Now().Format(time.RFC3339)
				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					ft = "Some error occured"
				}
				a <- Ttype{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
			}
		}()
	}

	superChan := make(chan Ttype, 10)

	go taskCreturer(superChan)

	task_worker := func(a Ttype) Ttype {
		tt, _ := time.Parse(time.RFC3339, a.cT)
		if tt.After(time.Now().Add(-20 * time.Second)) {
			a.taskRESULT = []byte("task has been successed")
		} else {
			a.taskRESULT = []byte("something went wrong")
		}
		a.fT = time.Now().Format(time.RFC3339Nano)

		time.Sleep(time.Millisecond * 150)

		return a
	}

	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	tasksorter := func(t Ttype) {
		if string(t.taskRESULT[14:]) == "successed" {
			doneTasks <- t
		} else {
			undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
		}
	}

	go func() {
		// получение тасков
		for t := range superChan {
			t = task_worker(t)
			go tasksorter(t)
		}
		close(superChan)
	}()

	result := map[int]Ttype{}
	err := []error{}
	go func() {
		for r := range doneTasks {
			go func() {
				result[r.id] = r
			}()
		}
		for r := range undoneTasks {
			go func() {
				err = append(err, r)
			}()
		}
		close(doneTasks)
		close(undoneTasks)
	}()

	time.Sleep(time.Second * 3)

	println("Errors:")
	for r := range err {
		println(r)
	}

	println("Done tasks:")
	for r := range result {
		println(r)
	}
}
