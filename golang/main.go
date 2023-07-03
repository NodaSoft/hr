package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

type Task struct {
	id           uuid.UUID
	creationTime time.Time
}

func taskCreator(ctx context.Context, taskCh chan<- Task) {
	for {
		select {
		case <-ctx.Done():
			close(taskCh)
			return

		default:
			timeNow := time.Now()

			taskCh <- Task{creationTime: timeNow, id: uuid.New()}
			time.Sleep(time.Millisecond * 150)
		}

	}
}

func taskWorker(taskCh <-chan Task, wg *sync.WaitGroup, taskResults *TasksResults) {

	for t := range taskCh {
		switch {
		case t.creationTime.Nanosecond()%2 > 0:
			taskResults.errorsMu.Lock()
			taskResults.errors = append(taskResults.errors, fmt.Errorf("Task id %d time %s, error: Some error occured", t.id, t.creationTime))
			taskResults.errorsMu.Unlock()

		case t.creationTime.After(time.Now().Add(-20 * time.Second)):
			taskResults.doneMu.Lock()
			taskResults.doneTasksIds[t.id] = t
			taskResults.doneMu.Unlock()

		default:
			taskResults.errorsMu.Lock()
			taskResults.errors = append(taskResults.errors, fmt.Errorf("Task id %d time %s, error: Something went wrong", t.id, t.creationTime))
			taskResults.errorsMu.Unlock()
		}

	}

	wg.Done()
}

type TasksResults struct {
	doneTasksIds map[uuid.UUID]Task
	doneMu       sync.Mutex

	errors   []error
	errorsMu sync.Mutex
}

func main() {
	var (
		wg           sync.WaitGroup
		ctx, cancel  = context.WithTimeout(context.Background(), time.Second*2)
		workersCount = 10
		tasksCh      = make(chan Task, 10)
	)

	defer cancel()

	taskResults := TasksResults{
		doneTasksIds: make(map[uuid.UUID]Task),
	}

	go taskCreator(ctx, tasksCh)

	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go taskWorker(tasksCh, &wg, &taskResults)
	}

	wg.Wait()

	println("Errors:")
	for r := range taskResults.errors {
		println(r)
	}

	println("Done tasks:")
	for r := range taskResults.doneTasksIds {
		fmt.Println(r)
	}
}
