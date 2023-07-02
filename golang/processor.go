package main

import (
	"context"
	"errors"
	"sync"

	"go.uber.org/zap"
)

// Если пытаются добавить таску в обработку, но процессор уже завершился.
var ErrStopped = errors.New("processor is stopped")

// функция-обработчик задачи
type ProcessorFunc[T any] func(task T) (T, error)

// Processor создает пул воркеров и обрабатывает входящие таски
type Processor[T any] struct {
	params ProcessorParams[T]

	// superChan принимает все входящие задачи, из него читают запущеные воркеры.
	superChan chan<- T
	// stop при закрытии сигнализирует воркерам о необходимости завершиться.
	stop chan<- struct{}
	// done закрывается когда все воркеры завершились.
	done <-chan struct{}
}

type ProcessorParams[T any] struct {
	RunnerFunc ProcessorFunc[T]
	NumWorkers int
	InputCap   int
	ResCap     int
	ErrsCap    int
}

func NewProcessor[T any](params ProcessorParams[T]) *Processor[T] {
	return &Processor[T]{
		params: params,
	}
}

// Start запускает обработчик тасков.
func (p *Processor[T]) Start() (results <-chan T, errs <-chan error, doneCh <-chan struct{}) {
	// Эти каналы не закрываются никогда, сигнал о завершении работы - закрытие канала doneCh.
	superChan := make(chan T, p.params.InputCap)
	resChan := make(chan T, p.params.ResCap)
	errChan := make(chan error, p.params.ErrsCap)

	done := make(chan struct{})
	stop := make(chan struct{})

	p.superChan = superChan
	p.done = done
	p.stop = stop

	// Запускает N воркеров
	var running sync.WaitGroup
	running.Add(p.params.NumWorkers)
	for i := 0; i < p.params.NumWorkers; i++ {
		wid := i
		worker := worker[T]{
			runner:    p.params.RunnerFunc,
			superChan: superChan,
			resChan:   resChan,
			errChan:   errChan,
			stop:      stop,
		}

		log.Debug("RUN WORKER ", zap.Int("id", wid))
		go func() {
			defer running.Done()
			defer log.Debug("STOPPED WORKER ", zap.Int("id", wid))
			worker.run()
		}()
	}

	// Большой брат следит за воркерами
	go func() {
		defer close(done)
		running.Wait()
	}()

	return resChan, errChan, done
}

// worker описывает отдельного воркера, у котрого есть функция-обработчик,
// канал для входных данных, канал для выходных данных, канал для выходных ошибок и канал для остановки воркера.
type worker[T any] struct {
	runner    ProcessorFunc[T]
	superChan <-chan T
	resChan   chan<- T
	errChan   chan<- error
	stop      <-chan struct{}
}

// runWorker runs the worker.
func (w *worker[T]) run() {
	for {
		select {
		case task := <-w.superChan:
			res, err := w.runner(task)
			if err != nil {
				select {
				case w.errChan <- err:
					log.Debug("ERROR WRITTEN")
					continue
				case <-w.stop:
					return
				}
			}

			select {
			case w.resChan <- res:
				log.Debug("RESULT WRITTEN")
			case <-w.stop:
				return
			}
		case <-w.stop:
			return
		}
	}
}

// Stop останавливает обработчик тасков. Нельзя вызывать асинхронно.
func (p *Processor[T]) Stop(ctx context.Context) error {
	defer log.Debug("STOPPED")
	// foolproof
	if p.stop == nil {
		return nil
	}
	close(p.stop)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-p.done:
		return nil
	}
}

// AddTask добавляет таску в очередь обработки, но не гарантирует что таска будет обработана.
// Например если вызвать AddTask в момент вызова Stop
func (p *Processor[T]) AddTask(ctx context.Context, task T) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	// TODO on stop err
	case p.superChan <- task:
		log.Debug("TASK ADDED")
		return nil
	case <-p.done:
		return ErrStopped
	}
}
