package creator

import (
	"context"

	"github.com/Quantum12k/hr/golang/internal/task"
)

type Creator struct {
	NewTasksCh chan *task.Task
}

func New(ctx context.Context) *Creator {
	c := &Creator{
		NewTasksCh: make(chan *task.Task),
	}

	go c.run(ctx)

	return c
}

func (c *Creator) run(ctx context.Context) {
	defer c.cleanup()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			c.NewTasksCh <- task.New()
		}
	}
}

func (c *Creator) cleanup() {
	close(c.NewTasksCh)
}
