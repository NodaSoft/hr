package main

import (
	"fmt"
	"sync"
	"time"
)

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// После обработки тасков в течении 3 секунд приложение должно выводить накопленные к этому моменту успешные таски и отдельно ошибки обработки тасков.

// ЗАДАНИЕ: сделать из плохого кода хороший и рабочий - as best as you can.
// Важно сохранить логику появления ошибочных тасков.
// Важно оставить асинхронные генерацию и обработку тасков.
// Сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через pull-request в github
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id          int
	cT          string // время создания
	fT          string // время выполнения
	taskSuccess bool
	taskRESULT  string
}

func main() {
	var wg sync.WaitGroup
	stopTheWorld := make(chan bool)
	taskCreturer := func(a chan Ttype) {
		go func() {
			for {
				ft := time.Now().Format(time.RFC3339)
				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					ft = "Some error occured"
				}
				select {
				case <-stopTheWorld:
					close(a)
					return
				default:
					//Вот тут я не знаю, айдишник по идее лучше поменять, так как много тасок друг друга переписывают
					a <- Ttype{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
				}
			}
		}()
	}

	superChan := make(chan Ttype, 10)
	go taskCreturer(superChan)

	task_worker := func(a Ttype) Ttype {
		tt, _ := time.Parse(time.RFC3339, a.cT)
		if tt.After(time.Now().Add(-20 * time.Second)) {
			a.taskSuccess = true
			a.taskRESULT = "task has been successed"
		} else {
			a.taskSuccess = false
			a.taskRESULT = "something went wrong"
		}
		a.fT = time.Now().Format(time.RFC3339Nano)

		return a
	}

	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	tasksorter := func(t Ttype) {
		if t.taskSuccess {
			doneTasks <- t
		} else {
			undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
		}
	}

	go func() {
		// получение тасков
		for t := range superChan {
			valWorked := task_worker(t)
			tasksorter(valWorked)
		}
		close(doneTasks)
		close(undoneTasks)
	}()

	result := map[int]Ttype{}
	err := []error{}
	wg.Add(1)
	go func() {
		for r := range doneTasks {
			result[r.id] = r
		}
		for r := range undoneTasks {
			err = append(err, r)
		}
		wg.Done()
	}()

	time.Sleep(time.Second * 3)
	stopTheWorld <- true
	wg.Wait()

	println("Errors:")
	for r := range err {
		println(r)
	}

	println("Done tasks:")
	for r := range result {
		println(r)
	}
}
