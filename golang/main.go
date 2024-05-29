package main

import (
	"context"
	"errors"
	"fmt"
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

// A task represents a meaninglessness of our life
type task struct {
	id         int
	createdAt  time.Time // время создания
	finishedAt time.Time // время выполнения
	result     []byte
	err        error
}

func taskCreator(ctx context.Context, taskCh chan<- task) {
	defer close(taskCh)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			now := time.Now()
			t := task{id: int(time.Now().Unix()), createdAt: now}
			if now.Nanosecond()%2 > 0 {
				t.err = errors.New("some error occured")
			}
			taskCh <- t
			time.Sleep(time.Millisecond * 150) // Ограничиваем число созданных задач
		}
	}
}

func processTask(t task) task {
	if t.err != nil {
		return t
	}

	now := time.Now()
	if t.createdAt.Before(now.Add(-20 * time.Second)) {
		t.err = errors.New("something went wrong")
		return t
	}

	t.result = []byte("task has been successful")
	t.finishedAt = now

	return t
}

func sortTask(doneCh chan<- task, errCh chan<- error, t task) {
	if t.err != nil {
		errCh <- fmt.Errorf("task id %d time %s, error %w", t.id, t.createdAt, t.err)
	} else {
		doneCh <- t
	}
}

func taskWorker(taskCh <-chan task, doneCh chan<- task, errCh chan<- error) {
	defer close(doneCh)
	defer close(errCh)

	for t := range taskCh {
		t = processTask(t)
		sortTask(doneCh, errCh, t)
	}
}

func printTasks(results map[int]task, errs []error) {
	fmt.Println("Errors:")
	for _, e := range errs {
		fmt.Println(e)
	}
	fmt.Println("Done tasks:")
	for k := range results {
		fmt.Println(k)
	}
}

func main() {
	// Контекст, останавливающий создание задач через 10 секунд
	creatorCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Каналы будут закрыты внутри taskCreator и taskWorker
	superChan := make(chan task, 10)
	doneTasks := make(chan task)
	undoneTasks := make(chan error)

	go taskCreator(creatorCtx, superChan)
	go taskWorker(superChan, doneTasks, undoneTasks)

	results := map[int]task{}
	errs := []error{}

	// Выводим задачи каждые 3 секунды
	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			printTasks(results, errs)
			if doneTasks == nil && undoneTasks == nil {
				return
			}
		case t, ok := <-doneTasks:
			if ok {
				results[t.id] = t
			} else {
				doneTasks = nil
			}
		case e, ok := <-undoneTasks:
			if ok {
				errs = append(errs, e)
			} else {
				undoneTasks = nil
			}
		}
	}
}
