package main

import (
	"fmt"
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
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

func main() {

	result := map[int]Ttype{}
	err := []error{}

	superChan := make(chan Ttype, 10)
	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	go newTask(superChan)

	go receiveTasks(superChan, doneTasks, undoneTasks)

	go calculateResult(doneTasks, undoneTasks, result, err)

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

func newTask(a chan Ttype) {
	// go func() {
	for {
		ft := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			ft = "Some error occured"
		}
		a <- Ttype{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
	}
	// }()
}

func newWorker(a Ttype) Ttype {
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

func taskSort(t Ttype, doneTasks chan Ttype, undoneTasks chan error) {
	if string(t.taskRESULT[14:]) == "successed" {
		doneTasks <- t
	} else {
		undoneTasks <- fmt.Errorf("task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
	}
}

func receiveTasks(superChan chan Ttype, doneTasks chan Ttype, undoneTasks chan error) {
	for t := range superChan {
		t = newWorker(t)
		go taskSort(t, doneTasks, undoneTasks)
	}
	close(superChan)
}

func calculateResult(doneTasks chan Ttype, undoneTasks chan error, result map[int]Ttype, err []error) {
	for r := range doneTasks {
		r := r
		go func() {
			result[r.id] = r
		}()
	}
	for r := range undoneTasks {
		r := r
		go func() {
			err = append(err, r)
		}()
	}
	close(doneTasks)
	close(undoneTasks)
}
