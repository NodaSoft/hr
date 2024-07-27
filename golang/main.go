package main

import (
	"fmt"
	"sync"
	"time"
)

type Ttype struct {
	id         int
	cT         string
	fT         string
	taskRESULT []byte
}

func main() {
	taskCreature := func(a chan Ttype) {
		go func() {
			taskID := 0
			startTime := time.Now()
			for time.Since(startTime) < 10*time.Second {
				taskID++
				ft := time.Now().Format(time.RFC3339)
				if time.Now().Nanosecond()%2 > 0 {
					ft = "Some error occurred"
				}
				a <- Ttype{id: taskID, cT: ft}
				time.Sleep(time.Second)
			}
			close(a)
		}()
	}

	superChan := make(chan Ttype, 10)
	doneTasks := make(chan Ttype, 10)
	undoneTasks := make(chan error, 10)

	var wg sync.WaitGroup
	taskCreature(superChan)

	taskWorker := func(a Ttype) Ttype {
		tt, err := time.Parse(time.RFC3339, a.cT)
		if err == nil && tt.After(time.Now().Add(-20*time.Second)) {
			a.taskRESULT = []byte("task has been successes")
		} else {
			a.taskRESULT = []byte("something went wrong")
		}
		a.fT = time.Now().Format(time.RFC3339Nano)
		time.Sleep(time.Millisecond * 150)
		return a
	}

	taskSorter := func(t Ttype) {
		if string(t.taskRESULT) == "task has been successes" {
			doneTasks <- t
		} else {
			undoneTasks <- fmt.Errorf("task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
		}
	}

	go func() {
		for t := range superChan {
			wg.Add(1)
			go func(task Ttype) {
				defer wg.Done()
				task = taskWorker(task)
				taskSorter(task)
			}(t)
		}
		wg.Wait()
		close(doneTasks)
		close(undoneTasks)
	}()

	go func() {
		for {
			time.Sleep(3 * time.Second)
			fmt.Println("Errors:")
			for len(undoneTasks) > 0 {
				fmt.Println(<-undoneTasks)
			}

			fmt.Println("Done tasks:")
			for len(doneTasks) > 0 {
				task := <-doneTasks
				fmt.Printf("Task id %d finished at %s with result %s\n", task.id, task.fT, task.taskRESULT)
			}
		}
	}()

	wg.Wait()

}
