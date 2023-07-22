package main

import (
	"bytes"
	"context"
	"fmt"
	"strings"
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
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskResult []byte
}

var (
	badStatus, successStatus []byte = []byte("something went wrong"), []byte("task has been successes")
)

const (
	errSome string = "some error occurred"
)

var (
	taskCreature = func(ctx context.Context, a chan Ttype, wg *sync.WaitGroup) {
		go func() {
			defer wg.Done()
			for {
				deadLine, _ := ctx.Deadline()
				if time.Now().After(deadLine) {
					break
				}
				ft := time.Now().Format(time.RFC3339)
				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					ft = "some error occurred"
				}
				a <- Ttype{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
			}
		}()
	}

	taskWorker = func(a Ttype) Ttype {
		if a.cT == errSome {
			a.taskResult = badStatus
		} else {
			tt, err := time.Parse(time.RFC3339, a.cT)
			if err != nil {
				a.taskResult = badStatus
			} else if tt.After(time.Now().Add(-20 * time.Second)) {
				a.taskResult = successStatus
			} else {
				a.taskResult = badStatus
			}
		}

		a.fT = time.Now().Format(time.RFC3339Nano)

		time.Sleep(time.Millisecond * 150)

		return a
	}
)

func main() {
	var (
		superChan   = make(chan Ttype, 10)
		doneTasks   = make(chan Ttype)
		undoneTasks = make(chan error)
		wg          = sync.WaitGroup{}
	)
	wg.Add(1)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	go taskCreature(ctx, superChan, &wg)

	taskSorter := func(t Ttype) {
		if bytes.Equal(t.taskResult, successStatus) {
			doneTasks <- t
		} else {
			undoneTasks <- fmt.Errorf("task id %d time %s, error %s", t.id, t.cT, string(t.taskResult))
		}
	}

	go func() {
		// получение тасков
		for t := range superChan {
			t = taskWorker(t)
			go taskSorter(t)
		}
		close(superChan)
	}()

	//------Handling Data--------------------

	result, errors := dataHandling(doneTasks, undoneTasks)
	wg.Wait()

	//------Result Printing--------------------

	resultPrinting(result, errors)

}

func dataHandling(doneTasks chan Ttype, undoneTasks chan error) (map[int]Ttype, []error) {
	var (
		result = make(map[int]Ttype, len(doneTasks))
		errors []error
		mutex  sync.RWMutex
	)
	go func() {
		for tsk := range doneTasks {
			go func(tsk Ttype) {
				mutex.Lock()
				result[tsk.id] = tsk
				mutex.Unlock()
			}(tsk)
		}
		close(doneTasks)
	}()

	go func() {
		for e := range undoneTasks {
			go func(e error) {
				errors = append(errors, e)
			}(e)
		}
		close(undoneTasks)
	}()
	return result, errors
}

func resultPrinting(result map[int]Ttype, errors []error) {
	var (
		sbErr, sbDone = strings.Builder{}, strings.Builder{}
		wgPrinting    = sync.WaitGroup{}
	)
	wgPrinting.Add(2)
	go func() {
		defer wgPrinting.Done()
		sbErr.WriteString("Errors:\n")
		for _, err := range errors {
			sbErr.WriteString(fmt.Sprintf("%s\n", err.Error()))
		}
	}()

	go func() {
		defer wgPrinting.Done()
		sbDone.WriteString("Done tasks:\n")
		for _, res := range result {
			sbDone.WriteString(fmt.Sprintf("%d\n", res.id))
		}
	}()
	wgPrinting.Wait()

	println(sbDone.String())
	println(sbErr.String())
}
