package main

import (
	"fmt"
	"sync"
	"test_tack/task"
	"time"
)

func main() {

	c := make(chan *task.Task, 10)

	wg := sync.WaitGroup{}

	doneTasks := make(chan *task.Task, 100)
	undoneTasks := make(chan error, 100)

	arDoneTasks := []*task.Task{}
	arUndoneTasks := []error{}

	go func() {
		for v := range c {
			wg.Done()
			v.Work()
			if string(v.Result) == "success" {
				doneTasks <- v
			} else {
				undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", v.ID, v.TimeStart, v.Result)
			}

		}
	}()

	go readChan(doneTasks, undoneTasks, arDoneTasks, arUndoneTasks)

	start := time.Now()
	i := 1
	for time.Since(start) < 10*time.Second {
		t := task.New(i)
		i++

		wg.Add(1)
		c <- t
	}

	close(c)

	wg.Wait()
}

func readChan(s chan *task.Task, e chan error, arDoneTasks []*task.Task, arUndoneTasks []error) {

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	defer close(s)
	defer close(e)
	for {
		select {
		case value := <-s:
			value.TimeClose = time.Now()
			arDoneTasks = append(arDoneTasks, value)
		case val := <-e:
			arUndoneTasks = append(arUndoneTasks, val)
		case <-ticker.C:
			
			fmt.Println("Успешные таски:")
			for _, v := range arDoneTasks {
				fmt.Println(v)
			}

			fmt.Println("не Успешные таски:")
			for _, v := range arUndoneTasks {
				fmt.Println(v)
			}

		}
	}
}
