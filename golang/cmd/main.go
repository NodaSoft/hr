package main

import (
	"task/internal/task"
	"time"
)

func main() {
	var errA []error // error array

	quit := time.After(time.Second * 10)

	inter := make(chan bool, 1) // interrupter

	c := make(chan task.Task)     // superChain
	result := map[int]task.Task{} // result

	// interrupter function
	go func() {
		for {
			select {
			case <-quit:
				return
			default:
				time.Sleep(time.Second * 3)
				inter <- true
			}
		}
	}()

	go task.Create(c, quit)

	// sort function
	go func() {
		for t := range c {
			t, err := t.Worker()
			if err != nil {
				errA = append(errA, err)
			} else {
				result[len(result)+1] = t
			}
		}
	}()

	// using main goroutine
	for {
		select {
		case <-inter:
			println("Errors:")
			for _, r := range errA {
				println(r.Error())
			}

			println("Done tasks:")
			for _, r := range result {
				println(r.Id)
			}
		case <-quit:
			return
		}
	}
}
