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
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskResult []byte
}

func main() {
	var (
		result      = make(map[int]Ttype)
		err         []error
		resultMutex sync.Mutex
		errMutex    sync.Mutex
	)
	taskCreturer := func(a chan Ttype, stop chan struct{}) {
		defer close(a)
		for {
			select {
			case <-stop:
				return
			default:
				ft := time.Now().Format(time.RFC3339)
				if time.Now().Nanosecond()%2 > 0 {
					ft = "Some error occured"
				}
				a <- Ttype{cT: ft, id: int(time.Now().UnixNano())}
			}
		}
	}
	stopTaskCr := make(chan struct{})
	superChan := make(chan Ttype, 10)

	go taskCreturer(superChan, stopTaskCr)

	taskWorker := func(a Ttype) Ttype {
		tt, _ := time.Parse(time.RFC3339, a.cT)
		if tt.After(time.Now().Add(-20 * time.Second)) {
			a.taskResult = []byte("task has been successed")
		} else {
			a.taskResult = []byte("something went wrong")
		}
		a.fT = time.Now().Format(time.RFC3339Nano)

		time.Sleep(time.Millisecond * 150)
		return a
	}

	taskSorter := func(t Ttype, sortWg *sync.WaitGroup) {
		if string(t.taskResult[14:]) == "successed" {
			resultMutex.Lock()
			result[t.id] = t
			resultMutex.Unlock()
		} else {
			errMutex.Lock()
			err = append(err, fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskResult))
			errMutex.Unlock()
		}
		sortWg.Done()
	}

	tasksSorted := make(chan struct{})
	countGotten := 0
	go func() {
		var sortWg sync.WaitGroup
		for t := range superChan {
			countGotten++
			t = taskWorker(t)
			sortWg.Add(1)
			go taskSorter(t, &sortWg)
		}
		sortWg.Wait()
		tasksSorted <- struct{}{}
	}()

	time.Sleep(time.Second * 3)
	stopTaskCr <- struct{}{}
	<-tasksSorted

	println("Errors:")
	for i := range err {
		println(i)
	}

	println("Done tasks:")
	for key := range result {
		println(key)
	}
}
