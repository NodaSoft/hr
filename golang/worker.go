package main

import (
	"runtime"
	"sync"
	"time"
)

var (
	// Note: typically we choose default coefficient in range [2, 8] based on type of workload
	ParallelWorkerGoroutines = 4 * runtime.NumCPU()
)

// Note: using unique types for different ids allow not to missmatch them at compile time
type TaskID int64

// Task represents ability to fill our life
// with anything we want (but we have to make sure
// we'll be able to process it later)
type Task[T any] struct {
	ID       TaskID
	CreateAt time.Time
	Input    T
}

func NewTask[T any](t T) Task[T] {
	return Task[T]{
		ID:       TaskID(Snowflake()),
		CreateAt: time.Now(),
		Input:    t,
	}
}

type Processor[T, R any] func(T) (R, error)

type TaskProcessor[T, R any] func(Task[T]) Result[T, R]

type Result[T, R any] struct {
	ID          TaskID
	Task        *Task[T]
	CompletedAt time.Time
	Value       R
	Error       error
}

// Note: there are many viable option for worker interface.
// We choose simplest form - function that takes task channel and returns result channel.
type Worker[T, R any] func(<-chan Task[T]) <-chan Result[T, R]

type Middleware[T, R any] func(f func(Task[T]) Result[T, R]) func(Task[T]) Result[T, R]

func emptyMiddleware[T, R any](f func(Task[T]) Result[T, R]) func(Task[T]) Result[T, R] {
	return func(t Task[T]) Result[T, R] {
		return f(t)
	}
}

func NewParallelWorker[T, R any](processor Processor[T, R], m Middleware[T, R], n int) Worker[T, R] {
	w := newParallelWorker[T, R](processor, m, n)

	return w.Run
}

func newParallelWorker[T, R any](p Processor[T, R], m Middleware[T, R], n int) parallelWorker[T, R] {
	if n <= 0 {
		n = ParallelWorkerGoroutines
	}

	if m == nil {
		m = emptyMiddleware[T, R]
	}

	runner := m(func(t Task[T]) Result[T, R] {
		r, err := p(t.Input)

		result := Result[T, R]{
			ID:          t.ID,
			Task:        &t,
			CompletedAt: time.Now(),
			Value:       r,
			Error:       err,
		}

		return result
	})

	tokens := make(chan struct{}, n)
	for i := 0; i < n; i++ {
		tokens <- struct{}{}
	}

	w := parallelWorker[T, R]{
		Tokens: tokens,
		Runner: runner,
	}

	return w
}

type parallelWorker[T, R any] struct {
	Tokens chan struct{}
	Runner func(Task[T]) Result[T, R]
}

func (w *parallelWorker[T, R]) Run(tasks <-chan Task[T]) <-chan Result[T, R] {
	results := make(chan Result[T, R])

	go func() {
		wg := sync.WaitGroup{}
		defer func() {
			wg.Wait()
			close(results)
		}()

		for t := range tasks {
			<-w.Tokens

			wg.Add(1)
			go func(t Task[T]) {
				result := w.Runner(t)

				results <- result

				w.Tokens <- struct{}{}
				wg.Done()
			}(t)
		}
	}()

	return results
}
