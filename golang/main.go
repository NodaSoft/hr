package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// * обновленный код отправить через pull-request.

// Приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме.
// Должно выводить успешные таски и ошибки по мере выполнения.
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

// A Ttype represents a meaninglessness of our life
type (
	Ttype struct {
		id         int
		cT         string // время создания
		taskRESULT string
	}

	taskMap struct {
		mu    sync.RWMutex
		tasks map[int]Ttype
	}

	errs struct {
		mu     sync.RWMutex
		errors []error
	}
)

const (
	taskSuccess = "success"
	taskFailed  = "failed"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	tasksCh := make(chan Ttype, 10)

	go createWorker(ctx, tasksCh)

	doneTasksCh := make(chan Ttype)
	undoneTasksCh := make(chan error)

	go taskSorter(ctx, tasksCh, doneTasksCh, undoneTasksCh)

	result := &taskMap{
		tasks: make(map[int]Ttype),
	}
	e := errs{
		errors: make([]error, 0),
	}

	go doneReader(doneTasksCh, result)
	go undoneReader(undoneTasksCh, &e)

	<-ctx.Done()

	fmt.Println("Errors:")
	for _, err := range e.read() {
		fmt.Println(err)
	}

	fmt.Println("Tasks:")
	for _, task := range result.readAll() {
		fmt.Println(task)
	}
}

func (t *taskMap) write(key int, value Ttype) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.tasks[key] = value
}

func (t *taskMap) read(key int) Ttype {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.tasks[key]
}

func (t *taskMap) readAll() []Ttype {
	t.mu.RLock()
	defer t.mu.RUnlock()

	tasks := make([]Ttype, 0, len(t.tasks))

	for _, task := range t.tasks {
		tasks = append(tasks, task)
	}

	return tasks
}

func (e *errs) write(err error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.errors = append(e.errors, err)
}

func (e *errs) read() []error {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.errors
}

func createWorker(ctx context.Context, ch chan Ttype) {
	for {
		select {
		case <-ctx.Done():
			close(ch)

			return
		default:
			ft := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				ft = "Some error occured"
			}
			ch <- Ttype{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
		}
	}
}

func taskSorter(ctx context.Context, tasksChan, doneCh chan Ttype, undoneCh chan error) {
	for {
		select {
		case <-ctx.Done():
			close(doneCh)
			close(undoneCh)

			return
		case t, ok := <-tasksChan:
			if !ok {
				return
			}

			t = taskResult(t)

			if t.taskRESULT == taskSuccess {
				doneCh <- t
			} else {
				undoneCh <- fmt.Errorf("task id %d time %s, result %s", t.id, t.cT, t.taskRESULT)
			}
		}
	}
}

func taskResult(task Ttype) Ttype {
	tt, _ := time.Parse(time.RFC3339, task.cT)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		task.taskRESULT = taskSuccess
	} else {
		task.taskRESULT = taskFailed
	}

	time.Sleep(time.Millisecond * 150)

	return task
}

func doneReader(doneTasksCh chan Ttype, taskMap *taskMap) {
	for r := range doneTasksCh {
		taskMap.write(r.id, r)
	}
}

func undoneReader(undoneTasksCh chan error, e *errs) {
	for err := range undoneTasksCh {
		e.write(err)
	}
}
