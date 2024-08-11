package main

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
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

var (
	errTask = errors.New("some error occured")
)

const (
	resultOK           = "task has been successed"
	resultErr          = "something went wrong"
	workInterval       = 150 * time.Millisecond
	generationInterval = 10 * time.Second
	printInterval      = 3 * time.Second
)

type task struct {
	id         int
	createdAt  time.Time
	finishedAt time.Time
	result     string
	err        error
}

func (t *task) error() string {
	return fmt.Sprintf("Id: %d created: %s, error: %s, result: %s\n",
		t.id, t.createdAt.Format(time.RFC3339), t.err, t.result)
}

func (t *task) success() string {
	return fmt.Sprintf("id: %d, created: %s, finished: %s, result: %s\n",
		t.id, t.createdAt.Format(time.RFC3339), t.finishedAt.Format(time.RFC3339), t.result)
}

func main() {
	done := make(chan struct{})
	tasks := tasksPipe(done)

	w := newWorker(printInterval, tasks)
	go w.work()
	go w.sort()
	go w.print()

	time.Sleep(generationInterval)
	close(done)
	w.stop()
}
func tasksPipe(
	done <-chan struct{},
) <-chan task {
	tasks := make(chan task)

	go func() {
		var i int
		for {
			select {
			case <-done:
				close(tasks)
				return
			default:
				i++
				task := task{id: i, createdAt: time.Now()}

				if task.createdAt.UnixMicro()%2 > 0 {
					task.err = errTask
				}

				tasks <- task
			}
		}
	}()

	return tasks
}

type worker struct {
	printInterval time.Duration

	tasks    <-chan task
	results  chan task
	shutdown chan struct{}
	wg       sync.WaitGroup

	errors    *bytes.Buffer
	successes *bytes.Buffer
}

func newWorker(
	printInterval time.Duration,
	tasks <-chan task,
) *worker {
	w := worker{
		printInterval: printInterval,
		tasks:         tasks,
		results:       make(chan task),
		shutdown:      make(chan struct{}),
		errors:        new(bytes.Buffer),
		successes:     new(bytes.Buffer),
		wg:            sync.WaitGroup{},
	}

	return &w
}

func (w *worker) work() {
	w.wg.Add(1)
	for {
		select {
		case <-w.shutdown:
			close(w.results)
			w.wg.Done()
			return
		case task := <-w.tasks:

			if task.err != nil {
				task.result = resultErr
			} else {
				task.result = resultOK
				task.finishedAt = time.Now()
			}

			w.results <- task
			time.Sleep(workInterval)
		}

	}
}

func (w *worker) sort() {
	w.wg.Add(1)
	for {
		select {
		case <-w.shutdown:
			w.wg.Done()
			return
		case task := <-w.results:
			if task.err != nil {
				w.errors.WriteString(task.error())
			} else {
				w.successes.WriteString(task.success())
			}
		}
	}
}

func (w *worker) print() {
	w.wg.Add(1)
	for {
		select {
		case <-w.shutdown:
			w.printResults()
			w.wg.Done()
			return
		case <-time.Tick(w.printInterval):
			w.printResults()
		}
	}
}

func (w *worker) printResults() {
	if w.successes.Len() > 0 {
		fmt.Println("Successes:")
		fmt.Println(w.successes.String())
	}
	if w.errors.Len() > 0 {
		fmt.Println("Errors:")
		fmt.Println(w.errors.String())
	}

	w.errors.Reset()
	w.successes.Reset()
}

func (w *worker) stop() {
	close(w.shutdown)
	w.wg.Wait()
}
