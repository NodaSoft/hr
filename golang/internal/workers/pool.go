package workers

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/NodaSoft/hr/golang/internal/tasks"
)

type Pool struct {
	workers    chan struct{}
	tasks      <-chan *tasks.Task
	processed  chan *tasks.Task
	inProgress atomic.Int32
	launched   atomic.Bool
	shutdown   chan struct{}
	onShutdown chan struct{}
	stopped    chan struct{}
	result     struct {
		failed map[*tasks.Task]error
		done   map[int]*tasks.Task
	}
}

func NewPool(size int, t <-chan *tasks.Task) *Pool {
	return &Pool{
		workers:    make(chan struct{}, size),
		tasks:      t,
		processed:  make(chan *tasks.Task, size),
		inProgress: atomic.Int32{},
		launched:   atomic.Bool{},
		shutdown:   make(chan struct{}, 1),
		onShutdown: make(chan struct{}, 1),
		stopped:    make(chan struct{}, 1),
	}
}

func (p *Pool) Execute() {
	if !p.launched.CompareAndSwap(false, true) { // already launched
		return
	}

	p.result.failed = make(map[*tasks.Task]error)
	p.result.done = make(map[int]*tasks.Task)

	// task processing (workers)
	go func() {
		var onShutdown bool
		for {
			select {
			case <-p.shutdown: // waiting for shutdown
				onShutdown = true
			case p.workers <- struct{}{}: // waiting for available worker
				select {
				case <-p.shutdown: // waiting for shutdown
					onShutdown = true
				case task := <-p.tasks: // waiting for a new task
					go p.launchWorker(task)
					p.inProgress.Add(1)
				}
			}

			if onShutdown {
				p.onShutdown <- struct{}{}
				return
			}
		}
	}()

	// result processing (sorting)
	go func() {
		var finishAndExit bool
		for {
			select {
			case <-p.onShutdown:
				finishAndExit = true
			case task := <-p.processed:
				if _, err := task.State(); err != nil {
					p.result.failed[task] = fmt.Errorf("id %d, created %s: %s", task.Id(), task.CreatedAt(), err)
				} else {
					p.result.done[task.Id()] = task
				}
				p.inProgress.Add(-1)
			}

			if finishAndExit && p.inProgress.Load() == 0 {
				p.stopped <- struct{}{}
				return
			}
		}
	}()
}

// Shutdown worker pool gracefully.
// All started workers will be finished and their results processed.
func (p *Pool) Shutdown(ctx context.Context) (err error) {
	if !p.launched.Load() {
		return nil
	}

	go func() {
		p.shutdown <- struct{}{}
	}()

	select {
	case <-p.stopped:
		p.launched.CompareAndSwap(true, false)
	case <-ctx.Done():
		err = fmt.Errorf("failed to shutdown worker pool gracefully: %w", ctx.Err())
	}

	return err
}

// Will return error if worker pool still running.
func (p *Pool) Result() (done map[int]*tasks.Task, failed map[*tasks.Task]error, err error) {
	if p.launched.Load() {
		return nil, nil, errors.New("worker pool is still running")
	}
	return p.result.done, p.result.failed, nil
}
