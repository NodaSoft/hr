package main

import (
	"fmt"
	"sync"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	ID         int
	CT         string // время создания
	FT         string // время выполнения
	taskResult []byte
}
type Res struct {
	sync.Mutex
	result map[int]Ttype
}
type ErrorSt struct {
	sync.Mutex
	err []error
}

func main() {
	superChan := make(chan Ttype, 10)

	go taskCreturer(superChan)

	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)
	go func() {
		// получение тасков
		for t := range superChan {
			task := taskWorker(t)
			go tasksorter(task, doneTasks, undoneTasks)
		}
		close(superChan)
		//		close(superChan)
	}()
	var res Res
	var errorSt ErrorSt

	res.result = map[int]Ttype{}
	errorSt.err = []error{}
	go func() {
		for r := range doneTasks {
			go func(r Ttype) {
				res.Lock()
				res.result[r.ID] = r
				res.Unlock()
			}(r)
		}
		for er := range undoneTasks {
			go func(er error) {
				errorSt.Lock()
				errorSt.err = append(errorSt.err, er)
				errorSt.Unlock()
			}(er)
		}

		close(doneTasks)
		close(undoneTasks)
	}()

	time.Sleep(time.Second * 3)
	println("Errors:")
	for r := range errorSt.err {
		println(r)
	}

	println("Done tasks:")
	for r := range res.result {
		println(r)
	}

}
func taskCreturer(a chan Ttype) {
	for {
		FT := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			FT = "Some error occured"
		}
		a <- Ttype{CT: FT, ID: int(time.Now().Unix())} // передаем таск на выполнение
	}
}
func taskWorker(a Ttype) Ttype {
	tt, _ := time.Parse(time.RFC3339, a.CT)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		a.taskResult = []byte("task has been successed")
	} else {
		a.taskResult = []byte("something went wrong")
	}
	a.FT = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)

	return a
}
func tasksorter(t Ttype, doneTasks chan Ttype, undoneTasks chan error) {
	if string(t.taskResult[14:]) == "successed" {
		doneTasks <- t
	} else {
		undoneTasks <- fmt.Errorf("task ID %d time %s, error %s", t.ID, t.CT, t.taskResult)
	}
}
