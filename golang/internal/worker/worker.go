package worker

import (
	"context"
	"sync"
	"time"

	"github.com/Quantum12k/hr/golang/internal/task"
)

const (
	sleepDuration  = 150 * time.Millisecond
	maxActiveTasks = 5
)

type Worker struct {
	pendingTasksCh chan *task.Task
	DoneTasksCh    chan *task.Task
}

func New(ctx context.Context, taskAcceptorCh chan *task.Task) *Worker {
	w := &Worker{
		pendingTasksCh: taskAcceptorCh,
		DoneTasksCh:    make(chan *task.Task),
	}

	go w.run(ctx)

	return w
}

func (w *Worker) run(ctx context.Context) {
	defer w.cleanup()

	activeTasks := make(chan struct{}, maxActiveTasks)
	wg := sync.WaitGroup{}

	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			return
		case t, ok := <-w.pendingTasksCh:
			if !ok {
				return
			}

			activeTasks <- struct{}{}
			wg.Add(1)

			go func() {
				t.Execute()
				w.DoneTasksCh <- t

				<-activeTasks
				wg.Done()
			}()
		}

		time.Sleep(sleepDuration)
	}
}

func (w *Worker) cleanup() {
	close(w.DoneTasksCh)
}
