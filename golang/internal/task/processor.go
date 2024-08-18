package task

import (
	"log"
	"sync"
)

func Process(wg *sync.WaitGroup, pending, completed chan *Task) {
	defer func() {
		close(completed)
		log.Print("closed completed channel. All the tasks are completed.")
		wg.Done()
	}()

	lwg := sync.WaitGroup{}

	for t := range pending {
		lwg.Add(1)
		go func(task *Task) {
			defer lwg.Done()
			completed <- Worker(task)
		}(t)
	}
	lwg.Wait()
}
