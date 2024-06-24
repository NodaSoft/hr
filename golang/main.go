package main

import (
	"fmt"
	"sync"
	"time"
)

type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

func main() {
	superChan := make(chan Ttype, 10)
	var wg sync.WaitGroup

	taskCreator := func(a chan Ttype) {
		defer close(a)
		end := time.After(10 * time.Second)
		for {
			select {
			case <-end:
				return
			default:
				ft := time.Now().Format(time.RFC3339)
				if time.Now().Nanosecond()%2 > 0 {
					ft = "Some error occured"
				}
				a <- Ttype{cT: ft, id: int(time.Now().Unix())}
				time.Sleep(500 * time.Millisecond)
			}
		}
	}

	taskProcessor := func(a chan Ttype, wg *sync.WaitGroup) {
		defer wg.Done()
		successfulTasks := make([]Ttype, 0)
		failedTasks := make([]Ttype, 0)
		for task := range a {
			if task.cT == "Some error occured" {
				failedTasks = append(failedTasks, task)
			} else {
				successfulTasks = append(successfulTasks, task)
			}
		}
		fmt.Println("Successful tasks:", len(successfulTasks))
		fmt.Println("Failed tasks:", len(failedTasks))
	}

	wg.Add(1)
	go taskCreator(superChan)
	go taskProcessor(superChan, &wg)

	ticker := time.NewTicker(3 * time.Second)
	go func() {
		for range ticker.C {
			fmt.Println("Processing...")
		}
	}()

	wg.Wait()
	ticker.Stop()
	fmt.Println("All tasks processed.")
}
