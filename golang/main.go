package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	ErrTaskFailed   = errors.New("task failed")
	ErrPastDeadline = errors.New("past deadline")
)

const (
	RunTime     = 3 * time.Second
	TaskTimeout = 1 * time.Second

	TaskResult = "success"
)

func Process(t struct{}) ([]byte, error) {
	time.Sleep(150 * time.Millisecond)

	failed := time.Now().Nanosecond()%2 > 0
	if failed {
		return nil, ErrTaskFailed
	} else {
		return []byte(TaskResult), nil
	}
}

func NewTimeoutMiddleware[T, R any](d time.Duration) Middleware[T, R] {
	return func(f func(Task[T]) Result[T, R]) func(Task[T]) Result[T, R] {
		return func(t Task[T]) Result[T, R] {
			deadline := time.Now().Add(-d)
			if t.CreateAt.Before(deadline) {
				result := Result[T, R]{
					ID:          t.ID,
					Task:        &t,
					CompletedAt: time.Now(),
					Error:       ErrPastDeadline,
				}

				return result
			}

			return f(t)
		}
	}
}

func main() {
	// There are multiple options to satisfy
	// "приложение эмулирует получение и обработку тасков..."
	// For example we can create App struct, config, etc. in order to emulate real life app.
	// or we can put everything in main:

	ctx, cancel := context.WithTimeout(context.Background(), RunTime)

	tasks := func(ctx context.Context) <-chan Task[struct{}] {
		tasks := make(chan Task[struct{}])

		go func() {
			for {
				task := Task[struct{}]{
					ID:       TaskID(Snowflake()),
					CreateAt: time.Now(),
					Input:    struct{}{},
				}

				select {
				case tasks <- task:

				case <-ctx.Done():
					close(tasks)
					return
				}
			}
		}()

		return tasks
	}(ctx)

	timeoutMiddleware := NewTimeoutMiddleware[struct{}, []byte](TaskTimeout)
	worker := NewParallelWorker(Process, timeoutMiddleware, 0)

	results := worker(tasks)

	completed := make(map[TaskID]Result[struct{}, []byte])
	failed := make([]Result[struct{}, []byte], 0)
loop:
	for {
		select {
		case r, ok := <-results:
			if !ok {
				break loop
			}

			if r.Error == nil {
				completed[r.ID] = r
			} else {
				failed = append(failed, r)
			}

		case <-terminationSignal():
			cancel()
		}
	}

	log.Println("Errors:")
	for _, r := range failed {
		createdAt := r.Task.CreateAt.Format(time.RFC3339)
		log.Printf("Task id %d time %s, error '%s'", r.ID, createdAt, r.Error)
	}

	log.Println("Done task IDs:")
	for id := range completed {
		log.Println(id)
	}
}

func terminationSignal() <-chan os.Signal {
	c := make(chan os.Signal, 2)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	return c
}
