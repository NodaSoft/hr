package task

import (
	"log"
	"sync"
)

func Sort(wg *sync.WaitGroup, completedQueue, doneTasks, errTasks chan *Task) {
	defer func() {
		close(doneTasks)
		close(errTasks)
		log.Println("closed done and err channels")
		wg.Done()
	}()

	for t := range completedQueue {
		err := t.Result.Error()
		if err != nil {
			errTasks <- t
		} else {
			doneTasks <- t
		}
	}
}
