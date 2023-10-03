package main

import (
	"bytes"
	"fmt"
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
	Loop() (map[int]Ttype, []error)
}

func NewProcessor(bufferSize int, createTask func() Ttype) Processor {
	return &processor{
		bufferSize:  bufferSize,
		buffer:      make(chan Ttype, bufferSize),
		doneTasks:   make(chan Ttype, bufferSize/2),
		undoneTasks: make(chan error, bufferSize/2),
		createTask:  createTask,
	}
}

type processor struct {
	bufferSize  int
	buffer      chan Ttype
	doneTasks   chan Ttype
	undoneTasks chan error
	createTask  func() Ttype
}

func (p *processor) Loop() (map[int]Ttype, []error) {
	result := make(map[int]Ttype, p.bufferSize/2)
	err := make([]error, 0, p.bufferSize/2)

	go p.taskCreator()
	go p.taskReciever()

	go func() {
		for r := range p.doneTasks {
			result[r.id] = r
		}
	}()

	go func() {
		for r := range p.undoneTasks {
			err = append(err, r)
		}
	}()

	time.Sleep(time.Second * 5)

	return result, err
}

func (p *processor) taskCreator() {
	for {
		p.buffer <- p.createTask()
	}
}

func (p *processor) taskReciever() {
	for t := range p.buffer {
		ct := t
		go func() {
			p.taskSorter(p.taskWorker(ct))
		}()
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

	time.Sleep(time.Millisecond * 150)

	a.taskRESULT = []byte("task has been successed")
	a.fT = time.Now().Format(time.RFC3339Nano)

	return a
}

func (p *processor) taskSorter(t Ttype) {
	if len(t.taskRESULT) > 14 && string(t.taskRESULT[14:]) == "successed" {
		p.doneTasks <- t
	} else {
		p.undoneTasks <- fmt.Errorf("task id [%d] time [%s], error [%s]", t.id, t.cT, t.taskRESULT)
	}
}
