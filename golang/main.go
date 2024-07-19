package main

import (
	"context"
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

type Task struct {
	id          int
	createdTime string // время создания
	finishTime  string // время выполнения
	result      []byte
}

type Results struct {
	items map[int]Task
	m     *sync.RWMutex
}

type Errors struct {
	items []error
	m     *sync.RWMutex
}

const (
	WORK_TIME       = time.Second * 10
	CREATE_INTERVAL = time.Millisecond * 150
)

func TaskProducer(ctx context.Context, w *sync.WaitGroup, queueToExecute chan Task) {
	defer func() {
		fmt.Println("Producer is done...")
		close(queueToExecute)
		w.Done()
	}()
	var (
		now         time.Time
		createdTime string
	)
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		default:
			now = time.Now()
			createdTime = now.Format(time.RFC3339Nano)
			if now.UnixMilli()%2 > 0 {
				createdTime = "!!!wrong format!!!"
			}
			queueToExecute <- Task{createdTime: createdTime, id: int(now.Unix())}
			time.Sleep(CREATE_INTERVAL)
		}
	}
}

func TaskWorker(ctx context.Context, w *sync.WaitGroup, queueToExecute chan Task, doneTasks chan Task, undoneTasks chan error) {
	defer fmt.Println("Worker is done...")
	defer close(doneTasks)
	defer close(undoneTasks)
	defer w.Done()
loop:
	for {
		select {
		case task, ok := <-queueToExecute:
			if !ok {
				break loop
			}
			_, err := time.Parse(time.RFC3339Nano, task.createdTime)
			if err != nil {
				undoneTasks <- fmt.Errorf("task id %d, time %s, error: %s", task.id, task.createdTime, err)
			} else {
				task.result = []byte("task has been succeed")
				task.finishTime = time.Now().Format(time.RFC3339Nano)
				doneTasks <- task
			}
		case <-ctx.Done():
			break loop
		}
	}
}

func main() {
	w := sync.WaitGroup{}
	w.Add(5)
	ctx := context.Background()
	ctx, closer := context.WithTimeout(ctx, WORK_TIME)
	defer closer()

	queueToExecute := make(chan Task, 10)

	go TaskProducer(ctx, &w, queueToExecute)

	doneTasks := make(chan Task)
	undoneTasks := make(chan error)
	go TaskWorker(ctx, &w, queueToExecute, doneTasks, undoneTasks)

	results := Results{m: &sync.RWMutex{}, items: make(map[int]Task)}
	errors := Errors{m: &sync.RWMutex{}, items: make([]error, 0)}
	go func() {
		defer func() {
			fmt.Println("Filler for done tasks is done...")
			w.Done()
		}()
		for task := range doneTasks {
			results.m.Lock()
			results.items[task.id] = task
			results.m.Unlock()
		}
	}()
	go func() {
		defer func() {
			fmt.Println("Filler for undone tasks is done...")
			w.Done()
		}()
		for err := range undoneTasks {
			errors.m.Lock()
			errors.items = append(errors.items, err)
			errors.m.Unlock()
		}
	}()
	go func() {
		defer func() {
			fmt.Println("Printer is done...")
			w.Done()
		}()
		ticker := time.NewTicker(time.Second * 3)
	loop:
		for {
			select {
			case <-ctx.Done():
				break loop
			case <-ticker.C:
				fmt.Println("--------------------------------------------------")
				fmt.Println("Errors:")
				errors.m.RLock()
				for _, err := range errors.items {
					fmt.Println(err)
				}
				errors.m.RUnlock()

				fmt.Println("Done tasks:")
				results.m.RLock()
				for res := range results.items {
					fmt.Println(res, ", result: ", string(results.items[res].result))
				}
				results.m.RUnlock()
				fmt.Println("--------------------------------------------------")
			}
		}
	}()
	w.Wait()
	fmt.Println("Finished")
}
