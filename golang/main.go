package main

import (
	"bytes"
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

var ErrorResultBytes = []byte("Some error occured")

func main() {
	taskCreturer := func(a chan Ttype) {
		for {
			ct := time.Now().Format(time.RFC3339)
			task := Ttype{cT: ct, id: int(time.Now().UnixNano())}
			if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				task.taskRESULT = ErrorResultBytes
			}
			a <- task // передаем таск на выполнение
		}
	}

	superChan := make(chan Ttype, 10)

	go taskCreturer(superChan)

	task_worker := func(a Ttype) Ttype {
		tt, err := time.Parse(time.RFC3339, a.cT)
		if err != nil {
			a.taskRESULT = []byte(fmt.Sprintf("time parse error [%v]", err))
			return a
		}

		if bytes.Equal(a.taskRESULT, ErrorResultBytes) {
			return a
		}

		if tt.Before(time.Now().Add(-20 * time.Second)) {
			a.taskRESULT = []byte("task outdated")
			return a
		}

		time.Sleep(time.Millisecond * 150)

		a.taskRESULT = []byte("task has been successed")
		a.fT = time.Now().Format(time.RFC3339Nano)

		return a
	}

	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	tasksorter := func(t Ttype) {
		if len(t.taskRESULT) > 14 && string(t.taskRESULT[14:]) == "successed" {
			doneTasks <- t
		} else {
			undoneTasks <- fmt.Errorf("task id [%d] time [%s], error [%s]", t.id, t.cT, t.taskRESULT)
		}
	}

	go func() {
		// получение тасков
		for t := range superChan {
			t = task_worker(t)
			go tasksorter(t)
		}
	}()

	result := map[int]Ttype{}
	err := []error{}

	go func() {
		for r := range doneTasks {
			result[r.id] = r
		}
	}()

	go func() {
		for r := range undoneTasks {
			err = append(err, r)
		}
	}()

	time.Sleep(time.Second * 5)

	println("Errors:")
	for _, r := range err {
		fmt.Printf("%v\n", r)
	}

	println("Done tasks:")
	for _, r := range result {
		fmt.Printf("%d %s\n", r.id, string(r.taskRESULT))
	}
}
