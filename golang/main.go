package main

import (
	"fmt"
	"time"
)

type Ttype struct {
	id         int
	cT         string
	fT         string
	taskRESULT string
}

func main() {
	taskCreator := func(a chan Ttype) {
		go func() {
			for {
				ft := time.Now().Format(time.RFC3339)
				var taskResult string
				if time.Now().Nanosecond()%2 > 0 {
					taskResult = "something went wrong"
				} else {
					taskResult = "task has been succeeded"
				}
				a <- Ttype{cT: ft, id: int(time.Now().Unix()), taskRESULT: taskResult}
			}
		}()
	}

	superChan := make(chan Ttype, 10)
	go taskCreator(superChan)

	taskWorker := func(a Ttype) Ttype {
		tt, _ := time.Parse(time.RFC3339, a.cT)
		if tt.After(time.Now().Add(-20 * time.Second)) {
			a.taskRESULT = "task has been succeeded"
		} else {
			a.taskRESULT = "something went wrong"
		}
		a.fT = time.Now().Format(time.RFC3339Nano)

		time.Sleep(time.Millisecond * 150)
		return a
	}

	doneTasks := make(chan Ttype)
	undoneTasks := make(chan Ttype)

	taskSorter := func(t Ttype) {
		if t.taskRESULT == "task has been succeeded" {
			doneTasks <- t
		} else {
			undoneTasks <- t
		}
	}

	go func() {
		for t := range superChan {
			t = taskWorker(t)
			go taskSorter(t)
		}
	}()

	result := map[int]Ttype{}
	errors := []error{}

	go func() {
		for r := range doneTasks {
			result[r.id] = r
		}
		close(doneTasks)
	}()

	go func() {
		for r := range undoneTasks {
			errors = append(errors,fmt.Errorf("task id %d time %s, error %s", r.id, r.cT, r.taskRESULT))
		}
		close(undoneTasks)
	}()

	time.Sleep(time.Second * 3)

	fmt.Println("Errors:")
	for _, r := range errors {
		fmt.Println(r)
	}

	fmt.Println("Done tasks:")
	for _, r := range result {
		fmt.Println(r)
	}
}
