package main

import (
	"fmt"
	"sync"
	"time"
)

type Ttype struct {
	id         int
	cT         time.Time
	fT         time.Time
	taskRESULT []byte
}

func main() {
	var wg sync.WaitGroup
	taskCreturer := func(a chan Ttype) {
		for {
			ft := time.Now()
			if time.Now().Nanosecond()%2 > 0 { 
				ft = time.Time{}
			}
			a <- Ttype{cT: ft, id: int(time.Now().Unix())}
		}
	}

	superChan := make(chan Ttype, 10)
	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	go taskCreturer(superChan)

	go func() {
		// receive tasks
		for t := range superChan {
			wg.Add(1)
			go func(t Ttype) {
				defer wg.Done()
				if t.cT.IsZero() {
					undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, "Some error occured", "")
				} else if time.Since(t.cT) < 20*time.Second {
					t.taskRESULT = []byte("task has been successed")
					t.fT = time.Now()
					doneTasks <- t
				} else {
					undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, "something went wrong")
				}
				time.Sleep(time.Millisecond * 150)
			}(t)
		}
		close(superChan)
	}()

	var result = map[int]Ttype{}
	var errs = []error{}
	go func() {
		for r := range doneTasks {
			result[r.id] = r
		}
		for r := range undoneTasks {
			errs = append(errs, r)
		}
		close(doneTasks)
		close(undoneTasks)
	}()

	wg.Wait()

	fmt.Println("Errors:")
	for _, r := range errs {
		fmt.Println(r)
	}

	fmt.Println("Done tasks:")
	for _, r := range result {
		fmt.Println(r)
	}
}
