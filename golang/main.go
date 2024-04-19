package main

import (
	"fmt"
	"time"
)

type Ttype struct {
	id         int
	cT         time.Time 
	fT         time.Time 
	taskRESULT string
}

func main() {
	superChan := make(chan Ttype)
	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	go taskCreator(superChan)

	for i := 0; i < 5; i++ {
		go taskWorker(superChan, doneTasks, undoneTasks)
	}

	go resultCollector(doneTasks, undoneTasks)

	time.Sleep(3 * time.Second)
	close(superChan)
	close(doneTasks)
	close(undoneTasks)

	fmt.Println("Processing complete.")
}

func taskCreator(tasks chan<- Ttype) {
	for {
		if time.Now().Nanosecond()%2 > 0 {
			
			tasks <- Ttype{id: int(time.Now().Unix()), cT: time.Now(), taskRESULT: "error"}
		} else {
			tasks <- Ttype{id: int(time.Now().Unix()), cT: time.Now()}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func taskWorker(tasks <-chan Ttype, done chan<- Ttype, errors chan<- error) {
	for t := range tasks {
		if t.taskRESULT == "error" {
			errors <- fmt.Errorf("Task ID %d: Error occurred", t.id)
		} else {
			
			time.Sleep(150 * time.Millisecond)
			t.fT = time.Now()
			t.taskRESULT = "success"
			done <- t
		}
	}
}

func resultCollector(done <-chan Ttype, errors <-chan error) {
	for {
		select {
		case task, ok := <-done:
			if ok {
				fmt.Printf("Task ID %d completed at %s\n", task.id, task.fT.Format(time.RFC3339))
			}
		case err, ok := <-errors:
			if ok {
				fmt.Println(err)
			}
		case <-time.After(3 * time.Second):
			fmt.Println("Stopping result collection...")
			return
		}
	}
}
