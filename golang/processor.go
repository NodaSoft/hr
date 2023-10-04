package main

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"
)

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

type Processor interface {
	Loop(ctx context.Context) (map[int]Ttype, []error)
}

func NewProcessor(goroutines int, bufferSize int, createTask func() Ttype) Processor {
	return &processor{
		goroutines:  goroutines,
		bufferSize:  bufferSize,
		buffer:      make(chan Ttype, bufferSize),
		doneTasks:   make(chan Ttype, bufferSize/2),
		undoneTasks: make(chan error, bufferSize/2),
		createTask:  createTask,
	}
}

type processor struct {
	goroutines  int
	bufferSize  int
	buffer      chan Ttype
	doneTasks   chan Ttype
	undoneTasks chan error
	createTask  func() Ttype
}

func (p *processor) Loop(ctx context.Context) (map[int]Ttype, []error) {
	result := make(map[int]Ttype, p.bufferSize/2)
	err := make([]error, 0, p.bufferSize/2)

	var wg sync.WaitGroup
	wg.Add(3)

	go p.taskCreator(ctx, &wg)
	go p.taskReciever(ctx, &wg)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case r := <-p.doneTasks:
				result[r.id] = r
			case r := <-p.undoneTasks:
				err = append(err, r)
			}
		}
	}()

	wg.Wait()

	return result, err
}

func (p *processor) taskCreator(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case p.buffer <- p.createTask():
		}
	}
}

func (p *processor) taskReciever(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	limit := make(chan struct{}, p.goroutines)

	for {
		select {
		case <-ctx.Done():
			return
		case t := <-p.buffer:
			wg.Add(1)
			limit <- struct{}{}
			go func() {
				defer wg.Done()
				defer func() { <-limit }()
				p.taskSorter(ctx, p.taskWorker(t))
			}()

		}
	}
}

func (p *processor) taskWorker(a Ttype) Ttype {
	tt, err := time.Parse(time.RFC3339, a.cT)
	if err != nil {
		a.taskRESULT = []byte(fmt.Sprintf("time parse error [%v]", err))
		return a
	}

	if bytes.Equal(a.taskRESULT, ErrorResultBytes) {
		return a
	}

	if tt.Before(time.Now().Add(-20 * time.Second)) {
		a.taskRESULT = []byte("task outdated")
		return a
	}

	time.Sleep(time.Millisecond * 150) // doing hard work

	a.taskRESULT = []byte("task has been successed")
	a.fT = time.Now().Format(time.RFC3339Nano)

	return a
}

func (p *processor) taskSorter(ctx context.Context, t Ttype) {
	if len(t.taskRESULT) > 14 && string(t.taskRESULT[14:]) == "successed" {
		select {
		case <-ctx.Done():
			return
		case p.doneTasks <- t:
			return
		}
	}

	select {
	case <-ctx.Done():
		return
	case p.undoneTasks <- fmt.Errorf("task id [%d] time [%s], error [%s]", t.id, t.cT, t.taskRESULT):
		return
	}
}
