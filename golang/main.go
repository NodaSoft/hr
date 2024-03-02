package main

import (
	"fmt"
	"time"
	"sync"
	"bytes"
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
	cT         time.Time // время создания
	fT         time.Time // время выполнения
	taskRESULT []byte
}


func (a *Ttype) task_worker() {
	if a.cT.Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
		a.taskRESULT = []byte("Some error occured")
		return
	}
	if a.cT.After(time.Now().Add(-20 * time.Second)) {
		a.taskRESULT = []byte("task has been successed")
	} else {
		a.taskRESULT = []byte("something went wrong")
	}
	a.fT = time.Now()
	time.Sleep(time.Millisecond * 150)
}


func main() {
	superChan := make(chan *Ttype, 10)

	go func() {
		
		for t, i := time.Now(),0; time.Since(t) < time.Second * 3; i++  {
			t := time.Now()
			superChan <- &Ttype{
				cT: t,
				id: i,
			}
                }
		close(superChan)
	}()

	doneTasks := make(chan *Ttype)
	undoneTasks := make(chan error)

	var wg sync.WaitGroup
	for i:=0; i< 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// получение тасков
			for t := range superChan {
				t.task_worker()
				if bytes.HasSuffix(t.taskRESULT,[]byte("successed")) {
					doneTasks <- t
				} else {
					undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT.Format(time.RFC3339), t.taskRESULT)
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(doneTasks)
		close(undoneTasks)
	}()

	result := map[int]*Ttype{}
	errs := make([]error,0,1000)
	rok, eok := true, true
	var r *Ttype
	var e error
	for rok || eok {
		select {
		case r, rok = <- doneTasks:
			if rok {
				result[r.id] = r
			}
		case e, eok = <- undoneTasks:
			if eok {
				errs = append(errs, e)
			}
		}
	}


	fmt.Println("Errors:")
	for _,e := range errs {
		fmt.Println(e)
	}


	fmt.Println("Done tasks:")
	for r,_ := range result {
		fmt.Println(r)
	}
}
