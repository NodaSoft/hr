package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// Ttype represents a task with creation and finish times
type Ttype struct {
	id         uint64
	cT         string // creation time
	fT         string // finish time
	taskRESULT []byte
}

type TtypeWithError struct {
	t   Ttype
	err error
}

var someError error = errors.New("Some error occurred")

func main() {
	// taskCreator := func(a chan Ttype, ctxCreator context.Context) {
	// 	var nextIdx uint64 = 0
	// 	for {
	// 		ft := time.Now().Format(time.RFC3339)
	// 		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
	// 			ft = "Some error occured"
	// 		}
	// 		nextIdx++
	// 		select {
	// 		case <-ctxCreator.Done():
	// 			close(a)
	// 			return
	// 		case a <- Ttype{cT: ft, id: nextIdx}: // передаем таск на выполнение
	// 			//
	// 		}
	// 	}
	// }
	taskCreator := func(a chan Ttype, ctxCreator context.Context) {
		var nextIdx uint64 = 0
		for {
			ft := time.Now().Format(time.RFC3339)
			// Условие появления ошибочных тасков с использованием рандома
			if rand.Intn(2) == 0 {
				ft = "Some error occurred"
			}
			atomic.AddUint64(&nextIdx, 1)
			select {
			case <-ctxCreator.Done():
				close(a)
				return
			case a <- Ttype{cT: ft, id: atomic.LoadUint64(&nextIdx)}: // передаем таск на выполнение
				//
			}
		}
	}

	taskWorker := func(a Ttype) (*Ttype, error) {
		taskFinished, err := time.Parse(time.RFC3339, a.cT)
		if err == nil {
			if taskFinished.After(time.Now().Add(-20 * time.Second)) {
				a.taskRESULT = []byte("task has been successed")
			} else {
				a.taskRESULT = []byte("something went wrong")
				err = someError
			}
		} else {
			a.taskRESULT = []byte("parsing error")
			err = someError
		}
		a.fT = time.Now().Format(time.RFC3339Nano)
		time.Sleep(time.Millisecond * 150) // simulate processing time
		return &a, err
	}

	superChan := make(chan Ttype, 1000)
	ctxCreator, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	defer cancel()
	go taskCreator(superChan, ctxCreator)

	doneTasks := make(chan *Ttype)
	undoneTasks := make(chan error)

	go func() {
		var wg sync.WaitGroup
		for task := range superChan {
			wg.Add(1)
			go func(task Ttype) {
				defer wg.Done()
				t, err := taskWorker(task)
				if err == nil {
					doneTasks <- t
				} else {
					undoneTasks <- fmt.Errorf("task id %d time %s, error %s: %w", t.id, t.cT, string(t.taskRESULT), err)
				}
			}(task)
		}
		wg.Wait()
		close(doneTasks)
		close(undoneTasks)
	}()

	result := make(map[uint64]*Ttype)
	errors := []error{}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for res := range doneTasks {
			fmt.Println("pep")
			fmt.Println(res.id)
			result[res.id] = res
		}
	}()

	go func() {
		defer wg.Done()
		for err := range undoneTasks {
			errors = append(errors, err)
		}
	}()

	wg.Wait()

	fmt.Println("Errors:")
	for i, err := range errors {
		fmt.Printf("error-%d: %v\n", i, err)
	}

	fmt.Println("Done tasks:")
	for id, res := range result {
		fmt.Printf("res-%d: %s\n", id, res.taskRESULT)
	}
}
