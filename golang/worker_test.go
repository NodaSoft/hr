package main

import (
	"context"
	"testing"
	"time"
)

func Test_parallelWorkerProcessInParallel(t *testing.T) {
	// test example
	// depending on project's / company's code style and code culture
	// we may choose write or not to write tests like this

	n := 3

	done := make(chan any)
	processing := make(chan any, n)

	worker := newParallelWorker[any, any](func(t any) (any, error) {
		processing <- t
		<-done

		return nil, nil
	}, nil, n)

	tasks := make(chan Task[any], n)
	for i := 0; i < n; i++ {
		tasks <- NewTask[any](nil)
	}

	results := worker.Run(tasks)

	timeout, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	processingCount := 0
	for processingCount < n {
		select {
		case <-processing:
			processingCount++

		case <-timeout.Done():
			t.Errorf("worker processes less than %d tasks", n)
			return
		}
	}

	select {
	case <-worker.Tokens:
		t.Errorf("worker can process more than %d tasks", n)
		return
	default:
	}

	close(done)
	close(tasks)

	completed := 0
	timeout, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
loop_result:
	for {
		select {
		case _, ok := <-results:
			if !ok {
				break loop_result
			} else {
				completed++
			}

		case <-timeout.Done():
			t.Errorf("result channel wasn't closed")
		}
	}
	if completed != n {
		t.Errorf("results = %d, expected %d", completed, n)
	}
}

func Benchmark_parallelWorker_Run(b *testing.B) {
	// pretty meaningless in this case

	n := b.N

	tasks := func() chan Task[any] {
		tasks := make(chan Task[any])
		go func() {
			for i := 0; i < n; i++ {
				tasks <- NewTask[any](nil)
			}
			close(tasks)
		}()
		return tasks
	}()

	worker := NewParallelWorker[any, any](func(t any) (any, error) { return nil, nil }, nil, 0)

	b.ResetTimer()

	results := worker(tasks)

	for range results {
	}
}
