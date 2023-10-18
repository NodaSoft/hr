package tasks

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

type Generator struct {
	hasStarted atomic.Bool
	tasks      chan *Task
	stop       chan struct{}
	stopped    chan struct{}
}

func NewGenerator(maxQueue int) *Generator {
	return &Generator{
		hasStarted: atomic.Bool{},
		tasks:      make(chan *Task, maxQueue),
		stop:       make(chan struct{}, 1),
		stopped:    make(chan struct{}, 1),
	}
}

func (g *Generator) Start() {
	if !g.hasStarted.CompareAndSwap(false, true) {
		return // already started
	}

	go func() {
		var id int
		for {
			id++ // for simplicity and readability, since generator is mock anyway (could be UUID)

			// this is part of the task (error simulation)
			createdAt := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 {
				createdAt = "Some error occurred"
			}

			select {
			case <-g.stop:
				g.stopped <- struct{}{}
				return
			case g.tasks <- New(id, createdAt):
				continue
			}
		}
	}()
}

func (g *Generator) Stop(ctx context.Context) (err error) {
	if !g.hasStarted.Load() {
		return nil
	}

	go func() {
		g.stop <- struct{}{}
	}()

	select {
	case <-g.stopped:
		g.hasStarted.CompareAndSwap(true, false)
	case <-ctx.Done():
		err = fmt.Errorf("failed to stop task generator: %w", ctx.Err())
	}

	return err
}

func (g *Generator) Tasks() <-chan *Task {
	return g.tasks
}
