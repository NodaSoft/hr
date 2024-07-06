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
				if time.Now().UnixMicro()%2 > 0 { // вот такое условие появления ошибочных тасков
					// если в условии использовать Nanosecond(), то условие никогда не выполняется т.к. последние 2 цифры всегда ноль.
					ft = "Some error occured"
				}
				a <- Ttype{cT: ft, id: int(time.Now().UnixMicro())} // передаем таск на выполнение
				// если использовать Unix() то за секунду создается несколько тасков с одинаковым id, как будто такого не должно быть
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
			result[r.id] = r
		}
		close(doneTasks)
	}()

	go func() {
		for r := range undoneTasks {
			err = append(err, r)
		}
		close(undoneTasks)
	}()

	duration := 10 * time.Second
	messageInterval := 3 * time.Second

	startTime := time.Now()
	nextMessageTime := startTime.Add(messageInterval)

	for {
		currentTime := time.Now()

		if currentTime.Sub(startTime) >= duration {
			break
		}

		if currentTime.After(nextMessageTime) {

			println("Errors:")
			for _, r := range err {
				println(r.Error())
			}

			println("Done tasks:")
			for r := range result {
				println(r)
			}

			nextMessageTime = nextMessageTime.Add(messageInterval)
		}
	}
}
