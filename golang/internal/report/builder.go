package report

import (
	"context"
	"log"
	"sync"
	"task/internal/task"
)

type Builder struct {
	m         sync.Mutex
	succeeded map[int]string
	errors    map[int]string
}

func NewBuilder() *Builder {
	return &Builder{
		succeeded: make(map[int]string),
		errors:    make(map[int]string),
	}
}

func (b *Builder) AddSucceeded(t *task.Task) {
	b.m.Lock()
	defer b.m.Unlock()
	b.succeeded[t.Id] = t.String()
}

func (b *Builder) AddError(t *task.Task) {
	b.m.Lock()
	defer b.m.Unlock()
	b.errors[t.Id] = t.String()
}

func (b *Builder) Report() Report {
	b.m.Lock()
	defer b.m.Unlock()
	report := NewReport(b.succeeded, b.errors)
	b.unsafeReset()
	return report
}

func (b *Builder) unsafeReset() {
	b.succeeded = make(map[int]string)
	b.errors = make(map[int]string)
}

func BuildReport(wg *sync.WaitGroup, b *Builder, doneTasks, errTasks chan *task.Task, buildStop context.CancelFunc) {
	defer func() {
		log.Println("stop building reports")
		buildStop()
		wg.Done()
	}()

	lwg := sync.WaitGroup{}

	lwg.Add(1)
	go func() {
		defer lwg.Done()

		for t := range doneTasks {
			b.AddSucceeded(t)
		}
	}()

	lwg.Add(1)
	go func() {
		defer lwg.Done()

		for e := range errTasks {
			b.AddError(e)
		}
	}()

	lwg.Wait()
}
