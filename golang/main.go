package main

import (
	"errors"
	"fmt"
	"time"
)

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// Должно выводить успешные таски и ошибки обработки тасков по мере выполнения.

// Важно сохранить логику появления ошибочных тасков.
// Сделать правильную мультипоточность обработки заданий.

var (
	ErrTaskExecute = errors.New("Some error occured")
)

type Ttype struct {
	id         int64
	cT         time.Time // время создания
	fT         time.Time // время выполнения
	taskRESULT string
	taskERROR  error
}

func main() {
	taskCreturer := func(a chan *Ttype) {
		go func() {
			for {
				ct := time.Now()
				task := &Ttype{
					id: time.Now().Unix(),
					cT: ct,
				}

				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					task.taskERROR = ErrTaskExecute
				}
				a <- task
			}
		}()
	}

	superChan := make(chan *Ttype, 10)

	go taskCreturer(superChan)

	task_worker := func(a *Ttype) *Ttype {
		if a.taskERROR == nil {
			a.taskRESULT = "successed"
		}
		a.fT = time.Now()

		time.Sleep(time.Millisecond * 150)

		return a
	}

	doneTasks := make(chan *Ttype)
	undoneTasks := make(chan error)

	tasksorter := func(t *Ttype) {
		if t.taskERROR == nil {
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

	result := make([]Ttype, len(doneTasks))
	err := make([]error, len(undoneTasks))

	go func() {
		for r := range doneTasks {
			donnedTask := r

			go func() {
				result = append(result, *donnedTask)
			}()
		}
	}()

	go func() {

		for r := range undoneTasks {
			undonnedTask := r

			go func() {
				err = append(err, undonnedTask)
			}()
		}
	}()

	sleepTime := time.Second * 1
	time.Sleep(sleepTime)

	close(doneTasks)
	close(undoneTasks)

	println("Errors:")
	for r := range err {
		println(err[r].Error())
	}

	println("Done tasks:")
	for r := range result {
		println(result[r].id)
	}
}
