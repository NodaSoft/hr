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

type Task struct {
	id                       int
	created, carried, result string
}

type CMap struct {
	sync.Mutex
	res map[string][]interface{}
}

func tasksCreator(a chan<- Task, l int) {
	for ; l > 0; l-- {
		ft := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			ft = "Some error occured"
		}
		a <- Task{created: ft, id: int(time.Now().UnixNano())} // передаем таск на выполнение
	}
	close(a)
}

func tasksWorker(s <-chan Task, done chan<- bool, r *CMap) {
	for {
		task, open := <-s

		if open {
			t := handleTask(task)

			if string(t.result[14:]) == "successed" {
				r.Lock()
				r.res["Done tasks"] = append(r.res["Done tasks"], t.id)
				r.Unlock()
			} else {
				r.Lock()
				r.res["Errors"] = append(
					r.res["Errors"],
					fmt.Errorf("task id %d time %s, error %s",
						t.id, t.created, t.result),
				)
				r.Unlock()
			}

		} else {
			done <- true
			return
		}
	}
}

func handleTask(a Task) Task {
	tt, _ := time.Parse(time.RFC3339, a.created)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		a.result = "task has been successed"
	} else {
		a.result = "something went wrong"
	}
	a.carried = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)

	return a
}

func main() {
	tasksLen := 10
	created, done := make(chan Task), make(chan bool)

	result := CMap{
		Mutex: sync.Mutex{},
		res: map[string][]interface{}{
			"Errors":     make([]interface{}, 0),
			"Done tasks": make([]interface{}, 0),
		},
	}

	go tasksCreator(created, tasksLen)
	go tasksWorker(created, done, &result)

	<-done

	for k, v := range result.res {
		fmt.Println(k + ":")
		for _, el := range v {
			fmt.Println(el)
		}
	}
}
