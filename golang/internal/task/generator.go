package task

import (
	"context"
	"log"
	"sync"
	"time"
)

func Generate(ctx context.Context, wg *sync.WaitGroup, pending chan *Task) {
	defer func() {
		close(pending)
		log.Print("pending channel closed ")
		wg.Done()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(40 * time.Millisecond)
			pending <- NewTask()
		}
	}

}
