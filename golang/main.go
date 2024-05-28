package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// Приложение должно генерировать таски 10 сек. Каждые 3 секунды должно выводить в консоль результат всех обработанных к этому моменту тасков (отдельно успешные и отдельно с ошибками).

// ЗАДАНИЕ: сделать из плохого кода хороший и рабочий - as best as you can.
// Важно сохранить логику появления ошибочных тасков.
// Важно оставить асинхронные генерацию и обработку тасков.
// Сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через pull-request в github
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	tasksChan, runCreator := NewTaskCreator()
	go runCreator(ctx)
	go func() {
		time.Sleep(time.Second * 3)
		cancel()
	}()
	runHandlersPool(ctx, tasksChan, 100).print()
}

func runHandlersPool(
	ctx context.Context,
	tasksChan <-chan Task,
	count int,
) _Result {
	var r _Result
	r.doneTasks = make(map[int]Task)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < count; i++ {
		doneTasksChan, errChan, handler := NewTaskHandler(tasksChan)
		go handler(ctx)
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case t := <-doneTasksChan:
					mu.Lock()
					r.doneTasks[t.ID] = t
					mu.Unlock()
				case err := <-errChan:
					mu.Lock()
					r.errs = append(r.errs, err)
					mu.Unlock()
				}
			}
		}()
	}
	wg.Wait()
	return r
}

type _Result struct {
	errs      []error
	doneTasks map[int]Task
}

func (r _Result) print() {
	fmt.Println("Errors:")
	for _, err := range r.errs {
		fmt.Printf("%v\n", err)
	}
	fmt.Println("Done tasks:")
	for _, task := range r.doneTasks {
		fmt.Printf("%v\n", task)
	}
}

// task_handler.go

type Task struct {
	ID      int
	Created time.Time // время создания
	Handled time.Time // время выполнения
	Err     error
	Result  []byte
}

func (t Task) String() string {
	return fmt.Sprintf(
		"Task(ID=%v, Created=%v, Handled=%v, Err=%v, _Result=%v)",
		t.ID,
		t.Created.Format(time.RFC3339),
		t.Handled.Format(time.RFC3339),
		t.Err,
		string(t.Result),
	)
}

var ErrOnCreateTask = errors.New("err on create task")

func NewTaskCreator() (<-chan Task, func(ctx context.Context)) {
	tasksChan := make(chan Task)
	started := atomic.Bool{}
	return tasksChan, func(ctx context.Context) {
		if started.Swap(true) {
			log.Printf("[ERROR] run creator: attempt to start creator again")
			return
		} else if ctx.Err() != nil {
			log.Printf("[ERROR] run creator: ctx already done")
			return
		}
		for {
			t := Task{ID: time.Now().Nanosecond()}
			if time.Now().Nanosecond()%2 > 0 { // фейлятся таски в нечетные наносекунды
				t.Err = ErrOnCreateTask
			}
			t.Created = time.Now()
			select {
			case <-ctx.Done():
				close(tasksChan)
				return
			case tasksChan <- t:
			}
		}
	}
}

// task_handler.go

func NewTaskHandler(tasksCh <-chan Task) (<-chan Task, <-chan error, func(ctx context.Context)) {
	doneTasksChan := make(chan Task)
	errChan := make(chan error)
	started := atomic.Bool{}
	return doneTasksChan, errChan, func(ctx context.Context) {
		if started.Swap(true) {
			log.Printf("[ERROR] run handler: attempt to start handler again")
			return
		} else if ctx.Err() != nil {
			log.Printf("[ERROR] run handler: ctx already done")
			return
		}
		onCtxDone := func() {
			close(doneTasksChan)
			close(errChan)
		}
		for t := range tasksCh {
			t = handleTask(t)
			if t.Err != nil {
				select {
				case <-ctx.Done():
					onCtxDone()
					return
				case errChan <- fmt.Errorf("handle task=%v: %w", t, t.Err):
				}
			} else {
				select {
				case <-ctx.Done():
					onCtxDone()
					return
				case doneTasksChan <- t:
				}
			}
		}
	}
}

var ErrOnHandleTask = errors.New("err on handle task")

func handleTask(task Task) Task {
	if task.Err != nil {
		// return task without changes
	} else if task.Created.After(time.Now().Add(-20 * time.Second)) {
		task.Result = []byte("task has been successes")
	} else {
		task.Result = []byte("something went wrong")
		task.Err = ErrOnHandleTask
	}
	time.Sleep(time.Millisecond * 150) // симуляция процесса работы?
	task.Handled = time.Now()
	return task
}
