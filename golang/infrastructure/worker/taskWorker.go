package worker

import (
	"TestTask/infrastructure/creator"
	"context"
	"sync"
)

type TaskWorker[TMesValue creator.Stringable, TIn creator.TaskMessage[TMesValue], TOut any] struct {
	workerCount int
	executor    func(TIn) TOut
}

func (w TaskWorker[TMesValue, TIn, TOut]) Work(ctx context.Context, inCh <-chan TIn, outCh chan<- TOut, complete chan<- bool) {
	wg := sync.WaitGroup{}
	wg.Add(w.workerCount)
	for i := 0; i < w.workerCount; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case task, ok := <-inCh:
					if !ok {
						return
					}
					select {
					case <-ctx.Done():
						return
					case outCh <- w.executor(task):
					}
				}
			}
		}()
	}
	wg.Wait()
	complete <- true
}

func GetTaskWorker[TMesValue creator.Stringable, TOut any](workerCount int, executor func(creator.TaskMessage[TMesValue]) TOut) Worker[creator.TaskMessage[TMesValue], TOut] {
	return &TaskWorker[TMesValue, creator.TaskMessage[TMesValue], TOut]{
		workerCount: workerCount,
		executor:    executor,
	}
}
