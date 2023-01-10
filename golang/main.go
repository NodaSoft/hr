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

func task_creturer(superChan chan<- Ttype) {
	for {
		ft := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			ft = "Some error occured"
		}
		superChan <- Ttype{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
	}
}

func task_worker(superChan <-chan Ttype, sort_ch chan<- Ttype) {
	for a := range superChan {
		tt, _ := time.Parse(time.RFC3339, a.cT)
		if tt.After(time.Now().Add(-20 * time.Second)) {
			a.taskRESULT = []byte("task has been successed")
		} else {
			a.taskRESULT = []byte("something went wrong")
		}
		a.fT = time.Now().Format(time.RFC3339Nano)

		time.Sleep(time.Millisecond * 150)

		sort_ch <- a
	}
	close(sort_ch)
}

func task_sorter(
	sort_ch <-chan Ttype,
	doneTasks chan<- Ttype,
	undoneTasks chan<- error,
) {
	for t := range sort_ch {
		if string(t.taskRESULT[14:]) == "successed" {
			doneTasks <- t
		} else {
			undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
		}
	}
	close(doneTasks)
	close(undoneTasks)
}

func main() {
	superChan := make(chan Ttype, 10)
	sort_ch := make(chan Ttype, 10)
	doneTasks_ch := make(chan Ttype)
	undoneTasks_ch := make(chan error)

	go task_creturer(superChan)
	go task_worker(superChan, sort_ch)
	go task_sorter(sort_ch, doneTasks_ch, undoneTasks_ch)

	result := map[int]Ttype{}
	go func() {
		for t := range doneTasks_ch {
			result[t.id] = t
		}
	}()

	err := []error{}
	go func() {
		for t := range undoneTasks_ch {
			err = append(err, t)
		}
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
