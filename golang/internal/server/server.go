package server

import (
	"log"
	"sync"
	"taskprocessor/internal/task"
)

type sC chan task.Ttype

type S struct {
	log            *log.Logger
	errorChanel    chan error
	superChanel    sC
	superChanelOut sC
	doneChanel     sC
	result         struct {
		success map[int]task.Ttype
		errors  []error
	}
}

func New(l *log.Logger, chSize int) *S {

	s := S{
		log:            l,
		superChanel:    make(chan task.Ttype, chSize),
		superChanelOut: make(chan task.Ttype, chSize),
		doneChanel:     make(chan task.Ttype, chSize),
		errorChanel:    make(chan error, chSize),
	}
	s.result.success = make(map[int]task.Ttype)

	return &s
}

func (s *S) Process() {
	var mutex sync.Mutex

	go func() {
		go task.TaskGenerator(s.superChanel)
		go task.TaskProcessor(s.superChanel, s.superChanelOut)

		for {
			select {
			case t := <-s.superChanelOut:
				if ok, err := t.IsSucceed(); ok {
					s.doneChanel <- t
				} else {
					s.errorChanel <- err
				}
			case t := <-s.doneChanel:
				mutex.Lock()
				s.result.success[t.Id()] = t
				mutex.Unlock()
			case err := <-s.errorChanel:
				mutex.Lock()
				s.result.errors = append(s.result.errors, err)
				mutex.Unlock()
			}
		}
	}()
}

func (s *S) Stop() {
	if _, ok := <-s.superChanel; !ok {
		close(s.superChanel)
	}

	if _, ok := <-s.superChanelOut; !ok {
		close(s.superChanelOut)
	}

	if _, ok := <-s.errorChanel; !ok {
		close(s.errorChanel)
	}

	if _, ok := <-s.doneChanel; !ok {
		close(s.doneChanel)
	}
}

func (s *S) ShowResult() {
	if len(s.result.errors) > 0 {
		println("Errors:")
		for _, v := range s.result.errors {
			println(v.Error())
		}
	}

	println("Done tasks:")
	for _, v := range s.result.success {
		println(v.Id())
	}
}
