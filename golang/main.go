package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

const (
	statusSuccess = 1
	statusFailed  = 2
)

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id         int
	cT         time.Time // время создания
	fT         time.Time // время выполнения
	taskRESULT int
	err        error
}

func main() {
	superChan, cancel := CreateTasks()
	doneTasks, undoneTasks, done := TasksResult(superChan)

	result := make(map[int]Ttype)
	errs := make([]error, 0)

	go func() {
		for {
			select {
			case t := <-doneTasks:
				result[t.id] = t
			case err := <-undoneTasks:
				errs = append(errs, err)
			case <-done:
				return
			}
		}
	}()

	time.Sleep(time.Second * 25)

	cancel()
	close(done)

	println("Errors:")
	for _, err := range errs {
		println(err.Error())
	}

	println("Done tasks:")
	for r := range result {
		println(r)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
}

func CreateTasks() (<-chan Ttype, func()) {
	ch := make(chan Ttype)
	done := make(chan struct{})
	cancel := func() {
		close(done)
	}
	go func() {
		for {
			select {
			case <-done:
				break
			default:
				t := time.Now()
				var err error
				if (t.Nanosecond()/1000)%2 > 0 { // вот такое условие появления ошибочных тасков
					err = errors.New("some error occured")
				}
				ch <- Ttype{id: int(t.Unix()), cT: t, err: err} // передаем таск на выполнение
			}
		}
		close(ch)
	}()
	return ch, cancel
}

func TasksResult(superChan <-chan Ttype) (chan Ttype, chan error, chan struct{}) {
	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)
	done := make(chan struct{})

	go func() {
		for t := range superChan {
			t = SetTasksResult(t)
			go func(t Ttype) {
				select {
				case <-done:
					return
				default:
					if t.taskRESULT == statusSuccess {
						doneTasks <- t
					} else {
						undoneTasks <- fmt.Errorf("task id %d time %s, error %s", t.id, t.cT, t.err)
					}
				}
			}(t)
		}
		close(doneTasks)
		close(undoneTasks)
	}()

	return doneTasks, undoneTasks, done
}

func SetTasksResult(a Ttype) Ttype {
	if a.err == nil {
		a.taskRESULT = statusSuccess
	} else {
		a.taskRESULT = statusFailed
	}
	a.fT = time.Now()

	time.Sleep(time.Millisecond * 150)

	return a
}
