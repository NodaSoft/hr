package main

import (
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
	errTask = fmt.Errorf("something went wrong")

	timeFormat   = time.RFC3339Nano
	timeToWork   = time.Second * 10
	timeInterval = time.Second * 3
)

// A Ttype represents a meaninglessness of our life
type Task struct {
	Id           int
	CreationTime string // время создания
	FinishTime   string // время выполнения
	Result       []byte
	Error        error
}

type Work struct {
	newTaskQueue    chan Task
	doneTasksQueue  chan Task
	undoneTaskQueue chan Task
	wg              sync.WaitGroup
	done            chan struct{}
}

func New() *Work {
	return &Work{
		wg:   sync.WaitGroup{},
		done: make(chan struct{}),
	}
}

// Run запускает всю логику обработки тасков tasks
func (w *Work) Run(tasks chan Task) {
	w.newTaskQueue = tasks
	w.startWork()
	// w.PrintAll()
	w.PrintNew()

	w.wg.Wait()

}

func (w *Work) startWork() {
	w.doneTasksQueue = make(chan Task)
	w.undoneTaskQueue = make(chan Task)
	go func() {
		for t := range w.newTaskQueue {
			t = taskWorker(t)
			w.sortTask(t)
		}

		close(w.doneTasksQueue)
		close(w.undoneTaskQueue)
	}()
}

// sortTask сортирует таски
func (w *Work) sortTask(t Task) {
	if t.Error != nil {
		w.undoneTaskQueue <- t
		return
	}
	w.doneTasksQueue <- t
}

// PrintAll выводить в консоль результат всех обработанных к этому моменту тасков
// (отдельно успешные и отдельно с ошибками).
func (w *Work) PrintAll() {
	doneTaskClosed := false
	undoneTaskClosed := false
	errMu := sync.RWMutex{}
	doneMu := sync.RWMutex{}
	result := []Task{}
	errors := []error{}

	print := func() {
		fmt.Println("Errors:")
		errMu.RLock()
		for _, e := range errors {
			fmt.Println(e)
		}
		errMu.RUnlock()

		fmt.Println("Done tasks:")
		doneMu.RLock()
		for _, r := range result {
			fmt.Printf("id: %v, start: %v, finish: %v result: %s\n", r.Id, r.CreationTime, r.FinishTime, r.Result)
		}
		doneMu.RUnlock()
	}

	w.wg.Add(2)
	go func() {
		for t := range w.doneTasksQueue {
			doneMu.Lock()
			result = append(result, t)
			doneMu.Unlock()
		}
		doneTaskClosed = true
		w.wg.Done()
	}()
	go func() {
		for t := range w.undoneTaskQueue {
			errMu.Lock()
			errors = append(errors, fmt.Errorf("task id %d, undoneTaskstime %s, error %w", t.Id, t.CreationTime, t.Error))
			errMu.Unlock()
		}
		undoneTaskClosed = true
		w.wg.Done()
	}()

	w.wg.Add(1)
	go func() {
		tick := time.NewTicker(timeInterval)
		for {
			select {
			case <-tick.C:
				print()
			default:
				if doneTaskClosed && undoneTaskClosed {
					print()
					w.wg.Done()
					return
				}
			}
		}
	}()
}

// PrintNew выводить в консоль результат новых обработанных к этому моменту тасков
// (отдельно успешные и отдельно с ошибками).
//
// На тот случай, если я неправильно понял задание вывода в консоль
func (w *Work) PrintNew() {
	doneTaskClosed := false
	undoneTaskClosed := false
	errMu := sync.Mutex{}
	doneMu := sync.Mutex{}
	result := []Task{}
	errors := []error{}

	print := func() {
		fmt.Println("Errors:")
		errMu.Lock()
		for _, e := range errors {
			fmt.Println(e)
		}
		errors = errors[:0]
		errMu.Unlock()

		fmt.Println("Done tasks:")
		doneMu.Lock()
		for _, r := range result {
			fmt.Printf("id: %v, start: %v, finish: %v result: %s\n", r.Id, r.CreationTime, r.FinishTime, r.Result)
		}
		result = result[:0]
		doneMu.Unlock()
	}

	w.wg.Add(2)
	go func() {
		for t := range w.doneTasksQueue {
			doneMu.Lock()
			result = append(result, t)
			doneMu.Unlock()
		}
		doneTaskClosed = true
		w.wg.Done()
	}()
	go func() {
		for t := range w.undoneTaskQueue {
			errMu.Lock()
			errors = append(errors, fmt.Errorf("task id: %d, undoneTaskstime: %s, error: %w", t.Id, t.CreationTime, t.Error))
			errMu.Unlock()
		}
		undoneTaskClosed = true
		w.wg.Done()
	}()

	w.wg.Add(1)
	go func() {
		tick := time.NewTicker(timeInterval)
		for {
			select {
			case <-tick.C:
				print()
			default:
				if doneTaskClosed && undoneTaskClosed {
					print()
					w.wg.Done()
					return
				}
			}
		}
	}()
}

// Генератора тасков.
// по условию 10 секунд генерирует таски и закрывает канал.
func Generator() chan Task {
	out := make(chan Task, 10)
	ticker := time.NewTicker(timeToWork)
	id := 1
	go func() {
		for {
			select {
			case <-ticker.C:
				close(out)
				return
			default:
				// Важно сохранить логику появления ошибочных тасков. Поэтому ничего не трогал
				startTime := time.Now().Format(timeFormat)
				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					startTime = "Some error occured"
				}
				out <- Task{CreationTime: startTime, Id: id} // передаем таск на выполнение
				id++                                         // предыдущая могла присвоить одинаковые айдишники
			}
		}
	}()

	return out
}

// имитация работы
func taskWorker(task Task) Task {
	tt, _ := time.Parse(timeFormat, task.CreationTime)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		task.Result = []byte("task has been successed")
	} else {
		task.Error = errTask
	}
	task.FinishTime = time.Now().Format(timeFormat)

	time.Sleep(time.Millisecond * 150)

	return task
}

func main() {
	w := New()
	w.Run(Generator())
}
