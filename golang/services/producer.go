package services

import (
	"context"
	"nodasoft-golang/domain"
)

type producer struct {
	ctx   context.Context
	tasks chan *domain.Task
}

func newProducer(ctx context.Context, generationRate int) *producer {
	return &producer{
		ctx:   ctx,
		tasks: make(chan *domain.Task, generationRate),
	}
}

func (p *producer) getTasks() <-chan *domain.Task {
	return p.tasks
}

func (p *producer) generate() {
	defer close(p.tasks)
	defer println("generating channel closed")
	for {
		select {
		case p.tasks <- domain.NewTask():
			continue
		case <-p.ctx.Done():
			return
		}
	}
}
