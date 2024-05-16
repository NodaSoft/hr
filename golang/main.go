package main

import (
	"fmt"
	"time"
	"sync"
)

type Ttype struct {
	id         int
	cT         string
	fT         string
	taskRESULT []byte
}

func main() {
	superChan := make(chan Ttype, 10)
	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)
	var wg sync.WaitGroup

	taskCreturer := func(a chan Ttype) {
		for {
			ft := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 {
				ft = "Some error occured"
			}
			a <- Ttype{cT: ft, id: int(time.Now().Unix())}
			time.Sleep(100 * time.Millisecond)
		}
	}

	task_worker := func(a Ttype) Ttype {
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

	tasksorter := func(t Ttype) {
		if string(t.taskRESULT[14:]) == "successed" {
			doneTasks <- t
		} else {
			undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
		}
	}

	go taskCreturer(superChan)

	go func() {
		for t := range superChan {
			wg.Add(1)
			go func(task Ttype) {
				defer wg.Done()
				task = task_worker(task)
				tasksorter(task)
			}(t)
		}
	}()

	go func() {
		wg.Wait()
		close(doneTasks)
		close(undoneTasks)
	}()

	go func() {
		time.Sleep(time.Second * 3)
		close(superChan)
	}()

	result := map[int]Ttype{}
	err := []error{}

	for {
		select {
		case r, ok := <-doneTasks:
			if ok {
				result[r.id] = r
			}
		case r, ok := <-undoneTasks:
			if ok {
				err = append(err, r)
			}
		case <-time.After(time.Second * 3):
			goto DONE
		}
	}
