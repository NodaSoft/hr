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

type result string

const (
	unknown            result = ""
	succeeded          result = "succeeded"
	errorOccurred      result = "some error occurred"
	somethingWentWrong result = "something went wrong"
)

// Task represents a meaninglessness of our life
type Task struct {
	id            int
	creationTime  time.Time // время создания
	executionTime time.Time // время выполнения
	result
}

func taskCreaturer(ctx context.Context, ch chan Task) {
	for {
		select {
		case <-ctx.Done():
			close(ch)
			return
		default:
			t := Task{
				id:           int(time.Now().Unix()),
				creationTime: time.Now().UTC(),
			}

			if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				t.result = errorOccurred
			}
			ch <- t // передаем таск на выполнение
			time.Sleep(time.Millisecond * 150)
		}
	}
}

func taskWorker(t Task) Task {
	if t.result == unknown && t.creationTime.After(time.Now().UTC().Add(time.Second*-20)) {
		t.result = succeeded
	} else if t.result == unknown {
		t.result = somethingWentWrong
	}
	t.executionTime = time.Now().UTC()
	return t
}

func taskSorter(t Task, doneTasks chan Task, undoneTasks chan error) {
	if t.result == succeeded {
		doneTasks <- t
	} else {
		undoneTasks <- fmt.Errorf("task id %d time %s, error %s", t.id, t.creationTime.Format(time.RFC3339), t.result)
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	superChan := make(chan Task, 10)

	go taskCreaturer(ctx, superChan)

	doneTasks := make(chan Task)
	undoneTasks := make(chan error)

	go func() {
		// получение тасков
		for t := range superChan {
			t = taskWorker(t)
			taskSorter(t, doneTasks, undoneTasks)
		}
		close(doneTasks)
		close(undoneTasks)
	}()

	resultDoneTasks := map[int]Task{}
	var errs []error

	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		defer wg.Done()
		for t := range doneTasks {
			resultDoneTasks[t.id] = t
		}
	}()
	go func() {
		defer wg.Done()
		for err := range undoneTasks {
			errs = append(errs, err)
		}
	}()
	wg.Wait()

	println("Errors:")
	for _, err := range errs {
		println(err.Error())
	}

	println("Done tasks:")
	for id := range resultDoneTasks {
		println(id)
	}
}
