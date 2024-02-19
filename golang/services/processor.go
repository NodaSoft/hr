package services

import (
	"context"
	"math/rand"
	"nodasoft-golang/domain"
	"time"
)

type processor struct {
	ctx           context.Context
	incomingTasks <-chan *domain.Task
	success       chan *domain.Task
	errored       chan *domain.Task
}

func newProcessor(ctx context.Context, tasks <-chan *domain.Task) *processor {
	return &processor{
		ctx:           ctx,
		incomingTasks: tasks,
		success:       make(chan *domain.Task),
		errored:       make(chan *domain.Task),
	}
}

func (p *processor) resultChannels() (success, errors <-chan *domain.Task) {
	return p.success, p.errored
}

func (p *processor) doWork() {
	for task := range p.incomingTasks {
		p.processTask(task)
		p.routeToOutgoingChannels(task)
	}
	defer func() {
		close(p.success)
		close(p.errored)
	}()
}

func (p *processor) processTask(task *domain.Task) {
	task.MarkAsInProgress()
	timeToSleep := time.Duration(rand.Intn(300)) * time.Millisecond
	select {
	case <-time.After(timeToSleep):
		if p.canMarkTaskAsErrored(task) {
			task.MarkAsErrored()
			return
		}
		task.MarkAsSuccessfullyCompleted()
	case <-p.ctx.Done():
		// Ну например так. Лучше, чем просто оставить таску в промежуточном состоянии
		task.MarkAsErrored()
		return
	}
}

// [canMarkTaskAsErrored] Условие появления ошибочных тасков.
// Логично помечать таску ошибочной при обработке, а не при генерации,
// как задании
func (p *processor) canMarkTaskAsErrored(task *domain.Task) bool {
	return time.Now().Nanosecond()%2 > 0
}

func (p *processor) routeToOutgoingChannels(task *domain.Task) {
	if task.IsSuccessful() {
		p.success <- task
	} else {
		p.errored <- task
	}
}
