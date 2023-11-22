// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// Приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме.
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков.

package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

const (
	executionTime = time.Second * 3
	taskQueueLen  = 10
)

var (
	errCreationFailed   = errors.New("creation failed")
	errProcessingFailed = errors.New("processing failed")
)

type taskStatus int

const (
	taskStatusUnknown taskStatus = iota
	taskStatusSuccess
	taskStatusFailed
)

// Task represents a meaninglessness of our life.
type task struct {
	id         int64
	createdAt  time.Time // время создания
	executedAt time.Time // время выполнения
	result     taskStatus
	err        error
}

func (t task) Error() string {
	return fmt.Sprintf("task id %d time %s, error: %v", t.id, t.createdAt, t.err)
}

func (t task) Id() string {
	return strconv.FormatInt(t.id, 10)
}

func create(ctx context.Context) <-chan task {
	out := make(chan task, taskQueueLen)

	go func() {
		defer close(out)
		for {
			var err error
			now := time.Now()
			id := now.UnixNano() / 1000
			status := taskStatusUnknown

			if id%2 == 1 { // time.Now().Nanosecond() всегда четна, генерация в микросекундах.
				err = errCreationFailed
				status = taskStatusFailed
			}

			current := task{
				id:        id,
				createdAt: now,
				result:    status,
				err:       err,
			}

			select {
			case out <- current: // передаем таск на выполнение
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

func proceed(in <-chan task) (<-chan task, <-chan task) {
	done := make(chan task, taskQueueLen)
	failed := make(chan task, taskQueueLen)

	go func() {
		defer func() {
			close(done)
			close(failed)
		}()
		// получение тасков
		for t := range in {
			if t.result != taskStatusFailed {
				t = work(t)
			}

			if t.result == taskStatusSuccess {
				done <- t
			} else {
				failed <- t
			}
		}
	}()

	return done, failed
}

func work(t task) task {
	// twentySecondsAgo := time.Now().Add(-20 * time.Second) // невыполняющееся условие
	// if t.createdAt.After(twentySecondsAgo) {
	if t.id%4 == 0 {
		t.result = taskStatusSuccess
	} else {
		t.err = errProcessingFailed
		t.result = taskStatusFailed
	}
	t.executedAt = time.Now()

	time.Sleep(time.Millisecond * 150)

	return t
}

func fill(wg *sync.WaitGroup, in <-chan task) *[]task {
	var out []task

	wg.Add(1)
	go func() {
		defer wg.Done()
		for t := range in {
			out = append(out, t)
		}
	}()

	return &out
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), executionTime)
	defer cancel()

	tasksChan := create(ctx)
	done, failed := proceed(tasksChan)

	var wg sync.WaitGroup
	results := fill(&wg, done)
	fails := fill(&wg, failed)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	select {
	case <-ctx.Done():
	case <-sigCh:
		cancel()
	}

	wg.Wait()
	println("Errors:")
	for _, t := range *fails {
		println(t.Error())
	}
	println("Done tasks:")
	for _, t := range *results {
		println(t.Id())
	}
}
