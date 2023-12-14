package main

import (
	"errors"
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
	id           int64
	cT           time.Time // время создания
	fT           time.Time // время выполнения
	cErr         error
	taskResultOk bool
}

type TResult struct {
	result map[int64]Ttype
	m      sync.Mutex
}

type TError struct {
	errs []error
	m    sync.Mutex
}

const (
	taskCount = 10
	workers   = 10
)

func taskHandler(allTasks <-chan Ttype, doneTasks chan<- Ttype, undoneTasks chan<- error, endSignal chan struct{}, wg *sync.WaitGroup) {
	for {
		select {
		case t, ok := <-allTasks:
			if ok {
				t.taskResultOk = t.cErr == nil && t.cT.After(time.Now().Add(-20*time.Second))

				time.Sleep(time.Millisecond * 150) // эмуляция времени обработки задачи
				t.fT = time.Now()

				if t.taskResultOk {
					doneTasks <- t
				} else {
					undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT.Format(time.RFC3339), "something went wrong")
				}
			}
		case <-endSignal:
			wg.Done()
			return
		}
	}
}

func taskResults(doneTasks <-chan Ttype, undoneTasks <-chan error, r *TResult, es *TError, endSignal <-chan struct{}, wg *sync.WaitGroup) {
	for {
		select {
		case done, ok := <-doneTasks:
			if ok {
				r.m.Lock()
				(*r).result[done.id] = done
				r.m.Unlock()
			}
		case unDone, ok := <-undoneTasks:
			if ok {
				es.m.Lock()
				es.errs = append(es.errs, unDone)
				es.m.Unlock()
			}
		case <-endSignal:
			wg.Done()
			return
		}
	}
}

func main() {
	var (
		receiveW, resultW        sync.WaitGroup
		stopReceive, stopResults = make(chan struct{}), make(chan struct{})
		superChan                = make(chan Ttype)
		doneTasks                = make(chan Ttype)
		undoneTasks              = make(chan error)
		result                   TResult
		errs                     TError
	)
	result.result = map[int64]Ttype{}
	receiveW.Add(workers)
	resultW.Add(workers)

	for i := 0; i < workers; i++ {
		go taskHandler(superChan, doneTasks, undoneTasks, stopReceive, &receiveW)
		go taskResults(doneTasks, undoneTasks, &result, &errs, stopResults, &resultW)
	}

	for i := 0; i < taskCount; i++ {
		var (
			ct = time.Now()
			t  = Ttype{id: ct.UnixNano(), cT: ct}
		)
		if ct.Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			t.cErr = errors.New("Some error occured")
		}
		superChan <- t // передаем таск на выполнение
	}

	close(superChan)
	close(stopReceive)

	receiveW.Wait() // ожидание получения всех задач

	close(doneTasks)
	close(undoneTasks)
	close(stopResults)

	resultW.Wait() // ожидание обработки всех задач
	//time.Sleep(time.Second * 3)

	println("Errors:")
	for _, r := range errs.errs {
		println(r.Error())
	}

	println("Done tasks:")
	for r := range result.result {
		println(r)
	}
}
