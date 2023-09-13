package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

const (
	defaultDuration = 3 * time.Second
	afterDuration   = -20 * time.Second
	delayDuration   = time.Millisecond * 150
	defaultCapacity = 1024
	maxGoroutines   = 10
)

var ErrTimeBefore = errors.New("time before")

func op() error {
	n := time.Now().Nanosecond()
	if n%2 > 0 {
		return fmt.Errorf("bad task code %d", n)
	}
	return nil
}

type Request struct {
	Task Task
	fn   func() error
}

func (r *Request) Do(delay time.Duration) error {
	defer time.Sleep(delay)
	if err := r.fn(); err != nil {
		return err
	}
	afterTime := time.Now().Add(afterDuration)
	if r.Task.After(afterTime) {
		return nil
	}
	return ErrTimeBefore
}

type Result struct {
	id  uint32
	err error
}

type Task struct {
	id     uint32
	create time.Time
}

func (t *Task) Identifier() uint32 {
	return t.id
}

func (t *Task) After(u time.Time) bool {
	return t.create.After(u)
}

func requester(ctx context.Context) <-chan *Request {
	event := uint32(1)
	r := make(chan *Request)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				r <- &Request{Task: Task{id: event, create: time.Now()}, fn: op}
				atomic.AddUint32(&event, 1)
			}
		}
	}()
	return r
}

func worker(ctx context.Context, in <-chan *Request, out chan<- Result) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			select {
			case <-ctx.Done():
				return
			default:
				r := <-in
				id := r.Task.Identifier()
				err := r.Do(delayDuration)
				out <- Result{id: id, err: err}
			}

		}
	}

}

type safeInts struct {
	mu  sync.RWMutex
	arr []uint32
}

func (si *safeInts) append(i uint32) {
	si.mu.Lock()
	defer si.mu.Unlock()
	si.arr = append(si.arr, i)
}

func (si *safeInts) Len() int {
	si.mu.RLock()
	defer si.mu.RUnlock()
	return len(si.arr)
}

func (si *safeInts) String() string {
	si.mu.RLock()
	defer si.mu.RUnlock()
	var b []byte
	for i, n := range si.arr {
		if i > 0 {
			b = append(b, '\n')
		}
		b = append(b, strconv.Itoa(int(n))...)
	}
	return string(b)
}

type safeErrors struct {
	mu   sync.RWMutex
	errs []error
}

func (se *safeErrors) append(err error) {
	se.mu.Lock()
	defer se.mu.Unlock()
	se.errs = append(se.errs, err)
}

func (se *safeErrors) Len() int {
	se.mu.RLock()
	defer se.mu.RUnlock()
	return len(se.errs)
}

func (se *safeErrors) String() string {
	return errors.Join(se.errs...).Error()
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDuration)
	defer cancel()

	result := make(chan Result, maxGoroutines)

	requests := requester(ctx)
	for i := 0; i < maxGoroutines; i++ {
		go worker(ctx, requests, result)
	}

	errs := safeErrors{errs: make([]error, 0, defaultCapacity)}
	ids := safeInts{arr: make([]uint32, 0, defaultCapacity)}

MainLoop:
	for {
		select {
		case <-ctx.Done():
			break MainLoop
		case r := <-result:
			if r.err != nil {
				errs.append(r.err)
			} else {
				ids.append(r.id)
			}
		}
	}

	// display result
	fmt.Printf("Errors:\n%s\n", errs.String())
	fmt.Printf("Done tasks:\n%s\n", ids.String())
	fmt.Printf("Total tasks: %d\n", errs.Len()+ids.Len())
	fmt.Printf("Done.\n")
}
