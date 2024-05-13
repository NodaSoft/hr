package worker

import (
	"context"
	"fmt"
	"main/task"
	"maps"
	"sync"
	"time"
)

type Worker struct {
	wg                sync.WaitGroup  // Группа ожидания завершения работы
	ctx               context.Context // Контекст для завершения работы
	tasksQueueChan    chan task.Task  // Канал для передачи тасков
	receivedTasksChan chan task.Task  // Канал для получения выполненных тасков
	result            Result          // Результат выполнения тасков
}

type Result struct {
	done   map[int]task.Task // Выполненные таски
	undone map[int]task.Task // Невыполненные таски
	mutex  sync.Mutex        // Мьютекс для защиты карт
}

var ErrSomeErrorOccurred = fmt.Errorf("some error occurred") // Ошибка для эмуляции ошибочных тасков

// NewWorker создает новый экземпляр Worker.
func NewWorker(ctx context.Context) *Worker {
	return &Worker{
		ctx:               ctx,
		tasksQueueChan:    make(chan task.Task, 10),
		receivedTasksChan: make(chan task.Task),
		result: Result{
			done:   make(map[int]task.Task),
			undone: make(map[int]task.Task),
		},
	}
}

// Start запускает воркеры.
func (w *Worker) Start(workers int) {
	w.wg.Add(3 + workers)
	go w.creator()
	for range workers {
		go w.worker()
	}
	go w.sorter()
	go w.writer()
}

// Wait ожидает завершения работы.
func (w *Worker) Wait() {
	defer close(w.receivedTasksChan)
	w.wg.Wait()
}

// creator эмулирует создание тасков.
func (w *Worker) creator() {
	defer close(w.tasksQueueChan)
	defer w.wg.Done()

	for {
		if w.ctx.Err() != nil {
			return
		}

		creationTime := time.Now()

		var err error
		if creationTime.Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			err = ErrSomeErrorOccurred
		}

		select {
		case <-w.ctx.Done():
			return
		case w.tasksQueueChan <- task.Task{Id: creationTime.Nanosecond(), CreatedAt: creationTime, Error: err}:
		}
	}
}

// worker эмулирует выполнение тасков.
func (w *Worker) worker() {
	defer w.wg.Done()

	for {
		select {
		case <-w.ctx.Done():
			return
		case t, ok := <-w.tasksQueueChan:
			if !ok {
				return
			}

			select {
			case <-w.ctx.Done():
				return
			case w.receivedTasksChan <- t.Do():
			}
		}
	}
}

// sorter сортирует выполненные таски.
func (w *Worker) sorter() {
	defer w.wg.Done()

	for {
		select {
		case <-w.ctx.Done():
			return
		case t, ok := <-w.receivedTasksChan:
			if !ok {
				return
			}

			w.result.mutex.Lock()

			if t.Error == nil {
				w.result.done[t.Id] = t
			} else {
				w.result.undone[t.Id] = t
			}

			w.result.mutex.Unlock()
		}
	}
}

// writer выводит результаты выполнения тасков.
func (w *Worker) writer() {
	defer w.wg.Done()

	for {
		select {
		case <-w.ctx.Done():
			return
		case <-time.After(time.Second * 3):
			w.result.mutex.Lock()

			// Копирование карт для вывода
			done := maps.Clone(w.result.done)
			undone := maps.Clone(w.result.undone)

			// Очистка оригинальных карт
			w.result.done = make(map[int]task.Task)
			w.result.undone = make(map[int]task.Task)

			w.result.mutex.Unlock()

			// Вывод результатов
			write("Errors:", undone)
			write("Done tasks:", done)
		}
	}
}

// write выводит результаты выполнения тасков.
func write(title string, tasks map[int]task.Task) {
	fmt.Println(title, len(tasks))
	for _, t := range tasks {
		fmt.Println(t)
	}
}
