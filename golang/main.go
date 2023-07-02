package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// Приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме.
// В конце должно выводить успешные таски и ошибки при выполнении остальных тасков.

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id           int
	creationTime string // время создания
	finishTime   string // время выполнения
	taskResult   []byte
}

// spamTasks яростно спамит тасками в Processor.
func spamTasks(ctx context.Context, p *Processor[Ttype]) {
	for {
		now := time.Now()
		ft := now.Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			ft = "Some error occured"
		}
		// передаем таск на выполнение
		if err := p.AddTask(ctx, Ttype{creationTime: ft, id: int(now.Unix())}); err != nil {
			// Ну эту ошибку я считай сам придумал, поэтому особо умно обрабатывать её не нужно.
			if errors.Is(err, ErrStopping) {
				return
			}
			if ctx.Err() != nil {
				return
			}
			// all other unexpected err
			panic(err)
		}
	}
}

func main() {
	var (
		p = NewProcessor[Ttype](ProcessorParams[Ttype]{
			RunnerFunc: StupidWorker,
			NumWorkers: 10,
			InputCap:   10,
			ResCap:     10,
			ErrsCap:    10,
		})

		wg sync.WaitGroup

		taskResults = map[int]Ttype{}
		taskErrors  = []error{}
	)

	spamCtx, spamCancel := context.WithCancel(context.Background())
	defer spamCancel()

	// Запуск обрабатывателя.
	results, errs := p.Start()

	// Запуск спаммера.
	wg.Add(1)
	go func() {
		spamTasks(spamCtx, p)
	}()

	// Обработка результата.
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			// Если оба канала закрыты, то можно выходить
			if results == nil && errs == nil {
				return
			}
			select {
			case res, ok := <-results:
				if !ok {
					results = nil // nil канал блокируется на чтение
				}
				taskResults[res.id] = res
			case err, ok := <-errs:
				if !ok {
					errs = nil // nil канал блокируется на чтение
				}
				taskErrors = append(taskErrors, err)
			}
		}
	}()

	time.Sleep(time.Second * 3) // Ждун
	spamCancel()                // отмена спаммера

	// Задачи блокируются на 150мс, я решил их подождать 200мс.
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer stopCancel()
	if err := p.Stop(stopCtx); err != nil {
		panic("unable to stop properly: " + err.Error())
	}

	// Дожидаемся пока обработается результат и ошибки
	wg.Wait()

	// Я подумал что важно выводить информацию именно после окончания работы, а не в процессе.
	println("Errors:")
	for r := range taskErrors {
		println(r)
	}

	println("Done tasks:")
	for r := range taskResults {
		println(r)
	}
}

// функция-обработчик задачи
type ProcessorFunc[T any] func(task T) (T, error)

type Processor[T any] struct {
	params ProcessorParams[T]

	// superChan принимает все входящие задачи, из него читают запущеные воркеры.
	superChan chan T
	// resultChan отдает задачи, которые были обработаны воркерами без ошибок.
	resultChan chan T
	// errChan возвращает ошибки, полученные в процессе обработки задач.
	errChan chan error

	// stop сигнализирует воркерам о необходимости завершиться.
	stop chan struct{}
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
		params:    params,
		superChan: make(chan T, params.InputCap),
	}
}

// Start запускает обработчик тасков.
func (p *Processor[T]) Start() (resultCh <-chan T, errCh <-chan error) {
	results := make(chan T, p.params.ResCap)
	errs := make(chan error, p.params.ErrsCap)
	done := make(chan struct{})
	stop := make(chan struct{})

	p.resultChan = results
	p.errChan = errs
	p.done = done
	p.stop = stop

	// Запускает N воркеров
	var running sync.WaitGroup
	running.Add(p.params.NumWorkers)
	for i := 0; i < p.params.NumWorkers; i++ {
		go func() {
			running.Done()
			p.runWorker()
		}()
	}

	// Большой брат следит за воркерами
	go func() {
		defer close(done)
		defer close(errs)
		defer close(results)
		running.Wait()
	}()

	return results, errCh
}

// runWorker runs the worker.
func (p *Processor[T]) runWorker() {
	for {
		select {
		case task := <-p.superChan:
			res, err := p.params.RunnerFunc(task)
			if err != nil {
				select {
				case p.errChan <- err:
				case <-p.stop:
					return
				}
				continue
			}

			select {
			case p.resultChan <- res:
			case <-p.stop:
				return
			}
		case <-p.stop:
			return
		}
	}
}

// Stop останавливает обработчик тасков. Нельзя вызывать асинхронно.
func (p *Processor[T]) Stop(ctx context.Context) error {
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

var ErrStopping = errors.New("processor is stopping now")

// AddTask добавляет таску в очередь обработки, но не гарантирует что таска будет обработана.
// Например если вызвать AddTask в момент вызова Stop
func (p *Processor[T]) AddTask(ctx context.Context, task T) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	// TODO on stop err
	case p.superChan <- task:
		return nil
	case <-p.done:
		return ErrStopping
	}
}

// StupidWorker пинает SmartWorker чтобы он работал, а потом тупит 150мс.
func StupidWorker(data Ttype) (Ttype, error) {
	res, err := SmartWorker(data)

	// В любом случае это костыль чтобы проц не сгорел, поэтому без разницы куда его пихать.
	time.Sleep(time.Millisecond * 150)

	return res, err
}

// SmartWorker описывает какую-то умную логику обработки задачи.
func SmartWorker(t Ttype) (_ Ttype, err error) {
	// Запись времени окончания обработки даже в случае ошибки
	defer func() {
		t.finishTime = time.Now().Format(time.RFC3339Nano)
		// Я перенес форматирование ошибки в defer только ради того, чтобы эту ошибку можно было записать в taskResult.
		if err != nil {
			t.taskResult = []byte(err.Error())
			// не люблю кидать просто %s в строку, обычно использую %q, либо в формате прописываю кавычки.
			err = fmt.Errorf("Task id[%d], time[%s], error[%s]", t.id, t.creationTime, err)
		}
	}()

	tt, err := time.Parse(time.RFC3339, t.creationTime)
	if err != nil {
		return t, err
	}

	if tt.IsZero() || time.Since(tt) > 20*time.Second {
		return t, errors.New("something went wrong")
	}

	t.taskResult = []byte("task has been successed")
	return t, nil
}
